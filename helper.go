package main

import (
	"strings"
)

// ReadInput is now a method bound to the SchoolSystem struct.
// It uses the system's internal console scanner cleanly without needing parameters.
func (ss *SchoolSystem) ReadInput() string {
	if ss.ConsoleScanner.Scan() {
		return ss.ConsoleScanner.Text()
	}
	return ""
}

// SanitizeField is now a method bound to the SchoolSystem struct.
// It intercepts database-breaking characters before they can corrupt your flat files.
func (ss *SchoolSystem) SanitizeField(input string) string {
	trimmed := strings.TrimSpace(input)
	
	// Strip out colons (:) and pipes (|) so they don't break the data file columns
	escaped := strings.ReplaceAll(trimmed, ":", "")
	escaped = strings.ReplaceAll(escaped, "|", "")
	
	return escaped
}