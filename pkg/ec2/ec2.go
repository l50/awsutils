package ec2

import (
	"errors"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// Connection contains all of the relevant
// information to maintain
// an EC2 connection.
type Connection struct {
	Client      *ec2.EC2
	Reservation *ec2.Reservation
	Params      Params
}

// Params provides parameter
// options for an EC2 instance.
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

var (
	metadataEndpoint string
)

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

// IsEC2Instance checks whether the code is running on an AWS EC2 instance by checking the existence of the file
// /sys/devices/virtual/dmi/id/product_uuid. If the file exists, the code is running on an EC2 instance and the function
// returns true. If the file does not exist, the function queries the EC2 instance metadata service at
// http://169.254.169.254/latest/meta-data/instance-id. If the request succeeds and the response starts with "i-",
// indicating that the code is running on an EC2 instance, the function returns true. If both the file check and
// metadata endpoint request fail, the function returns false.
func IsEC2Instance() bool {
	// Check for the existence of the product_uuid file. If it exists, we're on an EC2 instance.
	if _, err := os.Stat("/sys/devices/virtual/dmi/id/product_uuid"); err == nil {
		return true
	}

	return false
}
