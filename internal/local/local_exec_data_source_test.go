package local

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/patrikkj/terraform-provider-tf/internal/provider"
	"github.com/patrikkj/terraform-provider-tf/internal/utils"
)

func TestAccLocalExecDataSource(t *testing.T) {
	create := utils.Heredoc(`
		data "tf_local_exec" "basic" {
			command = "echo hi"
		}

		data "tf_local_exec" "nonzero_allowed" {
			command = "false"
			fail_if_nonzero = false
		}

		data "tf_local_exec" "multiline" {
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
					resource.TestCheckResourceAttr("data.tf_local_exec.basic", "command", "echo hi"),
					resource.TestCheckResourceAttr("data.tf_local_exec.basic", "exit_code", "0"),
					resource.TestCheckResourceAttr("data.tf_local_exec.basic", "output", "hi\n"),

					// Non-zero exit with fail_if_nonzero = false
					resource.TestCheckResourceAttr("data.tf_local_exec.nonzero_allowed", "command", "false"),
					resource.TestCheckResourceAttr("data.tf_local_exec.nonzero_allowed", "exit_code", "1"),
					resource.TestCheckResourceAttr("data.tf_local_exec.nonzero_allowed", "output", ""),

					// Multiline command
					resource.TestCheckResourceAttr("data.tf_local_exec.multiline", "exit_code", "0"),
					resource.TestCheckResourceAttr("data.tf_local_exec.multiline", "output", "Line 1\nLine 2\n"),
				),
			},
		},
	})
}

// Test for expected failure when command returns non-zero with fail_if_nonzero = true
func TestAccLocalExecDataSource_FailIfNonZero(t *testing.T) {
	create := utils.Heredoc(`
		data "tf_local_exec" "nonzero_fail" {
			command = "false"
			fail_if_nonzero = true
		}
	`)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      create,
				ExpectError: regexp.MustCompile(`.*`),
			},
		},
	})
}
