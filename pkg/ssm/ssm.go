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

// Connection contains all of the relevant
// information to maintain
// an SSM connection.
type Connection struct {
	Client  ssmiface.SSMAPI
	Session *session.Session
	Params  Params
}

// Params provides parameter
// options for SSM.
type Params struct {
	Name      string
	Value     string
	Type      string
	Overwrite bool
}

// createClient is a helper function that
// returns a new ssm session.
func createClient() (ssmiface.SSMAPI, *session.Session) {
	sess := session.Must(session.NewSessionWithOptions(
		session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))

	// Create SSM service client
	svc := ssm.New(sess)

	return svc, sess
}

// AgentReady checks if an SSM agent is ready.
// Inputs:
//
//	svc is an Amazon SSM service client
//	name is the name of the parameter
//
// Output:
//
//	If success, return true and nil
//	Otherwise, return false and an error from the call to DescribeInstanceInformation
func AgentReady(svc ssmiface.SSMAPI, instanceID string, waitTime time.Duration) (bool, error) {
	// Timeout after waitSeconds seconds
	timeout := time.After(waitTime * time.Second)
	// This for a test. We're fine.
	//nolint:all
	ticker := time.Tick(500 * time.Millisecond)
	input := &ssm.DescribeInstanceInformationInput{}

	for {
		select {
		case <-timeout:
			return false, errors.New("timed out")
		case <-ticker:
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

// CreateConnection creates a connection
// with SSM and returns it.
func CreateConnection() Connection {
	ssmConnection := Connection{}
	ssmConnection.Client, ssmConnection.Session = createClient()

	return ssmConnection
}

// DeleteParam deletes a parameter in SSM
// Inputs:
//
//	svc is an Amazon SSM service client
//	name is the name of the parameter
//
// Output:
//
//	If success, information about the parameter and nil
//	Otherwise, nil and an error from the call to DeleteParam
func DeleteParam(svc ssmiface.SSMAPI, name string) error {
	_, err := svc.DeleteParameter(&ssm.DeleteParameterInput{
		Name: aws.String(name),
	})

	return err
}

// PutParam creates a parameter in SSM
// Inputs:
//
//	svc is an Amazon SSM service client
//	name is the name of the parameter
//	value is the value of the parameter
//	type is the type of parameter
//	overwrite sets the flag to rewrite
//	a parameter value
//
// Output:
//
//	If success, information about the parameter and nil
//	Otherwise, nil and an error from the call to PutParam
func PutParam(svc ssmiface.SSMAPI, name string, value string, paramType string, overwrite bool) error {
	_, err := svc.PutParameter(&ssm.PutParameterInput{
		Name:      aws.String(name),
		Value:     aws.String(value),
		Type:      aws.String(paramType),
		Overwrite: aws.Bool(overwrite),
	})

	return err
}

// GetParam fetches details of a parameter in SSM
// Inputs:
//
//	svc is an Amazon SSM service client
//	name is the name of the parameter
//	value is the value of the parameter
//	paramType is the type of parameter
//
// Output:
//
//	If success, information about the parameter and nil
//	Otherwise, nil and an error from the call to GetParam
func GetParam(svc ssmiface.SSMAPI, name string) (string, error) {

	results, err := svc.GetParameter(&ssm.GetParameterInput{
		Name: aws.String(name),
	})

	if err != nil {
		return "", err
	}

	return *results.Parameter.Value, err
}

// RunCommand runs an input command using SSM.
// Inputs:
//
//	svc is an Amazon SSM service client
//	instanceID is the instance to run the command on
//	command is the command to run
//
// Output:
//
//	If successful, the command output and nil will be returned.
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
	fmt.Printf("Now running %s", command[0])
	// Get output and check it for twenty iterations
	var i int
	for i = 0; i < 20; i++ {

		// Sleep for five seconds before attempting to
		// retrieve output.
		time.Sleep(5 * time.Second)
		output, _ := svc.GetCommandInvocation(&ssm.GetCommandInvocationInput{
			CommandId:  aws.String(commandID),
			InstanceId: aws.String(instanceID),
		})

		// Return command output if it's available
		if output.Status != nil {
			if *output.Status != "InProgress" {
				return *output.StandardOutputContent, nil
			}
		}

		i++
	}

	return "", errors.New("failed to run command")
}
