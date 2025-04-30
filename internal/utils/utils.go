package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/fs"
	"regexp"
	"strings"
	"time"
)

// ParseFileMode parses a file mode string into a fs.FileMode
func ParseFileMode(mode string) fs.FileMode {
	var result uint32
	if _, err := fmt.Sscanf(mode, "%o", &result); err != nil {
		return 0644 // Default to 0644 if parsing fails
	}
	return fs.FileMode(result)
}

// GenerateFileID creates a unique identifier for a file based on its path and timestamp
func GenerateFileID(path string, timestamp time.Time) string {
	h := md5.New()
	h.Write([]byte(path))
	h.Write([]byte(timestamp.UTC().Format(time.RFC3339)))
	return hex.EncodeToString(h.Sum(nil))
}

// GenerateExecID creates a unique identifier for a command based on its command and timestamp
func GenerateExecID(command string, timestamp time.Time) string {
	h := md5.New()
	h.Write([]byte(command))
	h.Write([]byte(timestamp.UTC().Format(time.RFC3339)))
	return hex.EncodeToString(h.Sum(nil))
}

// heredoc removes exactly one leading and trailing newline and dedents the string based on
// the indentation of the first line.
func Heredoc(s string) string {
	s = strings.TrimPrefix(s, "\n")
	s = strings.TrimSuffix(s, "\n")

	lines := strings.Split(s, "\n")
	if len(lines) == 0 {
		return ""
	}

	// Get indentation from first line using regex
	indentation := ""
	if match := regexp.MustCompile(`^[ \t]+`).FindString(lines[0]); match != "" {
		indentation = match
	}

	// Remove indentation from each line
	for i, line := range lines {
		lines[i] = strings.TrimPrefix(line, indentation)
	}
	return strings.Join(lines, "\n")
}
