package dynamodb

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

/** DynamoDBClient **/
type DynamoDBClient interface {
	GetItem(input dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
	PutItem(input dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
	DeleteItem(input dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error)
	Query(input dynamodb.QueryInput) (*dynamodb.QueryOutput, error)
	Scan(input dynamodb.ScanInput) (*dynamodb.ScanOutput, error)
}
