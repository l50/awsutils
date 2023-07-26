# awsutils/main

The `main` package is a collection of utility functions
designed to simplify common main tasks.

Table of contents:

- [Functions](#functions)
- [Installation](#installation)
- [Usage](#usage)
- [Tests](#tests)
- [Contributing](#contributing)
- [License](#license)

---

## Functions

### FindExportedFuncsWithoutTests(string)

```go
FindExportedFuncsWithoutTests(string) []string, error
```

FindExportedFuncsWithoutTests finds exported functions without tests

---

### GeneratePackageDocs()

```go
GeneratePackageDocs() error
```

GeneratePackageDocs generates package documentation
for packages in the current directory and its subdirectories.

---

### InstallDeps()

```go
InstallDeps() error
```

InstallDeps Installs go dependencies

---

### RunPreCommit()

```go
RunPreCommit() error
```

RunPreCommit runs all pre-commit hooks locally

---

### RunTests(string)

```go
RunTests(string) error
```

RunTests runs all of the unit tests

---

### UpdateDocs()

```go
UpdateDocs() error
```

UpdateDocs updates the package documentation
for packages in the current directory and its subdirectories.

---

### UpdateMirror(string)

```go
UpdateMirror(string) error
```

UpdateMirror updates pkg.go.dev and proxy.golang.org with the
release associated with the input tag

---

### UseFixCodeBlocks(string, string)

```go
UseFixCodeBlocks(string, string) error
```

UseFixCodeBlocks fixes code blocks for the input filepath
using the input language.

**Parameters:**

filepath: the path to the file or directory to fix
language: the language of the code blocks to fix

**Returns:**

error: an error if one occurred

Example:

```go
mage fixcodeblocks docs/docGeneration.go go
```

---

## Installation

To use the awsutils/main package, you first need to install it.
Follow the steps below to install via go get.

```bash
go get github.com/awsutils/l50/main
```

---

## Usage

After installation, you can import the package in your Go project
using the following import statement:

```go
import "github.com/awsutils/l50/main"
```

---

## Tests

To ensure the package is working correctly, run the following
command to execute the tests for `main/awsutils`:

```bash
go test -v
```

---

## Contributing

Pull requests are welcome. For major changes,
please open an issue first to discuss what
you would like to change.

---

## License

This project is licensed under the MIT
License - see the [LICENSE](../LICENSE)
file for details.
