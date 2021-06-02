package dynamodb

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const DEFAULT_AWS_REGION = "eu-central-1"
const DEFAULT_TABLENAME = "<DynamoDbTableNotSet>"

// Add will create a new item or update an existing item in DynamoDb.
func (r *DynamoDbRepository) Add(item ItemKey) error {

	r.logger.Debug("Add Item: ", identifierAsString(item))

	dynamoDb, err := r.newDynamoDb()
	if err != nil {
		return err
	}

	av, err := dynamodbattribute.MarshalMap(item)
	r.logger.Debugf("AttributeVAlue: %+v", av)
	if err != nil {
		r.logger.Errorf("MarshalMap Error: %s", err)
		return err
	}

	input := &dynamodb.PutItemInput{Item: av, TableName: r.TableName}
	_, err = dynamoDb.PutItem(input)
	if err != nil {
		r.logger.Errorf("PutItem Error: %s", err)
	}
	return err
}

// Get will try to read an item from DynamDb by passed item key.
// Passed item have to be a pointer, because it will unmarshal DynamiDb item values into it.
func (r *DynamoDbRepository) Get(item ItemKey) error {

	r.logger.Debug("Get item: ", identifierAsString(item))

	if reflect.ValueOf(item).Kind() != reflect.Ptr {
		msg := fmt.Sprintf("Expect pointer receiver for: %s\n", identifierAsString(item))
		r.logger.Error(msg)
		return errors.New(msg)
	}

	dynamoDb, err := r.newDynamoDb()
	if err != nil {
		return err
	}

	input, err := r.newGetItemInput(item)
	if err != nil {
		r.logger.Error(err)
		return err
	}

	result, err := dynamoDb.GetItem(input)
	if err != nil {
		r.logger.Error(err)
		return err
	}
	r.logger.Debugf("DynamoDb Response for %s is: %+v", identifierAsString(item), result.Item)

	if len(result.Item) == 0 {
		msg := fmt.Sprintf("Not found: %s\n", identifierAsString(item))
		r.logger.Info(msg)
		return errors.New(msg)
	}
	return dynamodbattribute.UnmarshalMap(result.Item, item)

}

// Delete will try to delete an item from DynamoDb item identified by passed item key.
func (r *DynamoDbRepository) Delete(item ItemKey) error {

	r.logger.Debug("Delete Item: ", identifierAsString(item))

	dynamoDb, err := r.newDynamoDb()
	if err != nil {
		return err
	}

	input, err := r.newDeleteItemInput(item)
	if err != nil {
		r.logger.Error(err)
		return err
	}

	_, err = dynamoDb.DeleteItem(input)
	if err != nil {
		r.logger.Error(err)
	}
	return err
}

// NewDynamoDb creates a DynamoDb client. Uses a singleton pattern which creates the client only once.
func (r *DynamoDbRepository) newDynamoDb() (*dynamodb.DynamoDB, error) {

	if r.dynamoDbClient == nil {
		sess, err := session.NewSession(r.Config)
		if err != nil {
			r.logger.Error("Unable to create AWS session: %s", err)
			return nil, err
		}
		r.dynamoDbClient = dynamodb.New(sess)
	}
	return r.dynamoDbClient, nil
}

// newGetItemInput creates a new DynamoDb GetItemInout for passed item.
func (r *DynamoDbRepository) newGetItemInput(item ItemKey) (*dynamodb.GetItemInput, error) {

	dynamodbKey, err := dynamodbattribute.MarshalMap(NewItemIdentifier(item.GetId(), item.GetObjectType()))
	if err != nil {
		return nil, err
	}
	return &dynamodb.GetItemInput{
		Key:       dynamodbKey,
		TableName: r.TableName,
	}, nil
}

// newDeleteItemInput creates a new DynamoDb DeleteItemInout for passed item.
func (r *DynamoDbRepository) newDeleteItemInput(item ItemKey) (*dynamodb.DeleteItemInput, error) {

	dynamodbKey, err := dynamodbattribute.MarshalMap(NewItemIdentifier(item.GetId(), item.GetObjectType()))
	if err != nil {
		return nil, err
	}
	return &dynamodb.DeleteItemInput{
		Key:       dynamodbKey,
		TableName: r.TableName,
	}, nil
}
