package provider

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLocalFileDataSource(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "test-file-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	// Write some content to the file
	content := "Hello, World!"
	if _, err := tempFile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	tempFile.Close()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccLocalFileDataSourceConfig(tempFile.Name()),
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

func testAccLocalFileDataSourceConfig(filePath string) string {
	return fmt.Sprintf(`
data "tf_local_file" "test" {
	path = "%s"
}

data "tf_local_file" "missing_optional" {
	path = "/nonexistent/file"
	fail_if_absent = false
}
`, filePath)
}

// Test for expected failure when reading non-existent file with fail_if_absent = true
func TestAccLocalFileDataSource_FailIfAbsent(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccLocalFileDataSourceConfigFailIfAbsent(),
				ExpectError: regexp.MustCompile(`Failed to read file`),
			},
		},
	})
}

func testAccLocalFileDataSourceConfigFailIfAbsent() string {
	return `
data "tf_local_file" "missing_required" {
	path           = "/path/to/nonexistent/file"
	fail_if_absent = true
}
`
}
