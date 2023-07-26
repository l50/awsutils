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

AgentReady checks if an SSM agent is ready on the instance.

**Parameters:**

svc: AWS SSM service client.
instanceID: AWS EC2 instance ID to check.
waitTime: Maximum wait time before timing out.

**Returns:**

bool: True if the agent is ready, false otherwise.
error: An error if any issue occurs while checking the agent.

---

### CheckAWSCLIInstalled(ssmiface.SSMAPI, string)

```go
CheckAWSCLIInstalled(ssmiface.SSMAPI, string) bool, error
```

CheckAWSCLIInstalled checks if AWS CLI is installed on the instance.

**Parameters:**

svc: AWS SSM service client.
instanceID: AWS EC2 instance ID to check.

**Returns:**

bool: True if AWS CLI is installed, false otherwise.
error: An error if any issue occurs while checking the installation.

---

### CreateConnection()

```go
CreateConnection() Connection
```

CreateConnection establishes a connection with AWS SSM.

**Returns:**

Connection: Struct with a connected SSM client and session.

---

### DeleteParam(ssmiface.SSMAPI, string)

```go
DeleteParam(ssmiface.SSMAPI, string) error
```

DeleteParam removes a parameter from AWS SSM.

**Parameters:**

svc: AWS SSM service client.
name: Name of the parameter to delete.

**Returns:**

error: An error if any issue occurs while deleting the parameter.

---

### GetParam(ssmiface.SSMAPI, string)

```go
GetParam(ssmiface.SSMAPI, string) string, error
```

GetParam retrieves a parameter from AWS SSM.

**Parameters:**

svc: AWS SSM service client.
name: Name of the parameter.

**Returns:**

string: Value of the parameter.
error: An error if any issue occurs while fetching the parameter.

---

### ListAllParameters(ssmiface.SSMAPI)

```go
ListAllParameters(ssmiface.SSMAPI) []*ssm.ParameterMetadata, error
```

ListAllParameters retrieves all parameters in the AWS SSM.

**Parameters:**

svc: AWS SSM service client.

**Returns:**

([]*ssm.ParameterMetadata): List of all parameters' metadata.
error: An error if any issue occurs while fetching the parameters.

---

### PutParam(ssmiface.SSMAPI, string, string, string, bool)

```go
PutParam(ssmiface.SSMAPI, string, string, string, bool) error
```

PutParam creates or updates a parameter in AWS SSM.

**Parameters:**

svc: AWS SSM service client.
name: Name of the parameter.
value: Value of the parameter.
paramType: Type of the parameter.
overwrite: Flag to overwrite an existing parameter.

**Returns:**

error: An error if any issue occurs while creating or updating the parameter.

---

### RunCommand(ssmiface.SSMAPI, string, []string)

```go
RunCommand(ssmiface.SSMAPI, string, []string) string, error
```

RunCommand executes an input command on an AWS instance via SSM.

**Parameters:**

svc: AWS SSM service client.
instanceID: AWS EC2 instance ID where the command should run.
command: List of command strings to be run.

**Returns:**

string: Output of the command execution.
error: An error if any issue occurs while executing the command.

---

## Installation

To use the awsutils/ssm package, you first need to install it.
Follow the steps below to install via go get.

```bash
go get github.com/l50/awsutils/ssm
```

---

## Usage

After installation, you can import the package in your Go project
using the following import statement:

```go
import "github.com/l50/awsutils/ssm"
```

---

## Tests

To ensure the package is working correctly, run the following
command to execute the tests for `ssm/awsutils`:

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
