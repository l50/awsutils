# awsutils/ec2

The `ec2` package is a collection of utility functions
designed to simplify common ec2 tasks.

Table of contents:

- [Functions](#functions)
- [Installation](#installation)
- [Usage](#usage)
- [Tests](#tests)
- [Contributing](#contributing)
- [License](#license)

---

## Functions

### CheckInstanceExists(\*ec2.EC2, string)

```go
CheckInstanceExists(*ec2.EC2, string) error
```

CheckInstanceExists checks if an EC2 instance with the given instance ID exists.

---

### CreateConnection()

```go
CreateConnection() Connection
```

CreateConnection creates a connection
with EC2 and returns a Connection.

---

### CreateInstance(\*ec2.EC2, Params)

```go
CreateInstance(*ec2.EC2, Params) *ec2.Reservation, error
```

CreateInstance returns an ec2 reservation for an instance
that is created with the input ec2Params.

---

### DestroyInstance(\*ec2.EC2, string)

```go
DestroyInstance(*ec2.EC2, string) error
```

DestroyInstance terminates the ec2 instance associated with
the input instanceID.

---

### GetInstanceID(\*ec2.Instance)

```go
GetInstanceID(*ec2.Instance) string
```

GetInstanceID returns the instance ID
from an input instanceReservation.

---

### GetInstancePublicIP(\*ec2.EC2, string)

```go
GetInstancePublicIP(*ec2.EC2, string) string, error
```

GetInstancePublicIP returns the public IP address
of the input instanceID.

---

### GetInstanceState(\*ec2.EC2, string)

```go
GetInstanceState(*ec2.EC2, string) string, error
```

GetInstanceState returns the state of the ec2
instance associated with the input instanceID.

---

### GetInstances(*ec2.EC2, []*ec2.Filter)

```go
GetInstances(*ec2.EC2, []*ec2.Filter) []*ec2.Instance, error
```

GetInstances returns ec2 instances that the
input client has access to.
If no filters are provided, all ec2 instances will
be returned by default.

---

### GetInstancesRunningForMoreThan24Hours(\*ec2.EC2)

```go
GetInstancesRunningForMoreThan24Hours(*ec2.EC2) []*ec2.Instance, error
```

GetInstancesRunningForMoreThan24Hours returns a list of all EC2 instances running
for more than 24 hours.

---

### GetLatestAMI(AMIInfo)

```go
GetLatestAMI(AMIInfo) string, error
```

GetLatestAMI retrieves the latest Amazon Machine Image (AMI) for a
specified distribution, version and architecture. It utilizes AWS SDK
to query AWS EC2 for the AMIs matching the provided pattern and returns
the latest one based on the creation date.

**Parameters:**

info: An AMIInfo struct containing necessary details like Distro,
Version, Architecture, and Region for which the AMI needs to be retrieved.

**Returns:**

string: The ID of the latest AMI found based on the provided information.

error: An error if any issue occurs while trying to get the latest AMI.

---

### GetRegion(\*ec2.EC2)

```go
GetRegion(*ec2.EC2) string, error
```

GetRegion returns the region associated with the input
ec2 client.

---

### GetRunningInstances(\*ec2.EC2)

```go
GetRunningInstances(*ec2.EC2) *ec2.DescribeInstancesOutput, error
```

GetRunningInstances returns all ec2 instances with a state of running.

---

### IsEC2Instance()

```go
IsEC2Instance() bool
```

IsEC2Instance checks whether the code is running on an AWS
EC2 instance by checking the existence of the file
/sys/devices/virtual/dmi/id/product_uuid. If the file exists,
the code is running on an EC2 instance, and the function
returns true. If the file does not exist, the function returns false,
indicating that the code is not running on an EC2 instance.

Example usage:

    isEC2 := IsEC2Instance()
    if isEC2 {
        fmt.Println("Running on an EC2 instance")
    } else {
        fmt.Println("Not running on an EC2 instance")
    }

Returns:

bool: A boolean value that indicates whether the code is running on an EC2 instance.

---

### TagInstance(\*ec2.EC2, string, string, string)

```go
TagInstance(*ec2.EC2, string, string, string) error
```

TagInstance tags the instance tied to the input ID with the specified tag.

---

### WaitForInstance(\*ec2.EC2, string)

```go
WaitForInstance(*ec2.EC2, string) error
```

WaitForInstance waits for the input instanceID to get to
a running state.

---

## Installation

To use the awsutils/ec2 package, you first need to install it.
Follow the steps below to install via go get.

```bash
go get github.com/awsutils/l50/ec2
```

---

## Usage

After installation, you can import the package in your Go project
using the following import statement:

```go
import "github.com/awsutils/l50/ec2"
```

---

## Tests

To ensure the package is working correctly, run the following
command to execute the tests for `awsutils/ec2`:

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
