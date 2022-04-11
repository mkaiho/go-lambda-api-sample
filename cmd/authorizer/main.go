package main

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mkaiho/go-lambda-api-sample/entity"
	"github.com/mkaiho/go-lambda-api-sample/infrastructure"
	"github.com/mkaiho/go-lambda-api-sample/util"
)

var (
	issure          = "https://cognito-idp.ap-northeast-1.amazonaws.com/ap-northeast-1_pV4ISzjlY"
	audience        = []string{"1fb3tjglh346btqasbliodqdj2"}
	validitySeconds = 360
)

var (
	initErr      error
	tokenManager infrastructure.JWTManager
)

func init() {
	// logger
	var (
		loggerConf util.LoggerConfig
	)
	{
		loggerConf, initErr = util.NewLoggerConfig(util.InfoLevel)
		if initErr != nil {
			return
		}
		infrastructure.InjectZapLogger(loggerConf)
	}

	// ID
	var (
		idgen entity.IDGenerator
	)
	{
		idgen = infrastructure.NewULIDGenerator()
	}

	// token
	var (
		tokenManagerConf infrastructure.TokenManagerConfig
	)
	{
		tokenManagerConf = infrastructure.NewTokenManagerConfig(
			issure,
			audience,
			time.Second*time.Duration(validitySeconds),
		)
		tokenManager = infrastructure.NewJWTManager(tokenManagerConf, idgen)
	}
}
func main() {
	lambda.Start(handle)
}

func handle(
	ctx context.Context,
	event events.APIGatewayCustomAuthorizerRequest,
) (events.APIGatewayCustomAuthorizerResponse, error) {
	var token = strings.Replace(event.AuthorizationToken, "Bearer ", "", 1)
	var resp = events.APIGatewayCustomAuthorizerResponse{}

	if initErr != nil {
		return resp, nil
	}

	publicKey, err := tokenManager.ReadJWKSetFile("jwks.json")
	if err != nil {
		return resp, err
	}

	claims, err := tokenManager.VerifyWithKeySet([]byte(token), publicKey)
	if err != nil {
		return resp, nil
	}

	resp.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
		Version: "2012-10-17",
		Statement: []events.IAMPolicyStatement{
			{
				Action:   []string{"execute-api:Invoke"},
				Effect:   "Allow",
				Resource: []string{event.MethodArn},
			},
		},
	}
	resp.Context = map[string]interface{}{
		"Claims": string(claims),
	}

	return resp, nil
}
