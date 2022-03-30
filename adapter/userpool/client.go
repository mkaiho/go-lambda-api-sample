package userpool

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

type AttributeName string

func (n AttributeName) String() string {
	return string(n)
}

const (
	AttributeNameSub           AttributeName = "sub"
	AttributeNameEmailVerified AttributeName = "email_verified"
	AttributeNameEmail         AttributeName = "email"
	AttributeNameFamilyName    AttributeName = "family_name"
	AttributeNameGivenName     AttributeName = "given_name"
)

type UpdateAttributeValue struct {
	Value      *string
	NullUpdate bool
}

func (v *UpdateAttributeValue) IsUpdate() bool {
	return v.NullUpdate || v.Value != nil
}

func (v *UpdateAttributeValue) ToRow(name string) *cognitoidentityprovider.AttributeType {
	if !v.IsUpdate() {
		return nil
	}

	return &cognitoidentityprovider.AttributeType{
		Name:  &name,
		Value: v.Value,
	}
}

type UpdateAttributeValues []*UpdateAttributeValue

type UpdateAttribute struct {
	EmailVerified UpdateAttributeValue
	Email         UpdateAttributeValue
	FamilyName    UpdateAttributeValue
	GivenName     UpdateAttributeValue
}

func (a *UpdateAttribute) ToRows() []*cognitoidentityprovider.AttributeType {
	var rows []*cognitoidentityprovider.AttributeType

	if attribute := a.EmailVerified; attribute.IsUpdate() {
		rows = append(rows, attribute.ToRow(string(AttributeNameEmailVerified)))
	}
	if attribute := a.Email; attribute.IsUpdate() {
		rows = append(rows, attribute.ToRow(string(AttributeNameEmail)))
	}
	if attribute := a.FamilyName; attribute.IsUpdate() {
		rows = append(rows, attribute.ToRow(string(AttributeNameFamilyName)))
	}
	if attribute := a.GivenName; attribute.IsUpdate() {
		rows = append(rows, attribute.ToRow(string(AttributeNameGivenName)))
	}

	return rows
}

/** SignUp Input/Output **/
type (
	SignUpInput struct {
		Email      string
		Password   string
		FamilyName string
		GivenName  string
	}
	SignUpOutput struct {
		Sub string
	}
)

/** ConfirmSignUp Input/Output **/
type (
	ConfirmSignUpInput struct {
		Email            string
		ConfirmationCode string
	}
)

/** ResendConfirmationCode Input/Output **/
type (
	ResendConfirmationCodeInput struct {
		Email string
	}
)

/** AdminCreateUser Input/Output **/
type (
	AdminCreateUserInput struct {
		UserName string
		Email    string
	}
	AdminCreateUserOutput struct {
		UserName string
		Email    string
		Sub      string
	}
)

/** AdminUpdateUser Input **/
type (
	AdminUpdateUserInput struct {
		UserAttribute UpdateAttribute
	}
)

/** AdminGetUser Input/Output **/
type (
	AdminGetUserInput struct {
		Email string
	}
	AdminGetUserOutput struct {
		Sub           string
		EmailVerified string
		FamilyName    string
		GivenName     string
		Email         string
	}
)

func (i *AdminGetUserOutput) IsEmailVerified() bool {
	if i == nil {
		return false
	}
	return i.EmailVerified == "true"
}

func (in *AdminCreateUserInput) Validate() error {
	if len(in.UserName) == 0 {
		return errors.New("CreateUserInput.UserName is required")
	}
	if len(in.UserName) == 0 {
		return errors.New("CreateUserInput.Email is required")
	}
	return nil
}

/** Client **/
type Client interface {
	SignUp(ctx context.Context, input SignUpInput) (*SignUpOutput, error)
	ConfirmSignUp(ctx context.Context, input ConfirmSignUpInput) error
	ResendConfirmationCode(ctx context.Context, input ResendConfirmationCodeInput) error
	AdminCreateUser(ctx context.Context, input AdminCreateUserInput) (*AdminCreateUserOutput, error)
	AdminUpdateUser(ctx context.Context, input AdminUpdateUserInput) error
	AdminGetUser(ctx context.Context, input AdminGetUserInput) (*AdminGetUserOutput, error)
}
