package ec2

import (
	"fmt"
	"os"
	"testing"

	utils "github.com/l50/goutils"
)

var (
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
)

func TestCreateInstance(t *testing.T) {
	ec2client := CreateClient()

	ec2Reservation, err := CreateInstance(ec2client, ec2Params)
	if err != nil {
		t.Fatalf(
			"error running CreateInstance(): %v", err)
	}

	ec2Params.InstanceID = *ec2Reservation.Instances[0].InstanceId
	fmt.Printf("Successfully created instance: %s\n", ec2Params.InstanceID)

	err = DestroyInstance(ec2client, ec2Params.InstanceID)
	if err != nil {
		t.Fatalf(
			"error running DestroyInstance(): %v", err)
	}
}

func TestTagInstance(t *testing.T) {
	ec2client := CreateClient()

	ec2Reservation, err := CreateInstance(ec2client, ec2Params)
	if err != nil {
		t.Fatalf(
			"error running CreateInstance(): %v", err)
	}

	ec2Params.InstanceID = *ec2Reservation.Instances[0].InstanceId

	err = TagInstance(ec2client, ec2Params.InstanceID, "Env", "Prod")
	if err != nil {
		t.Fatalf(
			"error running TagInstance(): %v", err)
	}

	err = DestroyInstance(ec2client, ec2Params.InstanceID)
	if err != nil {
		t.Fatalf(
			"error running DestroyInstance(): %v", err)
	}
}

func TestDestroyInstance(t *testing.T) {
	ec2client := CreateClient()

	ec2Reservation, err := CreateInstance(ec2client, ec2Params)
	if err != nil {
		t.Fatalf(
			"error running CreateInstance(): %v", err)
	}

	ec2Params.InstanceID = *ec2Reservation.Instances[0].InstanceId

	err = DestroyInstance(ec2client, ec2Params.InstanceID)
	if err != nil {
		t.Fatalf(
			"error running DestroyInstance(): %v", err)
	}
}

func TestGetRunningInstances(t *testing.T) {
	ec2client := CreateClient()

	result, err := GetRunningInstances(ec2client)
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			fmt.Printf("Instance ID for running instance: %v\n",
				*instance.InstanceId)
		}
	}
	if err != nil {
		t.Fatalf(
			"error running GetRunningInstance(): %v", err)
	}
}

func TestGetInstancePublicIP(t *testing.T) {
	ec2client := CreateClient()

	ec2Reservation, err := CreateInstance(ec2client, ec2Params)
	if err != nil {
		t.Fatalf(
			"error running CreateInstance(): %v", err)
	}

	ec2Params.InstanceID = *ec2Reservation.Instances[0].InstanceId

	err = WaitForInstance(ec2client, ec2Params.InstanceID)
	if err != nil {
		t.Fatalf(
			"error running WaitForInstance(): %v", err)
	}

	ec2Params.PublicIP, err = GetInstancePublicIP(ec2client,
		ec2Params.InstanceID)
	if err != nil {
		t.Fatalf(
			"error running GetInstancePublicIP(): %v", err)
	}

	fmt.Printf("Successfully grabbed public IP: %s\n",
		ec2Params.PublicIP)

	err = DestroyInstance(ec2client, ec2Params.InstanceID)
	if err != nil {
		t.Fatalf(
			"error running DestroyInstance(): %v", err)
	}
}

func TestGetRegion(t *testing.T) {
	ec2client := CreateClient()

	_, err := GetRegion(ec2client)
	if err != nil {
		t.Fatalf(
			"error running GetRegion(): %v", err)
	}
}
