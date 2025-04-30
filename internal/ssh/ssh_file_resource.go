package ssh

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrikkj/terraform-provider-tf/internal/provider"
	"github.com/patrikkj/terraform-provider-tf/internal/utils"
)

type SSHFileResourceModel struct {
	Path            types.String `tfsdk:"path"`
	Content         types.String `tfsdk:"content"`
	Permissions     types.String `tfsdk:"permissions"`
	FailIfAbsent    types.Bool   `tfsdk:"fail_if_absent"`
	DeleteOnDestroy types.Bool   `tfsdk:"delete_on_destroy"`
	Id              types.String `tfsdk:"id"`

	// Connection details
	SSHConnectionModel
	UseProviderAsBastion types.Bool          `tfsdk:"use_provider_as_bastion"`
	Bastion              *SSHConnectionModel `tfsdk:"bastion"`
}

var SSHFileResourceSchema = schema.Schema{
	Description: "Manage files over SSH with potential side effects",
	Attributes: map[string]schema.Attribute{
		"path":              schema.StringAttribute{Required: true, Description: "Path to the file"},
		"content":           schema.StringAttribute{Required: true, Description: "Content of the file"},
		"permissions":       schema.StringAttribute{Optional: true, Computed: true, Default: stringdefault.StaticString("0644"), Description: "File permissions (e.g., '0644')"},
		"fail_if_absent":    schema.BoolAttribute{Optional: true, Description: "Whether to fail if the file does not exist"},
		"delete_on_destroy": schema.BoolAttribute{Optional: true, Computed: true, Default: booldefault.StaticBool(true), Description: "Whether to delete the file when the resource is destroyed. Defaults to true."},
		"id":                schema.StringAttribute{Computed: true, Description: "Unique identifier for this file"},

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

var _ resource.Resource = &SSHFileResource{}

func NewSSHFileResource() resource.Resource {
	return &SSHFileResource{}
}

type SSHFileResource struct {
	manager *SSHManager
}

func (r *SSHFileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ssh_file"
}

func (r *SSHFileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = SSHFileResourceSchema
}

func (r *SSHFileResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	state, ok := req.ProviderData.(provider.ProviderState)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected provider.ProviderState, got: %T", req.ProviderData),
		)
		return
	}

	sshState, ok := state["ssh"].(*SSHProviderState)
	if !ok {
		resp.Diagnostics.AddError(
			"Missing SSH Provider State",
			"SSH provider state not found in provider data",
		)
		return
	}

	r.manager = sshState.manager
}

func (r *SSHFileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SSHFileResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate a unique, stable ID before writing the file
	data.Id = types.StringValue(utils.GenerateID(data.Path.ValueString(), time.Now()))

	client, err := r.manager.GetClient(
		*data.SSHConnectionModel.toConfig(),
		data.UseProviderAsBastion.ValueBool(),
		data.Bastion.toConfig(),
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError("Failed to get SSH client", err.Error())
		return
	}

	if err := writeFile(ctx, client, data.Path.ValueString(), data.Content.ValueString(), data.Permissions.ValueString()); err != nil {
		resp.Diagnostics.AddError("Failed to write file", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SSHFileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SSHFileResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.manager.GetClient(
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
		resp.State.RemoveResource(ctx)
		return
	}

	data.Content = types.StringValue(content)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SSHFileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SSHFileResourceModel

	// Get the current state
	var state SSHFileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the planned changes
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve the original ID from state
	data.Id = state.Id

	client, err := r.manager.GetClient(
		*data.SSHConnectionModel.toConfig(),
		data.UseProviderAsBastion.ValueBool(),
		data.Bastion.toConfig(),
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError("Failed to get SSH client", err.Error())
		return
	}

	if err := writeFile(ctx, client, data.Path.ValueString(), data.Content.ValueString(), data.Permissions.ValueString()); err != nil {
		resp.Diagnostics.AddError("Failed to update file", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SSHFileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SSHFileResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if delete_on_destroy is set to false
	if !data.DeleteOnDestroy.IsNull() && !data.DeleteOnDestroy.ValueBool() {
		// Skip deletion if delete_on_destroy is false
		return
	}

	client, err := r.manager.GetClient(
		*data.SSHConnectionModel.toConfig(),
		data.UseProviderAsBastion.ValueBool(),
		data.Bastion.toConfig(),
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError("Failed to get SSH client", err.Error())
		return
	}

	if err := deleteFile(client, data.Path.ValueString()); err != nil {
		resp.Diagnostics.AddError("Failed to delete file", err.Error())
		return
	}
}
