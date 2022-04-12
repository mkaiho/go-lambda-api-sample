package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mkaiho/go-lambda-api-sample/adapter/dynamodb"
	"github.com/mkaiho/go-lambda-api-sample/adapter/web"
	"github.com/mkaiho/go-lambda-api-sample/entity"
	"github.com/mkaiho/go-lambda-api-sample/infrastructure"
	"github.com/mkaiho/go-lambda-api-sample/usecase"
)

type ResponseBody struct {
	Result *web.UserDetail  `json:"user,omitempty"`
	Error  *web.ErrorResult `json:"error,omitempty"`
}

type Response events.APIGatewayProxyResponse

func NewResponse(resp *web.CreateUserResponse) (*events.APIGatewayProxyResponse, error) {
	body := &ResponseBody{
		Result: resp.User,
		Error:  resp.Error,
	}
	bs, err := body.bodyString()
	if err != nil {
		return nil, err
	}
	return &events.APIGatewayProxyResponse{
		StatusCode: resp.Status.Int(),
		Body:       bs,
	}, nil
}

func (body *ResponseBody) bodyString() (string, error) {
	if body == nil {
		return "", nil
	}
	b, err := json.Marshal(*body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func handle(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var err error

	// DynamoDB
	var (
		dynamodbClient               dynamodb.DynamoDBClient
		dynamodbAttributeValueMapper dynamodb.AttributeValueMapper
	)
	{
		dynamodbClient, err = infrastructure.NewDynamoDBClientFromEnv(os.Getenv("AWS_DEFAULT_REGION"))
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: web.ResponseStatusInternalServerError.Int(),
				Body:       err.Error(),
			}, err
		}
		dynamodbAttributeValueMapper = infrastructure.NewDynamoDBAttributeValueMapper()
	}
	// repository
	var (
		idValidator entity.IDValidator = infrastructure.NewULIDValidator()
		idGenerator entity.IDGenerator = infrastructure.NewULIDGenerator()
		idManager                      = entity.NewIDManager(idValidator, idGenerator)
		usersWriter entity.UsersWriter = dynamodb.NewUsersWriter(idManager, dynamodbClient, dynamodbAttributeValueMapper)
	)
	// usecase
	var (
		createUserUseCase usecase.CreateUserUseCase = usecase.NewCreateUserUseCase(usersWriter)
	)
	// web handler
	var (
		createUserHandler web.CreateUserHandler = web.NewCreateUserHandler(idManager, createUserUseCase)
	)

	req := new(web.CreateUserRequest)
	err = json.Unmarshal([]byte(event.Body), req)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: web.ResponseStatusInternalServerError.Int(),
			Body:       err.Error(),
		}, err
	}

	result := createUserHandler.Handle(*req)
	resp, err := NewResponse(result)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: web.ResponseStatusInternalServerError.Int(),
			Body:       err.Error(),
		}, err
	}
	return *resp, nil
}

func main() {
	lambda.Start(handle)
}
