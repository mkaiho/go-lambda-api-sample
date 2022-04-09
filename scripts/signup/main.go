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
				Usage:    "sign up user email",
				Aliases:  []string{"e"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "password",
				Usage:    "sign up user password",
				Aliases:  []string{"p"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "family-name",
				Usage:    "sign up user family name",
				Aliases:  []string{"f"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "given-name",
				Usage:    "sign up user given name",
				Aliases:  []string{"g"},
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
		familyName string
		givenName  string
	)
	{
		userPoolID = cliCtx.String("user-pool-id")
		clientID = cliCtx.String("client-id")
		email = cliCtx.String("email")
		password = cliCtx.String("password")
		familyName = cliCtx.String("family-name")
		givenName = cliCtx.String("given-name")
	}
	userpoolClient := infrastructure.NewCognitoUserPoolClient(region, userPoolID, clientID)

	output, err := userpoolClient.SignUp(ctx, userpool.SignUpInput{
		Email:      email,
		Password:   password,
		FamilyName: familyName,
		GivenName:  givenName,
	})
	if err != nil {
		logger.Error(err, "failed to sign up")
		return err
	}

	logger.Info(
		"created user",
		"sub", output.Sub,
		"email", email,
		"familyName", familyName,
		"givenName", givenName,
	)
	return nil
}
