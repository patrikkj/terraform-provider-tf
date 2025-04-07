package provider

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLocalFileResource(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test-dir-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test multiple file operations in a single step
			{
				Config: testAccLocalFileResourceConfig(tempDir),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Basic file checks
					resource.TestCheckResourceAttr("tf_local_file.test", "path", filepath.Join(tempDir, "test.txt")),
					resource.TestCheckResourceAttr("tf_local_file.test", "content", "hello world"),

					// File with permissions checks
					resource.TestCheckResourceAttr("tf_local_file.test_perms", "path", filepath.Join(tempDir, "test_perms.txt")),
					resource.TestCheckResourceAttr("tf_local_file.test_perms", "content", "secure content"),
					resource.TestCheckResourceAttr("tf_local_file.test_perms", "permissions", "0600"),

					// Nested file checks
					resource.TestCheckResourceAttr("tf_local_file.test_nested", "path", filepath.Join(tempDir, "nested/dir/test.txt")),
					resource.TestCheckResourceAttr("tf_local_file.test_nested", "content", "nested file content"),

					// File with specific permissions
					resource.TestCheckResourceAttr("tf_local_file.test_write", "path", filepath.Join(tempDir, "test_write.txt")),
					resource.TestCheckResourceAttr("tf_local_file.test_write", "content", "Hello from Terraform!\nThis is a test file."),
					resource.TestCheckResourceAttr("tf_local_file.test_write", "permissions", "0644"),
				),
			},
			// Test updates to multiple files
			{
				Config: testAccLocalFileResourceConfigUpdates(tempDir),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tf_local_file.test", "path", filepath.Join(tempDir, "test.txt")),
					resource.TestCheckResourceAttr("tf_local_file.test", "content", "updated content"),

					resource.TestCheckResourceAttr("tf_local_file.test_write", "path", filepath.Join(tempDir, "test_write.txt")),
					resource.TestCheckResourceAttr("tf_local_file.test_write", "content", "Updated content!\nThe file has been modified."),
					resource.TestCheckResourceAttr("tf_local_file.test_write", "permissions", "0644"),
				),
			},
		},
	})
}

func testAccLocalFileResourceConfig(tempDir string) string {
	return fmt.Sprintf(`
resource "tf_local_file" "test" {
  path    = "%s/test.txt"
  content = "hello world"
}

resource "tf_local_file" "test_perms" {
  path        = "%s/test_perms.txt"
  content     = "secure content"
  permissions = "0600"
}

resource "tf_local_file" "test_nested" {
  path    = "%s/nested/dir/test.txt"
  content = "nested file content"
}

resource "tf_local_file" "test_write" {
  path        = "%s/test_write.txt"
  content     = "Hello from Terraform!\nThis is a test file."
  permissions = "0644"
}
`, tempDir, tempDir, tempDir, tempDir)
}

func testAccLocalFileResourceConfigUpdates(tempDir string) string {
	return fmt.Sprintf(`
resource "tf_local_file" "test" {
  path    = "%s/test.txt"
  content = "updated content"
}

resource "tf_local_file" "test_perms" {
  path        = "%s/test_perms.txt"
  content     = "secure content"
  permissions = "0600"
}

resource "tf_local_file" "test_nested" {
  path    = "%s/nested/dir/test.txt"
  content = "nested file content"
}

resource "tf_local_file" "test_write" {
  path        = "%s/test_write.txt"
  content     = "Updated content!\nThe file has been modified."
  permissions = "0644"
}
`, tempDir, tempDir, tempDir, tempDir)
}
