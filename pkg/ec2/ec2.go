package ec2

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// Params represents a particular EC2 instance
type Params struct {
	ImageID          string
	InstanceType     string
	MinCount         int
	MaxCount         int
	SecurityGroupIDs []string
	KeyName          string
	SubnetID         string
	VolumeSize       int64
	InstanceID       string
	InstanceName     string
	PublicIP         string
}

// CreateClient returns a new ec2 session.
func CreateClient() *ec2.EC2 {
	sess := session.Must(session.NewSessionWithOptions(
		session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))

	// Create EC2 service client
	svc := ec2.New(sess)

	return svc
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
		ImageId:          aws.String(ec2Params.ImageID),
		InstanceType:     aws.String(ec2Params.InstanceType),
		MinCount:         aws.Int64(int64(ec2Params.MinCount)),
		MaxCount:         aws.Int64(int64(ec2Params.MaxCount)),
		SecurityGroupIds: aws.StringSlice(ec2Params.SecurityGroupIDs),
		// Omitted in favor of enforcing use of SSM.
		// KeyName:    aws.String(ec2Params.KeyName),
		SubnetId: aws.String(ec2Params.SubnetID),
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
		InstanceIds: []*string{&instanceID},
	})

	if err != nil {
		return "", err
	}

	return *result.Reservations[0].Instances[0].NetworkInterfaces[0].Association.PublicIp, nil
}

// GetRegion returns the region associted with the input
// ec2 client.
func GetRegion(client *ec2.EC2) (string, error) {
	region := client.Config.Region

	if region == nil {
		return "", errors.New("failed to retrieve region")
	}

	return *region, nil
}
