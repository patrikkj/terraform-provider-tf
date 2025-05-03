package ssh

import (
	"crypto/md5"
	"encoding/hex"

	"golang.org/x/crypto/ssh"
)

// ConnectionKey represents the unique identifying parts of a connection
type ConnectionKey struct {
	Host               string
	User               string
	Port               int64
	PasswordHash       string
	PrivateKeyHash     string
	UseProviderBastion bool
	Bastion            *ConnectionKey
	FromClient         string
}

// NewConnectionKey creates a new ConnectionKey from connection parameters
func NewConnectionKey(config SSHConnectionConfig, useProviderAsBastion bool, bastion *SSHConnectionConfig, fromClient *ssh.Client) ConnectionKey {
	key := ConnectionKey{
		UseProviderBastion: useProviderAsBastion,
		Host:               stringValue(config.Host),
		User:               config.User,
		Port:               config.Port,
		PasswordHash:       hashSensitive(stringValue(config.Password)),
		PrivateKeyHash:     hashSensitive(stringValue(config.PrivateKey)),
	}

	// Add bastion details if present
	if bastion != nil {
		bastionKey := NewConnectionKey(*bastion, false, nil, fromClient)
		key.Bastion = &bastionKey
	}

	// Add fromClient address if present
	if fromClient != nil {
		key.FromClient = fromClient.RemoteAddr().String()
	}

	return key
}

// stringValue safely converts a pointer to string to its value or empty string if nil
func stringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// hashSensitive takes a sensitive string and returns its MD5 hash
func hashSensitive(s string) string {
	if s == "" {
		return s
	}
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}
