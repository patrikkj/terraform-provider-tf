package ssh

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/patrikkj/terraform-provider-tf/internal/utils"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// readFile reads a file's contents over SFTP
func readFile(client *ssh.Client, path string) (string, error) {
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return "", fmt.Errorf("failed to create SFTP client: %w", err)
	}
	defer sftpClient.Close()

	f, err := sftpClient.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return "", fmt.Errorf("failed to read file contents: %w", err)
	}

	// Return content as-is, preserving newlines
	return string(content), nil
}

// writeFile writes content to a file over SFTP
func writeFile(ctx context.Context, client *ssh.Client, path, content string, permissions string) error {
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return fmt.Errorf("failed to create SFTP client: %w", err)
	}
	defer sftpClient.Close()

	// Create directory if needed
	dirPath := filepath.Dir(path)
	if err := sftpClient.MkdirAll(dirPath); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create or truncate the file
	f, err := sftpClient.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	// Write content as-is, without modifying newlines
	if _, err := f.Write([]byte(content)); err != nil {
		return fmt.Errorf("failed to write file content: %w", err)
	}

	// Close the file before changing permissions
	f.Close()

	// Try to set permissions, but don't fail if it doesn't work
	mode := utils.ParseFileMode(permissions)
	tflog.Debug(ctx, fmt.Sprintf("Attempting to chmod %s to %s", path, mode))

	if err := sftpClient.Chmod(path, mode); err != nil {
		// Log the permission error but don't fail the operation
		tflog.Warn(ctx, fmt.Sprintf("Warning: Could not set permissions on %s to %s: %s",
			path, mode, err))
	}

	return nil
}

// deleteFile deletes a file over SFTP
func deleteFile(client *ssh.Client, path string) error {
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return fmt.Errorf("failed to create SFTP client: %w", err)
	}
	defer sftpClient.Close()

	if err := sftpClient.Remove(path); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}
