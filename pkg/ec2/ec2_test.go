package ec2

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	utils "github.com/l50/goutils"
)

var (
	err           error
	verbose       = false
	volumeSize, _ = utils.StringToInt64(os.Getenv("VOLUME_SIZE"))
	ec2Params     = Params{
		ImageID:          os.Getenv("AMI"),
		InstanceName:     os.Getenv("INST_NAME"),
		InstanceType:     os.Getenv("INST_TYPE"),
		MinCount:         1,
		MaxCount:         1,
		SecurityGroupIDs: []string{os.Getenv("SEC_GRP_ID")},
		SubnetID:         os.Getenv("SUBNET_ID"),
		VolumeSize:       volumeSize,
	}
	ec2Connection = Connection{}
)

func init() {
	ec2Connection.Client = createClient()
	ec2Connection.Params = ec2Params
	ec2Connection.Reservation, err = CreateInstance(
		ec2Connection.Client,
		ec2Connection.Params,
	)
	if err != nil {
		log.Fatalf(
			"error running CreateInstance(): %v",
			err,
		)
	}

	ec2Connection.Params.InstanceID = GetInstanceID(
		ec2Connection.Reservation.Instances[0],
	)

	log.Println(
		"Waiting for test instance to finish initialization - please wait",
	)

	// Wait for instance to finish
	// initialization.
	err = WaitForInstance(
		ec2Connection.Client,
		ec2Connection.Params.InstanceID,
	)
	if err != nil {
		log.Fatalf(
			"error running WaitForInstance(): %v",
			err,
		)
	}

	fmt.Printf("Successfully created instance: %s\n",
		ec2Connection.Params.InstanceID)
}

func TestTagInstance(t *testing.T) {
	err = TagInstance(
		ec2Connection.Client,
		ec2Connection.Params.InstanceID,
		"Env",
		"Prod",
	)

	if err != nil {
		t.Fatalf(
			"error running TagInstance(): %v", err)
	}
}

func TestGetRunningInstances(t *testing.T) {
	result, err := GetRunningInstances(
		ec2Connection.Client)
	log.Println("Running instance IDs:")
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			fmt.Println(*instance.InstanceId)
		}
	}

	if err != nil {
		t.Fatalf(
			"error running GetRunningInstance(): %v", err)
	}
}

func TestGetInstancePublicIP(t *testing.T) {
	ec2Connection.Params.PublicIP, err =
		GetInstancePublicIP(
			ec2Connection.Client,
			ec2Connection.Params.InstanceID,
		)

	if err != nil {
		t.Fatalf(
			"error running GetInstancePublicIP(): %v",
			err,
		)
	}

	fmt.Printf(
		"Successfully grabbed public IP: %s\n",
		ec2Connection.Params.PublicIP)
}

func TestGetRegion(t *testing.T) {
	_, err := GetRegion(ec2Connection.Client)
	if err != nil {
		t.Fatalf(
			"error running GetRegion(): %v",
			err,
		)
	}
}

func TestGetInstances(t *testing.T) {
	// Test with no filters
	instances, err := GetInstances(ec2Connection.Client, nil)
	if err != nil {
		t.Fatalf(
			"error running GetInstances() with no filters: %v",
			err,
		)
	}

	if verbose {
		log.Println("The following instances were found: ")
		for _, instance := range instances {
			fmt.Println(*instance.InstanceId)
		}
	}

	// Test with filters
	filters := []*ec2.Filter{
		{
			Name: aws.String("tag:Name"),
			Values: []*string{
				aws.String("goInstance"),
			},
		},
	}

	instances, err = GetInstances(ec2Connection.Client, filters)
	if err != nil {
		t.Fatalf(
			"error running GetInstances() with filters: %v",
			err,
		)
	}

	if verbose {
		log.Println("Using filters, the following instances were found: ")
		for _, instance := range instances {
			fmt.Println(*instance.InstanceId)
		}
	}
}

func TestDestroyInstance(t *testing.T) {
	t.Cleanup(func() {
		err = DestroyInstance(
			ec2Connection.Client,
			ec2Connection.Params.InstanceID,
		)
		if err != nil {
			t.Fatalf(
				"error running DestroyInstance(): %v",
				err,
			)
		}
	})
}
