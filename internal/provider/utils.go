package provider

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/fs"
	"time"
)

// Helper function to parse file mode
func parseFileMode(mode string) fs.FileMode {
	var result uint32
	if _, err := fmt.Sscanf(mode, "%o", &result); err != nil {
		return 0644 // Default to 0644 if parsing fails
	}
	return fs.FileMode(result)
}

// generateFileID creates a unique identifier for a file based on its path and timestamp
func generateFileID(path string, timestamp time.Time) string {
	h := md5.New()
	h.Write([]byte(path))
	h.Write([]byte(timestamp.UTC().Format(time.RFC3339)))
	return hex.EncodeToString(h.Sum(nil))
}

// generateExecID creates a unique identifier for a command based on its command and timestamp
func generateExecID(command string, timestamp time.Time) string {
	h := md5.New()
	h.Write([]byte(command))
	h.Write([]byte(timestamp.UTC().Format(time.RFC3339)))
	return hex.EncodeToString(h.Sum(nil))
}
