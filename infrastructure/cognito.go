package infrastructure

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/mkaiho/go-lambda-api-sample/adapter/userpool"
	"github.com/mkaiho/go-lambda-api-sample/entity"
	"github.com/mkaiho/go-lambda-api-sample/util"
)

var _ userpool.Client = (*cognitoUserPoolClient)(nil)

type cognitoUserPoolClient struct {
	userPoolID string
	clientID   string
	cognito    *cognitoidentityprovider.CognitoIdentityProvider
}

func NewCognitoUserPoolClient(region string, userPoolID string, clientID string) userpool.Client {
	conf := aws.NewConfig().
		WithRegion(region).
		WithCredentials(credentials.NewEnvCredentials())
	sess, err := session.NewSession(conf)
	if err != nil {
		panic(err)
	}
	return &cognitoUserPoolClient{
		userPoolID: userPoolID,
		clientID:   clientID,
		cognito:    cognitoidentityprovider.New(sess),
	}
}

func (c *cognitoUserPoolClient) SignUp(ctx context.Context, input userpool.SignUpInput) (*userpool.SignUpOutput, error) {
	resp, err := c.cognito.SignUpWithContext(ctx, &cognitoidentityprovider.SignUpInput{
		ClientId: aws.String(c.clientID),
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			{
				Name:  aws.String(userpool.AttributeNameEmail.String()),
				Value: aws.String(input.Email),
			},
			{
				Name:  aws.String(userpool.AttributeNameGivenName.String()),
				Value: aws.String(input.GivenName),
			},
			{
				Name:  aws.String(userpool.AttributeNameFamilyName.String()),
				Value: aws.String(input.FamilyName),
			},
		},
		Username: aws.String(input.Email),
		Password: aws.String(input.Password),
	})
	if err != nil {
		return nil, err
	}

	out := userpool.SignUpOutput{
		Sub: *resp.UserSub,
	}
	return &out, nil
}

func (c *cognitoUserPoolClient) ConfirmSignUp(ctx context.Context, input userpool.ConfirmSignUpInput) error {
	_, err := c.cognito.ConfirmSignUpWithContext(ctx, &cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         aws.String(c.clientID),
		Username:         &input.Email,
		ConfirmationCode: &input.ConfirmationCode,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *cognitoUserPoolClient) InitiateAuth(ctx context.Context, input userpool.InitiateAuthInput) (*userpool.InitiateAuthOutput, error) {
	resp, err := c.cognito.InitiateAuthWithContext(ctx, &cognitoidentityprovider.InitiateAuthInput{
		ClientId: &c.clientID,
		AuthFlow: aws.String(cognitoidentityprovider.AuthFlowTypeUserPasswordAuth),
		AuthParameters: map[string]*string{
			"USERNAME": &input.Email,
			"PASSWORD": &input.Password,
		},
	})
	if err != nil {
		return nil, err
	}

	out := userpool.InitiateAuthOutput{
		IDToken:      resp.AuthenticationResult.IdToken,
		AccessToken:  resp.AuthenticationResult.AccessToken,
		RefreshToken: resp.AuthenticationResult.RefreshToken,
	}
	return &out, nil
}

func (c *cognitoUserPoolClient) ResendConfirmationCode(ctx context.Context, input userpool.ResendConfirmationCodeInput) error {
	_, err := c.cognito.ResendConfirmationCodeWithContext(ctx, &cognitoidentityprovider.ResendConfirmationCodeInput{
		ClientId: aws.String(c.clientID),
		Username: &input.Email,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *cognitoUserPoolClient) AdminCreateUser(ctx context.Context, input userpool.AdminCreateUserInput) (*userpool.AdminCreateUserOutput, error) {
	resp, err := c.cognito.AdminCreateUser(&cognitoidentityprovider.AdminCreateUserInput{
		UserPoolId: &c.userPoolID,
		Username:   aws.String(input.UserName),
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			{
				Name:  aws.String(userpool.AttributeNameEmail.String()),
				Value: aws.String(input.Email),
			},
		},
		DesiredDeliveryMediums: []*string{
			aws.String("EMAIL"),
		},
	})
	if err != nil {
		return nil, err
	}

	out := userpool.AdminCreateUserOutput{
		UserName: *resp.User.Username,
	}
	log.Println(resp.User.Attributes)
	for _, attr := range resp.User.Attributes {
		switch *attr.Name {
		case userpool.AttributeNameEmail.String():
			out.Email = *attr.Value
		case userpool.AttributeNameSub.String():
			out.Sub = *attr.Value
		}
	}

	return &out, nil
}

func (c *cognitoUserPoolClient) AdminUpdateUser(ctx context.Context, input userpool.AdminUpdateUserInput) error {
	_, err := c.cognito.AdminUpdateUserAttributesWithContext(ctx, &cognitoidentityprovider.AdminUpdateUserAttributesInput{
		UserPoolId:     &c.userPoolID,
		Username:       input.UserAttribute.Email.Value,
		UserAttributes: input.UserAttribute.ToRows(),
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *cognitoUserPoolClient) AdminGetUser(ctx context.Context, input userpool.AdminGetUserInput) (*userpool.AdminGetUserOutput, error) {
	resp, err := c.cognito.AdminGetUserWithContext(ctx, &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: &c.userPoolID,
		Username:   &input.Email,
	})
	if err != nil {
		var une error = new(cognitoidentityprovider.UserNotFoundException)
		if errors.As(err, &une) {
			return nil, entity.ErrResourceNotFound
		}
		return nil, err
	}

	out := userpool.AdminGetUserOutput{}
	for _, attr := range resp.UserAttributes {
		if attr.Name == nil {
			util.GetLogger().Info("skipped because the attribute name is nil")
			continue
		}
		if attr.Value == nil {
			continue
		}
		switch *attr.Name {
		case userpool.AttributeNameSub.String():
			out.Sub = *attr.Value
		case userpool.AttributeNameEmailVerified.String():
			out.EmailVerified = *attr.Value
		case userpool.AttributeNameGivenName.String():
			out.GivenName = *attr.Value
		case userpool.AttributeNameFamilyName.String():
			out.FamilyName = *attr.Value
		case userpool.AttributeNameEmail.String():
			out.Email = *attr.Value
		}
	}

	return &out, nil
}
