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

// Connection contains all of the relevant
// information to maintain
// an EC2 connection.
//
// **Attributes:**
//
// Client: an EC2 session
// Reservation: an EC2 reservation
// Params: parameters for an EC2 instance
type Connection struct {
	Client      *ec2.EC2
	Reservation *ec2.Reservation
	Params      Params
}

// Params provides parameter
// options for an EC2 instance.
//
// **Attributes:**
//
// AssociatePublicIPAddress: whether or not to associate a public IP address
// ImageID: the AMI ID to use
// InstanceProfile: the IAM instance profile to use
// InstanceType: the instance type to use
// MinCount: the minimum number of instances to launch
// MaxCount: the maximum number of instances to launch
// SecurityGroupIDs: the security group IDs to use
// KeyName: the key name to use
// SubnetID: the subnet ID to use
// VolumeSize: the volume size to use
// InstanceID: the instance ID to use
// InstanceName: the instance name to use
// PublicIP: the public IP to use
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
	InstanceID               string
	InstanceName             string
	PublicIP                 string
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

// createClient is a helper function that
// returns a new ec2 session.
func createClient() *ec2.EC2 {
	sess := session.Must(session.NewSessionWithOptions(
		session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))

	// Create EC2 service client
	svc := ec2.New(sess)

	return svc
}

// CreateConnection creates a connection
// with EC2 and returns a Connection.
func CreateConnection() Connection {
	ec2Connection := Connection{}
	ec2Connection.Client = createClient()

	return ec2Connection
}

// CreateInstance returns an ec2 reservation for an instance
// that is created with the input ec2Params.
func CreateInstance(client *ec2.EC2, ec2Params Params) (*ec2.Reservation, error) {
	result, err := client.RunInstances(&ec2.RunInstancesInput{
		BlockDeviceMappings: []*ec2.BlockDeviceMapping{
			{
				DeviceName: aws.String("/dev/sdh"),
				Ebs: &ec2.EbsBlockDevice{
					VolumeSize: aws.Int64(ec2Params.VolumeSize),
				},
			},
		},
		IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
			Name: aws.String(ec2Params.InstanceProfile),
		},
		ImageId:      aws.String(ec2Params.ImageID),
		InstanceType: aws.String(ec2Params.InstanceType),
		MinCount:     aws.Int64(int64(ec2Params.MinCount)),
		MaxCount:     aws.Int64(int64(ec2Params.MaxCount)),
		// Omitted in favor of enforcing use of SSM.
		// KeyName:    aws.String(ec2Params.KeyName),
		NetworkInterfaces: []*ec2.InstanceNetworkInterfaceSpecification{
			{
				AssociatePublicIpAddress: aws.Bool(ec2Params.AssociatePublicIPAddress),
				DeviceIndex:              aws.Int64(int64(0)),
				SubnetId:                 aws.String(ec2Params.SubnetID),
				Groups:                   aws.StringSlice(ec2Params.SecurityGroupIDs),
			},
		},
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("instance"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String(ec2Params.InstanceName),
					},
				},
			},
		},
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// CheckInstanceExists checks if an EC2 instance with the given instance ID exists.
func CheckInstanceExists(client *ec2.EC2, instanceID string) error {
	instances, err := GetInstances(client, nil)
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

// TagInstance tags the instance tied to the input ID with the specified tag.
func TagInstance(client *ec2.EC2, instanceID string, tagKey string, tagValue string) error {
	_, err := client.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{&instanceID},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String(tagKey),
				Value: aws.String(tagValue),
			},
		},
	})

	if err != nil {
		return err
	}

	return nil
}

// DestroyInstance terminates the ec2 instance associated with
// the input instanceID.
func DestroyInstance(client *ec2.EC2, instanceID string) error {
	_, err := client.TerminateInstances(&ec2.TerminateInstancesInput{
		InstanceIds: []*string{&instanceID},
	})

	if err != nil {
		return err
	}

	return nil
}

// GetRunningInstances returns all ec2 instances with a state of running.
func GetRunningInstances(client *ec2.EC2) (*ec2.DescribeInstancesOutput, error) {
	result, err := client.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String("running")},
			},
		},
	})

	if err != nil {
		return nil, err
	}

	return result, err
}

// WaitForInstance waits for the input instanceID to get to
// a running state.
func WaitForInstance(client *ec2.EC2, instanceID string) error {
	instanceStatus := &ec2.DescribeInstanceStatusInput{
		InstanceIds: []*string{&instanceID},
	}
	err := client.WaitUntilInstanceStatusOk(instanceStatus)
	if err != nil {
		return err
	}

	return nil
}

// GetInstancePublicIP returns the public IP address
// of the input instanceID.
func GetInstancePublicIP(client *ec2.EC2, instanceID string) (string, error) {
	result, err := client.DescribeInstances(&ec2.DescribeInstancesInput{
		InstanceIds: []*string{aws.String(instanceID)},
	})

	if err != nil {
		return "", err
	}

	return *result.Reservations[0].
		Instances[0].
		NetworkInterfaces[0].
		Association.
		PublicIp, nil
}

// GetRegion returns the region associated with the input
// ec2 client.
func GetRegion(client *ec2.EC2) (string, error) {
	region := client.Config.Region

	if region == nil {
		return "", errors.New("failed to retrieve region")
	}

	return *region, nil
}

// GetInstanceID returns the instance ID
// from an input instanceReservation.
func GetInstanceID(instanceReservation *ec2.Instance) string {
	return *instanceReservation.InstanceId
}

// GetInstances returns ec2 instances that the
// input client has access to.
// If no filters are provided, all ec2 instances will
// be returned by default.
func GetInstances(client *ec2.EC2, filters []*ec2.Filter) (
	[]*ec2.Instance, error) {

	instances := []*ec2.Instance{}

	result, err := client.DescribeInstances(
		&ec2.DescribeInstancesInput{
			Filters: filters,
		})

	if err != nil {
		return instances, err
	}

	// Get instances from reservations and add
	// to the instances output.
	for _, reservation := range result.Reservations {
		instances = append(instances, reservation.Instances...)
	}

	return instances, nil
}

// GetInstanceState returns the state of the ec2
// instance associated with the input instanceID.
func GetInstanceState(client *ec2.EC2, instanceID string) (string, error) {
	result, err := client.DescribeInstances(
		&ec2.DescribeInstancesInput{
			InstanceIds: []*string{aws.String(instanceID)},
		})

	if err != nil {
		return "", err
	}

	return *result.Reservations[0].
		Instances[0].
		State.Name, nil
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

// GetInstancesRunningForMoreThan24Hours returns a list of all EC2 instances running
// for more than 24 hours.
func GetInstancesRunningForMoreThan24Hours(client *ec2.EC2) ([]*ec2.Instance, error) {
	// get all instances
	instances, err := GetInstances(client, nil)
	if err != nil {
		return nil, err
	}

	// filter out instances running for more than 24 hours
	var instancesOver24Hours []*ec2.Instance
	for _, instance := range instances {
		// Check if instance's LaunchTime is more than 24 hours ago
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
func GetLatestAMI(info AMIInfo) (string, error) {
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

	fmt.Println(amiNamePattern)
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
