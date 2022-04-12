package web

import (
	"github.com/mkaiho/go-lambda-api-sample/entity"
	"github.com/mkaiho/go-lambda-api-sample/usecase"
)

type UserDetail struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

/** List users **/
type (
	ListUsersResponse struct {
		Status ResponseStatus
		Users  []*UserDetail `json:"users,omitempty"`
		Error  *ErrorResult  `json:"error,omitempty"`
	}
	ListUsersHandler interface {
		Handle() *ListUsersResponse
	}
)

func NewListUsersHandler(listUsersUseCase usecase.ListUsersUseCase) ListUsersHandler {
	return &listUsersHandler{
		listUsersUseCase: listUsersUseCase,
	}
}

type listUsersHandler struct {
	listUsersUseCase usecase.ListUsersUseCase
}

func (h *listUsersHandler) Handle() *ListUsersResponse {
	users, err := h.listUsersUseCase.List()
	if err != nil {
		errResult, status := makeErrorResult(err)
		return &ListUsersResponse{
			Status: status,
			Error:  errResult,
		}
	}

	details := make([]*UserDetail, len(users))
	for i, user := range users {
		details[i] = &UserDetail{
			ID:    user.UserID().Value(),
			Name:  user.Name(),
			Email: user.Email().Value(),
		}
	}

	return &ListUsersResponse{
		Status: ResponseStatusOK,
		Users:  details,
	}
}

/** Get user by ID **/
type (
	GetUserRequest struct {
		ID string `json:"id"`
	}
	GetUserResponse struct {
		Status ResponseStatus
		User   *UserDetail  `json:"user,omitempty"`
		Error  *ErrorResult `json:"error,omitempty"`
	}
	GetUserHandler interface {
		Handle(req GetUserRequest) *GetUserResponse
	}
)

func NewGetUserHandler(idm entity.IDManager, getUserUseCase usecase.GetUserUseCase) GetUserHandler {
	return &getUserHandler{
		idm:            idm,
		getUserUseCase: getUserUseCase,
	}
}

type getUserHandler struct {
	idm            entity.IDManager
	getUserUseCase usecase.GetUserUseCase
}

func (h *getUserHandler) Handle(req GetUserRequest) *GetUserResponse {
	id, err := h.idm.From(req.ID)
	if err != nil {
		errResult, status := makeErrorResult(err)
		return &GetUserResponse{
			Status: status,
			Error:  errResult,
		}
	}
	user, err := h.getUserUseCase.Get(id)
	if err != nil {
		errResult, status := makeErrorResult(err)
		return &GetUserResponse{
			Status: status,
			Error:  errResult,
		}
	}
	if user == nil {
		return nil
	}
	return &GetUserResponse{
		Status: ResponseStatusOK,
		User: &UserDetail{
			ID:    user.UserID().Value(),
			Name:  user.Name(),
			Email: user.Email().Value(),
		},
	}
}

/** Create users **/
type (
	CreateUserRequest struct {
		User *UserDetail `json:"user"`
	}
	CreateUserResponse struct {
		Status ResponseStatus
		User   *UserDetail  `json:"user,omitempty"`
		Error  *ErrorResult `json:"error,omitempty"`
	}
	CreateUserHandler interface {
		Handle(req CreateUserRequest) *CreateUserResponse
	}
)

func NewCreateUserHandler(idm entity.IDManager, createUserUseCase usecase.CreateUserUseCase) CreateUserHandler {
	return &createUserHandler{
		idm:               idm,
		createUserUseCase: createUserUseCase,
	}
}

type createUserHandler struct {
	idm               entity.IDManager
	createUserUseCase usecase.CreateUserUseCase
}

func (h *createUserHandler) Handle(req CreateUserRequest) *CreateUserResponse {
	email, err := entity.NewEmail(req.User.Email)
	if err != nil {
		errResult, status := makeErrorResult(err)
		return &CreateUserResponse{
			Status: status,
			Error:  errResult,
		}
	}
	user, err := h.createUserUseCase.Create(entity.NewUser(h.idm.Generate(), req.User.Name, email))
	if err != nil {
		errResult, status := makeErrorResult(err)
		return &CreateUserResponse{
			Status: status,
			Error:  errResult,
		}
	}
	return &CreateUserResponse{
		Status: ResponseStatusCreated,
		User: &UserDetail{
			ID:    user.UserID().Value(),
			Name:  user.Name(),
			Email: user.Email().Value(),
		},
	}
}

/** Delete user by ID **/
type (
	DeleteUserRequest struct {
		ID string `json:"id"`
	}
	DeleteUserResponse struct {
		Status ResponseStatus
		Error  *ErrorResult `json:"error,omitempty"`
	}
	DeleteUserHandler interface {
		Handle(req DeleteUserRequest) *DeleteUserResponse
	}
)

func NewDeleteUserHandler(idm entity.IDManager, deleteUserUseCase usecase.DeleteUserUseCase) DeleteUserHandler {
	return &deleteUserHandler{
		idm:               idm,
		deleteUserUseCase: deleteUserUseCase,
	}
}

type deleteUserHandler struct {
	idm               entity.IDManager
	deleteUserUseCase usecase.DeleteUserUseCase
}

func (h *deleteUserHandler) Handle(req DeleteUserRequest) *DeleteUserResponse {
	id, err := h.idm.From(req.ID)
	if err != nil {
		errResult, status := makeErrorResult(err)
		return &DeleteUserResponse{
			Status: status,
			Error:  errResult,
		}
	}
	err = h.deleteUserUseCase.Delete(id)
	if err != nil {
		errResult, status := makeErrorResult(err)
		return &DeleteUserResponse{
			Status: status,
			Error:  errResult,
		}
	}
	return &DeleteUserResponse{
		Status: ResponseStatusNoContent,
	}
}
