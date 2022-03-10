package dynamodb

/** DynamoDB Credentials **/
type DynamoDBCredentials interface {
	AccessKeyID() string
	SecretAccessKey() string
}
