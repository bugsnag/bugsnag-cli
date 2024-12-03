package utils

import "time"

// GetTodaysDate returns the current date formatted according to the specified layout string.
// If no format is provided (empty string), it defaults to the "2006-01-02" format.
//
// Parameters:
//   - format (string): A layout string for formatting the date. See Go's time package documentation
//     for formatting details.
//
// Returns:
// - string: The formatted current date.
func GetTodaysDate(format string) string {
	// Default format if no format is provided
	const defaultFormat = "2006-01-02"

	// Use the provided format if it's not empty; otherwise, use the default format
	if format == "" {
		format = defaultFormat
	}

	return time.Now().Format(format)
}
