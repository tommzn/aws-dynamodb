package dynamodb

// NewItemIdentifier returns a new ItemIdentifier with passed id and object type.
func NewItemIdentifier(id, objectType string) ItemKey {
	return &ItemIdentifier{Id: id, ObjectType: objectType}
}

// GetId returns the id of a DynamoDb item.
func (id *ItemIdentifier) GetId() string {
	return id.Id
}

// GetObjectType returns the object type of a DynamoDb item.
func (id *ItemIdentifier) GetObjectType() string {
	return id.ObjectType
}
