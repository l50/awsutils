package ec2_test

import (
	"encoding/json"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	ec2utils "github.com/l50/awsutils/ec2"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Ensure region is set to us-west-1
	region, err := ec2utils.NewConnection().GetRegion()
	if err != nil {
		panic(err)
	}
	if region != "us-west-1" {
		panic("region must be set to us-west-1 for these tests")
	}
}

func TestIsSubnetPublic(t *testing.T) {
	c := ec2utils.NewConnection()
	vpcs, err := c.ListVPCs()
	assert.NoError(t, err)
	subnets, err := c.ListVPCSubnets(*vpcs[0].VpcId, "private")
	assert.NoError(t, err)
	assert.True(t, len(subnets) > 0)

	tests := []struct {
		name      string
		subnetID  string
		expectErr bool
	}{
		{
			name:      "Publicly Routed Subnet ID",
			subnetID:  *subnets[0].SubnetId,
			expectErr: false,
		},
		{
			name:      "Non-existent Subnet ID",
			subnetID:  "subnet-notrealatall",
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := ec2utils.NewConnection()
			_, gotError := c.IsSubnetPublic(tc.subnetID)

			if tc.expectErr {
				assert.Error(t, gotError)
			} else {
				assert.NoError(t, gotError)
			}
		})
	}
}

func TestGetSubnetID(t *testing.T) {
	c := ec2utils.NewConnection()
	vpcs, err := c.ListVPCs()
	assert.NoError(t, err)

	allSubnets, err := c.ListVPCSubnets(*vpcs[0].VpcId, "all")
	assert.NoError(t, err)
	assert.True(t, len(allSubnets) > 0, "There should be at least one subnet in 'all' category")

	publicSubnets, err := c.ListVPCSubnets(*vpcs[0].VpcId, "public")
	assert.NoError(t, err)

	privateSubnets, err := c.ListVPCSubnets(*vpcs[0].VpcId, "private")
	assert.NoError(t, err)

	// Test if the first subnet from allSubnets has a valid 'Name' tag
	var validSubnetName string
	if len(allSubnets) > 0 {
		for _, tag := range allSubnets[0].Tags {
			if *tag.Key == "Name" {
				validSubnetName = *tag.Value
				break
			}
		}
	}
	assert.NotEmpty(t, validSubnetName, "Subnet must have a 'Name' tag")

	tests := []struct {
		name       string
		subnetName string
		expectErr  bool
	}{
		{
			name:       "Valid Subnet Name",
			subnetName: validSubnetName,
			expectErr:  false,
		},
		{
			name:       "Invalid Subnet Name",
			subnetName: "InvalidSubnet",
			expectErr:  true,
		},
	}

	// Add a public subnet test case if available
	if len(publicSubnets) > 0 {
		publicSubnetName := getSubnetNameFromTag(publicSubnets[0])
		tests = append(tests, struct {
			name       string
			subnetName string
			expectErr  bool
		}{
			name:       "Public Subnet",
			subnetName: publicSubnetName,
			expectErr:  false,
		})
	}

	// Add a private subnet test case if available
	if len(privateSubnets) > 0 {
		privateSubnetName := getSubnetNameFromTag(privateSubnets[0])
		tests = append(tests, struct {
			name       string
			subnetName string
			expectErr  bool
		}{
			name:       "Private Subnet",
			subnetName: privateSubnetName,
			expectErr:  false,
		})
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := ec2utils.NewConnection()
			_, gotError := c.GetSubnetID(tc.subnetName)

			if tc.expectErr {
				assert.Error(t, gotError)
			} else {
				assert.NoError(t, gotError)
			}
		})
	}
}

// Helper function to extract the name of a subnet from its tags
func getSubnetNameFromTag(subnet *ec2.Subnet) string {
	for _, tag := range subnet.Tags {
		if *tag.Key == "Name" {
			return *tag.Value
		}
	}
	return ""
}

func TestGetVPCID(t *testing.T) {
	c := ec2utils.NewConnection()
	vpcs, err := c.ListVPCs()
	assert.NoError(t, err)
	assert.True(t, len(vpcs) > 0)

	// Ensure your VPCs have 'Name' tags and fetch the name of the first VPC
	var validVPCName string
	for _, tag := range vpcs[0].Tags {
		if *tag.Key == "Name" {
			validVPCName = *tag.Value
			break
		}
	}
	assert.NotEmpty(t, validVPCName, "VPC must have a 'Name' tag")

	tests := []struct {
		name      string
		vpcName   string
		expectErr bool
	}{
		{
			name:      "Valid VPC Name",
			vpcName:   validVPCName,
			expectErr: false,
		},
		{
			name:      "Default VPC",
			vpcName:   "default",
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
			c := ec2utils.NewConnection()
			id, gotError := c.GetVPCID(tc.vpcName)

			if tc.expectErr {
				assert.Error(t, gotError)
				assert.Empty(t, id)
			} else {
				assert.NoError(t, gotError)
				assert.NotEmpty(t, id)
			}
		})
	}
}

func TestListSecurityGroupsForSubnet(t *testing.T) {
	c := ec2utils.NewConnection()
	validSubnetID, err := c.GetSubnetID("test-subnet-2")
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
			_, gotError := c.ListSecurityGroupsForSubnet(tc.subnetID)

			if tc.expectErr {
				assert.Error(t, gotError)
			} else {
				assert.NoError(t, gotError)
			}
		})
	}
}

func TestListSecurityGroupsForVpc(t *testing.T) {
	c := ec2utils.NewConnection()
	vpcs, err := c.ListVPCs()
	assert.NoError(t, err)
	assert.True(t, len(vpcs) > 0)
	validVPCID := *vpcs[0].VpcId

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
			_, gotError := c.ListSecurityGroupsForVpc(tc.vpcID)

			if tc.expectErr {
				assert.Error(t, gotError)
			} else {
				assert.NoError(t, gotError)
			}
		})
	}
}

func TestListVPCSubnets(t *testing.T) {
	tests := []struct {
		name           string
		subnetLocation string
		wantSubnetIds  []string
		wantErr        error
	}{
		{
			name:           "valid request with all subnets",
			subnetLocation: "all",
		},
	}

	subnetIDRegex := regexp.MustCompile(`^subnet-[0-9a-f]+$`)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := ec2utils.NewConnection()
			vpcs, err := c.ListVPCs()
			assert.NoError(t, err)

			gotSubnets, gotError := c.ListVPCSubnets(*vpcs[0].VpcId, tc.subnetLocation)

			if tc.wantErr != nil {
				if gotError == nil {
					t.Errorf("expected an error but got none")
				}
				if gotError.Error() != tc.wantErr.Error() {
					t.Errorf("expected error %q, got %q", tc.wantErr, gotError)
				}
				return
			}

			if gotError != nil {
				t.Fatalf("unexpected error: %s", gotError)
			}

			gotSubnetsBytes, err := json.Marshal(gotSubnets)
			if err != nil {
				t.Fatalf("failed to marshal subnet output: %s", err)
			}

			var gotSubnetsSlice []struct {
				SubnetID *string `json:"SubnetId"`
			}
			if err := json.Unmarshal(gotSubnetsBytes, &gotSubnetsSlice); err != nil {
				t.Fatalf("failed to unmarshal subnet output: %s", err)
			}

			for _, subnet := range gotSubnetsSlice {
				if subnet.SubnetID == nil || !subnetIDRegex.MatchString(*subnet.SubnetID) {
					t.Errorf("received invalid subnet ID: %v", subnet.SubnetID)
				}
			}
		})
	}
}

func TestListVPCs(t *testing.T) {
	tests := []struct {
		name      string
		expectErr bool
	}{
		{
			name:      "Valid Request",
			expectErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := ec2utils.NewConnection()
			_, gotError := c.ListVPCs()

			if tc.expectErr {
				assert.Error(t, gotError)
			} else {
				assert.NoError(t, gotError)
			}
		})
	}
}
