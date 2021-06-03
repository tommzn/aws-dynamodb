package testing

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// SetupTableForTest will create a new DynamoDb table with passed name
// and a composed primary key with an ID as hash key and an object type as sort key.
func SetupTableForTest(tablename, region, endpoint *string) error {

	dynamoDb, err := dynamoDbClient(region, endpoint)
	if err != nil {
		return err
	}

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

	_, err = dynamoDb.CreateTable(createTableInput)
	if err != nil {
		return err
	}
	return nil
}

// TearDownTableForTest will delete passed table from DynamoDb.
// Attention: If you use wrong region and endpoint settings this may drop an undesired table!
func TearDownTableForTest(tablename, region, endpoint *string) error {

	dynamoDb, err := dynamoDbClient(region, endpoint)
	if err != nil {
		return err
	}

	deleteTableInput := &dynamodb.DeleteTableInput{
		TableName: tablename,
	}

	_, err = dynamoDb.DeleteTable(deleteTableInput)
	if err != nil {
		return err
	}
	return nil
}

// dynamoDbClient creates a new dynamoDb client with passed config.
func dynamoDbClient(region, endpoint *string) (*dynamodb.DynamoDB, error) {

	awsConfig := &aws.Config{
		Region:   region,
		Endpoint: endpoint,
	}
	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, err
	}
	d := dynamodb.New(sess)
	return d, nil
}

// listTables returns all available DynamoDb tables.
func listTables(region, endpoint *string) (*[]*string, error) {

	dynamoDb, err := dynamoDbClient(region, endpoint)
	if err != nil {
		return nil, err
	}

	res, err := dynamoDb.ListTables(&dynamodb.ListTablesInput{})
	if err == nil {
		return &res.TableNames, nil
	}
	return nil, err
}
