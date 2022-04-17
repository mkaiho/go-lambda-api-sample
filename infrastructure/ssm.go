package infrastructure

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	adapter "github.com/mkaiho/go-lambda-api-sample/adapter/ssm"
)

var _ adapter.Credentials = (*ssmCredentials)(nil)
var _ adapter.Client = (*ssmClient)(nil)

/** SSM Credentials **/
type ssmCredentials struct {
	accessKeyID     string
	secretAccessKey string
}

func (c *ssmCredentials) AccessKeyID() string {
	return c.accessKeyID
}

func (c *ssmCredentials) SecretAccessKey() string {
	return c.secretAccessKey
}

/** SSM Client **/
type ssmClient struct {
	ssm *ssm.SSM
}

func NewSSMClient(creds adapter.Credentials, region string) (adapter.Client, error) {
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

	return &ssmClient{
		ssm: ssm.New(sess, c),
	}, nil
}

func NewSSMClientFromEnv(region string) (adapter.Client, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	c := aws.NewConfig().
		WithRegion(region).
		WithCredentials(credentials.NewEnvCredentials())

	return &ssmClient{
		ssm: ssm.New(sess, c),
	}, nil
}

func (c *ssmClient) GetValue(key string) (string, error) {
	return c.getParameter(key, false)
}

func (c *ssmClient) GetSecret(key string) (string, error) {
	return c.getParameter(key, true)
}

func (c *ssmClient) getParameter(key string, withDecryption bool) (string, error) {
	resp, err := c.ssm.GetParameter(&ssm.GetParameterInput{
		Name:           &key,
		WithDecryption: &withDecryption,
	})
	if err != nil {
		return "nil", err
	}
	if resp.Parameter.Value == nil {
		return "", nil
	}

	return *resp.Parameter.Value, nil
}
