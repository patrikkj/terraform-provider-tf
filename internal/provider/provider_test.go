package provider

import (
	"fmt"

	"github.com/joho/godotenv"
)

func init() {
	// Load the .env file from multiple possible locations
	envFiles := []string{
		".env",
		"../.env",
		"../../.env",
		"../../../.env",
	}

	for _, file := range envFiles {
		if err := godotenv.Load(file); err == nil {
			fmt.Printf("Successfully loaded env file: %s\n", file)
			break
		}
	}
}
