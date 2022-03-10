package usecase

import (
	"github.com/mkaiho/go-lambda-api-sample/entity"
)

/** CreateUser **/
type CreateUserUseCase interface {
	Create(user entity.User) (entity.User, error)
}

func NewCreateUserUseCase(usersWriter entity.UsersWriter) CreateUserUseCase {
	return &createUserUseCase{
		usersWriter: usersWriter,
	}
}

type createUserUseCase struct {
	usersWriter entity.UsersWriter
}

func (u *createUserUseCase) Create(user entity.User) (entity.User, error) {
	user, err := u.usersWriter.Insert(user)
	if err != nil {
		if IsErrUseCase(err) {
			return nil, err
		}
		return nil, NewErrFatal(err.Error())
	}
	return user, nil
}
