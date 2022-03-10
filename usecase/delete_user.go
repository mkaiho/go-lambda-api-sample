package usecase

import (
	"github.com/mkaiho/go-lambda-api-sample/entity"
)

/** DeleteUser **/
type DeleteUserUseCase interface {
	Delete(id entity.UserID) error
}

func NewDeleteUserUseCase(usersWriter entity.UsersWriter) DeleteUserUseCase {
	return &deleteUserUseCase{
		usersWriter: usersWriter,
	}
}

type deleteUserUseCase struct {
	usersWriter entity.UsersWriter
}

func (u *deleteUserUseCase) Delete(id entity.UserID) error {
	err := u.usersWriter.Delete(id)
	if err != nil {
		return NewErrFatal(err.Error())
	}

	return nil
}
