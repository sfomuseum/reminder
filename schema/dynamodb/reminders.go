package dynamodb

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var DynamoDBRemindersTable = &dynamodb.CreateTableInput{
	KeySchema: []types.KeySchemaElement{
		{
			AttributeName: aws.String("Id"), // partition key
			KeyType:       "HASH",
		},
	},
	AttributeDefinitions: []types.AttributeDefinition{
		{
			AttributeName: aws.String("Id"),
			AttributeType: "N",
		},
	},
	/*
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String("name"),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("Name"),
						KeyType:       "HASH",
					},
				},
				Projection: &types.Projection{
					ProjectionType: "ALL",
				},
			},
			{
				IndexName: aws.String("by_created"),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("Created"),
						KeyType:       "HASH",
					},
				},
				Projection: &types.Projection{
					ProjectionType: "KEYS_ONLY",
				},
			},
		},
	*/
	BillingMode: BILLING_MODE,
	TableName:   &REMINDERS_TABLE_NAME,
}
