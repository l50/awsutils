# awsutils/magefiles

`magefiles` provides utilities that would normally be managed
and executed with a `Makefile`. Instead of being written in the make language,
magefiles are crafted in Go and leverage the [Mage](https://magefile.org/) library.

---

## Table of contents

- [Functions](#functions)
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

GeneratePackageDocs creates documentation for the various packages
in the project.

Example usage:

```go
mage generatepackagedocs
```

**Returns:**

error: An error if any issue occurs during documentation generation.

---

### InstallDeps()

```go
InstallDeps() error
```

InstallDeps installs the Go dependencies necessary for developing
on the project.

Example usage:

```go
mage installdeps
```

**Returns:**

error: An error if any issue occurs while trying to
install the dependencies.

---

### Run()

```go
Run() error
```

Run runs the unit tests and extracts failing functions and their tests.

---

### RunPreCommit()

```go
RunPreCommit() error
```

RunPreCommit updates, clears, and executes all pre-commit hooks
locally. The function follows a three-step process:

First, it updates the pre-commit hooks.
Next, it clears the pre-commit cache to ensure a clean environment.
Lastly, it executes all pre-commit hooks locally.

Example usage:

```go
mage runprecommit
```

**Returns:**

error: An error if any issue occurs at any of the three stages
of the process.

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

UpdateMirror updates pkg.go.dev with the release associated with the
input tag

Example usage:

```go
mage updatemirror v2.0.1
```

**Parameters:**

tag: the tag to update pkg.go.dev with

**Returns:**

error: An error if any issue occurs while updating pkg.go.dev

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
