package web

import (
	"github.com/mkaiho/go-lambda-api-sample/entity"
	"github.com/mkaiho/go-lambda-api-sample/usecase"
)

type UserDetail struct {
	ID   string `json:"id"`
	Name string `json:"name"`
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
			ID:   user.UserID().Value(),
			Name: user.Name(),
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
			ID:   user.UserID().Value(),
			Name: user.Name(),
		},
	}
}
