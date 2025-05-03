package ssh

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

// Optional is a utility type for handling optional values
type Optional[T any] struct {
	value *T
}

// NewOptional creates a new Optional with a value
func NewOptional[T any](value T) Optional[T] {
	return Optional[T]{value: &value}
}

// Value returns the value or the zero value if nil
func (o Optional[T]) Value() T {
	if o.value == nil {
		var zero T
		return zero
	}
	return *o.value
}

// String returns a string representation of the value or "<nil>" if nil
func (o Optional[T]) String() string {
	if o.value == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%v", *o.value)
}

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
		Host:               NewOptional(config.Host).String(),
		User:               NewOptional(config.User).String(),
		Port:               NewOptional(config.Port).String(),
		PasswordHash:       hashSensitive(NewOptional(config.Password).String()),
		PrivateKeyHash:     hashSensitive(NewOptional(config.PrivateKey).String()),
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
