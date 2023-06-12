package ssm_test

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	ec2utils "github.com/l50/awsutils/ec2"
	ssmutils "github.com/l50/awsutils/ssm"
	"github.com/l50/goutils/v2/str"
)

var (
	err       error
	ssmParams = ssmutils.Params{
		Name:      "TestParam",
		Value:     "123456",
		Type:      "String",
		Overwrite: true,
	}
	ssmConnection = ssmutils.CreateConnection()
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
	if err := ssmutils.PutParam(ssmConnection.Client,
		ssmParams.Name, ssmParams.Value,
		ssmParams.Type, ssmParams.Overwrite); err != nil {
		log.Fatalf(
			"error running CreateSSMParam(): %v",
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

func TestDeleteParam(t *testing.T) {
	if err := ssmutils.PutParam(ssmConnection.Client,
		ssmParams.Name, ssmParams.Value,
		ssmParams.Type, ssmParams.Overwrite); err != nil {
		log.Fatalf(
			"error running CreateSSMParam(): %v",
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

type MockSSMAPI struct {
	ssmiface.SSMAPI
}

func (m *MockSSMAPI) SendCommand(input *ssm.SendCommandInput) (*ssm.SendCommandOutput, error) {
	return &ssm.SendCommandOutput{
		Command: &ssm.Command{
			CommandId: aws.String("fakeCommandID"),
			Status:    aws.String("Success"),
		},
	}, nil
}

func (m *MockSSMAPI) GetCommandInvocation(input *ssm.GetCommandInvocationInput) (*ssm.GetCommandInvocationOutput, error) {
	// Mock response based on the input, return non-nil pointers
	return &ssm.GetCommandInvocationOutput{
		CommandId:             aws.String("fakeCommandID"),
		InstanceId:            aws.String("fakeInstanceID"),
		Status:                aws.String("Success"),
		StandardOutputContent: aws.String(""),
	}, nil
}

// func TestCheckAWSCLIInstalled(t *testing.T) {
// 	tests := []struct {
// 		name       string
// 		instanceID string
// 		commandID  string
// 		want       bool
// 		err        error
// 	}{
// 		{
// 			name:       "AWS CLI is installed",
// 			instanceID: "instanceID1",
// 			commandID:  "commandID-AWSCLIInstalled",
// 			want:       true,
// 			err:        nil,
// 		},
// 		{
// 			name:       "AWS CLI is not installed",
// 			instanceID: "instanceID2",
// 			commandID:  "commandID-AWSCLINotInstalled",
// 			want:       false,
// 			err:        errors.New("AWS CLI is not installed on the instance"),
// 		},
// 	}

// 	mockClient := &mockSSMClient{}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := ssmutils.CheckAWSCLIInstalled(mockClient, tt.instanceID)
// 			if err != nil && err.Error() != tt.err.Error() {
// 				t.Errorf("CheckAWSCLIInstalled() error = %v, wantErr %v", err, tt.err)
// 				return
// 			}
// 			if got != tt.want {
// 				t.Errorf("CheckAWSCLIInstalled() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

func TestCheckAWSCLIInstalled(t *testing.T) {
	mockSvc := &MockSSMAPI{}

	// Testing the CheckAWSCLIInstalled function with the mock service
	_, err := ssmutils.CheckAWSCLIInstalled(mockSvc, "i-1234567890abcdef0")
	if err != nil {
		if err.Error() != "AWS CLI is not installed on the instance" {
			t.Errorf("Unexpected error: %s", err)
		}
	}
}

func TestRunCommand(t *testing.T) {
	timeout := time.Duration(60 * time.Second)
	ssmConnection = ssmutils.CreateConnection()
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

	agentStatus, err := ssmutils.AgentReady(ssmConnection.Client, ec2Connection.Params.InstanceID, timeout)
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

	result, err := ssmutils.RunCommand(ssmConnection.Client,
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
