package ssh

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"

	"golang.org/x/crypto/ssh"
)

// ConnectionKey represents the unique identifying parts of a connection
type ConnectionKey struct {
	Host               string
	User               string
	Port               string
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
		User:               stringValue(config.User),
		Port:               int64Value(config.Port),
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

// String returns a string representation of the ConnectionKey for debugging
func (k ConnectionKey) String() string {
	bastionStr := "<nil>"
	if k.Bastion != nil {
		bastionStr = k.Bastion.String()
	}
	return fmt.Sprintf("host=%s|user=%s|port=%s|pwd=%s|key=%s|useProviderBastion=%v|bastion=%s|from=%s",
		k.Host, k.User, k.Port, k.PasswordHash, k.PrivateKeyHash, k.UseProviderBastion, bastionStr, k.FromClient)
}

// stringValue safely converts a pointer to string to its value or empty string if nil
func stringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// int64Value safely converts a pointer to int64 to its string value or empty string if nil
func int64Value(i *int64) string {
	if i == nil {
		return ""
	}
	return strconv.FormatInt(*i, 10)
}

// hashSensitive takes a sensitive string and returns its MD5 hash
func hashSensitive(s string) string {
	if s == "" {
		return s
	}
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}
