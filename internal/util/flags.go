package util

import (
	"os"
	"strings"
)

// DebugEnabled returns true / false if DEBUG=true is set
func DebugEnabled() bool {
	debugVal := os.Getenv("DEBUG")
	if strings.ToLower(debugVal) == "true" {
		return true
	}

	return false
}
