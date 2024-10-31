package dynamodb

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var REMINDERS_TABLE_NAME = "reminders"

var BILLING_MODE = types.BillingModePayPerRequest

var DynamoDBTables = map[string]*dynamodb.CreateTableInput{
	REMINDERS_TABLE_NAME: DynamoDBRemindersTable,
}
