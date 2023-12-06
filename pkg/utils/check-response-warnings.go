package utils

import (
	"encoding/json"
	"fmt"
)

// CheckResponseWarnings takes a JSON-encoded response body as a byte slice and extracts
// the list of warnings from it. The function expects the JSON structure to have a key
// "warnings" containing an array of warning messages.
//
// Parameters:
//   - body: A byte slice representing the JSON-encoded response body.
//
// Returns:
//   - []interface{}: A slice of interfaces representing the list of warnings.
//   - error: An error if there was an issue decoding the JSON.
//     If there are no errors, the error will be nil.
func CheckResponseWarnings(body []byte) ([]interface{}, error) {
	// Unmarshal the JSON response into a map
	var responseMap map[string]interface{}
	err := json.Unmarshal(body, &responseMap)

	// Check for JSON decoding errors
	if err != nil {
		return nil, fmt.Errorf("Error decoding response JSON: %s", err.Error())
	}

	// Extract the list of warnings from the "warnings" key
	warnings, ok := responseMap["warnings"].([]interface{})

	// Check if the "warnings" key contains a valid array of interfaces
	if ok {
		return warnings, nil
	}

	// If the "warnings" key is not present or does not contain a valid array, return an error
	return nil, nil
}
