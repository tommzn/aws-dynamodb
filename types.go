package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	log "github.com/tommzn/go-log"
)

// ItemIdentifier can be used in objects which should be persisted in DynamoDb.
type ItemIdentifier struct {
	Id         string `json:"Id"`
	ObjectType string `json:"ObjectType"`
}

// DynamoDbRepository is a wrapper to AWS DynamoDb SDK.
type DynamoDbRepository struct {

	// Config contains the AWS config to access DynamoDb.
	Config *aws.Config

	// TableName defines the DynamoDb table which should be used.
	TableName *string

	// Logger will write logs for errors and and other messages depending pn used log level.
	logger log.Logger

	// dynamoDbClient is a used to access DynamoDb apis.
	dynamoDbClient *dynamodb.DynamoDB
}
