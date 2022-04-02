package api

import (
	"encoding/json"
	"net/http"
)

type ErrorObject struct {
	ErrorCategory    string           `json:"errorCategory"`
	ErrorCode        string           `json:"errorCode"`
	ErrorDescription string           `json:"errorDescription,omitempty"`
}

type Error interface {
	error
	json.Marshaler
	GetCode() int
}

type ErrorBadRequest struct {
	InternalErrorDescription string
}

func (e *ErrorBadRequest) Error() string { return e.InternalErrorDescription }
func (e *ErrorBadRequest) GetCode() int  { return http.StatusBadRequest }
func (e *ErrorBadRequest) MarshalJSON() ([]byte, error) {
	errObj := ErrorObject{
		ErrorCategory:    "internal",
		ErrorCode:        "BAD_REQUEST",
		ErrorDescription: e.InternalErrorDescription,
	}
	obj, err := json.Marshal(errObj)
	return obj, err
}

type ErrorInternalServerError struct {
	InternalErrorDescription string
}

func (e *ErrorInternalServerError) Error() string { return e.InternalErrorDescription }
func (e *ErrorInternalServerError) GetCode() int  { return http.StatusInternalServerError }
func (e *ErrorInternalServerError) MarshalJSON() ([]byte, error) {
	errObj := ErrorObject{
		ErrorCategory:    "internal",
		ErrorCode:        "INTERNAL_ERROR",
		ErrorDescription: e.InternalErrorDescription,
	}
	obj, err := json.Marshal(errObj)
	return obj, err
}

type ErrorNotFound struct {}
func (e *ErrorNotFound) Error() string { return "Not found" }
func (e *ErrorNotFound) GetCode() int { return http.StatusNotFound }
func (e *ErrorNotFound) MarshalJSON() ([]byte, error) {
	errObj := ErrorObject{
		ErrorCategory: "internal",
		ErrorCode: "NOT_FOUND",
		ErrorDescription: "Not found"}
	obj, err := json.Marshal(errObj)
	return obj, err
}
