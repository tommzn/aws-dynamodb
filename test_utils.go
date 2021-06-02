package dynamodb

import config "github.com/tommzn/go-config"

// dynamoDbSettings extracts DynamoDb settings from passed config.
func dynamoDbSettings(conf config.Config) (*string, *string, *string) {

	tablename := conf.Get("aws.dynamodb.tablename", nil)
	region := conf.Get("aws.dynamodb.awsregion", nil)
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
