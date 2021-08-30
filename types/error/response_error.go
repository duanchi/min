package error

import "net/http"

type ResponseError struct {
	Message string
	ErrorCode int
	StatusCode int
	ErrorData interface{}
}

func (this ResponseError) Error() string {
	return this.Message
}

func (this ResponseError) Code() int {
	return this.ErrorCode
}

func (this ResponseError) Status() int {
	return this.StatusCode
}

func (this ResponseError) Data() interface{} {
	return this.ErrorData
}

func NewResponseError(message string, code int) ResponseError {
	return ResponseError{
		StatusCode: http.StatusInternalServerError,
		Message: message,
		ErrorCode: code,
	}
}

func NewResponseErrorWithData(message string, code int, data interface{}) ResponseError {
	return ResponseError{
		StatusCode: http.StatusInternalServerError,
		Message: message,
		ErrorCode: code,
		ErrorData: data,
	}
}