package spriocError

import "fmt"

type SpriocError struct {
	Err     error
	Message string
	Code    int
}

func New(err error, message string, code int) SpriocError {
	return SpriocError{Err: err, Message: message, Code: code}
}

func (err SpriocError) Error() string {
	return fmt.Sprintf("Code: %d; %s\n%s", err.Code, err.Message, err.Err.Error())
}
