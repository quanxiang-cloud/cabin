package error

import "fmt"

type Error struct {
	Code    int64  `json:"code"`
	Message string `json:"msg,omitempty"`
}

func New(code int64, params ...interface{}) Error {
	if len(params) > 0 {
		return Error{
			Code:    code,
			Message: fmt.Sprintf(Translation(code), params...),
		}
	}
	return Error{
		Code: code,
	}
}

func (e Error) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return Translation(e.Code)
}

func NewErrorWithString(code int64, message string) Error {
	return Error{
		Code:    code,
		Message: message,
	}
}
