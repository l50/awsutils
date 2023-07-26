package ec2

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// Connection provides a connection
// to AWS EC2.
//
// **Attributes:**
//
// Client: the EC2 client
type Connection struct {
	Client *ec2.EC2
}

// Params provides information
// about an EC2 instance.
//
// **Attributes:**
//
// AssociatePublicIPAddress: whether to associate a public IP address
// ImageID: the ID of the AMI to use
// InstanceProfile: the name of the instance profile to use
// InstanceType: the type of the instance to use
// MinCount: the minimum number of instances to launch
// MaxCount: the maximum number of instances to launch
// SecurityGroupIDs: the IDs of the security groups to use
// KeyName: the name of the key pair to use
// SubnetID: the ID of the subnet to use
// VolumeSize: the size of the volume to use
// InstanceName: the name of the instance to use
type Params struct {
	AssociatePublicIPAddress bool
	ImageID                  string
	InstanceProfile          string
	InstanceType             string
	MinCount                 int
	MaxCount                 int
	SecurityGroupIDs         []string
	KeyName                  string
	SubnetID                 string
	VolumeSize               int64
	InstanceName             string
}

// AMIInfo provides information
// about an AMI.
//
// **Attributes:**
//
// Distro: the distro to use
// Version: the version to use
// Architecture: the architecture to use
// Region: the region to use
type AMIInfo struct {
	Distro       string
	Version      string
	Architecture string
	Region       string
}

// NewConnection creates a new connection
// to AWS EC2.
//
// **Returns:**
//
// *Connection: a new connection to AWS EC2
func NewConnection() *Connection {
	sess := session.Must(session.NewSessionWithOptions(
		session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))

	svc := ec2.New(sess)

	return &Connection{Client: svc}
}

// IsEC2Instance checks whether the code is running on an AWS
// EC2 instance by checking the existence of the file
// /sys/devices/virtual/dmi/id/product_uuid. If the file exists,
// the code is running on an EC2 instance, and the function
// returns true. If the file does not exist, the function returns false,
// indicating that the code is not running on an EC2 instance.
//
// **Returns:**
//
// bool: A boolean value that indicates whether the code is running on an EC2 instance.
func IsEC2Instance() bool {
	// Check for the existence of the product_uuid file. If it exists, we're on an EC2 instance.
	if _, err := os.Stat("/sys/devices/virtual/dmi/id/product_uuid"); err == nil {
		return true
	}

	return false
}

// CreateInstance creates a new EC2 instance
// with the provided parameters.
//
// **Parameters:**
//
// ec2Params: the parameters to use
//
// **Returns:**
//
// *ec2.Reservation: the reservation of the created instance
//
// error: an error if any issue occurs while trying to create the instance
func (c *Connection) CreateInstance(ec2Params Params) (*ec2.Reservation, error) {
	input := &ec2.RunInstancesInput{
		BlockDeviceMappings: c.getBlockDeviceMappings(ec2Params),
		IamInstanceProfile:  c.getIAMInstanceProfile(ec2Params),
		ImageId:             aws.String(ec2Params.ImageID),
		InstanceType:        aws.String(ec2Params.InstanceType),
		MinCount:            aws.Int64(int64(ec2Params.MinCount)),
		MaxCount:            aws.Int64(int64(ec2Params.MaxCount)),
		NetworkInterfaces:   c.getNetworkInterfaces(ec2Params),
		TagSpecifications:   c.getTagSpecifications(ec2Params),
	}

	result, err := c.Client.RunInstances(input)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// CheckInstanceExists checks whether an instance
// with the provided ID exists.
//
// **Parameters:**
//
// instanceID: the ID of the instance to check
//
// **Returns:**
//
// error: an error if any issue occurs while trying to check the instance
func (c *Connection) CheckInstanceExists(instanceID string) error {
	instances, err := c.GetInstances(nil)
	if err != nil {
		return err
	}

	for _, instance := range instances {
		if *instance.InstanceId == instanceID {
			return nil
		}
	}

	return fmt.Errorf("instance %s does not exist", instanceID)
}

// TagInstance tags an instance with the provided key and value.
//
// **Parameters:**
//
// instanceID: the ID of the instance to tag
//
// tagKey: the key of the tag to use
//
// tagValue: the value of the tag to use
//
// **Returns:**
//
// error: an error if any issue occurs while trying to tag the instance
func (c *Connection) TagInstance(instanceID string, tagKey string, tagValue string) error {
	input := &ec2.CreateTagsInput{
		Resources: []*string{&instanceID},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String(tagKey),
				Value: aws.String(tagValue),
			},
		},
	}

	_, err := c.Client.CreateTags(input)
	if err != nil {
		return err
	}

	return nil
}

// DestroyInstance destroys the instance with the provided ID.
//
// **Parameters:**
//
// instanceID: the ID of the instance to destroy
//
// **Returns:**
//
// error: an error if any issue occurs while trying to destroy the instance
func (c *Connection) DestroyInstance(instanceID string) error {
	input := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{&instanceID},
	}

	_, err := c.Client.TerminateInstances(input)
	if err != nil {
		return err
	}

	return nil
}

// GetRunningInstances retrieves all running instances.
//
// **Returns:**
//
// *ec2.DescribeInstancesOutput: the output of the DescribeInstances operation
//
// error: an error if any issue occurs while trying to retrieve the running instances
func (c *Connection) GetRunningInstances() (*ec2.DescribeInstancesOutput, error) {
	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String("running")},
			},
		},
	}

	result, err := c.Client.DescribeInstances(input)
	if err != nil {
		return nil, err
	}

	return result, err
}

// WaitForInstance waits until the instance with the provided ID
// is in the running state.
//
// **Parameters:**
//
// instanceID: the ID of the instance to wait for
//
// **Returns:**
//
// error: an error if any issue occurs while trying to wait for the instance
func (c *Connection) WaitForInstance(instanceID string) error {
	input := &ec2.DescribeInstanceStatusInput{
		InstanceIds: []*string{&instanceID},
	}

	err := c.Client.WaitUntilInstanceStatusOk(input)
	if err != nil {
		return err
	}

	return nil
}

// GetInstancePublicIP retrieves the public IP address of the instance
// with the provided ID.
//
// **Parameters:**
//
// instanceID: the ID of the instance to use
//
// **Returns:**
//
// string: the public IP address of the instance
//
// error: an error if any issue occurs while trying to retrieve the public IP address
func (c *Connection) GetInstancePublicIP(instanceID string) (string, error) {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{aws.String(instanceID)},
	}

	result, err := c.Client.DescribeInstances(input)
	if err != nil {
		return "", err
	}

	return *result.Reservations[0].
		Instances[0].
		NetworkInterfaces[0].
		Association.
		PublicIp, nil
}

// GetRegion retrieves the region of the connection.
//
// **Returns:**
//
// string: the region of the connection
//
// error: an error if any issue occurs while trying to retrieve the region
func (c *Connection) GetRegion() (string, error) {
	region := c.Client.Config.Region
	if region == nil {
		return "", errors.New("failed to retrieve region")
	}

	return *region, nil
}

// GetInstances retrieves all instances matching the provided filters.
//
// **Parameters:**
//
// filters: the filters to use
//
// **Returns:**
//
// []*ec2.Instance: the instances matching the provided filters
//
// error: an error if any issue occurs while trying to retrieve the instances
func (c *Connection) GetInstances(filters []*ec2.Filter) ([]*ec2.Instance, error) {
	input := &ec2.DescribeInstancesInput{
		Filters: filters,
	}

	result, err := c.Client.DescribeInstances(input)
	if err != nil {
		return nil, err
	}

	var instances []*ec2.Instance
	for _, reservation := range result.Reservations {
		instances = append(instances, reservation.Instances...)
	}

	return instances, nil
}

// GetInstanceState retrieves the state of the instance with the provided ID.
//
// **Parameters:**
//
// instanceID: the ID of the instance to use
//
// **Returns:**
//
// string: the state of the instance
//
// error: an error if any issue occurs while trying to retrieve the state
func (c *Connection) GetInstanceState(instanceID string) (string, error) {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{aws.String(instanceID)},
	}

	result, err := c.Client.DescribeInstances(input)
	if err != nil {
		return "", err
	}

	return *result.Reservations[0].Instances[0].State.Name, nil
}

// GetInstancesRunningForMoreThan24Hours retrieves all instances
// that have been running for more than 24 hours.
//
// **Returns:**
//
// []*ec2.Instance: the instances that have been running for more than 24 hours
//
// error: an error if any issue occurs while trying to retrieve the instances
func (c *Connection) GetInstancesRunningForMoreThan24Hours() ([]*ec2.Instance, error) {
	instances, err := c.GetInstances(nil)
	if err != nil {
		return nil, err
	}

	var instancesOver24Hours []*ec2.Instance
	for _, instance := range instances {
		if instance.LaunchTime.Before(time.Now().Add(-24 * time.Hour)) {
			instancesOver24Hours = append(instancesOver24Hours, instance)
		}
	}

	return instancesOver24Hours, nil
}

// GetLatestAMI retrieves the latest Amazon Machine Image (AMI) for a
// specified distribution, version and architecture. It utilizes AWS SDK
// to query AWS EC2 for the AMIs matching the provided pattern and returns
// the latest one based on the creation date.
//
// **Parameters:**
//
// info: An AMIInfo struct containing necessary details like Distro,
// Version, Architecture, and Region for which the AMI needs to be retrieved.
//
// **Returns:**
//
// string: The ID of the latest AMI found based on the provided information.
//
// error: An error if any issue occurs while trying to get the latest AMI.
func (c *Connection) GetLatestAMI(info AMIInfo) (string, error) {
	versionToAMIName := map[string]map[string]map[string]string{
		"ubuntu": {
			"22.04": {
				"amd64": "ubuntu/images/hvm-ssd/ubuntu-jammy-%s-amd64-server-*",
				"arm64": "ubuntu/images/hvm-ssd/ubuntu-jammy-%s-arm64-server-*",
			},
			"20.04": {
				"amd64": "ubuntu/images/hvm-ssd/ubuntu-focal-%s-amd64-server-*",
				"arm64": "ubuntu/images/hvm-ssd/ubuntu-focal-%s-arm64-server-*",
			},
			"18.04": {
				"amd64": "ubuntu/images/hvm-ssd/ubuntu-bionic-%s-amd64-server-*",
				"arm64": "ubuntu/images/hvm-ssd/ubuntu-bionic-%s-arm64-server-*",
			},
		},
		"centos": {
			"7": {
				"x86_64": "CentOS Linux %s x86_64 HVM EBS*",
				"arm64":  "CentOS Linux %s arm64 HVM EBS*",
			},
			"8": {
				"x86_64": "CentOS %s AMI*",
				"arm64":  "CentOS %s ARM64 AMI*",
			},
		},
		"debian": {
			"10": {
				"amd64": "debian-%s-buster-hvm-amd64-gp2*",
				"arm64": "debian-%s-buster-hvm-arm64-gp2*",
			},
		},
		"kali": {
			"2023.1": {
				"amd64": "kali-linux-%s-amd64*",
				"arm64": "kali-linux-%s-arm64*",
			},
		},
	}

	distToOwner := map[string]string{
		"ubuntu": "099720109477", // Canonical
		"debian": "136693071363", // Debian
		"kali":   "679593333241", // Kali Linux
	}

	owner, ok := distToOwner[info.Distro]
	if !ok {
		return "", fmt.Errorf("unsupported distribution: %s", info.Distro)
	}

	amiNamePattern, ok := versionToAMIName[info.Distro][info.Version][info.Architecture]
	if !ok {
		return "", fmt.Errorf("unsupported distribution/version/architecture: %s/%s/%s", info.Distro, info.Version, info.Architecture)
	}

	amiNamePattern = fmt.Sprintf(amiNamePattern, info.Version)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(info.Region),
	})

	if err != nil {
		return "", err
	}

	svc := ec2.New(sess)

	input := &ec2.DescribeImagesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("name"),
				Values: []*string{aws.String(amiNamePattern + "*")},
			},
		},
		Owners: []*string{aws.String(owner)},
	}

	result, err := svc.DescribeImages(input)
	if err != nil {
		return "", err
	}

	if len(result.Images) == 0 {
		return "", fmt.Errorf("no images found for distro: %s, version: %s, "+
			"architecture: %s", info.Distro, info.Version, info.Architecture)
	}

	// Sort images by CreationDate in descending order
	sort.Slice(result.Images, func(i, j int) bool {
		iTime, _ := time.Parse(time.RFC3339, *result.Images[i].CreationDate)
		jTime, _ := time.Parse(time.RFC3339, *result.Images[j].CreationDate)
		return iTime.After(jTime)
	})

	// Get the latest image (first image after sorting in descending order)
	latestImage := result.Images[0]

	return *latestImage.ImageId, nil
}

// FindOverlyPermissiveInboundRules checks if a specific security group permits all inbound traffic.
// Specifically, it checks if the security group has an inbound rule with the IP protocol set to "-1",
// which allows all IP traffic. This is useful for identifying security groups
// that are configured with lenient security rules, especially in testing environments.
// The function uses AWS SDK to describe security groups in AWS EC2 and checks their inbound rules.
//
// **Parameters:**
//
// secGrpID: A string containing the ID of the security group which needs to be checked for the all traffic inbound rule.
//
// **Returns:**
//
// bool: A boolean value indicating whether the security group permits all inbound traffic or not.
//
// error: An error if any issue occurs while trying to describe the security group or check its inbound rules.
func (c *Connection) FindOverlyPermissiveInboundRules(secGrpID string) (bool, error) {
	input := &ec2.DescribeSecurityGroupsInput{
		GroupIds: []*string{aws.String(secGrpID)},
	}

	resp, err := c.Client.DescribeSecurityGroups(input)
	if err != nil {
		return false, err
	}

	for _, group := range resp.SecurityGroups {
		for _, permission := range group.IpPermissions {
			if *permission.IpProtocol == "-1" {
				return true, nil
			}
		}
	}

	return false, nil
}

// ListSecurityGroups lists all security groups.
//
// **Returns:**
//
// []*ec2.SecurityGroup: all security groups
//
// error: an error if any issue occurs while trying to list the security groups
func (c *Connection) ListSecurityGroups() ([]*ec2.SecurityGroup, error) {
	input := &ec2.DescribeSecurityGroupsInput{}

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
// []string: the IDs of the security groups for the provided subnet ID
//
// error: an error if any issue occurs while trying to list the security groups
func (c *Connection) ListSecurityGroupsForSubnet(subnetID string) ([]string, error) {
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

	var groupIDs []string
	for _, group := range result.SecurityGroups {
		groupIDs = append(groupIDs, *group.GroupId)
	}

	return groupIDs, nil
}

// ListSecurityGroupsForVpc lists all security groups for the provided VPC ID.
//
// **Parameters:**
//
// vpcID: the ID of the VPC to use
//
// **Returns:**
//
// []string: the IDs of the security groups for the provided VPC ID
//
// error: an error if any issue occurs while trying to list the security groups
func (c *Connection) ListSecurityGroupsForVpc(vpcID string) ([]string, error) {
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

	var groupIDs []string
	for _, group := range result.SecurityGroups {
		groupIDs = append(groupIDs, *group.GroupId)
	}

	return groupIDs, nil
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
// error: an error if any issue occurs while trying to check whether the provided subnet ID is publicly routable
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
// error: an error if any issue occurs while trying to retrieve the ID of the subnet with the provided name
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
// vpcName: the name of the VPC to use
//
// **Returns:**
//
// string: the ID of the VPC with the provided name
//
// error: an error if any issue occurs while trying to retrieve the ID of the VPC with the provided name
func (c *Connection) GetVPCID(vpcName string) (string, error) {
	input := &ec2.DescribeVpcsInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:Name"),
				Values: []*string{
					aws.String(vpcName),
				},
			},
		},
	}

	result, err := c.Client.DescribeVpcs(input)
	if err != nil {
		return "", err
	}

	if len(result.Vpcs) == 0 {
		return "", errors.New("no VPC found with the provided name")
	}

	return *result.Vpcs[0].VpcId, nil
}

// Helper function to check if a given resource exists
func (c *Connection) checkResourceExistence(resourceName, resourceID string) error {
	switch resourceName {
	case "subnet":
		input := &ec2.DescribeSubnetsInput{
			SubnetIds: []*string{aws.String(resourceID)},
		}
		_, err := c.Client.DescribeSubnets(input)
		return err
	case "vpc":
		input := &ec2.DescribeVpcsInput{
			VpcIds: []*string{aws.String(resourceID)},
		}
		_, err := c.Client.DescribeVpcs(input)
		return err
	default:
		return errors.New("unsupported resource type")
	}
}

func (c *Connection) getBlockDeviceMappings(ec2Params Params) []*ec2.BlockDeviceMapping {
	return []*ec2.BlockDeviceMapping{
		{
			DeviceName: aws.String("/dev/sdh"),
			Ebs: &ec2.EbsBlockDevice{
				VolumeSize: aws.Int64(ec2Params.VolumeSize),
			},
		},
	}
}

func (c *Connection) getIAMInstanceProfile(ec2Params Params) *ec2.IamInstanceProfileSpecification {
	return &ec2.IamInstanceProfileSpecification{
		Name: aws.String(ec2Params.InstanceProfile),
	}
}

func (c *Connection) getNetworkInterfaces(ec2Params Params) []*ec2.InstanceNetworkInterfaceSpecification {
	return []*ec2.InstanceNetworkInterfaceSpecification{
		{
			AssociatePublicIpAddress: aws.Bool(ec2Params.AssociatePublicIPAddress),
			DeviceIndex:              aws.Int64(int64(0)),
			SubnetId:                 aws.String(ec2Params.SubnetID),
			Groups:                   aws.StringSlice(ec2Params.SecurityGroupIDs),
		},
	}
}

func (c *Connection) getTagSpecifications(ec2Params Params) []*ec2.TagSpecification {
	return []*ec2.TagSpecification{
		{
			ResourceType: aws.String("instance"),
			Tags: []*ec2.Tag{
				{
					Key:   aws.String("Name"),
					Value: aws.String(ec2Params.InstanceName),
				},
			},
		},
	}
}
