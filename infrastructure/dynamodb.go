package infrastructure

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	adapter "github.com/mkaiho/go-lambda-api-sample/adapter/dynamodb"
)

var _ adapter.DynamoDBCredentials = (*dynamoDBCredentials)(nil)
var _ adapter.DynamoDBClient = (*dynamoDBClient)(nil)
var _ adapter.AttributeValueMapper = (*dynamoDBAttributeValueMapper)(nil)

/** DynamoDB Credentials **/
type dynamoDBCredentials struct {
	accessKeyID     string
	secretAccessKey string
}

func (c *dynamoDBCredentials) AccessKeyID() string {
	return c.accessKeyID
}

func (c *dynamoDBCredentials) SecretAccessKey() string {
	return c.secretAccessKey
}

func NewDynamoDBCredentials(accessKeyID string, secretAccessKey string) *dynamoDBCredentials {
	return &dynamoDBCredentials{
		accessKeyID:     accessKeyID,
		secretAccessKey: secretAccessKey,
	}
}

/** DynamoDB Client **/
type dynamoDBClient struct {
	ds *dynamodb.DynamoDB
}

func NewDynamoDBClient(creds adapter.DynamoDBCredentials, region string) (*dynamoDBClient, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	c := aws.NewConfig().
		WithRegion(region).
		WithCredentials(credentials.NewCredentials(&credentials.StaticProvider{
			Value: credentials.Value{
				AccessKeyID:     creds.AccessKeyID(),
				SecretAccessKey: creds.SecretAccessKey(),
			},
		}))
	ds := dynamodb.New(sess, c)
	return &dynamoDBClient{
		ds: ds,
	}, nil
}

func NewDynamoDBClientFromEnv(region string) (*dynamoDBClient, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	c := aws.NewConfig().
		WithRegion(region).
		WithCredentials(credentials.NewEnvCredentials())
	ds := dynamodb.New(sess, c)
	return &dynamoDBClient{
		ds: ds,
	}, nil
}

func (c *dynamoDBClient) GetItem(input dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	return c.ds.GetItem(&input)
}

func (c *dynamoDBClient) PutItem(input dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return c.ds.PutItem(&input)
}

func (c *dynamoDBClient) DeleteItem(input dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	return c.ds.DeleteItem(&input)
}

func (c *dynamoDBClient) Query(input dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	return c.ds.Query(&input)
}

// TODO: Consideration about exceeding 1MB
func (c *dynamoDBClient) Scan(input dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
	return c.ds.Scan(&input)
}

/** AttributeValueMapper **/
type dynamoDBAttributeValueMapper struct{}

func NewDynamoDBAttributeValueMapper() adapter.AttributeValueMapper {
	return &dynamoDBAttributeValueMapper{}
}

func (am *dynamoDBAttributeValueMapper) MarshalMap(in interface{}) (map[string]*dynamodb.AttributeValue, error) {
	return dynamodbattribute.MarshalMap(in)
}

func (am *dynamoDBAttributeValueMapper) UnmarshalMap(m map[string]*dynamodb.AttributeValue, out interface{}) error {
	return dynamodbattribute.UnmarshalMap(m, out)
}
