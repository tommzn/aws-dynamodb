package dynamodb

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type IdentifierTestSuite struct {
	suite.Suite
}

func TestIdentifierTestSuite(t *testing.T) {
	suite.Run(t, new(IdentifierTestSuite))
}

func (suite *IdentifierTestSuite) TestCreateItemKey() {

	id := "id-1"
	objectType := "TestItem"
	itemKey := NewItemIdentifier(id, objectType)
	suite.Equal(id, itemKey.GetId())
	suite.Equal(objectType, itemKey.GetObjectType())
}
