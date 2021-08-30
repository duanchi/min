package error

import "net/http"

type ForbiddenError struct {
	Message string
	ErrorCode int
	StatusCode int
	ErrorData interface{}
}

func (this ForbiddenError) Error() string {
	return this.Message
}

func (this ForbiddenError) Code() int {
	return this.ErrorCode
}

func (this ForbiddenError) Status() int {
	return this.StatusCode
}

func (this ForbiddenError) Data() interface{} {
	return this.ErrorData
}

func NewForbiddenError(message string, code int) ForbiddenError {
	return ForbiddenError{
		StatusCode: http.StatusForbidden,
		Message: message,
		ErrorCode: code,
	}
}

func NewForbiddenErrorWithData(message string, code int, data interface{}) ForbiddenError {
	return ForbiddenError{
		StatusCode: http.StatusForbidden,
		Message: message,
		ErrorCode: code,
		ErrorData: data,
	}
}