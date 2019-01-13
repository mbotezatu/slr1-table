package grammar

import "strings"

// IsTerminal returns true if x is a terminal symbol or empty string
func IsTerminal(x string) bool {
	if x == "" {
		return true
	}

	firstChar := string([]rune(x)[0])
	if firstChar == strings.ToLower(firstChar) {
		return true
	}

	return false
}
