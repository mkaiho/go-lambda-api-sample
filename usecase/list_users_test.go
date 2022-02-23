package usecase

import (
	"errors"
	"fmt"
	"testing"

	"github.com/mkaiho/go-lambda-api-sample/entity"
	mockentity "github.com/mkaiho/go-lambda-api-sample/mocks/entity"
	"github.com/stretchr/testify/assert"
)

func Test_listUsersUseCase_List(t *testing.T) {
	type fields struct {
		usersReader entity.UsersReader
	}
	tests := []struct {
		name    string
		fields  fields
		want    assert.ValueAssertionFunc
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Return users",
			fields: fields{
				usersReader: func() entity.UsersReader {
					const size = 3
					users := createDummyUsers(size)
					reader := new(mockentity.UsersReader)
					reader.On("FindAll").Return(users, nil)
					return reader
				}(),
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				for i, user := range i1.([]entity.User) {
					wantUserID := fmt.Sprintf("%03d", i+1)
					if got := user.UserID().Value(); got != wantUserID {
						t.Errorf("listUsersUseCase.List() got.UserID().Value() = %v, want.UserID().Value() %v", got, wantUserID)
						return false
					}
					wantUserName := fmt.Sprintf("testuser_%03d", i+1)
					if got := user.Name(); got != wantUserName {
						t.Errorf("listUsersUseCase.List() got.Name() = %v, want.Name() %v", got, wantUserName)
						return false
					}
				}
				return true
			},
			wantErr: assert.NoError,
		},
		{
			name: "Return error when UsersReader return error",
			fields: fields{
				usersReader: func() entity.UsersReader {
					reader := new(mockentity.UsersReader)
					reader.On("FindAll").Return(nil, errors.New("dummy error"))
					return reader
				}(),
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				return assert.Nil(tt, i1, "listUsersUseCase.List() = %v, want %v", i1, nil)
			},
			wantErr: func(tt assert.TestingT, e error, i ...interface{}) bool {
				wantErr := "dummy error"
				return assert.EqualError(tt, e, wantErr, "listUsersUseCase.List() error = %v, wantErr %v", e, wantErr)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &listUsersUseCase{
				usersReader: tt.fields.usersReader,
			}
			got, err := u.List()
			if !tt.wantErr(t, err) {
				return
			}
			tt.want(t, got)
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
