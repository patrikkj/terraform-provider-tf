package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/joho/godotenv"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"tf": providerserver.NewProtocol6WithError(New("test")()),
}

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

func testAccPreCheck(t *testing.T) {
	// No pre-check needed for local provider
}
