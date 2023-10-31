# awsutils/iam

The `iam` package is a collection of utility functions
designed to simplify common iam tasks.

---

## Table of contents

- [Functions](#functions)
- [Installation](#installation)
- [Usage](#usage)
- [Tests](#tests)
- [Contributing](#contributing)
- [License](#license)

---

## Functions

### GetIamInstanceProfile(string)

```go
GetIamInstanceProfile(string) *iam.InstanceProfile, error
```

GetIamInstanceProfile retrieves the IAM instance profile by its name

---

## Installation

To use the awsutils/iam package, you first need to install it.
Follow the steps below to install via go get.

```bash
go get github.com/l50/awsutils/iam
```

---

## Usage

After installation, you can import the package in your Go project
using the following import statement:

```go
import "github.com/l50/awsutils/iam"
```

---

## Tests

To ensure the package is working correctly, run the following
command to execute the tests for `awsutils/iam`:

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
