package dynamodb

import (
	"time"

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
	config *aws.Config

	// TableName defines the DynamoDb table which should be used.
	tableName *string

	// Logger will write logs for errors and and other messages depending pn used log level.
	logger log.Logger

	// dynamoDbClient is a used to access DynamoDb apis.
	dynamoDbClient *dynamodb.DynamoDB

	// lockTtl defines the life time of a lock.
	lockTtl time.Duration
}

// QueryRequest is used to query items for a partition key.
type QueryRequest struct {

	// ObjectType items should be read for.
	ObjectType string

	// Items will contain a pointer to a slice of desired items.
	// It's used to define the type of items which should be returned.
	Items interface{}
}

// ItemLock is a lock for an item in DynamoDb.
type ItemLock struct {

	// ItemIdentifier for locked item.
	*ItemIdentifier

	// ExpiresAt is the life time of a lock in epoch seconds.
	ExpiresAt int64

	// LockId is an id to identify a lock.
	LockId string
}
