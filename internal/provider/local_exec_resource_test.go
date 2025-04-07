package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLocalExecResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccLocalExecResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Basic command execution
					resource.TestCheckResourceAttr("tf_local_exec.basic", "command", "echo 'hello world'"),
					resource.TestCheckResourceAttr("tf_local_exec.basic", "exit_code", "0"),
					resource.TestCheckResourceAttr("tf_local_exec.basic", "output", "hello world\n"),

					// Non-zero exit with fail_if_nonzero = false
					resource.TestCheckResourceAttr("tf_local_exec.nonzero_allowed", "command", "false"),
					resource.TestCheckResourceAttr("tf_local_exec.nonzero_allowed", "exit_code", "1"),

					// Whoami command
					resource.TestCheckResourceAttr("tf_local_exec.whoami", "command", "whoami"),
					resource.TestCheckResourceAttrSet("tf_local_exec.whoami", "output"),
					resource.TestCheckResourceAttr("tf_local_exec.whoami", "exit_code", "0"),

					// Multiline command
					resource.TestCheckResourceAttr("tf_local_exec.multiline", "command", "echo \"Line 1\"\necho \"Line 2\"\n"),
					resource.TestCheckResourceAttr("tf_local_exec.multiline", "exit_code", "0"),
					resource.TestCheckResourceAttr("tf_local_exec.multiline", "output", "Line 1\nLine 2\n"),
				),
			},
			// Test updates to commands
			{
				Config: testAccLocalExecResourceConfigUpdates(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tf_local_exec.basic", "command", "echo 'updated'"),
					resource.TestCheckResourceAttr("tf_local_exec.basic", "exit_code", "0"),
					resource.TestCheckResourceAttr("tf_local_exec.basic", "output", "updated\n"),
				),
			},
		},
	})
}

func testAccLocalExecResourceConfig() string {
	return `
resource "tf_local_exec" "basic" {
  command = "echo 'hello world'"
}

resource "tf_local_exec" "on_destroy" {
  command = "echo 'hello world'"
  on_destroy = "echo 'on_destroy' > /tmp/on_destroy"
}

resource "tf_local_exec" "nonzero_allowed" {
  command        = "false"
  fail_if_nonzero = false
}

resource "tf_local_exec" "whoami" {
  command = "whoami"
}

resource "tf_local_exec" "multiline" {
  command = <<-EOF
    echo "Line 1"
    echo "Line 2"
  EOF
}
`
}

func testAccLocalExecResourceConfigUpdates() string {
	return `
resource "tf_local_exec" "basic" {
  command = "echo 'updated'"
}

resource "tf_local_exec" "nonzero_allowed" {
  command        = "false"
  fail_if_nonzero = false
}

resource "tf_local_exec" "whoami" {
  command = "whoami"
}

resource "tf_local_exec" "multiline" {
  command = <<-EOF
    echo "Line 1"
    echo "Line 2"
  EOF
}
`
}
