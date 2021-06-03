package testing

import (
	test "testing"

	"github.com/stretchr/testify/suite"
)

type DynamoDbTestSuite struct {
	suite.Suite
	tablename string
	region    string
	endpoint  string
}

func TestDynamoDbTestSuite(t *test.T) {
	suite.Run(t, new(DynamoDbTestSuite))
}

func (suite *DynamoDbTestSuite) SetupTest() {
	suite.tablename = "TestTable"
	suite.region = "eu-central-5"
	suite.endpoint = "http://localhost:8000"
}

func (suite *DynamoDbTestSuite) TestSetupAndTearDownTable() {

	suite.False(suite.tableExists())

	suite.Nil(SetupTableForTest(&suite.tablename, &suite.region, &suite.endpoint))
	suite.True(suite.tableExists())

	suite.Nil(TearDownTableForTest(&suite.tablename, &suite.region, &suite.endpoint))
	suite.False(suite.tableExists())
}

func (suite *DynamoDbTestSuite) tableExists() bool {

	tables, err := listTables(&suite.region, &suite.endpoint)
	if err != nil || len(tables) == 0 {
		return false
	}

	for _, existingTable := range tables {
		if *existingTable == suite.tablename {
			return true
		}
	}
	return false
}
