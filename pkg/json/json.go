package json

import (
	"encoding/json"
	"strings"
)

// ParseStruct parses a given struct into and returns a JSON string that can be interpreted by nftlb.
func ParseStruct(data interface{}) (string, error) {
	indentedJSON, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return "", err
	}
	// Add spaces before every :, it's required by nftlb
	return strings.Replace(string(indentedJSON), "\":", "\" :", -1), nil
}
