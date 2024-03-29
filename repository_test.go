package dynamodb

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	testutils "github.com/tommzn/aws-dynamodb/testing"
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	utils "github.com/tommzn/go-utils"
)

type RepositoryTestSuite struct {
	suite.Suite
	logLevel log.LogLevel
	conf     config.Config
	repo     Repository
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

func (suite *RepositoryTestSuite) SetupTest() {
	suite.logLevel = log.Error
	suite.conf = loadConfigForTest()
	suite.repo = NewRepository(suite.conf, loggerForTest(suite.logLevel))
	tablename, region, endpoint := dynamoDbSettings(suite.conf)
	suite.Nil(testutils.SetupTableForTest(tablename, region, endpoint))
}

func (suite *RepositoryTestSuite) TearDownTest() {
	tablename, region, endpoint := dynamoDbSettings(suite.conf)
	suite.Nil(testutils.TearDownTableForTest(tablename, region, endpoint))
}

func (suite *RepositoryTestSuite) TestCrudActions() {

	item := newItemForTest()
	suite.Nil(suite.repo.Add(item))

	item2 := newTestItemWithoutValues(item)
	suite.Nil(suite.repo.Get(item2))
	suite.Equal(item.Val1, item2.Val1)
	suite.Equal(item.Val2, item2.Val2)
	suite.Equal(item.Val3.Format(time.RFC3339), item2.Val3.Format(time.RFC3339))

	item.Val1 = "yYy"
	item.Val2 = 9887654
	item.Val3 = time.Now().Add(-5 * time.Minute)
	suite.Nil(suite.repo.Add(item))

	item3 := newTestItemWithoutValues(item)
	suite.Nil(suite.repo.Get(item3))
	suite.Equal(item.Val1, item3.Val1)
	suite.Equal(item.Val2, item3.Val2)
	suite.Equal(item.Val3.Format(time.RFC3339), item3.Val3.Format(time.RFC3339))

	suite.Nil(suite.repo.Delete(item))
	suite.NotNil(suite.repo.Get(item))

	suite.NotNil(suite.repo.Get(*item))

	item.ItemIdentifier.Id = "xxx"
	err := suite.repo.Get(item)
	suite.NotNil(err)
}

func (suite *RepositoryTestSuite) TestObjectLock() {

	item := newItemForTest()
	itemLock, err := suite.repo.Lock(item)
	suite.Nil(err)
	suite.NotNil(itemLock)

	suite.NotNil(suite.repo.Add(itemLock))

	expiresAt := itemLock.ExpiresAt
	time.Sleep(1 * time.Second)
	_, err1 := suite.repo.Renew(itemLock)
	suite.Nil(err1)
	suite.True(itemLock.ExpiresAt > expiresAt)

	itemLock.LockId = utils.NewId()
	_, err1_1 := suite.repo.Renew(itemLock)
	suite.NotNil(err1_1)

	itemLock2, err2 := suite.repo.Lock(item)
	suite.NotNil(err2)
	suite.Nil(itemLock2)

	suite.Nil(suite.repo.Unlock(itemLock))

	_, err3 := suite.repo.Renew(itemLock)
	suite.NotNil(err3)

	itemLock4, err4 := suite.repo.Lock(item)
	suite.Nil(err4)
	suite.NotNil(itemLock4)

}

func (suite *RepositoryTestSuite) TestLockAfterExpiration() {

	suite.repo.(*DynamoDbRepository).lockTtl = 1 * time.Second

	item := newItemForTest()
	itemLock, err := suite.repo.Lock(item)
	suite.Nil(err)
	suite.NotNil(itemLock)

	_, err1 := suite.repo.Lock(item)
	suite.NotNil(err1)

	time.Sleep(2 * time.Second)
	itemLock2, err2 := suite.repo.Lock(item)
	suite.Nil(err2)
	suite.NotNil(itemLock2)
}

func (suite *RepositoryTestSuite) TestWithErrors() {

	item := newItemForTest()
	suite.repo.(*DynamoDbRepository).tableName = nil
	suite.NotNil(suite.repo.Get(item))
}

func (suite *RepositoryTestSuite) TestQueryItems() {

	for i := 1; i <= 3; i++ {
		suite.Nil(suite.repo.Add(newItemForTest()))
	}

	items := []testItem{}
	suite.Nil(suite.repo.Query("TestItems", &items))
	suite.Len(items, 3)

	items2 := []testItem{}
	suite.Nil(suite.repo.Query("XXX", &items2))
	suite.Len(items2, 0)

	suite.NotNil(suite.repo.Query("XXX", []testItem{}))
}
