package web

import (
	"net/http"
)

type ResponseStatus int

func (s ResponseStatus) Int() int {
	return int(s)
}

func (s ResponseStatus) String() string {
	return http.StatusText(int(s))
}

const (
	// 2XX
	ResponseStatusOK       ResponseStatus = http.StatusOK
	ResponseStatusCreated  ResponseStatus = http.StatusCreated
	ResponseStatusAccepted ResponseStatus = http.StatusAccepted
	// 4XX
	ResponseStatusBadRequest   ResponseStatus = http.StatusBadRequest
	ResponseStatusUnauthorized ResponseStatus = http.StatusUnauthorized
	ResponseStatusForbidden    ResponseStatus = http.StatusForbidden
	ResponseStatusNotFound     ResponseStatus = http.StatusNotFound
	// 5XX
	ResponseStatusInternalServerError ResponseStatus = http.StatusInternalServerError
)
