terraform {
  required_version = ">= 1.0"

  required_providers {
    tf = {
      source  = "patrikkj/tf"
      version = "0.0.1"
    }
  }
}

provider "tf" {}

# Create a local file
resource "tf_tf_local_file" "test" {
  path    = "test.txt"
  content = "Hello, World!"
}

# Execute a local command
resource "tf_tf_local_exec" "test" {
  command = "echo 'Command executed' > command_output.txt"
}
resource "tf_tf_local_exec" "test2" {
  command = "echo 'Command executed'"
}
