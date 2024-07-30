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

### AWSService.AddRoleToInstanceProfile(string)

```go
AddRoleToInstanceProfile(string) *iam.AddRoleToInstanceProfileOutput, error
```

AddRoleToInstanceProfile adds a role to an instance profile.

**Parameters:**

profileName: The name of the instance profile to add the role to.
roleName: The name of the role to add to the instance profile.

**Returns:**

*iam.AddRoleToInstanceProfileOutput: A pointer to the AddRoleToInstanceProfileOutput.
error: An error if any issue occurs while trying to add the role to the instance profile.

---

### AWSService.AttachRolePolicy(string)

```go
AttachRolePolicy(string) *iam.AttachRolePolicyOutput, error
```

AttachRolePolicy attaches a policy to a role.

**Parameters:**

roleName: The name of the role to attach the policy to.
policyArn: The ARN of the policy to attach.

**Returns:**

*iam.AttachRolePolicyOutput: A pointer to the AttachRolePolicyOutput.
error: An error if any issue occurs while trying to attach the policy to the role.

---

### AWSService.CreateInstanceProfile(string)

```go
CreateInstanceProfile(string) *iam.CreateInstanceProfileOutput, error
```

CreateInstanceProfile creates a new instance profile with the given name.

**Parameters:**

profileName: The name of the instance profile to create.

**Returns:**

*iam.CreateInstanceProfileOutput: A pointer to the CreateInstanceProfileOutput.
error: An error if any issue occurs while trying to create the instance profile.

---

### AWSService.CreateRole(string)

```go
CreateRole(string) *iam.CreateRoleOutput, error
```

CreateRole creates a new role with the given name and assume role policy.

**Parameters:**

roleName: The name of the role to create.
assumeRolePolicy: The policy that the role will assume.

**Returns:**

*iam.CreateRoleOutput: A pointer to the CreateRoleOutput.
error: An error if any issue occurs while trying to create the role.

---

### AWSService.DeleteInstanceProfile(string)

```go
DeleteInstanceProfile(string) *iam.DeleteInstanceProfileOutput, error
```

DeleteInstanceProfile deletes an instance profile.

**Parameters:**

profileName: The name of the instance profile to delete.

**Returns:**

*iam.DeleteInstanceProfileOutput: A pointer to the DeleteInstanceProfileOutput.
error: An error if any issue occurs while trying to delete the instance profile.

---

### AWSService.DeleteRole(string)

```go
DeleteRole(string) *iam.DeleteRoleOutput, error
```

DeleteRole deletes a role.

**Parameters:**

roleName: The name of the role to delete.

**Returns:**

*iam.DeleteRoleOutput: A pointer to the DeleteRoleOutput.
error: An error if any issue occurs while trying to delete the role.

---

### AWSService.DeleteRolePolicy(string)

```go
DeleteRolePolicy(string) *iam.DeleteRolePolicyOutput, error
```

DeleteRolePolicy deletes a policy from a role.

**Parameters:**

roleName: The name of the role to delete the policy from.
policyName: The name of the policy to delete.

**Returns:**

*iam.DeleteRolePolicyOutput: A pointer to the DeleteRolePolicyOutput.
error: An error if any issue occurs while trying to delete the policy from the role.

---

### AWSService.DetachRolePolicy(string)

```go
DetachRolePolicy(string) *iam.DetachRolePolicyOutput, error
```

DetachRolePolicy detaches a policy from a role.

**Parameters:**

roleName: The name of the role to detach the policy from.
policyArn: The ARN of the policy to detach.

**Returns:**

*iam.DetachRolePolicyOutput: A pointer to the DetachRolePolicyOutput.
error: An error if any issue occurs while trying to detach the policy from the role.

---

### AWSService.GetAWSIdentity()

```go
GetAWSIdentity() *AWSIdentity, error
```

GetAWSIdentity retrieves the AWS identity of the caller.

**Returns:**

*AWSIdentity: A pointer to the AWSIdentity of the caller.
error: An error if any issue occurs while trying to get the AWS identity.

---

### AWSService.GetInstanceProfile(string)

```go
GetInstanceProfile(string) *types.InstanceProfile, error
```

GetInstanceProfile retrieves the instance profile for a given profile name.

**Parameters:**

profileName: The name of the profile to retrieve.

**Returns:**

*types.InstanceProfile: A pointer to the InstanceProfile.
error: An error if any issue occurs while trying to get the instance profile.

---

### AWSService.PutRolePolicy(string)

```go
PutRolePolicy(string) *iam.PutRolePolicyOutput, error
```

PutRolePolicy updates the policy for a role.

**Parameters:**

roleName: The name of the role to update the policy for.
policyName: The name of the policy to update.
policyDocument: The policy document to update.

**Returns:**

*iam.PutRolePolicyOutput: A pointer to the PutRolePolicyOutput.
error: An error if any issue occurs while trying to update the policy for the role.

---

### AWSService.RemoveRoleFromInstanceProfile(string)

```go
RemoveRoleFromInstanceProfile(string) *iam.RemoveRoleFromInstanceProfileOutput error
```

RemoveRoleFromInstanceProfile removes a role from an instance profile.

**Parameters:**

profileName: The name of the instance profile to remove the role from.
roleName: The name of the role to remove from the instance profile.

**Returns:**

*iam.RemoveRoleFromInstanceProfileOutput: A pointer to the RemoveRoleFromInstanceProfileOutput.
error: An error if any issue occurs while trying to remove the role from the instance profile.

---

### NewAWSService()

```go
NewAWSService() *AWSService, error
```

NewAWSService creates a new AWSService with the default AWS configuration.

**Returns:**

*AWSService: A pointer to the newly created AWSService.
error: An error if any issue occurs while trying to create the AWSService.

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
