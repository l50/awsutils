//go:build mage

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/l50/goutils/v2/dev/lint"
	mageutils "github.com/l50/goutils/v2/dev/mage"
	"github.com/l50/goutils/v2/docs"
	"github.com/l50/goutils/v2/git"
	"github.com/l50/goutils/v2/sys"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/spf13/afero"
)

func init() {
	os.Setenv("GO111MODULE", "on")
}

// InstallDeps installs the Go dependencies necessary for developing
// on the project.
//
// Example usage:
//
// ```go
// mage installdeps
// ```
//
// **Returns:**
//
// error: An error if any issue occurs while trying to
// install the dependencies.
func InstallDeps() error {
	fmt.Println(color.YellowString("Running go mod tidy."))

	if err := mageutils.Tidy(); err != nil {
		return fmt.Errorf("failed to install dependencies: %v", err)
	}

	fmt.Println(color.YellowString("Installing dependencies."))
	if err := lint.InstallGoPCDeps(); err != nil {
		return fmt.Errorf("failed to install pre-commit dependencies: %v", err)
	}

	if err := mageutils.InstallVSCodeModules(); err != nil {
		return fmt.Errorf(color.RedString(
			"failed to install vscode-go modules: %v", err))
	}

	return nil
}

// FindExportedFuncsWithoutTests finds exported functions without tests
func FindExportedFuncsWithoutTests(pkg string) ([]string, error) {
	funcs, err := mageutils.FindExportedFuncsWithoutTests(os.Args[1])

	if err != nil {
		log.Fatalf("failed to find exported functions without tests: %v", err)
	}

	for _, funcName := range funcs {
		fmt.Println(funcName)
	}

	return funcs, nil

}

// GeneratePackageDocs creates documentation for the various packages
// in the project.
//
// Example usage:
//
// ```go
// mage generatepackagedocs
// ```
//
// **Returns:**
//
// error: An error if any issue occurs during documentation generation.
func GeneratePackageDocs() error {
	fs := afero.NewOsFs()

	repoRoot, err := git.RepoRoot()
	if err != nil {
		return fmt.Errorf("failed to get repo root: %v", err)
	}
	sys.Cd(repoRoot)

	repo := docs.Repo{
		Owner: "l50",
		Name:  "awsutils",
	}

	templatePath := filepath.Join("templates", "README.md.tmpl")
	if err := docs.CreatePackageDocs(fs, repo, templatePath); err != nil {
		return fmt.Errorf("failed to create package docs: %v", err)
	}

	return nil
}

// RunPreCommit updates, clears, and executes all pre-commit hooks
// locally. The function follows a three-step process:
//
// First, it updates the pre-commit hooks.
// Next, it clears the pre-commit cache to ensure a clean environment.
// Lastly, it executes all pre-commit hooks locally.
//
// Example usage:
//
// ```go
// mage runprecommit
// ```
//
// **Returns:**
//
// error: An error if any issue occurs at any of the three stages
// of the process.
func RunPreCommit() error {
	if !sys.CmdExists("pre-commit") {
		return fmt.Errorf("pre-commit is not installed, please follow the " +
			"instructions in the dev doc: " +
			"https://github.com/facebookincubator/TTPForge/tree/main/docs/dev")
	}

	fmt.Println(color.YellowString("Updating pre-commit hooks."))
	if err := lint.UpdatePCHooks(); err != nil {
		return err
	}

	fmt.Println(color.YellowString("Clearing the pre-commit cache to ensure we have a fresh start."))
	if err := lint.ClearPCCache(); err != nil {
		return err
	}

	fmt.Println(color.YellowString("Running all pre-commit hooks locally."))
	if err := lint.RunPCHooks(); err != nil {
		return err
	}

	return nil
}

// RunTests runs all of the unit tests
func RunTests(testType string) error {
	if testType == "" {
		testType = "coverage"
	}
	fmt.Println("Running unit tests.")
	if err := sh.RunV(filepath.Join(".hooks", "go-unit-tests.sh"), testType); err != nil {
		return fmt.Errorf("failed to run unit tests: %v", err)
	}

	return nil
}

// UpdateMirror updates pkg.go.dev with the release associated with the
// input tag
//
// Example usage:
//
// ```go
// mage updatemirror v2.0.1
// ```
//
// **Parameters:**
//
// tag: the tag to update pkg.go.dev with
//
// **Returns:**
//
// error: An error if any issue occurs while updating pkg.go.dev
func UpdateMirror(tag string) error {
	var err error
	fmt.Printf("Updating pkg.go.dev with the new tag %s.", tag)

	repo := docs.Repo{
		Owner: "l50",
		Name:  "awsutils",
	}

	err = sh.RunV("curl", "--silent", fmt.Sprintf(
		"https://sum.golang.org/lookup/github.com/%s/%s@%s",
		repo.Owner, repo.Name, tag))
	if err != nil {
		return fmt.Errorf("failed to update proxy.golang.org: %w", err)
	}

	err = sh.RunV("curl", "--silent", fmt.Sprintf(
		"https://proxy.golang.org/github.com/%s/%s/@v/%s.info",
		repo.Owner, repo.Name, tag))
	if err != nil {
		return fmt.Errorf("failed to update pkg.go.dev: %w", err)
	}

	return nil
}

// UpdateDocs updates the package documentation
// for packages in the current directory and its subdirectories.
func UpdateDocs() error {
	repo := docs.Repo{
		Owner: "l50",
		Name:  "awsutils",
	}

	fs := afero.NewOsFs()

	templatePath := filepath.Join("templates", "README.md.tmpl")

	if err := docs.CreatePackageDocs(fs, repo, templatePath); err != nil {
		return fmt.Errorf("failed to update docs: %v", err)
	}

	return nil
}

// Run runs the unit tests and extracts failing functions and their tests.
func Run() error {
	reader := bufio.NewReader(os.Stdin)

	// Select package
	fmt.Println("Please select a package to test:")
	packages, _ := listPackages()
	// Append option for running all tests
	packages = append(packages, "Run all tests")

	for i, pkg := range packages {
		fmt.Printf("[%d] %s\n", i, pkg)
	}
	fmt.Print("Enter the number of the package: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)     // remove trailing newline
	inputNum, err := strconv.Atoi(input) // convert string to integer
	if err != nil {
		fmt.Println("Invalid input, please enter a number.")
		return err
	}

	if inputNum == len(packages)-1 { // Check if user selected last option (run all tests)
		// Select test type
		fmt.Println("Please select a test type:")
		testTypes := []string{"all", "coverage", "short"}
		for i, testType := range testTypes {
			fmt.Printf("[%d] %s\n", i, testType)
		}
		fmt.Print("Enter the number of the test type: ")
		input, _ = reader.ReadString('\n')
		input = strings.TrimSpace(input)
		inputNum, err = strconv.Atoi(input)
		if err != nil {
			fmt.Println("Invalid input, please enter a number.")
			return err
		}
		selectedTestType := testTypes[inputNum]

		mg.SerialDeps(RunTests(selectedTestType))
	} else {
		selectedPackage := packages[inputNum]
		mg.SerialDeps(runTests(selectedPackage), extractFailedFunctions(selectedPackage))
	}
	return nil
}

// extractFunctionName extracts the function name from a test output line.
func extractFunctionName(line string) string {
	// Assuming line is of the format: `=== RUN   TestFunctionName`
	parts := strings.Split(line, " ")
	if len(parts) >= 4 {
		return parts[3]
	}
	return ""
}

// findTestFunction finds the test function for a given function name.
func findTestFunction(functionName string) string {
	// Assuming each function `Foo` has a corresponding test function `TestFoo`
	return "Test" + functionName
}

// runTests executes go test.
func runTests(pkg string) error {
	fmt.Printf("Running tests for package %s...\n", pkg)
	if err := sh.Run("go", "test", "-v", pkg); err != nil {
		return fmt.Errorf("tests failed: %w", err)
	}
	return nil
}

// listPackages returns a slice of all package paths.
func listPackages() ([]string, error) {
	cmd := exec.Command("go", "list", "./...")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("unable to list packages: %w", err)
	}
	packages := strings.Split(string(output), "\n")

	// Filter out empty strings or strings containing only white spaces
	var validPackages []string
	for _, pkg := range packages {
		if strings.TrimSpace(pkg) != "" {
			validPackages = append(validPackages, pkg)
		}
	}
	return validPackages, nil
}

// TestEvent represents a test event.
type TestEvent struct {
	Action  string
	Package string
	Test    string
	Output  string
}

// extractFailedFunctions parses the test output and extracts failing functions and their tests.
func extractFailedFunctions(pkg string) error {
	fmt.Println("Extracting failed functions...")

	// Run the test with -json flag to parse the output easily.
	cmd := exec.Command("go", "test", "-json", pkg)
	output, _ := cmd.CombinedOutput()

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		var event TestEvent
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			return fmt.Errorf("unable to unmarshal test event: %w", err)
		}

		if event.Action == "fail" {
			fmt.Println("Failed function: ", event.Test)
		}
	}
	return nil
}
