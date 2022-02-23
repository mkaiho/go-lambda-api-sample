package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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

func NewResponse(resp *web.GetUserResponse) (*events.APIGatewayProxyResponse, error) {
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
	userID := event.PathParameters["user_id"]

	// repository
	var (
		idValidator entity.IDValidator = infrastructure.NewDummyIDValidator()
		idGenerator entity.IDGenerator = infrastructure.NewDummyIDGenerator()
		idManager                      = entity.NewIDManager(idValidator, idGenerator)
		usersReader entity.UsersReader = infrastructure.NewDummyUsersReader(idManager)
	)
	// usecase
	var (
		getUserUseCase usecase.GetUserUseCase = usecase.NewGetUserUseCase(usersReader)
	)
	// web handler
	var (
		getUserHandler web.GetUserHandler = web.NewGetUserHandler(idManager, getUserUseCase)
	)

	result := getUserHandler.Handle(web.GetUserRequest{
		ID: userID,
	})
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
