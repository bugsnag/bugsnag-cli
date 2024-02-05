package utils

// DisplayBlankIfEmpty Displays a string value as <blank> if it's an empty string
func DisplayBlankIfEmpty(value string) string {
	if value == "" {
		return "<blank>"
	}
	return value
}
