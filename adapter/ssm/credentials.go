package dynamodb

/** SSM Credentials **/
type Credentials interface {
	AccessKeyID() string
	SecretAccessKey() string
}
