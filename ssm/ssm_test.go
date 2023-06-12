package ssm_test

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	ec2utils "github.com/l50/awsutils/ec2"
	ssmutils "github.com/l50/awsutils/ssm"
	"github.com/l50/goutils/v2/str"
	"github.com/l50/goutils/v2/sys"
)

var (
	err       error
	ssmParams = ssmutils.Params{
		Name:      "TestParam",
		Value:     "123456",
		Type:      "String",
		Overwrite: true,
	}
	ssmConnection  = ssmutils.CreateConnection()
	verbose        bool
	testInstanceID string // shared instance ID for tests
	ec2Connection  = ec2utils.CreateConnection()
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
	volumeSize, _ := str.ToInt64(os.Getenv("VOLUME_SIZE"))
	ec2Params := ec2utils.Params{
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
	ec2Connection.Reservation, err = ec2utils.CreateInstance(ec2Connection.Client, ec2Params)
	if err != nil {
		log.Fatalf("error running CreateInstance(): %v", err)
	}

	testInstanceID = ec2utils.GetInstanceID(ec2Connection.Reservation.Instances[0])

	// Wait for the instance to be ready
	err = ec2utils.WaitForInstance(ec2Connection.Client, testInstanceID)
	if err != nil {
		log.Fatalf("error waiting for instance to be ready: %v", err)
	}

	// Double check that the instance is running
	state, err := ec2utils.GetInstanceState(ec2Connection.Client, testInstanceID)
	if err != nil {
		log.Fatalf("error getting instance state: %v", err)
	}
	if state != "running" {
		log.Fatalf("instance is not running: %v", err)
	}
}

func teardown() {
	ec2Connection = ec2utils.CreateConnection()
	err = ec2utils.DestroyInstance(ec2Connection.Client, testInstanceID)
	if err != nil {
		log.Fatalf("error running DestroyInstance(): %v", err)
	}
}

// The following test cases are defined
var tests = []struct {
	name string
}{
	{
		name: "TestGetParam",
	},
	{
		name: "TestDeleteParam",
	},
	{
		name: "TestCheckAWSCLIInstalled",
	},
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
			case "TestCheckAWSCLIInstalled":
				testCheckAWSCLIInstalled(t)
			case "TestRunCommand":
				testRunCommand(t)
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

func testCheckAWSCLIInstalled(t *testing.T) {
	if err := ec2utils.CheckInstanceExists(ec2Connection.Client, testInstanceID); err != nil {
		_, err := ssmutils.CheckAWSCLIInstalled(ssmConnection.Client, testInstanceID)
		if err != nil {
			if err.Error() != "AWS CLI is not installed on the instance" {
				t.Errorf("Unexpected error: %s", err)
			}
		}
	} else {
		t.Errorf("Instance %s does not exist", testInstanceID)
	}
}

func testRunCommand(t *testing.T) {
	timeout := time.Duration(60 * time.Second)

	agentStatus, err := ssmutils.AgentReady(ssmConnection.Client, testInstanceID, timeout)
	if err != nil {
		t.Fatalf("error running AgentReady(): %v", err)
	}

	if agentStatus {
		fmt.Printf("Successfully created SSM-managed instance: %s\n", testInstanceID)
	}

	command := []string{
		"echo",
		"Hello World!",
		"My name is $(whoami)",
	}

	result, err := ssmutils.RunCommand(ssmConnection.Client, testInstanceID, command)
	if err != nil {
		t.Fatalf("error running RunCommand(): %v", err)
	}

	if verbose {
		fmt.Println(result)
	}
}
