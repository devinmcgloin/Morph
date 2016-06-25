package handlers

import (
	"bytes"
	"encoding/json"
	"log"
)

func getJSON(m map[string]string) []byte {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(m)
	if err != nil {
		return []byte{}
	}
	log.Println(err, m)

	return buf.Bytes()
}
