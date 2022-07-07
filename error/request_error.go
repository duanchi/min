package error

import "net/http"

type RequestError struct {
	Message string
	ErrorCode int
	StatusCode int
	ErrorData interface{}
}

func (this RequestError) Error() string {
	return this.Message
}

func (this RequestError) Code() int {
	return this.ErrorCode
}

func (this RequestError) Status() int {
	return this.StatusCode
}

func (this RequestError) Data() interface{} {
	return this.ErrorData
}

func NewRequestError(message string, code int) RequestError {
	return RequestError{
		StatusCode: http.StatusBadRequest,
		Message: message,
		ErrorCode: code,
	}
}

func NewRequestErrorWithData(message string, code int, data interface{}) RequestError {
	return RequestError{
		StatusCode: http.StatusBadRequest,
		Message: message,
		ErrorCode: code,
		ErrorData: data,
	}
}