package dynamodb

import "fmt"

// identifierAsString returns a string represantation of an identifier.
func identifierAsString(id ItemKey) string {
	return fmt.Sprintf("%s:%s", id.GetId(), id.GetObjectType())
}
