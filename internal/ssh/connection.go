package ssh

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var SSHConnectionSchema = struct {
	Host                 schema.StringAttribute
	User                 schema.StringAttribute
	Password             schema.StringAttribute
	PrivateKey           schema.StringAttribute
	Port                 schema.Int64Attribute
	UseProviderAsBastion schema.BoolAttribute
	Bastion              schema.SingleNestedAttribute
}{
	Host:                 schema.StringAttribute{Description: "Override the provider's host configuration", Optional: true},
	User:                 schema.StringAttribute{Description: "Override the provider's user configuration", Optional: true},
	Password:             schema.StringAttribute{Description: "Override the provider's password configuration", Optional: true, Sensitive: true},
	PrivateKey:           schema.StringAttribute{Description: "Override the provider's private key configuration", Optional: true, Sensitive: true},
	Port:                 schema.Int64Attribute{Description: "The port number to connect to", Optional: true},
	UseProviderAsBastion: schema.BoolAttribute{Description: "Use the provider's connection as a bastion host", Optional: true},
	Bastion: schema.SingleNestedAttribute{
		Description: "Bastion host configuration",
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			"host":        schema.StringAttribute{Description: "The hostname or IP address of the bastion host", Required: true},
			"port":        schema.Int64Attribute{Description: "The port number of the bastion host", Optional: true},
			"user":        schema.StringAttribute{Description: "The username for bastion host authentication", Required: true},
			"password":    schema.StringAttribute{Description: "The password for bastion host authentication", Optional: true, Sensitive: true},
			"private_key": schema.StringAttribute{Description: "The private key for bastion host authentication", Optional: true, Sensitive: true},
		},
	},
}

// Common model for SSH connection configuration
type SSHConnectionModel struct {
	Host       types.String `tfsdk:"host"`
	User       types.String `tfsdk:"user"`
	Password   types.String `tfsdk:"password"`
	PrivateKey types.String `tfsdk:"private_key"`
	Port       types.Int64  `tfsdk:"port"`
}

type SSHConnectionConfig struct {
	Host       *string
	User       *string
	Password   *string
	PrivateKey *string
	Port       *int64
}

func (m *SSHConnectionModel) toConfig() *SSHConnectionConfig {
	// Handle nil receiver
	if m == nil {
		return nil
	}

	config := &SSHConnectionConfig{}

	if !m.Host.IsNull() {
		value := m.Host.ValueString()
		config.Host = &value
	}
	if !m.User.IsNull() {
		value := m.User.ValueString()
		config.User = &value
	}
	if !m.Password.IsNull() {
		value := m.Password.ValueString()
		config.Password = &value
	}
	if !m.PrivateKey.IsNull() {
		value := m.PrivateKey.ValueString()
		config.PrivateKey = &value
	}
	if !m.Port.IsNull() {
		value := m.Port.ValueInt64()
		config.Port = &value
	}

	return config
}
