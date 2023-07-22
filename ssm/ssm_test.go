package ssm_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	ec2utils "github.com/l50/awsutils/ec2"
	ssmutils "github.com/l50/awsutils/ssm"
	"github.com/l50/goutils/v2/str"
	"github.com/l50/goutils/v2/sys"
)

var (
	ssmParams = ssmutils.Params{
		Name:      "TestParam",
		Value:     "123456",
		Type:      "String",
		Overwrite: true,
	}
	ssmConnection     = ssmutils.CreateConnection()
	testEC2Connection *ec2utils.Connection
	testInstanceID    string
	reservation       *ec2.Reservation
	testParams        ec2utils.Params
)

func TestMain(m *testing.M) {
	requiredEnvVars := []string{
		"AMI",
		"VOLUME_SIZE",
		"INST_NAME",
		"INST_TYPE",
		"IAM_INSTANCE_PROFILE",
		"SEC_GRP_ID",
		"SUBNET_ID",
	}
	for _, ev := range requiredEnvVars {
		if err := sys.EnvVarSet(ev); err != nil {
			fmt.Printf("error setting environment variable %s: %v\n", ev, err)
			os.Exit(1)
		}
	}

	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	volumeSize, err := str.ToInt64(os.Getenv("VOLUME_SIZE"))
	if err != nil {
		log.Fatalf("error converting volume size to int64: %v", err)
	}
	testParams = ec2utils.Params{
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
	testEC2Connection = ec2utils.NewConnection()
	reservation, err = testEC2Connection.CreateInstance(testParams)
	if err != nil {
		fmt.Printf("failed to create instance: %v", err)
		os.Exit(1)
	}

	// Store the instance ID in a global variable for other tests to use
	testInstanceID = *reservation.Instances[0].InstanceId

	// Wait for the instance to be ready
	if err := testEC2Connection.WaitForInstance(testInstanceID); err != nil {
		log.Fatalf("error waiting for instance to be ready: %v", err)
	}

	// Double check that the instance is running
	state, err := testEC2Connection.GetInstanceState(testInstanceID)
	if err != nil {
		log.Fatalf("error getting instance state: %v", err)
	}
	if state != "running" {
		log.Fatalf("instance is not running: %v", err)
	}

	if len(reservation.Instances) == 0 {
		fmt.Println("No instances found in reservation")
		os.Exit(1)
	}

	// Schedule the instance to be destroyed after the test ends
	defer func() {
		err := testEC2Connection.DestroyInstance(testInstanceID)
		if err != nil {
			log.Fatalf("failed to destroy instance: %v", err)
		}
	}()
}

func teardown() {
	err := testEC2Connection.DestroyInstance(testInstanceID)
	if err != nil {
		fmt.Printf("failed to destroy instance: %v", err)
		os.Exit(1)
	}
}

var tests = []struct {
	name string
}{
	{
		name: "TestGetParam",
	},
	{
		name: "TestDeleteParam",
	},
	// {
	// 	name: "TestCheckAWSCLIInstalled",
	// },
	{
		name: "TestRunCommand",
	},
}

func TestAWSUtils(t *testing.T) {
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Run the setup and teardown logic for each subtest
			setup()
			defer teardown()

			switch tc.name {
			// Test each of these individually with:
			// go test -run TestAWSUtils/TestGetParam
			case "TestGetParam":
				testGetParam(t)
			case "TestDeleteParam":
				testDeleteParam(t)
			}
		})
	}
}

func testGetParam(t *testing.T) {
	if err := ssmutils.PutParam(ssmConnection.Client,
		ssmParams.Name, ssmParams.Value,
		ssmParams.Type, ssmParams.Overwrite); err != nil {
		t.Fatalf(
			"error running PutParam(): %v",
			err,
		)
	}

	result, err := ssmutils.GetParam(ssmConnection.Client,
		ssmParams.Name)
	if err != nil {
		t.Fatalf(
			"error running GetParam(): %v",
			err,
		)
	}
	fmt.Println(result)
}

func testDeleteParam(t *testing.T) {
	if err := ssmutils.PutParam(ssmConnection.Client,
		ssmParams.Name, ssmParams.Value,
		ssmParams.Type, ssmParams.Overwrite); err != nil {
		t.Fatalf(
			"error running PutParam(): %v",
			err,
		)
	}

	if err := ssmutils.DeleteParam(ssmConnection.Client, ssmParams.Name); err != nil {
		t.Fatalf(
			"error running DeleteParam(): %v",
			err,
		)
	}
}
