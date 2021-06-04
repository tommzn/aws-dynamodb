package dynamodb

// ItemKey is an interface each object have to fulfill to be persisted
// in DynamoDb.
type ItemKey interface {
	GetId() string
	GetObjectType() string
}

// Repository provides CRUD access for DynamoDb items.
type Repository interface {

	// Add or update an item in DynamoDb.
	Add(ItemKey) error

	// Get will try to read an item by specified key from DynamoDb.
	Get(ItemKey) error

	// Query will list all items for an object type.
	Query(string, interface{}) error

	// Delete will remove an item with specified key from DynamoDb.
	Delete(ItemKey) error

	// Lock will try to obtain a lock for an items identified by passed key.
	Lock(ItemKey) (*ItemLock, error)

	// Renew can be used to extend lease of an item lock.
	Renew(*ItemLock) (*ItemLock, error)

	// Unlock will delete passed object lock from DynamoDb.
	Unlock(*ItemLock) error
}
