package web

import "github.com/mkaiho/go-lambda-api-sample/usecase"

type ErrorResult struct {
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func makeErrorResult(err error) (*ErrorResult, ResponseStatus) {
	useCaseError := usecase.ToUseCaseError(err)
	switch useCaseError.Type() {
	case usecase.ErrorTypeNotFoundEntity:
		e := usecase.ToErrNotFoundEntity(useCaseError)
		return makeErrNotFound(e), ResponseStatusNotFound
	case usecase.ErrorTypeFatal:
		fallthrough
	default:
		e := usecase.ToErrFatal(useCaseError)
		return makeErrInternal(e), ResponseStatusInternalServerError
	}
}

func makeErrNotFound(err usecase.EntityNotFoundError) *ErrorResult {
	return &ErrorResult{
		Message: err.Error(),
		Details: nil,
	}
}

func makeErrInternal(err usecase.FatalError) *ErrorResult {
	return &ErrorResult{
		Message: err.Error(),
		Details: nil,
	}
}
