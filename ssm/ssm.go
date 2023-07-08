package ssm

import (
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

// Connection represents the necessary information to maintain
// an AWS Systems Manager (SSM) connection.
//
// **Attributes:**
//
// Client:  Amazon SSM service client interface.
// Session: AWS session from which the client is derived.
// Params:  Structure with parameters for the SSM service.
type Connection struct {
	Client  ssmiface.SSMAPI
	Session *session.Session
	Params  Params
}

// Params represents parameter options for SSM.
//
// **Attributes:**
//
// Name:      Parameter name.
// Value:     Parameter value.
// Type:      Parameter type.
// Overwrite: Flag to overwrite an existing parameter.
type Params struct {
	Name      string
	Value     string
	Type      string
	Overwrite bool
}

// createClient generates a new AWS session and an SSM service client.
//
// **Returns:**
//
// ssmiface.SSMAPI: Interface for Amazon SSM service client.
// *session.Session: AWS session.
func createClient() (ssmiface.SSMAPI, *session.Session) {
	sess := session.Must(session.NewSessionWithOptions(
		session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))

	svc := ssm.New(sess)

	return svc, sess
}

// AgentReady checks if an SSM agent is ready on the instance.
//
// **Parameters:**
//
// svc: AWS SSM service client.
// instanceID: AWS EC2 instance ID to check.
// waitTime: Maximum wait time before timing out.
//
// **Returns:**
//
// bool: True if the agent is ready, false otherwise.
// error: An error if any issue occurs while checking the agent.
func AgentReady(svc ssmiface.SSMAPI, instanceID string, waitTime time.Duration) (bool, error) {
	timeout := time.After(waitTime * time.Second)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	input := &ssm.DescribeInstanceInformationInput{}

	for {
		select {
		case <-timeout:
			return false, errors.New("timed out")
		case <-ticker.C:
			data, err := svc.DescribeInstanceInformation(input)
			if err != nil {
				return false, err
			}

			if len(data.InstanceInformationList) != 0 {
				for _, d := range data.InstanceInformationList {
					if *d.InstanceId == instanceID {
						if *d.PingStatus == "Online" {
							return true, nil
						}
					}
				}
			}
		}
	}
}

// CheckAWSCLIInstalled checks if AWS CLI is installed on the instance.
//
// **Parameters:**
//
// svc: AWS SSM service client.
// instanceID: AWS EC2 instance ID to check.
//
// **Returns:**
//
// bool: True if AWS CLI is installed, false otherwise.
// error: An error if any issue occurs while checking the installation.
func CheckAWSCLIInstalled(svc ssmiface.SSMAPI, instanceID string) (bool, error) {
	command := []string{"command -v aws"}
	output, err := RunCommand(svc, instanceID, command)
	if err != nil {
		return false, err
	}

	if output == "" {
		return false, errors.New("AWS CLI is not installed on the instance")
	}

	return true, nil
}

// CreateConnection establishes a connection with AWS SSM.
//
// **Returns:**
//
// Connection: Struct with a connected SSM client and session.
func CreateConnection() Connection {
	ssmConnection := Connection{}
	ssmConnection.Client, ssmConnection.Session = createClient()

	return ssmConnection
}

// DeleteParam removes a parameter from AWS SSM.
//
// **Parameters:**
//
// svc: AWS SSM service client.
// name: Name of the parameter to delete.
//
// **Returns:**
//
// error: An error if any issue occurs while deleting the parameter.
func DeleteParam(svc ssmiface.SSMAPI, name string) error {
	_, err := svc.DeleteParameter(&ssm.DeleteParameterInput{
		Name: aws.String(name),
	})

	return err
}

// PutParam creates or updates a parameter in AWS SSM.
//
// **Parameters:**
//
// svc: AWS SSM service client.
// name: Name of the parameter.
// value: Value of the parameter.
// paramType: Type of the parameter.
// overwrite: Flag to overwrite an existing parameter.
//
// **Returns:**
//
// error: An error if any issue occurs while creating or updating the parameter.
func PutParam(svc ssmiface.SSMAPI, name string, value string, paramType string, overwrite bool) error {
	_, err := svc.PutParameter(&ssm.PutParameterInput{
		Name:      aws.String(name),
		Value:     aws.String(value),
		Type:      aws.String(paramType),
		Overwrite: aws.Bool(overwrite),
	})

	return err
}

// GetParam retrieves a parameter from AWS SSM.
//
// **Parameters:**
//
// svc: AWS SSM service client.
// name: Name of the parameter.
//
// **Returns:**
//
// string: Value of the parameter.
// error: An error if any issue occurs while fetching the parameter.
func GetParam(svc ssmiface.SSMAPI, name string) (string, error) {

	results, err := svc.GetParameter(&ssm.GetParameterInput{
		Name: aws.String(name),
	})

	if err != nil {
		return "", err
	}

	return *results.Parameter.Value, err
}

// RunCommand executes an input command on an AWS instance via SSM.
//
// **Parameters:**
//
// svc: AWS SSM service client.
// instanceID: AWS EC2 instance ID where the command should run.
// command: List of command strings to be run.
//
// **Returns:**
//
// string: Output of the command execution.
// error: An error if any issue occurs while executing the command.
func RunCommand(svc ssmiface.SSMAPI, instanceID string, command []string) (string, error) {
	params := map[string][]*string{"commands": aws.StringSlice(command)}
	docName := "AWS-RunShellScript"

	cmdInput := &ssm.SendCommandInput{
		InstanceIds:  aws.StringSlice([]string{instanceID}),
		DocumentName: aws.String(docName),
		Parameters:   params,
	}

	inputResult, err := svc.SendCommand(cmdInput)
	if err != nil {
		return "", err
	}

	commandID := *inputResult.Command.CommandId
	fmt.Printf("Now running %s on %s\n", command[0], instanceID)

	for i := 0; i < 20; i++ {
		time.Sleep(5 * time.Second)
		output, _ := svc.GetCommandInvocation(&ssm.GetCommandInvocationInput{
			CommandId:  aws.String(commandID),
			InstanceId: aws.String(instanceID),
		})
		if *output.Status != "InProgress" {
			if *output.Status == "Success" {
				return *output.StandardOutputContent, nil
			}
			return "", errors.New(*output.StandardErrorContent)
		}
	}

	return "", errors.New("command timed out")
}

// ListAllParameters retrieves all parameters in the AWS SSM.
//
// **Parameters:**
//
// svc: AWS SSM service client.
//
// **Returns:**
//
// ([]*ssm.ParameterMetadata): List of all parameters' metadata.
// error: An error if any issue occurs while fetching the parameters.
func ListAllParameters(svc ssmiface.SSMAPI) ([]*ssm.ParameterMetadata, error) {
	input := &ssm.DescribeParametersInput{}
	result, err := svc.DescribeParameters(input)

	if err != nil {
		return nil, err
	}

	return result.Parameters, nil
}
