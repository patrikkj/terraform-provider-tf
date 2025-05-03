package ssh

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"sync"

	"golang.org/x/crypto/ssh"
)

// connectionKey represents the unique identifying parts of a connection
type connectionKey string

// hashSensitive takes a sensitive string and returns its MD5 hash
func hashSensitive(s string) string {
	if s == "<nil>" {
		return s
	}
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}

// SSHManager handles SSH connections for the provider
type SSHManager struct {
	providerConfig  *SSHConnectionConfig
	providerBastion *SSHConnectionConfig

	clientCache map[ConnectionKey]*ssh.Client
	lockMap     sync.Map // Map of mutexes per connection key
}

// getOrCreateLock returns a mutex for the given connection key
func (m *SSHManager) getOrCreateLock(key ConnectionKey) *sync.Mutex {
	actual, _ := m.lockMap.LoadOrStore(key, &sync.Mutex{})
	return actual.(*sync.Mutex)
}

// NewSSHManager creates a new SSH connection manager
func NewSSHManager(config *SSHConnectionConfig, bastion *SSHConnectionConfig) (*SSHManager, error) {
	return &SSHManager{
		providerConfig:  config,
		providerBastion: bastion,
		clientCache:     make(map[ConnectionKey]*ssh.Client),
	}, nil
}

// GetClient returns a cached SSH client or creates a new one if not found
func (m *SSHManager) GetClient(config SSHConnectionConfig, useProviderAsBastion bool, bastion *SSHConnectionConfig, fromClient *ssh.Client) (*ssh.Client, error) {
	key := NewConnectionKey(config, useProviderAsBastion, bastion, fromClient)

	// Get or create lock for this connection key
	lock := m.getOrCreateLock(key)
	lock.Lock()
	defer lock.Unlock()

	// Check if client exists in cache
	if client, ok := m.clientCache[key]; ok {
		return client, nil
	}

	// Create new client
	client, isNew, err := m.getClient(config, useProviderAsBastion, bastion, fromClient)
	if err != nil {
		return nil, err
	}

	// Cache the new client only if it was newly created
	if isNew {
		m.clientCache[key] = client
	}

	return client, nil
}

// getClient is the internal implementation that creates new SSH clients
func (m *SSHManager) getClient(config SSHConnectionConfig, useProviderAsBastion bool, bastion *SSHConnectionConfig, fromClient *ssh.Client) (*ssh.Client, bool, error) {
	// If useProviderAsBastion is true, use the provider client as bastion
	if useProviderAsBastion {
		providerClient, err := m.GetClient(*m.providerConfig, false, m.providerBastion, nil)
		if err != nil {
			return nil, false, fmt.Errorf("failed to connect to provider as bastion: %w", err)
		}
		client, err := m.GetClient(config, false, bastion, providerClient)
		return client, false, err
	}

	// If bastion is provided, get bastion client and recurse
	if bastion != nil {
		bastionClient, err := m.GetClient(*bastion, false, nil, fromClient)
		if err != nil {
			return nil, false, fmt.Errorf("failed to connect to bastion: %w", err)
		}
		client, err := m.GetClient(config, false, nil, bastionClient)
		return client, true, err
	}

	// If there is no configuration, fall back to provider client
	if fromClient == nil && config.Host == nil {
		providerClient, err := m.GetClient(*m.providerConfig, false, m.providerBastion, nil)
		if err != nil {
			return nil, false, fmt.Errorf("failed to connect to provider: %w", err)
		}
		return providerClient, true, nil
	}

	// Create ssh client configuration
	sshConfig, err := config.CreateSSHConfig()
	if err != nil {
		return nil, false, fmt.Errorf("failed to create ssh config: %w", err)
	}

	// Create target from port and host
	var port int64 = 22 // Default port
	if config.Port != nil {
		port = *config.Port
	}
	target := net.JoinHostPort(*config.Host, strconv.FormatInt(port, 10))

	// If there is no fromClient, return a new client using ssh.Dial
	if fromClient == nil {
		client, err := ssh.Dial("tcp", target, sshConfig)
		if err != nil {
			return nil, false, fmt.Errorf("failed to connect to target host: %w", err)
		}
		return client, true, nil
	}

	// Create new client through fromClient
	conn, err := fromClient.Dial("tcp", target)
	if err != nil {
		return nil, false, fmt.Errorf("failed to connect to target host through bastion: %w", err)
	}
	ncc, chans, reqs, err := ssh.NewClientConn(conn, target, sshConfig)
	if err != nil {
		return nil, false, fmt.Errorf("unable to create SSH connection through bastion: %w", err)
	}
	return ssh.NewClient(ncc, chans, reqs), true, nil
}
