package dynamodb

import "fmt"

// identifierAsString returns a string representation of an identifier.
func identifierAsString(id ItemKey) string {
	return fmt.Sprintf("%s:%s", id.GetObjectType(), id.GetId())
}
