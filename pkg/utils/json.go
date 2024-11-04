package utils

import (
	"bytes"
	"encoding/json"
)

// ToJSONString converts the given status to a JSON string.
func ToJSONString(status interface{}) string {
	if status == nil {
		return ""
	}
	// Create a buffer to write the JSON string to
	buf := new(bytes.Buffer)
	// Encode the state into the buffer
	err := json.NewEncoder(buf).Encode(status)

	if err != nil {
		panic(err)
	}
	// Return a string representation of the buffer
	return buf.String()
}
