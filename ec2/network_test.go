package ec2_test

import (
	"encoding/json"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	ec2utils "github.com/l50/awsutils/ec2"
	"github.com/stretchr/testify/assert"
)

func TestIsSubnetPubliclyRoutable(t *testing.T) {
	c := ec2utils.NewConnection()
	routableSubnetID, err := c.GetSubnetID("test-subnet-2")
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
			c := ec2utils.NewConnection()
			_, gotError := c.IsSubnetPubliclyRoutable(tc.subnetID)

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
			subnetName: "test-subnet-2",
			expectErr:  false,
		},
		{
			name:       "Invalid Subnet Name",
			subnetName: "InvalidSubnet",
			expectErr:  true,
		},
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
			c := ec2utils.NewConnection()
			_, gotError := c.GetVPCID(tc.vpcName)

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
	validVPCID, err := c.GetVPCID("test-vpc")
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
			c := ec2utils.NewConnection()
			_, gotError := c.ListSecurityGroupsForVpc(tc.vpcID)

			if tc.expectErr {
				assert.Error(t, gotError)
			} else {
				assert.NoError(t, gotError)
			}
		})
	}
}

func TestListSubnetsForVPC(t *testing.T) {
	tests := []struct {
		name             string
		vpcName          string
		subnetLocation   string
		mockVpcOutput    ec2.DescribeVpcsOutput
		mockSubnetOutput ec2.DescribeSubnetsOutput
		wantSubnetIds    []string
		wantErr          error
	}{
		{
			name:           "valid request with all subnets",
			vpcName:        "default",
			subnetLocation: "all",
			mockVpcOutput: ec2.DescribeVpcsOutput{
				Vpcs: []*ec2.Vpc{
					{
						VpcId: aws.String("vpc-12345"),
					},
				},
			},
			mockSubnetOutput: ec2.DescribeSubnetsOutput{
				Subnets: []*ec2.Subnet{
					{
						SubnetId: aws.String("subnet-eb0fc5b1"),
					},
					{
						SubnetId: aws.String("subnet-1b49f77d"),
					},
				},
			},
		},
	}

	// Regex to match subnet IDs (e.g., subnet-123abc)
	subnetIDRegex := regexp.MustCompile(`^subnet-[0-9a-f]+$`)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := ec2utils.NewConnection()
			gotSubnets, gotError := c.ListSubnetsForVPC(tc.vpcName, tc.subnetLocation)

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
