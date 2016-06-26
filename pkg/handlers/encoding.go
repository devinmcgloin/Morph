package handlers

// Encoding takes the internal response, and formats them properly for external
// consumption.

// TODO need to take dbrefs and form them as links.
// TODO need to change metadata return types to strings not ratios.
// TODO fill out owner information.

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
