package error

import "net/http"

type AuthorizeError struct {
	Message string
	ErrorCode int
	StatusCode int
	ErrorData interface{}
}

func (this AuthorizeError) Error() string {
	return this.Message
}

func (this AuthorizeError) Code() int {
	return this.ErrorCode
}

func (this AuthorizeError) Status() int {
	return this.StatusCode
}

func (this AuthorizeError) Data() interface{} {
	return this.ErrorData
}

func NewAuthorizeError(message string, code int) AuthorizeError {
	return AuthorizeError{
		StatusCode: http.StatusUnauthorized,
		Message: message,
		ErrorCode: code,
	}
}

func NewAuthorizeErrorWithData(message string, code int, data interface{}) AuthorizeError {
	return AuthorizeError{
		StatusCode: http.StatusUnauthorized,
		Message: message,
		ErrorCode: code,
		ErrorData: data,
	}
}