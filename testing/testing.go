package testing

import (
	"strings"
)

// StringAddress returns the address of the input string.
func StringAddress(str string) *string {
	return &str
}

// RightPad pads a string with spaces until the given limit, or it cuts the string to the given limit.
func RightPad(str string, limit int) string {
	str = str + strings.Repeat(" ", limit)
	return str[:limit]
}
