package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLocalExecDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccLocalExecDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Basic command execution
					resource.TestCheckResourceAttr("data.tf_local_exec.basic", "command", "echo hi"),
					resource.TestCheckResourceAttr("data.tf_local_exec.basic", "exit_code", "0"),
					resource.TestCheckResourceAttr("data.tf_local_exec.basic", "output", "hi\n"),

					// Whoami command
					resource.TestCheckResourceAttr("data.tf_local_exec.whoami", "command", "whoami"),
					resource.TestCheckResourceAttrSet("data.tf_local_exec.whoami", "output"),
					resource.TestCheckResourceAttr("data.tf_local_exec.whoami", "exit_code", "0"),

					// Non-zero exit with fail_if_nonzero = false
					resource.TestCheckResourceAttr("data.tf_local_exec.nonzero_allowed", "command", "false"),
					resource.TestCheckResourceAttr("data.tf_local_exec.nonzero_allowed", "exit_code", "1"),
					resource.TestCheckResourceAttr("data.tf_local_exec.nonzero_allowed", "output", ""),

					// Multiline command
					resource.TestCheckResourceAttr("data.tf_local_exec.multiline", "exit_code", "0"),
					resource.TestCheckResourceAttr("data.tf_local_exec.multiline", "output", "Line 1\nLine 2\n"),

					// Script command
					resource.TestCheckResourceAttr("data.tf_local_exec.script", "exit_code", "0"),
					resource.TestCheckResourceAttr("data.tf_local_exec.script", "output", "Hello\n"),
				),
			},
		},
	})
}

func testAccLocalExecDataSourceConfig() string {
	return `
data "tf_local_exec" "basic" {
  command = "echo hi"
}

data "tf_local_exec" "whoami" {
  command = "whoami"
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

data "tf_local_exec" "script" {
  command = "echo Hello"
}
`
}

// Test for expected failure when command returns non-zero with fail_if_nonzero = true
func TestAccLocalExecDataSource_FailIfNonZero(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccLocalExecDataSourceConfigFailIfNonZero(),
				ExpectError: regexp.MustCompile(`.*`),
			},
		},
	})
}

func testAccLocalExecDataSourceConfigFailIfNonZero() string {
	return `
data "tf_local_exec" "nonzero_fail" {
  command = "false"
  fail_if_nonzero = true
}
`
}
