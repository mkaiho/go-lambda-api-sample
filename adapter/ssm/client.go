package dynamodb

type Client interface {
	GetValue(key string) (string, error)
	GetSecret(key string) (string, error)
}
