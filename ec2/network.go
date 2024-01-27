package ec2

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// GetSubnetID retrieves the ID of the subnet with the provided name.
//
// **Parameters:**
//
// subnetName: the name of the subnet to use
//
// **Returns:**
//
// string: the ID of the subnet with the provided name
//
// error: an error if any issue occurs while trying to retrieve
// the ID of the subnet with the provided name
func (c *Connection) GetSubnetID(subnetName string) (string, error) {
	input := &ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:Name"),
				Values: []*string{
					aws.String(subnetName),
				},
			},
		},
	}

	result, err := c.Client.DescribeSubnets(input)
	if err != nil {
		return "", err
	}

	if len(result.Subnets) == 0 {
		return "", errors.New("no subnet found with the provided name")
	}

	subnetID := *result.Subnets[0].SubnetId
	if err := c.checkResourceExistence("subnet", subnetID); err != nil {
		return "", err
	}

	return subnetID, nil
}

// GetVPCID retrieves the ID of the VPC with the provided name.
//
// **Parameters:**
//
// vpcName: the name of the VPC to use. If "default" is provided, the function
// will return the ID of the default VPC.
//
// **Returns:**
//
// string: the ID of the VPC with the provided name
//
// error: an error if any issue occurs while trying to retrieve
// the ID of the VPC with the provided name
func (c *Connection) GetVPCID(vpcName string) (string, error) {
	var input *ec2.DescribeVpcsInput

	// Check if we're looking for the default VPC
	if vpcName == "default" {
		input = &ec2.DescribeVpcsInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("isDefault"),
					Values: []*string{aws.String("true")},
				},
			},
		}
	} else {
		input = &ec2.DescribeVpcsInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("tag:Name"),
					Values: []*string{aws.String(vpcName)},
				},
			},
		}
	}

	result, err := c.Client.DescribeVpcs(input)
	if err != nil {
		return "", err
	}

	if len(result.Vpcs) == 0 {
		if vpcName == "default" {
			return "", errors.New("no default VPC found")
		}
		return "", errors.New("no VPC found with the provided name")
	}

	return *result.Vpcs[0].VpcId, nil
}

// IsSubnetPubliclyRoutable checks whether the provided subnet ID
// is publicly routable.
//
// **Parameters:**
//
// subnetID: the ID of the subnet to use
//
// **Returns:**
//
// bool: a boolean value indicating whether the provided subnet ID is publicly routable
//
// error: an error if any issue occurs while trying to check whether the
// provided subnet ID is publicly routable
func (c *Connection) IsSubnetPubliclyRoutable(subnetID string) (bool, error) {
	if err := c.checkResourceExistence("subnet", subnetID); err != nil {
		return false, err
	}

	input := &ec2.DescribeRouteTablesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("association.subnet-id"),
				Values: []*string{aws.String(subnetID)},
			},
		},
	}
	result, err := c.Client.DescribeRouteTables(input)
	if err != nil {
		return false, err
	}
	for _, routeTable := range result.RouteTables {
		for _, route := range routeTable.Routes {
			if route.GatewayId != nil && *route.GatewayId != "local" && *route.DestinationCidrBlock == "0.0.0.0/0" {
				return true, nil
			}
		}
	}
	return false, nil
}

// ListSecurityGroupsForVpc lists all security groups for the provided VPC ID.
//
// **Parameters:**
//
// vpcID: the ID of the VPC to use
//
// **Returns:**
//
// []*ec2.SecurityGroup: all security groups for the provided VPC ID
//
// error: an error if any issue occurs while trying to list the security groups
func (c *Connection) ListSecurityGroupsForVpc(vpcID string) ([]*ec2.SecurityGroup, error) {
	if err := c.checkResourceExistence("vpc", vpcID); err != nil {
		return nil, err
	}
	input := &ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []*string{aws.String(vpcID)},
			},
		},
	}

	result, err := c.Client.DescribeSecurityGroups(input)
	if err != nil {
		return nil, err
	}

	return result.SecurityGroups, nil
}

// ListSubnetsForVPC lists subnets for the provided VPC name and subnet location.
//
// **Parameters:**
//
// vpcName: the name of the VPC to use. Returns subnets for the default VPC if "default" is provided.
// subnetLocation: the location of the subnet. Can be "public", "private", or "all".
//
// **Returns:**
//
// []*ec2.Subnet: the list of subnets for the provided VPC name and location
//
// error: an error if any issue occurs while trying to list the subnets
func (c *Connection) ListSubnetsForVPC(vpcName string, subnetLocation string) ([]*ec2.Subnet, error) {
	// Retrieve the VPC ID for the default VPC
	vpcID, err := c.GetVPCID(vpcName)
	if err != nil {
		return nil, err
	}

	// Validate subnetLocation
	if subnetLocation != "public" && subnetLocation != "private" && subnetLocation != "all" {
		return nil, errors.New("subnetLocation must be public, private, or all")
	}

	// Build the subnet filter based on subnetLocation
	var filters []*ec2.Filter
	if subnetLocation != "all" {
		filters = append(filters, &ec2.Filter{
			Name:   aws.String("tag:SubnetType"),
			Values: []*string{aws.String(subnetLocation)},
		})
	}

	// Always include the VPC ID in the filter
	filters = append(filters, &ec2.Filter{
		Name:   aws.String("vpc-id"),
		Values: []*string{aws.String(vpcID)},
	})

	// Describe subnets with the prepared filters
	input := &ec2.DescribeSubnetsInput{
		Filters: filters,
	}

	result, err := c.Client.DescribeSubnets(input)
	if err != nil {
		return nil, err
	}

	return result.Subnets, nil
}
