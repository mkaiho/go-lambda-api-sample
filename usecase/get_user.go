package usecase

import (
	"github.com/mkaiho/go-lambda-api-sample/entity"
)

/** GetUser **/
type GetUserUseCase interface {
	Get(id entity.UserID) (entity.User, error)
}

func NewGetUserUseCase(usersReader entity.UsersReader) GetUserUseCase {
	return &getUserUseCase{
		usersReader: usersReader,
	}
}

type getUserUseCase struct {
	usersReader entity.UsersReader
}

func (u *getUserUseCase) Get(id entity.UserID) (entity.User, error) {
	user, err := u.usersReader.FindByID(id)
	if err != nil {
		if IsErrUseCase(err) {
			return nil, err
		}
		return nil, NewErrFatal(err.Error())
	}
	if user == nil {
		return nil, NewErrNotFoundEntity("user", map[string]string{
			"id": id.Value(),
		})
	}
	return user, nil
}
