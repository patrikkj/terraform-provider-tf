package ssh

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrikkj/terraform-provider-tf/internal/provider"
	"github.com/patrikkj/terraform-provider-tf/internal/utils"
)

type SSHFileDataSourceModel struct {
	Path         types.String `tfsdk:"path"`
	Content      types.String `tfsdk:"content"`
	Permissions  types.String `tfsdk:"permissions"`
	FailIfAbsent types.Bool   `tfsdk:"fail_if_absent"`
	Id           types.String `tfsdk:"id"`

	// Connection details
	SSHConnectionModel
	UseProviderAsBastion types.Bool          `tfsdk:"use_provider_as_bastion"`
	Bastion              *SSHConnectionModel `tfsdk:"bastion"`
}

var SSHFileDataSourceSchema = schema.Schema{
	Description: "Read files over SSH",
	Attributes: map[string]schema.Attribute{
		"path":           schema.StringAttribute{Required: true, Description: "Path to the file"},
		"content":        schema.StringAttribute{Computed: true, Description: "Content of the file"},
		"permissions":    schema.StringAttribute{Computed: true, Optional: true, Description: "File permissions (e.g., '0644')"},
		"fail_if_absent": schema.BoolAttribute{Optional: true, Description: "Whether to fail if the file does not exist"},
		"id":             schema.StringAttribute{Computed: true, Description: "Unique identifier for this file"},

		// Common SSH connection attributes
		"host":                    SSHConnectionSchema.Host,
		"user":                    SSHConnectionSchema.User,
		"password":                SSHConnectionSchema.Password,
		"private_key":             SSHConnectionSchema.PrivateKey,
		"port":                    SSHConnectionSchema.Port,
		"use_provider_as_bastion": SSHConnectionSchema.UseProviderAsBastion,
		"bastion":                 SSHConnectionSchema.Bastion,
	},
}

var _ datasource.DataSource = &SSHFileDataSource{}

func NewSSHFileDataSource() datasource.DataSource {
	return &SSHFileDataSource{}
}

type SSHFileDataSource struct {
	manager *SSHManager
}

func (d *SSHFileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ssh_file"
}

func (d *SSHFileDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = SSHFileDataSourceSchema
}

func (d *SSHFileDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	state, _ := req.ProviderData.(provider.ProviderState)
	sshState, _ := state["ssh"].(*SSHProviderState)
	d.manager = sshState.manager
}

func (d *SSHFileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SSHFileDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate a unique ID early, based on the path
	data.Id = types.StringValue(utils.GenerateID(data.Path.ValueString(), time.Now()))

	client, err := d.manager.GetClient(
		*data.SSHConnectionModel.toConfig(),
		data.UseProviderAsBastion.ValueBool(),
		data.Bastion.toConfig(),
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError("Failed to get SSH client", err.Error())
		return
	}

	content, err := readFile(client, data.Path.ValueString())
	if err != nil {
		if data.FailIfAbsent.ValueBool() {
			resp.Diagnostics.AddError("Failed to read file", err.Error())
			return
		}
		// If fail_if_absent is false, return empty content
		data.Content = types.StringValue("")
	} else {
		data.Content = types.StringValue(content)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
