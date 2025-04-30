package local

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/patrikkj/terraform-provider-tf/internal/provider"
	"github.com/patrikkj/terraform-provider-tf/internal/utils"
)

func TestAccLocalExecResource(t *testing.T) {
	create := utils.Heredoc(`
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

		resource "tf_local_exec" "multiline" {
			command = <<-EOF
				echo "Line 1"
				echo "Line 2"
			EOF
		}
	`)

	update := utils.Heredoc(`
		resource "tf_local_exec" "basic" {
			command = "echo 'updated'"
		}

		resource "tf_local_exec" "nonzero_allowed" {
			command        = "false"
			fail_if_nonzero = false
		}

		resource "tf_local_exec" "multiline" {
			command = <<-EOF
				echo "Line 1"
				echo "Line 2"
			EOF
		}
	`)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: create,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Basic command execution
					resource.TestCheckResourceAttr("tf_local_exec.basic", "command", "echo 'hello world'"),
					resource.TestCheckResourceAttr("tf_local_exec.basic", "exit_code", "0"),
					resource.TestCheckResourceAttr("tf_local_exec.basic", "output", "hello world\n"),

					// Non-zero exit with fail_if_nonzero = false
					resource.TestCheckResourceAttr("tf_local_exec.nonzero_allowed", "command", "false"),
					resource.TestCheckResourceAttr("tf_local_exec.nonzero_allowed", "exit_code", "1"),

					// Multiline command
					resource.TestCheckResourceAttr("tf_local_exec.multiline", "command", "echo \"Line 1\"\necho \"Line 2\"\n"),
					resource.TestCheckResourceAttr("tf_local_exec.multiline", "exit_code", "0"),
					resource.TestCheckResourceAttr("tf_local_exec.multiline", "output", "Line 1\nLine 2\n"),
				),
			},
			// Test updates to commands
			{
				Config: update,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tf_local_exec.basic", "command", "echo 'updated'"),
					resource.TestCheckResourceAttr("tf_local_exec.basic", "exit_code", "0"),
					resource.TestCheckResourceAttr("tf_local_exec.basic", "output", "updated\n"),
				),
			},
		},
	})
}
