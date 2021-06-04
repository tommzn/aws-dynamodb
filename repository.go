package dynamodb

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	utils "github.com/tommzn/go-utils"
)

// DEFAULT_AWS_REGION defines the fallback AWS region if nothing has been specified by config.
const DEFAULT_AWS_REGION = "eu-central-1"

// DEFAULT_TABLENAME defines the default DynamoDb table name.
const DEFAULT_TABLENAME = "<DynamoDbTableNotSet>"

// lockObjectType is the object type used for locks.
const lockObjectType = "OBJECTLOCK"

// Add will create a new item or update an existing item in DynamoDb.
func (r *DynamoDbRepository) Add(item ItemKey) error {

	r.logger.Debug("Add Item: ", identifierAsString(item))

	if item.GetObjectType() == lockObjectType {
		return fmt.Errorf("Unsupported object type for Add: %s", item.GetObjectType())
	}

	av, err := dynamodbattribute.MarshalMap(item)
	r.logger.Debugf("AttributeValue: %+v", av)
	if err == nil {
		input := &dynamodb.PutItemInput{Item: av, TableName: r.tableName}
		_, err = r.dynamoDb().PutItem(input)
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

	result, err := r.dynamoDb().GetItem(r.newGetItemInput(item))
	if err == nil {

		r.logger.Debugf("DynamoDb Response for %s is: %+v", identifierAsString(item), result.Item)
		if len(result.Item) == 0 {
			msg := fmt.Sprintf("Not found: %s\n", identifierAsString(item))
			r.logger.Info(msg)
			return errors.New(msg)
		}
		return dynamodbattribute.UnmarshalMap(result.Item, item)

	}
	return err
}

// Delete will try to delete an item from DynamoDb item identified by passed item key.
func (r *DynamoDbRepository) Delete(item ItemKey) error {

	r.logger.Debug("Delete Item: ", identifierAsString(item))

	_, err := r.dynamoDb().DeleteItem(r.newDeleteItemInput(item))
	return err
}

// Query will list items for a specific object type. Receiver have to be a slice of expected type
// and will be used to unmarshal query result. So please pass it as a pointer.
// If there're no items for passed object type no error is returned, passed slice will stay empty.
func (r *DynamoDbRepository) Query(objectType string, receiver interface{}) error {

	r.logger.Debug("Query Items: ", objectType)

	if reflect.ValueOf(receiver).Kind() != reflect.Ptr {
		msg := "Expect pointer receiver for items."
		r.logger.Error(msg)
		return errors.New(msg)
	}

	result, err := r.dynamoDb().Query(r.newQueryInput(objectType))
	r.logger.Debugf("Query Result: %+v", result)
	if err == nil {
		err = dynamodbattribute.UnmarshalListOfMaps(result.Items, receiver)
		r.logger.Debugf("List Response: %+v", receiver)
	}
	return err
}

// Lock will try to obtain a lock passed item. Default life time of a lock is 5 min.
func (r *DynamoDbRepository) Lock(item ItemKey) (*ItemLock, error) {

	itemLock := r.newObjectLockForItem(item)
	input := r.newPutItemInputForLock(&itemLock)
	if _, err := r.dynamoDb().PutItem(input); err == nil {
		return &itemLock, nil
	} else {
		return nil, err
	}
}

// Renew can be used to extend life time of a lock.
func (r *DynamoDbRepository) Renew(itemLock *ItemLock) (*ItemLock, error) {

	itemLock.ExpiresAt = r.newLockExpiration()
	input := r.newPutItemInputForRenew(itemLock)
	if _, err := r.dynamoDb().PutItem(input); err == nil {
		return itemLock, nil
	} else {
		return nil, err
	}
}

// Unlock will remove given lock from DynamoDb.
func (r *DynamoDbRepository) Unlock(itemLock *ItemLock) error {

	return r.Delete(itemLock)
}

// dynamoDb creates a DynamoDb client. Uses a singleton pattern which creates the client only once.
func (r *DynamoDbRepository) dynamoDb() *dynamodb.DynamoDB {

	if r.dynamoDbClient == nil {
		sess := session.Must(session.NewSession(r.config))
		r.dynamoDbClient = dynamodb.New(sess)
	}
	return r.dynamoDbClient
}

// newGetItemInput creates a new DynamoDb GetItemInout for passed item.
func (r *DynamoDbRepository) newGetItemInput(item ItemKey) *dynamodb.GetItemInput {

	dynamodbKey, _ := dynamodbattribute.MarshalMap(NewItemIdentifier(item.GetId(), item.GetObjectType()))
	return &dynamodb.GetItemInput{
		Key:       dynamodbKey,
		TableName: r.tableName,
	}
}

// newDeleteItemInput creates a new DynamoDb DeleteItemInout for passed item.
func (r *DynamoDbRepository) newDeleteItemInput(item ItemKey) *dynamodb.DeleteItemInput {

	dynamodbKey, _ := dynamodbattribute.MarshalMap(NewItemIdentifier(item.GetId(), item.GetObjectType()))
	return &dynamodb.DeleteItemInput{
		Key:       dynamodbKey,
		TableName: r.tableName,
	}
}

// newQueryInput creates a new query input for AWS DynamoDb.
func (r *DynamoDbRepository) newQueryInput(objectType string) *dynamodb.QueryInput {

	keyCondition := expression.Key("ObjectType").Equal(expression.Value(objectType))
	expr, _ := expression.NewBuilder().WithKeyCondition(keyCondition).Build()
	r.logger.Debugf("Key expression: %+v", expr)

	expressionAttributeNames := expr.Names()
	expressionAttributeValues := expr.Values()
	r.logger.Debugf("Expr names: %+v", expressionAttributeNames)
	r.logger.Debugf("Expr values: %+v", expressionAttributeValues)

	return &dynamodb.QueryInput{
		ExpressionAttributeNames:  expressionAttributeNames,
		ExpressionAttributeValues: expressionAttributeValues,
		TableName:                 r.tableName,
		KeyConditionExpression:    expr.KeyCondition(),
	}
}

// newPutItemInputForLock creates a new conditional put item input for a lock item.
func (r *DynamoDbRepository) newPutItemInputForLock(itemLock *ItemLock) *dynamodb.PutItemInput {

	dynamodbLockData, _ := dynamodbattribute.MarshalMap(itemLock)

	expressionAttributeValues := make(map[string]*dynamodb.AttributeValue)
	ttlAttribute, _ := dynamodbattribute.Marshal(itemLock.ExpiresAt)
	expressionAttributeValues[":ExpiresAt"] = ttlAttribute
	return &dynamodb.PutItemInput{
		Item:                      dynamodbLockData,
		TableName:                 r.tableName,
		ConditionExpression:       aws.String("(attribute_not_exists(Id) AND attribute_not_exists(ObjectType)) or ExpiresAt < :ExpiresAt"),
		ExpressionAttributeValues: expressionAttributeValues,
	}
}

// newPutItemInputForLock creates a new conditional put item input for a lock item.
func (r *DynamoDbRepository) newPutItemInputForRenew(itemLock *ItemLock) *dynamodb.PutItemInput {

	dynamodbLockData, _ := dynamodbattribute.MarshalMap(itemLock)

	expressionAttributeValues := make(map[string]*dynamodb.AttributeValue)
	attrLockId, _ := dynamodbattribute.Marshal(itemLock.LockId)
	expressionAttributeValues[":LockId"] = attrLockId
	return &dynamodb.PutItemInput{
		Item:                      dynamodbLockData,
		TableName:                 r.tableName,
		ConditionExpression:       aws.String("(attribute_exists(Id) AND attribute_exists(ObjectType))  AND LockId = :LockId"),
		ExpressionAttributeValues: expressionAttributeValues,
	}
}

// newObjectLockForItem returns a lock object.
func (r *DynamoDbRepository) newObjectLockForItem(item ItemKey) ItemLock {
	return ItemLock{
		ItemIdentifier: NewItemIdentifier(identifierAsString(item), lockObjectType),
		ExpiresAt:      r.newLockExpiration(),
		LockId:         utils.NewId(),
	}
}

// newLockExpiration returns the new expiration time of a lock.
func (r *DynamoDbRepository) newLockExpiration() int64 {
	return time.Now().Add(r.lockTtl).Unix()
}
