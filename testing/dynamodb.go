package testing

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// SetupTableForTest will create a new DynamoDb table with passed name
// and a composed primary key with an ID as hash key and an object type as sort key.
func SetupTableForTest(tablename, region, endpoint *string) error {

	createTableInput := &dynamodb.CreateTableInput{
		TableName: tablename,
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			&dynamodb.AttributeDefinition{
				AttributeName: aws.String("ObjectType"),
				AttributeType: aws.String("S"),
			},
			&dynamodb.AttributeDefinition{
				AttributeName: aws.String("Id"),
				AttributeType: aws.String("S"),
			}},
		KeySchema: []*dynamodb.KeySchemaElement{
			&dynamodb.KeySchemaElement{
				AttributeName: aws.String("ObjectType"),
				KeyType:       aws.String("HASH"),
			},
			&dynamodb.KeySchemaElement{
				AttributeName: aws.String("Id"),
				KeyType:       aws.String("RANGE"),
			},
		},
	}

	_, err := dynamoDbClient(region, endpoint).CreateTable(createTableInput)
	return err
}

// TearDownTableForTest will delete passed table from DynamoDb.
// Attention: If you use wrong region and endpoint settings this may drop an undesired table!
func TearDownTableForTest(tablename, region, endpoint *string) error {

	deleteTableInput := &dynamodb.DeleteTableInput{
		TableName: tablename,
	}

	_, err := dynamoDbClient(region, endpoint).DeleteTable(deleteTableInput)
	return err
}

// listTables returns all available DynamoDb tables.
func listTables(region, endpoint *string) ([]*string, error) {

	res, err := dynamoDbClient(region, endpoint).ListTables(&dynamodb.ListTablesInput{})
	return res.TableNames, err
}

// dynamoDbClient creates a new dynamoDb client with passed config.
func dynamoDbClient(region, endpoint *string) *dynamodb.DynamoDB {

	awsConfig := &aws.Config{
		Region:   region,
		Endpoint: endpoint,
	}
	sess := session.Must(session.NewSession(awsConfig))
	return dynamodb.New(sess)
}
