package ec2_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	ec2utils "github.com/l50/awsutils/ec2"
	"github.com/stretchr/testify/assert"
)

var testEC2Connection *ec2utils.Connection
var testInstanceID string
var reservation *ec2.Reservation
var testParams ec2utils.Params

func init() {
	testParams = ec2utils.Params{
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
	testEC2Connection = ec2utils.NewConnection()
	var err error
	reservation, err = testEC2Connection.CreateInstance(testParams)
	if err != nil {
		fmt.Printf("failed to create instance: %v", err)
		os.Exit(1)
	}

	// Store the instance ID in a global variable for other tests to use
	testInstanceID = *reservation.Instances[0].InstanceId
}

func TestMain(m *testing.M) {
	if err := testEC2Connection.WaitForInstance(testInstanceID); err != nil {
		fmt.Printf("failed to wait for instance: %v", err)
		os.Exit(1)
	}

	if len(reservation.Instances) == 0 {
		fmt.Println("No instances found in reservation")
		os.Exit(1)
	}

	code := m.Run()

	err := testEC2Connection.DestroyInstance(testInstanceID)
	if err != nil {
		fmt.Printf("failed to destroy instance: %v", err)
		os.Exit(1)
	}

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
	assert.NoError(t, err)
	assert.NotNil(t, reservation)
	instanceID := *reservation.Instances[0].InstanceId

	// Schedule the instance to be destroyed after the test ends
	defer func() {
		err := c.DestroyInstance(instanceID)
		if err != nil {
			t.Fatalf("failed to destroy instance: %v", err)
		}
	}()
}

func TestCheckInstanceExists(t *testing.T) {
	// Test with the instance ID obtained from the setup
	err := testEC2Connection.CheckInstanceExists(testInstanceID)
	if err != nil {
		t.Fatalf("instance %s does not exist", testInstanceID)
	}
}

func TestTagInstance(t *testing.T) {
	// Test with the instance ID obtained from the setup
	err := testEC2Connection.TagInstance(testInstanceID, "key", "value")
	if err != nil {
		t.Fatalf("failed to tag instance %s: %v", testInstanceID, err)
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
	assert.NoError(t, err)
	assert.NotNil(t, reservation)
	instanceID := *reservation.Instances[0].InstanceId

	// Test the DestroyInstance function
	err = c.DestroyInstance(instanceID)
	if err != nil {
		t.Fatalf("failed to destroy instance %s: %v", instanceID, err)
	}
}

func TestGetRunningInstances(t *testing.T) {
	result, err := testEC2Connection.GetRunningInstances()
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestWaitForInstance(t *testing.T) {
	err := testEC2Connection.WaitForInstance(*reservation.Instances[0].InstanceId)
	assert.NoError(t, err)
}

func TestGetInstancePublicIP(t *testing.T) {
	ip, err := testEC2Connection.GetInstancePublicIP(*reservation.Instances[0].InstanceId)
	assert.NoError(t, err)
	assert.NotEmpty(t, ip)
}

func TestGetRegion(t *testing.T) {
	region, err := testEC2Connection.GetRegion()
	assert.NoError(t, err)
	assert.NotEmpty(t, region)
}

func TestGetInstances(t *testing.T) {
	instances, err := testEC2Connection.GetInstances(nil)
	assert.NoError(t, err)
	assert.NotNil(t, instances)
}

func TestGetInstanceState(t *testing.T) {
	state, err := testEC2Connection.GetInstanceState(*reservation.Instances[0].InstanceId)
	assert.NoError(t, err)
	assert.NotEmpty(t, state)
}

func TestGetInstancesRunningForMoreThan24Hours(t *testing.T) {
	instances, err := testEC2Connection.GetInstancesRunningForMoreThan24Hours()
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
			gotOutput, gotError := testEC2Connection.GetLatestAMI(tc.input)

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

func TestListSecurityGroupsForSubnet(t *testing.T) {
	validSubnetID, err := testEC2Connection.GetSubnetID("test-subnet-2")
	if err != nil {
		t.Fatalf("failed to get VPC ID: %v", err)
	}

	tests := []struct {
		name      string
		subnetID  string
		expectErr bool
	}{
		{
			name:      "Valid Subnet ID",
			subnetID:  validSubnetID,
			expectErr: false,
		},
		{
			name:      "Invalid Subnet ID",
			subnetID:  "subnet-invalid",
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, gotError := testEC2Connection.ListSecurityGroupsForSubnet(tc.subnetID)

			if tc.expectErr {
				assert.Error(t, gotError)
			} else {
				assert.NoError(t, gotError)
			}
		})
	}
}

func TestIsSubnetPubliclyRoutable(t *testing.T) {
	routableSubnetID, err := testEC2Connection.GetSubnetID("test-subnet-2")
	if err != nil {
		t.Fatalf("failed to get VPC ID: %v", err)
	}

	tests := []struct {
		name      string
		subnetID  string
		expectErr bool
	}{
		{
			name:      "Valid Subnet ID",
			subnetID:  routableSubnetID,
			expectErr: false,
		},
		{
			name:      "Invalid Subnet ID",
			subnetID:  "subnet-invalid",
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, gotError := testEC2Connection.IsSubnetPubliclyRoutable(tc.subnetID)

			if tc.expectErr {
				assert.Error(t, gotError)
			} else {
				assert.NoError(t, gotError)
			}
		})
	}
}

func TestListSecurityGroupsForVpc(t *testing.T) {
	validVPCID, err := testEC2Connection.GetVPCID("test-vpc")
	if err != nil {
		t.Fatalf("failed to get VPC ID: %v", err)
	}

	tests := []struct {
		name      string
		vpcID     string
		expectErr bool
	}{
		{
			name:      "Valid VPC ID",
			vpcID:     validVPCID,
			expectErr: false,
		},
		{
			name:      "Invalid VPC ID",
			vpcID:     "vpc-invalid",
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, gotError := testEC2Connection.ListSecurityGroupsForVpc(tc.vpcID)

			if tc.expectErr {
				assert.Error(t, gotError)
			} else {
				assert.NoError(t, gotError)
			}
		})
	}
}

func TestGetSubnetID(t *testing.T) {
	tests := []struct {
		name       string
		subnetName string
		expectErr  bool
	}{
		{
			name:       "Valid Subnet Name",
			subnetName: "test-subnet-2", // This should match a 'Name' tag of one of your subnets
			expectErr:  false,
		},
		{
			name:       "Invalid Subnet Name",
			subnetName: "InvalidSubnet", // This should not match any 'Name' tag of your subnets
			expectErr:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, gotError := testEC2Connection.GetSubnetID(tc.subnetName)

			if tc.expectErr {
				assert.Error(t, gotError)
			} else {
				assert.NoError(t, gotError)
			}
		})
	}
}

func TestGetVPCID(t *testing.T) {
	tests := []struct {
		name      string
		vpcName   string
		expectErr bool
	}{
		{
			name:      "Valid VPC Name",
			vpcName:   "test-vpc",
			expectErr: false,
		},
		{
			name:      "Invalid VPC Name",
			vpcName:   "InvalidVPC",
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, gotError := testEC2Connection.GetVPCID(tc.vpcName)

			if tc.expectErr {
				assert.Error(t, gotError)
			} else {
				assert.NoError(t, gotError)
			}
		})
	}
}
