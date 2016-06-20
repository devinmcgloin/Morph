package morphError

import "fmt"

type MorphError struct {
	Err     error
	Message string
	Code    int
}

func New(err error, message string, code int) MorphError {
	return MorphError{Err: err, Message: message, Code: code}
}

func (err MorphError) Error() string {
	return fmt.Sprintf("Code: %d; %s\n%s", err.Code, err.Message, err.Err.Error())
}
