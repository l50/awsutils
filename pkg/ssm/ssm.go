package ssm

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

// Connection contains all of the relevant
// information to maintain
// an SSM connection.
type Connection struct {
	Client  *ssm.SSM
	Session *session.Session
	Param   Param
}

// Param provides parameter
// options for SSM.
type Param struct {
	Name      string
	Value     string
	Type      string
	Overwrite bool
}

// createClient is a helper function that
// returns a new ssm session.
func createClient() (*ssm.SSM, *session.Session) {
	sess := session.Must(session.NewSessionWithOptions(
		session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))

	// Create SSM service client
	svc := ssm.New(sess)

	return svc, sess
}

// CreateConnection creates a connection
// with SSM and returns it.
func CreateConnection() Connection {
	ssmConnection := Connection{}
	ssmConnection.Client, ssmConnection.Session = createClient()

	return ssmConnection
}

// PutParam creates a parameter in SSM
// Inputs:
//     svc is an Amazon SSM service client
//     name is the name of the parameter
// Output:
//     If success, information about the parameter and nil
//     Otherwise, nil and an error from the call to GetParam
func PutParam(svc ssmiface.SSMAPI, name *string, value *string, paramType *string, overwrite *bool) (*ssm.PutParameterOutput, error) {
	results, err := svc.PutParameter(&ssm.PutParameterInput{
		Name:      name,
		Value:     value,
		Type:      paramType,
		Overwrite: overwrite,
	})

	return results, err
}

// GetParam fetches details of a parameter in SSM
// Inputs:
//     svc is an Amazon SSM service client
//     name is the name of the parameter
//     value is the value of the parameter
//     paramType is the type of parameter
// Output:
//     If success, information about the parameter and nil
//     Otherwise, nil and an error from the call to PutParam
func GetParam(svc ssmiface.SSMAPI, name *string) (*ssm.GetParameterOutput, error) {
	results, err := svc.GetParameter(&ssm.GetParameterInput{
		Name: name,
	})

	return results, err
}
