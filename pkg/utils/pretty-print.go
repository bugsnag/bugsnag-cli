package utils

import (
	"bytes"
	"encoding/json"
)

// PrettyPrintJson - Prints JSON with indentations
func PrettyPrintJson(str string) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "    "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}

// PrettyPrintMap
func PrettyPrintMap(item map[string]string) (string, error) {
	prettyJSON, err := json.MarshalIndent(item, "", "    ")
	if err != nil {
		return "", err
	}
	return string(prettyJSON), nil
}
