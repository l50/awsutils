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

### Connection.CheckInstanceExists(string)

```go
CheckInstanceExists(string) error
```

CheckInstanceExists checks whether an instance
with the provided ID exists.

**Parameters:**

instanceID: the ID of the instance to check

**Returns:**

error: an error if any issue occurs while trying to check the instance

---

### Connection.CreateInstance(Params)

```go
CreateInstance(Params) *ec2.Reservation, error
```

CreateInstance creates a new EC2 instance
with the provided parameters.

**Parameters:**

ec2Params: the parameters to use

**Returns:**

*ec2.Reservation: the reservation of the created instance

error: an error if any issue occurs while trying to create the instance

---

### Connection.DestroyInstance(string)

```go
DestroyInstance(string) error
```

DestroyInstance destroys the instance with the provided ID.

**Parameters:**

instanceID: the ID of the instance to destroy

**Returns:**

error: an error if any issue occurs while trying to destroy the instance

---

### Connection.FindOverlyPermissiveInboundRules(string)

```go
FindOverlyPermissiveInboundRules(string) bool, error
```

FindOverlyPermissiveInboundRules checks if a specific security group permits all inbound traffic.
Specifically, it checks if the security group has an inbound rule with the IP protocol set to "-1",
which allows all IP traffic. This is useful for identifying security groups
that are configured with lenient security rules, especially in testing environments.
The function uses AWS SDK to describe security groups in AWS EC2 and checks their inbound rules.

**Parameters:**

secGrpID: A string containing the ID of the security group which needs to be checked for the all traffic inbound rule.

**Returns:**

bool: A boolean value indicating whether the security group permits all inbound traffic or not.

error: An error if any issue occurs while trying to describe the security group or check its inbound rules.

---

### Connection.GetInstancePublicIP(string)

```go
GetInstancePublicIP(string) string, error
```

GetInstancePublicIP retrieves the public IP address of the instance
with the provided ID.

**Parameters:**

instanceID: the ID of the instance to use

**Returns:**

string: the public IP address of the instance

error: an error if any issue occurs while trying to retrieve the public IP address

---

### Connection.GetInstanceState(string)

```go
GetInstanceState(string) string, error
```

GetInstanceState retrieves the state of the instance with the provided ID.

**Parameters:**

instanceID: the ID of the instance to use

**Returns:**

string: the state of the instance

error: an error if any issue occurs while trying to retrieve the state

---

### Connection.GetInstances([]*ec2.Filter)

```go
GetInstances([]*ec2.Filter) []*ec2.Instance, error
```

GetInstances retrieves all instances matching the provided filters.

**Parameters:**

filters: the filters to use

**Returns:**

[]*ec2.Instance: the instances matching the provided filters

error: an error if any issue occurs while trying to retrieve the instances

---

### Connection.GetInstancesRunningForMoreThan24Hours()

```go
GetInstancesRunningForMoreThan24Hours() []*ec2.Instance, error
```

GetInstancesRunningForMoreThan24Hours retrieves all instances
that have been running for more than 24 hours.

**Returns:**

[]*ec2.Instance: the instances that have been running for more than 24 hours

error: an error if any issue occurs while trying to retrieve the instances

---

### Connection.GetLatestAMI(AMIInfo)

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

### Connection.GetRegion()

```go
GetRegion() string, error
```

GetRegion retrieves the region of the connection.

**Returns:**

string: the region of the connection

error: an error if any issue occurs while trying to retrieve the region

---

### Connection.GetRunningInstances()

```go
GetRunningInstances() *ec2.DescribeInstancesOutput, error
```

GetRunningInstances retrieves all running instances.

**Returns:**

*ec2.DescribeInstancesOutput: the output of the DescribeInstances operation

error: an error if any issue occurs while trying to retrieve the running instances

---

### Connection.GetSubnetID(string)

```go
GetSubnetID(string) string, error
```

GetSubnetID retrieves the ID of the subnet with the provided name.

**Parameters:**

subnetName: the name of the subnet to use

**Returns:**

string: the ID of the subnet with the provided name

error: an error if any issue occurs while trying to retrieve the ID of the subnet with the provided name

---

### Connection.GetVPCID(string)

```go
GetVPCID(string) string, error
```

GetVPCID retrieves the ID of the VPC with the provided name.

**Parameters:**

vpcName: the name of the VPC to use

**Returns:**

string: the ID of the VPC with the provided name

error: an error if any issue occurs while trying to retrieve the ID of the VPC with the provided name

---

### Connection.IsSubnetPubliclyRoutable(string)

```go
IsSubnetPubliclyRoutable(string) bool, error
```

IsSubnetPubliclyRoutable checks whether the provided subnet ID
is publicly routable.

**Parameters:**

subnetID: the ID of the subnet to use

**Returns:**

bool: a boolean value indicating whether the provided subnet ID is publicly routable

error: an error if any issue occurs while trying to check whether the provided subnet ID is publicly routable

---

### Connection.ListSecurityGroups()

```go
ListSecurityGroups() []*ec2.SecurityGroup, error
```

ListSecurityGroups lists all security groups.

**Returns:**

[]*ec2.SecurityGroup: all security groups

error: an error if any issue occurs while trying to list the security groups

---

### Connection.ListSecurityGroupsForSubnet(string)

```go
ListSecurityGroupsForSubnet(string) []string, error
```

ListSecurityGroupsForSubnet lists all security groups
for the provided subnet ID.

**Parameters:**

subnetID: the ID of the subnet to use

**Returns:**

[]string: the IDs of the security groups for the provided subnet ID

error: an error if any issue occurs while trying to list the security groups

---

### Connection.ListSecurityGroupsForVpc(string)

```go
ListSecurityGroupsForVpc(string) []string, error
```

ListSecurityGroupsForVpc lists all security groups for the provided VPC ID.

**Parameters:**

vpcID: the ID of the VPC to use

**Returns:**

[]string: the IDs of the security groups for the provided VPC ID

error: an error if any issue occurs while trying to list the security groups

---

### Connection.TagInstance(string, string, string)

```go
TagInstance(string, string, string) error
```

TagInstance tags an instance with the provided key and value.

**Parameters:**

instanceID: the ID of the instance to tag

tagKey: the key of the tag to use

tagValue: the value of the tag to use

**Returns:**

error: an error if any issue occurs while trying to tag the instance

---

### Connection.WaitForInstance(string)

```go
WaitForInstance(string) error
```

WaitForInstance waits until the instance with the provided ID
is in the running state.

**Parameters:**

instanceID: the ID of the instance to wait for

**Returns:**

error: an error if any issue occurs while trying to wait for the instance

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

**Returns:**

bool: A boolean value that indicates whether the code is running on an EC2 instance.

---

### NewConnection()

```go
NewConnection() *Connection
```

NewConnection creates a new connection
to AWS EC2.

**Returns:**

*Connection: a new connection to AWS EC2

---

## Installation

To use the awsutils/ec2 package, you first need to install it.
Follow the steps below to install via go get.

```bash
go get github.com/l50/awsutils/ec2
```

---

## Usage

After installation, you can import the package in your Go project
using the following import statement:

```go
import "github.com/l50/awsutils/ec2"
```

---

## Tests

To ensure the package is working correctly, run the following
command to execute the tests for `ec2/awsutils`:

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
