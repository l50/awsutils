# awsutils/secretsmanager

The `secretsmanager` package is a collection of utility functions
designed to simplify common secretsmanager tasks.

Table of contents:

- [Functions](#functions)
- [Installation](#installation)
- [Usage](#usage)
- [Tests](#tests)
- [Contributing](#contributing)
- [License](#license)

---

## Functions

### CreateConnection()

```go
CreateConnection() Connection
```

CreateConnection creates a connection
with Secrets Manager and returns it.

---

### CreateOrUpdateSecret(*secretsmanager.SecretsManager, string, string, string)

```go
CreateOrUpdateSecret(*secretsmanager.SecretsManager string string string) error
```

CreateOrUpdateSecret creates a new secret or updates an existing one.

---

### CreateSecret(*secretsmanager.SecretsManager, string, string, string)

```go
CreateSecret(*secretsmanager.SecretsManager, string, string, string) error
```

CreateSecret creates an input `secretName`
with the specified `secretValue`.

---

### DeleteSecret(*secretsmanager.SecretsManager, string, bool)

```go
DeleteSecret(*secretsmanager.SecretsManager, string, bool) error
```

DeleteSecret deletes an input `secretName`.
It will attempt to do so forcefully if `forceDelete`
is set to true.

---

### GetSecret(*secretsmanager.SecretsManager, string)

```go
GetSecret(*secretsmanager.SecretsManager, string) string, error
```

GetSecret returns the value of an input `secretName`.

---

### ReplicateSecret(Connection, string, string, []string)

```go
ReplicateSecret(Connection, string, string, []string) error
```

ReplicateSecret replicates a secret with the specified `secretName`
to multiple target regions.

---

### UpdateSecret(*secretsmanager.SecretsManager, string, string)

```go
UpdateSecret(*secretsmanager.SecretsManager, string, string) error
```

UpdateSecret updates an existing secret

---

## Installation

To use the awsutils/secretsmanager package, you first need to install it.
Follow the steps below to install via go get.

```bash
go get github.com/awsutils/l50/secretsmanager
```

---

## Usage

After installation, you can import the package in your Go project
using the following import statement:

```go
import "github.com/awsutils/l50/secretsmanager"
```

---

## Tests

To ensure the package is working correctly, run the following
command to execute the tests for `secretsmanager/awsutils`:

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
