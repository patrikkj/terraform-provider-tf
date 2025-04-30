package ssh

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/patrikkj/terraform-provider-tf/internal/provider"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSSHExecDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSSHExecDataSourceConfig(t),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Basic command execution
					resource.TestCheckResourceAttr("data.tf_ssh_exec.basic", "command", "echo hi"),
					resource.TestCheckResourceAttr("data.tf_ssh_exec.basic", "exit_code", "0"),
					resource.TestCheckResourceAttr("data.tf_ssh_exec.basic", "output", "hi\n"),

					// Whoami command
					resource.TestCheckResourceAttr("data.tf_ssh_exec.whoami", "command", "whoami"),
					resource.TestCheckResourceAttr("data.tf_ssh_exec.whoami", "output", getEnvVarOrSkip(t, "SSH_USER")+"\n"),
					resource.TestCheckResourceAttr("data.tf_ssh_exec.whoami", "exit_code", "0"),

					// Non-zero exit with fail_if_nonzero = false
					resource.TestCheckResourceAttr("data.tf_ssh_exec.nonzero_allowed", "command", "false"),
					resource.TestCheckResourceAttr("data.tf_ssh_exec.nonzero_allowed", "exit_code", "1"),
					resource.TestCheckResourceAttr("data.tf_ssh_exec.nonzero_allowed", "output", ""),

					// Multiline command
					resource.TestCheckResourceAttr("data.tf_ssh_exec.multiline", "exit_code", "0"),
					resource.TestCheckResourceAttr("data.tf_ssh_exec.multiline", "output", "Line 1\nLine 2\n"),

					// Script command
					resource.TestCheckResourceAttr("data.tf_ssh_exec.script", "exit_code", "0"),
					resource.TestMatchResourceAttr(
						"data.tf_ssh_exec.script",
						"output",
						regexp.MustCompile(`Hello\n-rw-.*\s+test.txt\n`),
					),
				),
			},
		},
	})
}

func testAccSSHExecDataSourceConfig(t *testing.T) string {
	return fmt.Sprintf(`
provider "tf" {
	ssh = {
		host     = "%s"
		user     = "%s"
		password = "%s"
	}
}

data "tf_ssh_exec" "basic" {
  command = "echo hi"
}

data "tf_ssh_exec" "whoami" {
  command = "whoami"
}

data "tf_ssh_exec" "nonzero_allowed" {
  command = "false"
  fail_if_nonzero = false
}

data "tf_ssh_exec" "multiline" {
  command = <<-EOF
	  echo "Line 1"
	  echo "Line 2"
	EOF
}

data "tf_ssh_exec" "script" {
  command = <<-EOT
      #!/bin/bash
      if [ ! -d "/tmp/test_dir" ]; then
          mkdir /tmp/test_dir
      fi
      cd /tmp/test_dir
      echo "Hello" > test.txt
      cat test.txt
      ls -l test.txt
    EOT
}
`, getEnvVarOrSkip(t, "SSH_HOST"), getEnvVarOrSkip(t, "SSH_USER"), getEnvVarOrSkip(t, "SSH_PASSWORD"))
}

// Test for expected failure when command returns non-zero with fail_if_nonzero = true
func TestAccSSHExecDataSource_FailIfNonZero(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccSSHExecDataSourceConfigFailIfNonZero(t),
				ExpectError: regexp.MustCompile(`command exited with non-zero status: 1`),
			},
		},
	})
}

func testAccSSHExecDataSourceConfigFailIfNonZero(t *testing.T) string {
	return fmt.Sprintf(`
provider "tf" {
	ssh = {
		host     = "%s"
		user     = "%s"
		password = "%s"
	}
}

data "tf_ssh_exec" "nonzero_fail" {
  command = "false"
  fail_if_nonzero = true
}

data "tf_ssh_exec" "nonzero_fail2" {
  command = "false"
  fail_if_nonzero = true
}
`, getEnvVarOrSkip(t, "SSH_HOST"), getEnvVarOrSkip(t, "SSH_USER"), getEnvVarOrSkip(t, "SSH_PASSWORD"))
}

// Add this new test function
func TestAccSSHExecDataSource_PrivateKey(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckPrivateKey(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSSHExecDataSourceConfigPrivateKey(t),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Basic command execution
					resource.TestCheckResourceAttr("data.tf_ssh_exec.basic", "command", "echo hi"),
					resource.TestCheckResourceAttr("data.tf_ssh_exec.basic", "exit_code", "0"),
					resource.TestCheckResourceAttr("data.tf_ssh_exec.basic", "output", "hi\n"),

					// Whoami command
					resource.TestCheckResourceAttr("data.tf_ssh_exec.whoami", "command", "whoami"),
					resource.TestCheckResourceAttr("data.tf_ssh_exec.whoami", "output", getEnvVarOrSkip(t, "SSH_USER")+"\n"),
					resource.TestCheckResourceAttr("data.tf_ssh_exec.whoami", "exit_code", "0"),
				),
			},
		},
	})
}

func testAccSSHExecDataSourceConfigPrivateKey(t *testing.T) string {
	return fmt.Sprintf(`
provider "tf" {
	ssh = {
		host        = "%s"
		user        = "%s"
		private_key = file("%s")
	}
}

data "tf_ssh_exec" "basic" {
	command = "echo hi"
}

data "tf_ssh_exec" "whoami" {
	command = "whoami"
}
`, getEnvVarOrSkip(t, "SSH_HOST"),
		getEnvVarOrSkip(t, "SSH_USER"),
		getEnvVarOrSkip(t, "SSH_PRIVATE_KEY_PATH"))
}
