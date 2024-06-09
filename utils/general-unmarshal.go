package utils

import "encoding/json"

// THis is simply for the conversion process - so that I can use JSON and Go structs in parallel
func GenericUnmarshal(jsonStr string, target interface{}) error {
	return json.Unmarshal([]byte(jsonStr), target)
}
