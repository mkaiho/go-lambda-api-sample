package usecase

import (
	"errors"
	"fmt"
	"strings"
)

var _ UseCaseError = (*useCaseError)(nil)
var _ UseCaseError = (*EntityNotFoundError)(nil)
var _ UseCaseError = (*EntityDuplicateError)(nil)
var _ UseCaseError = (*FatalError)(nil)

/** Error type **/
type ErrorType int

const (
	ErrorTypeNotFoundEntity  ErrorType = iota
	ErrorTypeDuplicateEntity ErrorType = iota
	ErrorTypeFatal           ErrorType = iota
)

/** Use case error **/
type UseCaseError interface {
	Type() ErrorType
	Error() string
}

type useCaseError struct {
	errType ErrorType
	message string
}

func IsErrUseCase(err error) bool {
	var errType UseCaseError
	return errors.As(err, &errType)
}

func ToUseCaseError(err error) UseCaseError {
	if !IsErrUseCase(err) {
		return NewErrFatal(err.Error())
	}
	return err.(UseCaseError)
}

func (e useCaseError) Type() ErrorType {
	return e.errType
}

func (e useCaseError) Error() string {
	return e.message
}

func newUseCaseError(errType ErrorType, message string) useCaseError {
	return useCaseError{
		errType: errType,
		message: message,
	}
}

/** EntityNotFoundError **/
type EntityNotFoundError struct {
	useCaseError
	entityName string
	params     map[string]string
}

func NewErrNotFoundEntity(entityName string, params map[string]string) EntityNotFoundError {
	var paramStrs []string
	for key, value := range params {
		paramStrs = append(paramStrs, fmt.Sprintf("%s=%s", key, value))
	}
	message := fmt.Sprintf("%s is not found. params: %s", entityName, strings.Join(paramStrs, ", "))
	return EntityNotFoundError{
		useCaseError: newUseCaseError(ErrorTypeNotFoundEntity, message),
		entityName:   entityName,
		params:       params,
	}
}

func ToErrNotFoundEntity(err UseCaseError) EntityNotFoundError {
	if !IsErrNotFound(err) {
		return NewErrNotFoundEntity("unknown", nil)
	}
	return err.(EntityNotFoundError)
}

func IsErrNotFound(err error) bool {
	var errType EntityNotFoundError
	return errors.As(err, &errType)
}

func (e EntityNotFoundError) EntityName() string {
	return e.entityName
}

func (e EntityNotFoundError) Params() map[string]string {
	return e.params
}

/** EntityDuplicateError **/
type EntityDuplicateError struct {
	useCaseError
	entityName string
	params     map[string]string
}

func NewErrDuplicateEntity(entityName string, params map[string]string) EntityDuplicateError {
	var paramStrs []string
	for key, value := range params {
		paramStrs = append(paramStrs, fmt.Sprintf("%s=%s", key, value))
	}
	message := fmt.Sprintf("%s is duplicate. params: %s", entityName, strings.Join(paramStrs, ", "))
	return EntityDuplicateError{
		useCaseError: newUseCaseError(ErrorTypeNotFoundEntity, message),
		entityName:   entityName,
		params:       params,
	}
}

func ToErrDuplicateEntity(err UseCaseError) EntityDuplicateError {
	if !IsErrNotFound(err) {
		return NewErrDuplicateEntity("unknown", nil)
	}
	return err.(EntityDuplicateError)
}

func IsErrDuplicate(err error) bool {
	var errType EntityDuplicateError
	return errors.As(err, &errType)
}

func (e EntityDuplicateError) EntityName() string {
	return e.entityName
}

func (e EntityDuplicateError) Params() map[string]string {
	return e.params
}

/** Fatal error **/
type FatalError struct {
	useCaseError
}

func NewErrFatal(message string) FatalError {
	return FatalError{
		useCaseError: newUseCaseError(ErrorTypeFatal, message),
	}
}

func IsErrFatal(err error) bool {
	var errType FatalError
	return errors.As(err, &errType)
}

func ToErrFatal(err UseCaseError) FatalError {
	if !IsErrFatal(err) {
		return NewErrFatal(err.Error())
	}
	return err.(FatalError)
}
