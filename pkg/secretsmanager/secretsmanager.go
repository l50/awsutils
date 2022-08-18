package secretsmanager

import (
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

// Connection contains all of the relevant
// information to maintain
// a Secrets Manager connection.
type Connection struct {
	Client  *secretsmanager.SecretsManager
	Session *session.Session
	Params  Params
}

// Params provides parameter
// options for Secrets Manager.
type Params struct {
	Name        string
	Description string
	Value       string
	Created     time.Time
	Modified    time.Time
}

// createClient is a helper function that
// returns a new secretsmanager session.
func createClient() (*secretsmanager.SecretsManager, *session.Session) {
	sess := session.Must(session.NewSessionWithOptions(
		session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))

	// Create secrets manager service client
	svc := secretsmanager.New(sess)

	return svc, sess
}

// CreateConnection creates a connection
// with Secrets Manager and returns it.
func CreateConnection() Connection {
	smConnection := Connection{}
	smConnection.Client, smConnection.Session = createClient()

	return smConnection
}

// CreateSecret creates an input `secretName`
// with the specified `secretValue`.
func CreateSecret(client *secretsmanager.SecretsManager,
	secretName string, secretDesc string, secretValue string) error {
	_, err := client.CreateSecret(&secretsmanager.CreateSecretInput{
		Name:         aws.String(secretName),
		Description:  aws.String(secretDesc),
		SecretString: aws.String(secretValue),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeInvalidParameterException:
				return errors.New(secretsmanager.ErrCodeInvalidParameterException)
			case secretsmanager.ErrCodeInvalidRequestException:
				return errors.New(secretsmanager.ErrCodeInvalidRequestException)
			case secretsmanager.ErrCodeLimitExceededException:
				return errors.New(secretsmanager.ErrCodeLimitExceededException)
			case secretsmanager.ErrCodeEncryptionFailure:
				return errors.New(secretsmanager.ErrCodeEncryptionFailure)
			case secretsmanager.ErrCodeResourceExistsException:
				return errors.New(secretsmanager.ErrCodeResourceExistsException)
			case secretsmanager.ErrCodeResourceNotFoundException:
				return errors.New(secretsmanager.ErrCodeResourceNotFoundException)
			case secretsmanager.ErrCodeMalformedPolicyDocumentException:
				return errors.New(secretsmanager.ErrCodeMalformedPolicyDocumentException)
			case secretsmanager.ErrCodeInternalServiceError:
				return errors.New(secretsmanager.ErrCodeInternalServiceError)
			case secretsmanager.ErrCodePreconditionNotMetException:
				return errors.New(secretsmanager.ErrCodePreconditionNotMetException)
			case secretsmanager.ErrCodeDecryptionFailure:
				return errors.New(secretsmanager.ErrCodeDecryptionFailure)
			default:
				return aerr
			}
		} else {
			return err
		}
	}
	return nil
}

// DeleteSecret deletes an input `secretName`.
// It will attempt to do so forcefully if `forceDelete`
// is set to true.
func DeleteSecret(client *secretsmanager.SecretsManager, secretName string, forceDelete bool) error {
	_, err := client.DeleteSecret(&secretsmanager.DeleteSecretInput{
		SecretId:                   aws.String(secretName),
		ForceDeleteWithoutRecovery: aws.Bool(true),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeResourceNotFoundException:
				return errors.New(secretsmanager.ErrCodeResourceNotFoundException)
			case secretsmanager.ErrCodeInvalidParameterException:
				return errors.New(secretsmanager.ErrCodeInvalidParameterException)
			case secretsmanager.ErrCodeInvalidRequestException:
				return errors.New(secretsmanager.ErrCodeInvalidRequestException)
			case secretsmanager.ErrCodeInternalServiceError:
				return errors.New(secretsmanager.ErrCodeInternalServiceError)
			default:
				return aerr
			}
		} else {
			return err
		}
	}

	return nil
}
