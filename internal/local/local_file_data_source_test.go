package local

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/patrikkj/terraform-provider-tf/internal/provider"
	"github.com/patrikkj/terraform-provider-tf/internal/utils"
)

func TestAccLocalFileDataSource(t *testing.T) {
	// Create a temporary file for testing
	tempFile, _ := os.CreateTemp("", "test-file-*.txt")
	defer os.Remove(tempFile.Name())

	// Write some content to the file
	content := "Hello, World!"
	tempFile.WriteString(content)
	tempFile.Close()

	create := utils.Heredoc(fmt.Sprintf(`
		data "tf_local_file" "test" {
			path = "%s"
		}

		data "tf_local_file" "missing_optional" {
			path = "/nonexistent/file"
			fail_if_absent = false
		}
	`, tempFile.Name()))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: create,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Test reading an existing file
					resource.TestCheckResourceAttr("data.tf_local_file.test", "path", tempFile.Name()),
					resource.TestCheckResourceAttr("data.tf_local_file.test", "content", content),

					// Test reading a non-existent file with fail_if_absent = false
					resource.TestCheckResourceAttr("data.tf_local_file.missing_optional", "path", "/nonexistent/file"),
					resource.TestCheckResourceAttr("data.tf_local_file.missing_optional", "content", ""),
				),
			},
		},
	})
}

// Test for expected failure when reading non-existent file with fail_if_absent = true
func TestAccLocalFileDataSource_FailIfAbsent(t *testing.T) {
	create := utils.Heredoc(`
		data "tf_local_file" "missing_required" {
			path           = "/path/to/nonexistent/file"
			fail_if_absent = true
		}
	`)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      create,
				ExpectError: regexp.MustCompile(`Failed to read file`),
			},
		},
	})
}
