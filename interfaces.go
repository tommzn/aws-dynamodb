package dynamodb

// ItemKey is an interface each object have to fulfill to be persisted
// in DynamoDb.
type ItemKey interface {
	GetId() string
	GetObjectType() string
}

// Repository provides CRUD access for DynamoDb items.
type Repository interface {
	Add(ItemKey) error
	Get(ItemKey) error
	Query(string, interface{}) error
	//Scan(ListItemRequest) (ListItemResponse, error)
	Delete(ItemKey) error
}
