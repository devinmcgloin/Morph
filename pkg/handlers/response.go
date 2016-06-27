package handlers

import (
	"encoding/json"
	"log"
)

type Response struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"-"`
}

func Resp(message string, code int) Response {
	return Response{Message: message, Code: code}
}

func (response Response) Error() string {
	byteResp, err := json.Marshal(response)
	if err != nil {
		log.Println(err, response)
		return ""
	}
	return string(byteResp)
}

func (response Response) Format() []byte {
	byteResp, err := json.Marshal(response)
	if err != nil {
		log.Println(err, response)
		return []byte{}
	}
	return byteResp
}
