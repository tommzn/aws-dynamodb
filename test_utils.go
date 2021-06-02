package dynamodb

import (
	"time"

	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	utils "github.com/tommzn/go-utils"
)

// dynamoDbSettings extracts DynamoDb settings from passed config.
func dynamoDbSettings(conf config.Config) (*string, *string, *string) {

	tablename := conf.Get("aws.dynamodb.tablename", nil)
	region := conf.Get("aws.dynamodb.region", nil)
	endpoint := conf.Get("aws.dynamodb.endpoint", nil)
	return tablename, region, endpoint
}

// loadConfigForTest returns test config from file testconfig.yml.
func loadConfigForTest() config.Config {

	configFile := "testconfig.yml"
	source := config.NewFileConfigSource(&configFile)
	conf, _ := source.Load()
	return conf
}

// testItem is a DynamoDb item used for package tests.
type testItem struct {
	*ItemIdentifier
	Val1 string
	Val2 int
	Val3 time.Time
}

// newItemForTest returns a test items with dummy values.
func newItemForTest() *testItem {
	return &testItem{
		ItemIdentifier: NewItemIdentifier(utils.NewId(), "TestItems"),
		Val1:           "xXx",
		Val2:           123445667,
		Val3:           time.Now(),
	}
}

// newTestItemWithoutValues copies item key from passed item into a new one.
func newTestItemWithoutValues(item ItemKey) *testItem {
	return &testItem{
		ItemIdentifier: NewItemIdentifier(item.GetId(), item.GetObjectType()),
	}
}

// loggerForTest returns a new stdout logger.
func loggerForTest(logLevel log.LogLevel) log.Logger {
	return log.NewLogger(logLevel, nil, nil)
}
