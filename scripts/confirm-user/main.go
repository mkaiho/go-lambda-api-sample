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
				Usage:    "confirmed user email",
				Aliases:  []string{"e"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "code",
				Usage:    "confirmation code",
				Aliases:  []string{"c"},
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
		code       string
	)
	{
		userPoolID = cliCtx.String("user-pool-id")
		clientID = cliCtx.String("client-id")
		email = cliCtx.String("email")
		code = cliCtx.String("code")
	}
	userpoolClient := infrastructure.NewCognitoUserPoolClient(region, userPoolID, clientID)

	err := userpoolClient.ConfirmSignUp(ctx, userpool.ConfirmSignUpInput{
		Email:            email,
		ConfirmationCode: code,
	})
	if err != nil {
		logger.Error(err, "failed to confirm user")
		return err
	}

	logger.Info(
		"confirmed user",
		"email", email,
	)
	return nil
}
