package dynamo

import (
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
)

// Connection contains all of the
// relevant information to maintain
// a connection with DynamoDB.
type Connection struct {
	Client *dynamodb.DynamoDB
	Params Params
}

// Params provides parameter
// options for a DynamoDB table.
type Params struct {
	ID       uuid.UUID
	Created  time.Time
	Modified time.Time
}

// createClient is a helper function that
// returns a new dynamo session.
func createClient() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSessionWithOptions(
		session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	return svc
}

// CreateConnection creates a connection
// with DynamoDB and returns a Connection.
func CreateConnection() Connection {
	dynamoConnection := Connection{}
	dynamoConnection.Client = createClient()

	return dynamoConnection
}
