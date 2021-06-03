// Package dynamodb provides a wrapper for AWS DynamoDb to run CRUD actions in items.
package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
)

// NewRepository creates a new DynamoDb repository by passed coinfig.
// By config you can defins the table name, region and endpoint for a local dynamodb.
func NewRepository(conf config.Config, logger log.Logger) Repository {

	tableName := conf.Get("aws.dynamodb.tablename", config.AsStringPtr(DEFAULT_TABLENAME))
	awsConfig := &aws.Config{
		Region:   conf.Get("aws.dynamodb.region", config.AsStringPtr(DEFAULT_AWS_REGION)),
		Endpoint: conf.Get("aws.dynamodb.endpoint", nil),
	}
	return &DynamoDbRepository{
		config:    awsConfig,
		tableName: tableName,
		logger:    logger,
	}
}
