package spriocError

import (
	"encoding/json"
	"log"
)

type SpriocError struct {
	Err     error  `json:"-"`
	Message string `json:"error"`
	Code    int    `json:"code"`
}

func New(err error, message string, code int) SpriocError {
	return SpriocError{Err: err, Message: message, Code: code}
}

func (spriocErr SpriocError) Error() string {
	byteError, err := json.Marshal(spriocErr)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(byteError)
}
