# awsutils/ssm

The `ssm` package is a collection of utility functions
designed to simplify common ssm tasks.

Table of contents:

- [Functions](#functions)
- [Installation](#installation)
- [Usage](#usage)
- [Tests](#tests)
- [Contributing](#contributing)
- [License](#license)

---

## Functions

### AgentReady(ssmiface.SSMAPI, string, time.Duration)

```go
AgentReady(ssmiface.SSMAPI, string, time.Duration) bool, error
```

AgentReady checks if an SSM agent is ready.
Inputs:

    svc is an Amazon SSM service client
    name is the name of the parameter

Output:

    If success, return true and nil
    Otherwise, return false and an error from the call to DescribeInstanceInformation

---

### CheckAWSCLIInstalled(ssmiface.SSMAPI, string)

```go
CheckAWSCLIInstalled(ssmiface.SSMAPI, string) bool, error
```

CheckAWSCLIInstalled checks if AWS CLI is installed on the instance.
Inputs:

    svc is an Amazon SSM service client
    instanceID is the instance to check

Output:

    If successful, return true and nil. If AWS CLI is not installed or an error occurred, return false and the error.

---

### CreateConnection()

```go
CreateConnection() Connection
```

CreateConnection creates a connection
with SSM and returns it.

---

### DeleteParam(ssmiface.SSMAPI, string)

```go
DeleteParam(ssmiface.SSMAPI, string) error
```

DeleteParam deletes a parameter in SSM
Inputs:

    svc is an Amazon SSM service client
    name is the name of the parameter

Output:

    If success, information about the parameter and nil
    Otherwise, nil and an error from the call to DeleteParam

---

### GetParam(ssmiface.SSMAPI, string)

```go
GetParam(ssmiface.SSMAPI, string) string, error
```

GetParam fetches details of a parameter in SSM
Inputs:

    svc is an Amazon SSM service client
    name is the name of the parameter
    value is the value of the parameter
    paramType is the type of parameter

Output:

    If success, information about the parameter and nil
    Otherwise, nil and an error from the call to GetParam

---

### PutParam(ssmiface.SSMAPI, string, string, string, bool)

```go
PutParam(ssmiface.SSMAPI, string, string, string, bool) error
```

PutParam creates a parameter in SSM
Inputs:

    svc is an Amazon SSM service client
    name is the name of the parameter
    value is the value of the parameter
    type is the type of parameter
    overwrite sets the flag to rewrite
    a parameter value

Output:

    If success, information about the parameter and nil
    Otherwise, nil and an error from the call to PutParam

---

### RunCommand(ssmiface.SSMAPI, string, []string)

```go
RunCommand(ssmiface.SSMAPI, string, []string) string, error
```

RunCommand runs an input command using SSM.
Inputs:

    svc is an Amazon SSM service client
    instanceID is the instance to run the command on
    command is the command to run

Output:

    If successful, the command output and nil will be returned.

---

## Installation

To use the awsutils/ssm package, you first need to install it.
Follow the steps below to install via go get.

```bash
go get github.com/awsutils/l50/ssm
```

---

## Usage

After installation, you can import the package in your Go project
using the following import statement:

```go
import "github.com/awsutils/l50/ssm"
```

---

## Tests

To ensure the package is working correctly, run the following
command to execute the tests for `awsutils/ssm`:

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
