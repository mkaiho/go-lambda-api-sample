package usecase

import "github.com/mkaiho/go-lambda-api-sample/entity"

/** ListUsers **/
type ListUsersUseCase interface {
	List() ([]entity.User, error)
}

func NewListUsersUseCase(usersReader entity.UsersReader) ListUsersUseCase {
	return &listUsersUseCase{
		usersReader: usersReader,
	}
}

type listUsersUseCase struct {
	usersReader entity.UsersReader
}

func (u *listUsersUseCase) List() ([]entity.User, error) {
	users, err := u.usersReader.FindAll()
	if err != nil {
		if IsErrUseCase(err) {
			return nil, err
		}
		return nil, NewErrFatal(err.Error())
	}

	return users, nil
}
