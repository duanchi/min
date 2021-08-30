package error

import "net/http"

type NotFoundError struct {
	Message string
	ErrorCode int
	StatusCode int
	ErrorData interface{}
}

func (this NotFoundError) Error() string {
	return this.Message
}

func (this NotFoundError) Code() int {
	return this.ErrorCode
}

func (this NotFoundError) Status() int {
	return this.StatusCode
}

func (this NotFoundError) Data() interface{} {
	return this.ErrorData
}

func NewNotFoundError(message string, code int) NotFoundError {
	return NotFoundError{
		StatusCode: http.StatusNotFound,
		Message: message,
		ErrorCode: code,
	}
}

func NewNotFoundErrorWithData(message string, code int, data interface{}) NotFoundError {
	return NotFoundError{
		StatusCode: http.StatusNotFound,
		Message: message,
		ErrorCode: code,
		ErrorData: data,
	}
}