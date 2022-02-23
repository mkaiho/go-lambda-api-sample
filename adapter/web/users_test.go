package web

import (
	"errors"
	"fmt"
	"testing"

	"github.com/mkaiho/go-lambda-api-sample/entity"
	mockentity "github.com/mkaiho/go-lambda-api-sample/mocks/entity"
	mockusecase "github.com/mkaiho/go-lambda-api-sample/mocks/usecase"
	"github.com/mkaiho/go-lambda-api-sample/usecase"
	"github.com/stretchr/testify/assert"
)

func Test_listUsersHandler_Handle(t *testing.T) {
	type fields struct {
		listUsersUseCase usecase.ListUsersUseCase
	}
	tests := []struct {
		name   string
		fields fields
		want   *ListUsersResponse
	}{
		{
			name: "Return users",
			fields: fields{
				listUsersUseCase: func() usecase.ListUsersUseCase {
					const size = 3
					users := createDummyUsers(size)
					listUsersUseCase := new(mockusecase.ListUsersUseCase)
					listUsersUseCase.On("List").Return(users, nil)
					return listUsersUseCase
				}(),
			},
			want: &ListUsersResponse{
				Status: ResponseStatusOK,
				Users: []*UserDetail{
					{
						ID:   "001",
						Name: "testuser_001",
					},
					{
						ID:   "002",
						Name: "testuser_002",
					},
					{
						ID:   "003",
						Name: "testuser_003",
					},
				},
			},
		},
		{
			name: "Return dummy error when repository return error",
			fields: fields{
				listUsersUseCase: func() usecase.ListUsersUseCase {
					listUsersUseCase := new(mockusecase.ListUsersUseCase)
					listUsersUseCase.On("List").Return(nil, errors.New("dummy error"))
					return listUsersUseCase
				}(),
			},
			want: &ListUsersResponse{
				Status: ResponseStatusInternalServerError,
				Error: &ErrorResult{
					Message: "dummy error",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &listUsersHandler{
				listUsersUseCase: tt.fields.listUsersUseCase,
			}
			got := h.Handle()
			assert.Equal(t, tt.want, got, "listUsersHandler.Handle() = %v, want %v", got, tt.want)
		})
	}
}

func createDummyUsers(size int) []entity.User {
	testUsers := make([]entity.User, size)
	for i := 0; i < size; i++ {
		testUserID := new(mockentity.UserID)
		testUserID.On("Value").Return(fmt.Sprintf("%03d", i+1))

		testUser := new(mockentity.User)
		testUser.On("UserID").Return(testUserID)
		testUser.On("Name").Return(fmt.Sprintf("testuser_%03d", i+1))

		testUsers[i] = testUser
	}

	return testUsers
}
