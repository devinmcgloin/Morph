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
	if err.Err != nil {
		return fmt.Sprintf("Code: %d; %s\n%s\n", err.Code, err.Message, err.Err.Error())
	}
	return fmt.Sprintf("Code: %d; %s\n", err.Code, err.Message)
}
