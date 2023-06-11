package ssm_test

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	ec2utils "github.com/l50/awsutils/ec2"
	"github.com/l50/awsutils/ssm"
	"github.com/l50/goutils/v2/str"
)

var (
	err       error
	ssmParams = ssm.Params{
		Name:      "TestParam",
		Value:     "123456",
		Type:      "String",
		Overwrite: true,
	}
	ssmConnection = ssm.CreateConnection()
	verbose       bool
)

func init() {
	verbose = false
	if err != nil {
		log.Fatalf(
			"error running createClient(): %v",
			err,
		)
	}
}

func TestGetParam(t *testing.T) {
	if err := ssm.PutParam(ssmConnection.Client,
		ssmParams.Name, ssmParams.Value,
		ssmParams.Type, ssmParams.Overwrite); err != nil {
		log.Fatalf(
			"error running CreateSSMParam(): %v",
			err,
		)
	}

	result, err := ssm.GetParam(ssmConnection.Client,
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
	if err := ssm.PutParam(ssmConnection.Client,
		ssmParams.Name, ssmParams.Value,
		ssmParams.Type, ssmParams.Overwrite); err != nil {
		log.Fatalf(
			"error running CreateSSMParam(): %v",
			err,
		)
	}

	if err := ssm.DeleteParam(ssmConnection.Client, ssmParams.Name); err != nil {
		t.Fatalf(
			"error running DeleteParam(): %v",
			err,
		)
	}
}

func TestRunCommand(t *testing.T) {
	timeout := time.Duration(60 * time.Second)
	ssmConnection = ssm.CreateConnection()
	if err != nil {
		log.Fatalf(
			"error running createClient(): %v",
			err,
		)
	}
	ec2Connection := ec2utils.CreateConnection()
	volumeSize, _ := str.ToInt64(os.Getenv("VOLUME_SIZE"))
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

	agentStatus, err := ssm.AgentReady(ssmConnection.Client, ec2Connection.Params.InstanceID, timeout)
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
		"echo",
		"Hello World!",
		"My name is $(whoami)",
	}

	result, err := ssm.RunCommand(ssmConnection.Client,
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
