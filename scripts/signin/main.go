package main

import (
	"context"
	"log"
	"os"

	"github.com/mkaiho/go-lambda-api-sample/adapter/userpool"
	"github.com/mkaiho/go-lambda-api-sample/infrastructure"
	"github.com/mkaiho/go-lambda-api-sample/util"
	"github.com/urfave/cli/v2"
)

const region = "ap-northeast-1"

var (
	initErr error
	ctx     context.Context
	logger  util.Logger
)

func init() {
	ctx = context.Background()
	var loggerConf util.LoggerConfig
	loggerConf, initErr = util.NewLoggerConfig(util.InfoLevel)
	infrastructure.InjectZapLogger(loggerConf)
	logger = util.GetLogger()
}

func main() {
	if initErr != nil {
		log.Fatal(initErr)
	}

	app := cli.App{
		Name:  "Sign up",
		Usage: "Sign up app",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "user-pool-id",
				Usage:    "User pool ID",
				EnvVars:  []string{"USER_POOL_ID"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "client-id",
				Usage:    "Client ID of User pool",
				EnvVars:  []string{"USER_POOL_CLIENT_ID"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "email",
				Usage:    "sign in user email",
				Aliases:  []string{"e"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "password",
				Usage:    "sign in user password",
				Aliases:  []string{"p"},
				Required: true,
			},
		},
		Action: runner,
	}
	app.Run(os.Args)
}

var runner cli.ActionFunc = func(cliCtx *cli.Context) error {
	var (
		userPoolID string
		clientID   string
		email      string
		password   string
	)
	{
		userPoolID = cliCtx.String("user-pool-id")
		clientID = cliCtx.String("client-id")
		email = cliCtx.String("email")
		password = cliCtx.String("password")
	}
	userpoolClient := infrastructure.NewCognitoUserPoolClient(region, userPoolID, clientID)

	output, err := userpoolClient.InitiateAuth(ctx, userpool.InitiateAuthInput{
		Email:    email,
		Password: password,
	})
	if err != nil {
		logger.Error(err, "failed to sign up")
		return err
	}

	logger.Info(
		"created user",
		"email", email,
		"idToken", output.IDToken,
		"accessToken", output.AccessToken,
		"refreshToken", output.RefreshToken,
	)
	return nil
}
