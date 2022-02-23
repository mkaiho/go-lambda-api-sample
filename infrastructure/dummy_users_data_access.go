package infrastructure

import (
	"fmt"

	"github.com/mkaiho/go-lambda-api-sample/entity"
	mockentity "github.com/mkaiho/go-lambda-api-sample/mocks/entity"
)

var _ entity.UsersReader = (*dummyUsersReader)(nil)

func NewDummyUsersReader(idManager entity.IDManager) entity.UsersReader {
	return &dummyUsersReader{
		idManager: idManager,
	}
}

type dummyUsersReader struct {
	idManager entity.IDManager
}

func (r *dummyUsersReader) FindAll() ([]entity.User, error) {
	const size = 3
	testUsers := createDummyUsers(size)
	return testUsers, nil
}

func (r *dummyUsersReader) FindByID(id entity.UserID) (entity.User, error) {
	const size = 3
	testUsers := createDummyUsers(size)
	for _, user := range testUsers {
		if id.Value() == user.UserID().Value() {
			return user, nil
		}
	}

	return nil, nil
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
