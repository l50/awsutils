# awsutils/dynamo

The `dynamo` package is a collection of utility functions
designed to simplify common dynamo tasks.

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
with DynamoDB and returns it.

---

### CreateTable(*dynamodb.DynamoDB, string)

```go
CreateTable(*dynamodb.DynamoDB, string) error
```

CreateTable creates a table with the input
`tableName`.

---

### DestroyTable(*dynamodb.DynamoDB, string)

```go
DestroyTable(*dynamodb.DynamoDB, string) error
```

DestroyTable destroys the input table.

---

### GetRegion(*dynamodb.DynamoDB)

```go
GetRegion(*dynamodb.DynamoDB) string, error
```

GetRegion returns the region associated with the input
dynamo client.

---

### GetTables(*dynamodb.DynamoDB)

```go
GetTables(*dynamodb.DynamoDB) []*string, error
```

GetTables returns all dynamoDB tables that the
input client has access to.
Resource:
https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/dynamo-example-list-tables.html

---

### WaitForTable(*dynamodb.DynamoDB, string)

```go
WaitForTable(*dynamodb.DynamoDB, string) error
```

WaitForTable waits for the creation process of the
input table to finish.

---

## Installation

To use the awsutils/dynamo package, you first need to install it.
Follow the steps below to install via go get.

```bash
go get github.com/awsutils/l50/dynamo
```

---

## Usage

After installation, you can import the package in your Go project
using the following import statement:

```go
import "github.com/awsutils/l50/dynamo"
```

---

## Tests

To ensure the package is working correctly, run the following
command to execute the tests for `dynamo/awsutils`:

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
