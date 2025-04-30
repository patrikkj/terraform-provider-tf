package ssh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/patrikkj/terraform-provider-tf/internal/provider"
)

func TestAccSSHFileResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test multiple file operations in a single step
			{
				Config: testAccSSHFileResourceConfig(t),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Basic file checks
					resource.TestCheckResourceAttr("tf_ssh_file.test", "path", "/tmp/test.txt"),
					resource.TestCheckResourceAttr("tf_ssh_file.test", "content", "hello world"),

					// File with permissions checks
					resource.TestCheckResourceAttr("tf_ssh_file.test_perms", "path", "/tmp/test_perms.txt"),
					resource.TestCheckResourceAttr("tf_ssh_file.test_perms", "content", "secure content"),
					resource.TestCheckResourceAttr("tf_ssh_file.test_perms", "permissions", "0600"),

					// Nested file checks
					resource.TestCheckResourceAttr("tf_ssh_file.test_nested", "path", "/tmp/nested/dir/test.txt"),
					resource.TestCheckResourceAttr("tf_ssh_file.test_nested", "content", "nested file content"),

					// File with specific permissions
					resource.TestCheckResourceAttr("tf_ssh_file.test_write", "path", "/tmp/test_write.txt"),
					resource.TestCheckResourceAttr("tf_ssh_file.test_write", "content", "Hello from Terraform!\nThis is a test file."),
					resource.TestCheckResourceAttr("tf_ssh_file.test_write", "permissions", "0644"),
				),
			},
			// Test updates to multiple files
			{
				Config: testAccSSHFileResourceConfigUpdates(t),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tf_ssh_file.test", "path", "/tmp/test.txt"),
					resource.TestCheckResourceAttr("tf_ssh_file.test", "content", "updated content"),

					resource.TestCheckResourceAttr("tf_ssh_file.test_write", "path", "/tmp/test_write.txt"),
					resource.TestCheckResourceAttr("tf_ssh_file.test_write", "content", "Updated content!\nThe file has been modified."),
					resource.TestCheckResourceAttr("tf_ssh_file.test_write", "permissions", "0644"),
				),
			},
		},
	})
}

func testAccSSHFileResourceConfig(t *testing.T) string {
	return fmt.Sprintf(`
provider "tf" {
	ssh = {
		host     = "%s"
		user     = "%s"
		password = "%s"
	}
}

resource "tf_ssh_file" "test" {
  path    = "/tmp/test.txt"
  content = "hello world"
}

resource "tf_ssh_file" "test_perms" {
  path        = "/tmp/test_perms.txt"
  content     = "secure content"
  permissions = "0600"
}

resource "tf_ssh_file" "test_nested" {
  path    = "/tmp/nested/dir/test.txt"
  content = "nested file content"
}

resource "tf_ssh_file" "test_write" {
  path        = "/tmp/test_write.txt"
  content     = "Hello from Terraform!\nThis is a test file."
  permissions = "0644"
}
`, getEnvVarOrSkip(t, "SSH_HOST"), getEnvVarOrSkip(t, "SSH_USER"), getEnvVarOrSkip(t, "SSH_PASSWORD"))
}

func testAccSSHFileResourceConfigUpdates(t *testing.T) string {
	return fmt.Sprintf(`
provider "tf" {
	ssh = {
		host     = "%s"
		user     = "%s"
		password = "%s"
	}
}

resource "tf_ssh_file" "test" {
  path    = "/tmp/test.txt"
  content = "updated content"
}

resource "tf_ssh_file" "test_perms" {
  path        = "/tmp/test_perms.txt"
  content     = "secure content"
  permissions = "0600"
}

resource "tf_ssh_file" "test_nested" {
  path    = "/tmp/nested/dir/test.txt"
  content = "nested file content"
}

resource "tf_ssh_file" "test_write" {
  path        = "/tmp/test_write.txt"
  content     = "Updated content!\nThe file has been modified."
  permissions = "0644"
}
`, getEnvVarOrSkip(t, "SSH_HOST"), getEnvVarOrSkip(t, "SSH_USER"), getEnvVarOrSkip(t, "SSH_PASSWORD"))
}
