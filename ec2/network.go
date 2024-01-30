package ec2

import (
	"errors"
	"fmt"
	"strings"

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
		return "", fmt.Errorf("error describing subnets: %v", err)
	}

	if len(result.Subnets) == 0 {
		return "", fmt.Errorf("no subnet found with the name: %s", subnetName)
	}

	subnetID := *result.Subnets[0].SubnetId
	if subnetID == "" {
		return "", fmt.Errorf("found subnet has empty ID for the name: %s", subnetName)
	}

	if err := c.checkResourceExistence("subnet", subnetID); err != nil {
		return "", fmt.Errorf("subnet with ID %s does not exist: %v", subnetID, err)
	}

	return subnetID, nil
}

// GetSubnetRouteTable retrieves the route table ID associated with a specific subnet.
func (c *Connection) GetSubnetRouteTable(subnetID string) (string, error) {
	if subnetID == "" {
		return "", errors.New("no subnet ID provided. Usage: GetSubnetRouteTable <subnet-id>")
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
		return "", fmt.Errorf("error fetching route table for subnet %s: %v", subnetID, err)
	}

	if len(result.RouteTables) == 0 {
		return "", fmt.Errorf("no route table found for subnet %s", subnetID)
	}

	return *result.RouteTables[0].RouteTableId, nil
}

// GetVPCID retrieves the information of a VPC with the provided name.
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

// IsSubnetPublic checks whether the provided subnet ID
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
func (c *Connection) IsSubnetPublic(subnetID string) (bool, error) {
	// Ensure the subnet exists before determining if it's public
	if err := c.checkResourceExistence("subnet", subnetID); err != nil {
		return false, err
	}

	routeTableID, err := c.GetSubnetRouteTable(subnetID)
	if err != nil {
		// Handle the case where there's no route table for the subnet
		if strings.Contains(err.Error(), "no route table found") {
			return false, nil
		}
		return false, err
	}

	input := &ec2.DescribeRouteTablesInput{
		RouteTableIds: []*string{aws.String(routeTableID)},
	}

	result, err := c.Client.DescribeRouteTables(input)
	if err != nil {
		return false, fmt.Errorf("error describing route table %s: %v", routeTableID, err)
	}

	// Check if result.RouteTables is not nil and has at least one entry
	if result.RouteTables == nil || len(result.RouteTables) == 0 {
		return false, fmt.Errorf("no route tables found for route table ID %s", routeTableID)
	}

	// Check if Routes is not nil and has at least one entry
	if result.RouteTables[0].Routes == nil || len(result.RouteTables[0].Routes) == 0 {
		return false, fmt.Errorf("no routes found in route table %s", routeTableID)
	}

	for _, route := range result.RouteTables[0].Routes {
		// Check if route.GatewayId is not nil before dereferencing
		if route.GatewayId != nil && strings.HasPrefix(*route.GatewayId, "igw-") {
			return true, nil
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

// ListSecurityGroupsForSubnet lists all security groups
// for the provided subnet ID.
//
// **Parameters:**
//
// subnetID: the ID of the subnet to use
//
// **Returns:**
//
// []*ec2.SecurityGroup: all security groups for the provided subnet ID
//
// error: an error if any issue occurs while trying to list the security groups
func (c *Connection) ListSecurityGroupsForSubnet(subnetID string) ([]*ec2.SecurityGroup, error) {
	if err := c.checkResourceExistence("subnet", subnetID); err != nil {
		return nil, err
	}
	input := &ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("ip-permission.cidr"),
				Values: []*string{
					aws.String(subnetID),
				},
			},
		},
	}

	result, err := c.Client.DescribeSecurityGroups(input)
	if err != nil {
		return nil, err
	}

	return result.SecurityGroups, nil
}

// ListVPCSubnets lists subnets for the provided VPC name and subnet location.
//
// **Parameters:**
//
// vpcID: the ID of the VPC to use.
// subnetLocation: the location of the subnet. Can be "public", "private", or "all".
//
// **Returns:**
//
// []*ec2.Subnet: the list of subnets for the provided VPC name and location
//
// error: an error if any issue occurs while trying to list the subnets
func (c *Connection) ListVPCSubnets(vpcID string, subnetLocation string) ([]*ec2.Subnet, error) {
	// Validate VPC existence
	if err := c.checkResourceExistence("vpc", vpcID); err != nil {
		return nil, err
	}

	// Validate subnetLocation
	if subnetLocation != "public" && subnetLocation != "private" && subnetLocation != "all" {
		return nil, errors.New("subnetLocation must be public, private, or all")
	}

	// Build the subnet filter based on subnetLocation
	var filters []*ec2.Filter

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

	var subnets []*ec2.Subnet
	for _, subnet := range result.Subnets {
		if subnetLocation == "all" {
			subnets = append(subnets, subnet)
			continue
		}

		isPublic, err := c.IsSubnetPublic(*subnet.SubnetId)
		if err != nil {
			if subnetLocation == "private" && isNoRouteTableError(err) {
				// Verify if the subnet is truly private by checking all route tables
				isReallyPrivate, verifyErr := c.isSubnetReallyPrivate(*subnet.SubnetId)
				if verifyErr != nil {
					return nil, verifyErr
				}
				if isReallyPrivate {
					subnets = append(subnets, subnet)
					continue
				}
			}
			return nil, fmt.Errorf("error checking if subnet %s is publicly routable: %v", *subnet.SubnetId, err)
		}

		if (subnetLocation == "public" && isPublic) || (subnetLocation == "private" && !isPublic) {
			subnets = append(subnets, subnet)
		}
	}

	return subnets, nil
}

// isSubnetReallyPrivate checks all route tables to confirm if a subnet is truly private.
func (c *Connection) isSubnetReallyPrivate(subnetID string) (bool, error) {
	routeTables, err := c.Client.DescribeRouteTables(&ec2.DescribeRouteTablesInput{})
	if err != nil {
		return false, fmt.Errorf("error describing route tables: %v", err)
	}

	for _, routeTable := range routeTables.RouteTables {
		for _, association := range routeTable.Associations {
			if association.SubnetId != nil && *association.SubnetId == subnetID {
				for _, route := range routeTable.Routes {
					if route.GatewayId != nil && strings.HasPrefix(*route.GatewayId, "igw-") {
						return false, nil // Subnet has a route to an IGW, so it's not private
					}
				}
			}
		}
	}

	return true, nil // No routes to an IGW found, subnet is private
}

// isNoRouteTableError checks if the error is due to a missing route table, which is a common scenario for private subnets
func isNoRouteTableError(err error) bool {
	// Adjust the condition to match the specific error message or error type you're receiving for no route table
	return strings.Contains(err.Error(), "no route table found")
}

// ListVPCs lists all VPCs.
//
// **Returns:**
//
// []*ec2.Vpc: all VPCs
//
// error: an error if any issue occurs while trying to list the VPCs
func (c *Connection) ListVPCs() ([]*ec2.Vpc, error) {
	input := &ec2.DescribeVpcsInput{}

	result, err := c.Client.DescribeVpcs(input)
	if err != nil {
		return nil, err
	}

	return result.Vpcs, nil
}
