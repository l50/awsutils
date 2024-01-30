package ec2_test

import (
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	ec2utils "github.com/l50/awsutils/ec2"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestNewConnection(t *testing.T) {
	c := ec2utils.NewConnection()
	assert.NotNil(t, c.Client)
}

func TestCreateInstance(t *testing.T) {
	testParams := ec2utils.Params{
		AssociatePublicIPAddress: true,
		ImageID:                  os.Getenv("AMI"),
		InstanceName:             os.Getenv("INST_NAME"),
		InstanceType:             os.Getenv("INST_TYPE"),
		InstanceProfile:          os.Getenv("IAM_INSTANCE_PROFILE"),
		MinCount:                 1,
		MaxCount:                 1,
		SecurityGroupIDs:         []string{os.Getenv("SEC_GRP_ID")},
		SubnetID:                 os.Getenv("SUBNET_ID"),
		VolumeSize:               8,
	}
	c := ec2utils.NewConnection()
	reservation, err := c.CreateInstance(testParams)
	// Schedule the instance to be destroyed after the test ends
	defer func() {
		err := c.DestroyInstance(*reservation.Instances[0].InstanceId)
		if err != nil {
			t.Fatalf("failed to destroy instance: %v", err)
		}
	}()
	assert.NoError(t, err)
	assert.NotNil(t, reservation)
}

func TestCheckInstanceExists(t *testing.T) {
	testParams := ec2utils.Params{
		AssociatePublicIPAddress: true,
		ImageID:                  os.Getenv("AMI"),
		InstanceName:             os.Getenv("INST_NAME"),
		InstanceType:             os.Getenv("INST_TYPE"),
		InstanceProfile:          os.Getenv("IAM_INSTANCE_PROFILE"),
		MinCount:                 1,
		MaxCount:                 1,
		SecurityGroupIDs:         []string{os.Getenv("SEC_GRP_ID")},
		SubnetID:                 os.Getenv("SUBNET_ID"),
		VolumeSize:               8,
	}
	c := ec2utils.NewConnection()
	reservation, err := c.CreateInstance(testParams)
	assert.NoError(t, err)
	assert.NotNil(t, reservation)

	err = c.CheckInstanceExists(*reservation.Instances[0].InstanceId)
	if err != nil {
		t.Fatalf("instance %s does not exist", *reservation.Instances[0].InstanceId)
	}

	// Schedule the instance to be destroyed after the test ends
	defer func() {
		err := c.DestroyInstance(*reservation.Instances[0].InstanceId)
		if err != nil {
			t.Fatalf("failed to destroy instance: %v", err)
		}
	}()
}

func TestTagInstance(t *testing.T) {
	testParams := ec2utils.Params{
		AssociatePublicIPAddress: true,
		ImageID:                  os.Getenv("AMI"),
		InstanceName:             os.Getenv("INST_NAME"),
		InstanceType:             os.Getenv("INST_TYPE"),
		InstanceProfile:          os.Getenv("IAM_INSTANCE_PROFILE"),
		MinCount:                 1,
		MaxCount:                 1,
		SecurityGroupIDs:         []string{os.Getenv("SEC_GRP_ID")},
		SubnetID:                 os.Getenv("SUBNET_ID"),
		VolumeSize:               8,
	}

	c := ec2utils.NewConnection()
	reservation, err := c.CreateInstance(testParams)
	// Schedule the instance to be destroyed after the test ends
	defer func() {
		err := c.DestroyInstance(*reservation.Instances[0].InstanceId)
		if err != nil {
			t.Fatalf("failed to destroy instance: %v", err)
		}
	}()
	assert.NoError(t, err)
	assert.NotNil(t, reservation)

	err = c.TagInstance(*reservation.Instances[0].InstanceId, "key", "value")
	if err != nil {
		t.Fatalf("failed to tag instance %s: %v", *reservation.Instances[0].InstanceId, err)
	}
}

func TestDestroyInstance(t *testing.T) {
	// Create an instance here instead of using a global one
	testParams := ec2utils.Params{
		AssociatePublicIPAddress: true,
		ImageID:                  os.Getenv("AMI"),
		InstanceName:             os.Getenv("INST_NAME"),
		InstanceType:             os.Getenv("INST_TYPE"),
		InstanceProfile:          os.Getenv("IAM_INSTANCE_PROFILE"),
		MinCount:                 1,
		MaxCount:                 1,
		SecurityGroupIDs:         []string{os.Getenv("SEC_GRP_ID")},
		SubnetID:                 os.Getenv("SUBNET_ID"),
		VolumeSize:               8,
	}
	c := ec2utils.NewConnection()
	reservation, err := c.CreateInstance(testParams)
	// Schedule the instance to be destroyed after the test ends
	defer func() {
		err := c.DestroyInstance(*reservation.Instances[0].InstanceId)
		if err != nil {
			t.Fatalf("failed to destroy instance: %v", err)
		}
	}()
	assert.NoError(t, err)
	assert.NotNil(t, reservation)
	if err := c.WaitForInstance(*reservation.Instances[0].InstanceId); err != nil {
		t.Fatalf("failed to wait for instance: %v", err)
	}
}

func TestGetRunningInstances(t *testing.T) {
	c := ec2utils.NewConnection()
	result, err := c.GetRunningInstances()
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestWaitForInstance(t *testing.T) {
	testParams := ec2utils.Params{
		AssociatePublicIPAddress: true,
		ImageID:                  os.Getenv("AMI"),
		InstanceName:             os.Getenv("INST_NAME"),
		InstanceType:             os.Getenv("INST_TYPE"),
		InstanceProfile:          os.Getenv("IAM_INSTANCE_PROFILE"),
		MinCount:                 1,
		MaxCount:                 1,
		SecurityGroupIDs:         []string{os.Getenv("SEC_GRP_ID")},
		SubnetID:                 os.Getenv("SUBNET_ID"),
		VolumeSize:               8,
	}
	c := ec2utils.NewConnection()
	reservation, err := c.CreateInstance(testParams)
	// Schedule the instance to be destroyed after the test ends
	defer func() {
		err := c.DestroyInstance(*reservation.Instances[0].InstanceId)
		if err != nil {
			t.Fatalf("failed to destroy instance: %v", err)
		}
	}()
	assert.NoError(t, err)
	assert.NotNil(t, reservation)
	err = c.WaitForInstance(*reservation.Instances[0].InstanceId)
	assert.NoError(t, err)
}

func TestGetInstancePublicIP(t *testing.T) {
	testParams := ec2utils.Params{
		AssociatePublicIPAddress: true,
		ImageID:                  os.Getenv("AMI"),
		InstanceName:             os.Getenv("INST_NAME"),
		InstanceType:             os.Getenv("INST_TYPE"),
		InstanceProfile:          os.Getenv("IAM_INSTANCE_PROFILE"),
		MinCount:                 1,
		MaxCount:                 1,
		SecurityGroupIDs:         []string{os.Getenv("SEC_GRP_ID")},
		SubnetID:                 os.Getenv("SUBNET_ID"),
		VolumeSize:               8,
	}
	c := ec2utils.NewConnection()
	reservation, err := c.CreateInstance(testParams)
	if err != nil {
		t.Fatalf("failed to create instance: %v", err)
	}
	if reservation == nil || len(reservation.Instances) == 0 || reservation.Instances[0] == nil {
		t.Fatalf("no instances were created")
	}
	instanceID := *reservation.Instances[0].InstanceId
	// Schedule the instance to be destroyed after the test ends
	defer func() {
		if err := c.DestroyInstance(instanceID); err != nil {
			t.Fatalf("failed to destroy instance: %v", err)
		}
	}()

	// Wait for the instance to be available before fetching the public IP
	if err := c.WaitForInstance(instanceID); err != nil {
		t.Fatalf("error waiting for instance: %v", err)
	}

	ip, err := c.GetInstancePublicIP(instanceID)
	assert.NoError(t, err)
	assert.NotEmpty(t, ip)
}

func TestGetRegion(t *testing.T) {
	c := ec2utils.NewConnection()
	region, err := c.GetRegion()
	assert.NoError(t, err)
	assert.NotEmpty(t, region)
}

func TestGetInstances(t *testing.T) {
	c := ec2utils.NewConnection()
	filters := []*ec2.Filter{
		{
			Name:   aws.String("instance-state-name"),
			Values: []*string{aws.String("running")},
		},
	}
	instances, err := c.GetInstances(filters)
	assert.NoError(t, err)
	assert.NotNil(t, instances)
}

func TestGetInstanceState(t *testing.T) {
	testParams := ec2utils.Params{
		AssociatePublicIPAddress: true,
		ImageID:                  os.Getenv("AMI"),
		InstanceName:             os.Getenv("INST_NAME"),
		InstanceType:             os.Getenv("INST_TYPE"),
		InstanceProfile:          os.Getenv("IAM_INSTANCE_PROFILE"),
		MinCount:                 1,
		MaxCount:                 1,
		SecurityGroupIDs:         []string{os.Getenv("SEC_GRP_ID")},
		SubnetID:                 os.Getenv("SUBNET_ID"),
		VolumeSize:               8,
	}
	c := ec2utils.NewConnection()
	reservation, err := c.CreateInstance(testParams)
	// Schedule the instance to be destroyed after the test ends
	defer func() {
		err := c.DestroyInstance(*reservation.Instances[0].InstanceId)
		if err != nil {
			t.Fatalf("failed to destroy instance: %v", err)
		}
	}()
	assert.NoError(t, err)
	assert.NotNil(t, reservation)
	state, err := c.GetInstanceState(*reservation.Instances[0].InstanceId)
	assert.NoError(t, err)
	assert.NotEmpty(t, state)
}

func TestGetInstancesRunningForMoreThan24Hours(t *testing.T) {
	c := ec2utils.NewConnection()
	instances, err := c.GetInstancesRunningForMoreThan24Hours()
	assert.NoError(t, err)
	assert.NotNil(t, instances)
}

func TestIsEC2Instance(t *testing.T) {
	result := ec2utils.IsEC2Instance()
	assert.NotNil(t, result)
}

func TestGetLatestAMI(t *testing.T) {
	tests := []struct {
		name      string
		input     ec2utils.AMIInfo
		expectErr bool
	}{
		{
			name: "Ubuntu 22.04 arm64",
			input: ec2utils.AMIInfo{
				Distro:       "ubuntu",
				Version:      "22.04",
				Architecture: "arm64",
				Region:       "us-west-1",
			},
			expectErr: false,
		},
		{
			name: "Ubuntu 20.04 amd64",
			input: ec2utils.AMIInfo{
				Distro:       "ubuntu",
				Version:      "20.04",
				Architecture: "amd64",
				Region:       "us-west-1",
			},
			expectErr: false,
		},
		{
			name: "Unsupported distro",
			input: ec2utils.AMIInfo{
				Distro:       "not-supported",
				Version:      "20.04",
				Architecture: "amd64",
				Region:       "us-west-1",
			},
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := ec2utils.NewConnection()
			gotOutput, gotError := c.GetLatestAMI(tc.input)

			if gotError != nil {
				assert.Error(t, gotError)
			} else {
				assert.NotEmpty(t, gotOutput)
				assert.True(t, strings.HasPrefix(gotOutput, "ami-"))
			}

			if tc.expectErr {
				assert.Error(t, gotError)
			} else {
				assert.NoError(t, gotError)
			}
		})
	}
}

func TestCreateSecurityGroup(t *testing.T) {
	c := ec2utils.NewConnection()
	vpcID, err := c.GetVPCID("default")
	if err != nil {
		t.Fatalf("failed to get VPC ID: %v", err)
	}

	tests := []struct {
		name            string
		groupName       string
		description     string
		vpcID           string
		expectedGroupID string
		expectErr       bool
	}{
		{
			name:            "Valid Input",
			groupName:       "test-group",
			description:     "test description",
			vpcID:           vpcID,
			expectedGroupID: "",
			expectErr:       false,
		},
		{
			name:            "Empty Group Name",
			groupName:       "",
			description:     "test description",
			vpcID:           vpcID,
			expectedGroupID: "",
			expectErr:       true,
		},
		{
			name:            "Empty Description",
			groupName:       "test-group",
			description:     "",
			vpcID:           vpcID,
			expectedGroupID: "",
			expectErr:       true,
		},
		{
			name:            "Empty VPC ID",
			groupName:       "test-group",
			description:     "test description",
			vpcID:           "non-existent-vpc-id",
			expectedGroupID: "",
			expectErr:       true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create Security Group first
			groupID, err := c.CreateSecurityGroup(tc.groupName, tc.description, tc.vpcID)

			// Check error for CreateSecurityGroup
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Ensure group id was populated
			if groupID == "" && !tc.expectErr {
				t.Errorf("groupID should not be empty, but was: %v", groupID)
			}

			// Defer the destruction of the security group if creation was successful
			if err == nil {
				defer func() {
					err := c.DestroySecurityGroup(groupID)
					if err != nil {
						t.Errorf("failed to destroy security group: %v", err)
					}
				}()
			}
		})
	}
}

func TestDestroySecurityGroup(t *testing.T) {
	c := ec2utils.NewConnection()
	vpcID, err := c.GetVPCID("default")
	if err != nil {
		t.Fatalf("failed to get VPC ID: %v", err)
	}
	tests := []struct {
		name             string
		groupName        string
		description      string
		vpcID            string
		expectCreateErr  bool
		expectDestroyErr bool
	}{
		{
			name:             "group-to-destroy",
			groupName:        "test-group",
			description:      "test description",
			vpcID:            vpcID,
			expectCreateErr:  false,
			expectDestroyErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create Security Group first
			groupID, createErr := c.CreateSecurityGroup(tc.groupName, tc.description, tc.vpcID)

			// Check error for CreateSecurityGroup
			if tc.expectCreateErr {
				assert.Error(t, createErr)
			} else {
				assert.NoError(t, createErr)
			}

			// If CreateSecurityGroup was successful, try to destroy it
			if createErr == nil {
				defer func() {
					destroyErr := c.DestroySecurityGroup(groupID)

					// Check error for DestroySecurityGroup
					if tc.expectDestroyErr {
						assert.Error(t, destroyErr)
					} else {
						assert.NoError(t, destroyErr)
					}
				}()
			}
		})
	}
}
