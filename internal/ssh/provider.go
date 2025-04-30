package ssh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	_provider "github.com/patrikkj/terraform-provider-tf/internal/provider"
)

type (
	SSHProvider      struct{}
	SSHProviderState struct {
		manager *SSHManager
	}
	SSHProviderModel struct {
		SSHConnectionModel
		Bastion *SSHConnectionModel `tfsdk:"bastion"`
	}
	FullSSHProviderModel struct {
		Ssh SSHProviderModel `tfsdk:"ssh"`
	}
)

func (p *SSHProvider) Name() string {
	return "ssh"
}

func (p *SSHProvider) Schema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host":        schema.StringAttribute{Description: "The hostname or IP address of the target SSH server", Optional: true},
			"port":        schema.Int64Attribute{Description: "The port number of the target SSH server", Optional: true},
			"user":        schema.StringAttribute{Description: "The username for SSH authentication", Optional: true},
			"password":    schema.StringAttribute{Description: "The password for SSH authentication", Optional: true, Sensitive: true},
			"private_key": schema.StringAttribute{Description: "The private key for SSH authentication", Optional: true, Sensitive: true},
			"bastion": schema.SingleNestedAttribute{
				Description: "Bastion host configuration",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"host":        SSHConnectionSchema.Host,
					"port":        SSHConnectionSchema.Port,
					"user":        SSHConnectionSchema.User,
					"password":    SSHConnectionSchema.Password,
					"private_key": SSHConnectionSchema.PrivateKey,
				},
			},
		},
	}
}

func (p *SSHProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var fullConfig FullSSHProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &fullConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}
	config := fullConfig.Ssh

	// Set default port if not specified
	if config.Port.IsNull() {
		config.Port = types.Int64Value(22)
	}

	// Create the SSH manager with provider configuration
	manager, err := NewSSHManager(config.SSHConnectionModel.toConfig(), config.Bastion.toConfig())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create SSH manager",
			err.Error(),
		)
		return
	}

	// Set the SSH state
	state := resp.DataSourceData.(_provider.ProviderState)
	state["ssh"] = &SSHProviderState{
		manager: manager,
	}
}

func (p *SSHProvider) DataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewSSHExecDataSource,
		NewSSHFileDataSource,
	}
}

func (p *SSHProvider) Resources() []func() resource.Resource {
	return []func() resource.Resource{
		NewSSHExecResource,
		NewSSHFileResource,
	}
}

func init() {
	_provider.RegisterProvider(&SSHProvider{})
}
