package dynamodb

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type RepositoryTestSuite struct {
	suite.Suite
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

/**
func (suite *RepositoryTestSuite) SetupTest() {
	suite.conf = loadConfigForTest()
	tablename, region, endpoint := dynamoDbSettings(suite.conf)
	suite.Nil(testing.SetupTableForTest(tablename, region, endpoint))
}

func (suite *RepositoryTestSuite) TearDownTest() {
	tablename, region, endpoint := dynamoDbSettings(suite.conf)
	suite.Nil(testing.TearDownTableForTest(tablename, region, endpoint))
}
*/

func (suite *RepositoryTestSuite) TestAddItem() {

}
