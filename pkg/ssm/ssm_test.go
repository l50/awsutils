package ssm

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	ec2utils "github.com/l50/awsutils/pkg/ec2"
	utils "github.com/l50/goutils"
)

var (
	err       error
	ssmParams = Params{
		Name:      "TestParam",
		Value:     "123456",
		Type:      "String",
		Overwrite: true,
	}
	ssmConnection = Connection{}
	verbose       bool
)

func init() {
	verbose = false
	ssmConnection.Client, ssmConnection.Session = createClient()
	if err != nil {
		log.Fatalf(
			"error running createClient(): %v",
			err,
		)
	}

	err := PutParam(ssmConnection.Client,
		ssmParams.Name, ssmParams.Value,
		ssmParams.Type, ssmParams.Overwrite)
	if err != nil {
		log.Fatalf(
			"error running CreateSSMParam(): %v",
			err,
		)
	}
}

func TestGetParam(t *testing.T) {
	result, err := GetParam(ssmConnection.Client,
		ssmParams.Name)
	if err != nil {
		t.Fatalf(
			"error running GetParam(): %v",
			err,
		)
	}
	fmt.Println(result)
}

func TestDeleteParam(t *testing.T) {
	err := DeleteParam(ssmConnection.Client, ssmParams.Name)
	if err != nil {
		t.Fatalf(
			"error running DeleteParam(): %v",
			err,
		)
	}
}

func TestRunCommand(t *testing.T) {
	timeout := time.Duration(60 * time.Second)
	ssmConnection.Client, ssmConnection.Session = createClient()
	if err != nil {
		log.Fatalf(
			"error running createClient(): %v",
			err,
		)
	}
	ec2Connection := ec2utils.CreateConnection()
	volumeSize, _ := utils.StringToInt64(os.Getenv("VOLUME_SIZE"))
	params := ec2utils.Params{
		AssociatePublicIPAddress: true,
		ImageID:                  os.Getenv("AMI"),
		InstanceName:             os.Getenv("INST_NAME"),
		InstanceType:             os.Getenv("INST_TYPE"),
		InstanceProfile:          os.Getenv("IAM_INSTANCE_PROFILE"),
		MinCount:                 1,
		MaxCount:                 1,
		SecurityGroupIDs:         []string{os.Getenv("SEC_GRP_ID")},
		SubnetID:                 os.Getenv("SUBNET_ID"),
		VolumeSize:               volumeSize,
	}
	ec2Connection.Reservation, err = ec2utils.CreateInstance(
		ec2Connection.Client,
		params,
	)
	if err != nil {
		log.Fatalf(
			"error running CreateInstance(): %v",
			err,
		)
	}
	ec2Connection.Params.InstanceID = ec2utils.GetInstanceID(
		ec2Connection.Reservation.Instances[0],
	)

	agentStatus, err := AgentReady(ssmConnection.Client, ec2Connection.Params.InstanceID, timeout)
	if err != nil {
		t.Fatalf(
			"error running AgentReady(): %v",
			err,
		)
	}

	if agentStatus {
		fmt.Printf("Successfully created SSM-managed instance: %s\n",
			ec2Connection.Params.InstanceID)
	}

	command := []string{
		"whoami",
	}

	result, err := RunCommand(ssmConnection.Client,
		ec2Connection.Params.InstanceID, command)
	if err != nil {
		t.Fatalf(
			"error running RunCommand(): %v",
			err,
		)
	}
	if verbose {
		fmt.Println(result)
	}

	err = ec2utils.DestroyInstance(
		ec2Connection.Client,
		ec2Connection.Params.InstanceID,
	)
	if err != nil {
		t.Fatalf(
			"error running DestroyInstance(): %v",
			err,
		)
	}
}
