package ssh

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	_provider "github.com/patrikkj/terraform-provider-tf/internal/provider"
	"github.com/patrikkj/terraform-provider-tf/internal/utils"
)

type SSHExecResourceModel struct {
	Command       types.String `tfsdk:"command"`
	Output        types.String `tfsdk:"output"`
	ExitCode      types.Int64  `tfsdk:"exit_code"`
	FailIfNonzero types.Bool   `tfsdk:"fail_if_nonzero"`
	OnDestroy     types.String `tfsdk:"on_destroy"`
	Id            types.String `tfsdk:"id"`

	// Connection details
	SSHConnectionModel
	UseProviderAsBastion types.Bool          `tfsdk:"use_provider_as_bastion"`
	Bastion              *SSHConnectionModel `tfsdk:"bastion"`
}

var SSHExecResourceSchema = schema.Schema{
	Description: "Execute commands over SSH with potential side effects",
	Attributes: map[string]schema.Attribute{
		"command":         schema.StringAttribute{Required: true, Description: "Command to execute"},
		"output":          schema.StringAttribute{Computed: true, Description: "Output of the command"},
		"exit_code":       schema.Int64Attribute{Computed: true, Description: "Exit code of the command"},
		"fail_if_nonzero": schema.BoolAttribute{Optional: true, Computed: true, Default: booldefault.StaticBool(true), Description: "Whether to fail if the command returns a non-zero exit code. Defaults to true if not specified."},
		"on_destroy":      schema.StringAttribute{Optional: true, Description: "Command to execute when the resource is destroyed"},
		"id":              schema.StringAttribute{Computed: true, Description: "Unique identifier for this execution"},

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

var _ resource.Resource = &SSHExecResource{}

func NewSSHExecResource() resource.Resource {
	return &SSHExecResource{}
}

type SSHExecResource struct {
	manager *SSHManager
}

func (r *SSHExecResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ssh_exec"
}

func (r *SSHExecResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = SSHExecResourceSchema
}

func (r *SSHExecResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	state, ok := req.ProviderData.(_provider.ProviderState)
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

func (r *SSHExecResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SSHExecResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set default values for computed fields
	if plan.Output.IsNull() {
		plan.Output = types.StringValue("")
	}
	if plan.ExitCode.IsNull() {
		plan.ExitCode = types.Int64Value(0)
	}

	// Generate a unique, stable ID before executing the command
	plan.Id = types.StringValue(utils.GenerateID(plan.Command.ValueString(), time.Now()))

	// Get SSH client
	client, err := r.manager.GetClient(
		*plan.SSHConnectionModel.toConfig(),
		plan.UseProviderAsBastion.ValueBool(),
		plan.Bastion.toConfig(),
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError("Failed to get SSH client", err.Error())
		return
	}

	// Execute the command
	output, exitCode, err := executeCommand(
		client,
		plan.Command.ValueString(),
		plan.FailIfNonzero.ValueBool(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Command execution failed", err.Error())
		return
	}
	plan.Output = types.StringValue(output)
	plan.ExitCode = types.Int64Value(exitCode)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SSHExecResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SSHExecResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// No need to re-run the command during read
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SSHExecResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state SSHExecResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve the original ID from state
	plan.Id = state.Id

	// Get SSH client
	client, err := r.manager.GetClient(
		*plan.SSHConnectionModel.toConfig(),
		plan.UseProviderAsBastion.ValueBool(),
		plan.Bastion.toConfig(),
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError("Failed to get SSH client", err.Error())
		return
	}

	// Execute the command
	output, exitCode, err := executeCommand(
		client,
		plan.Command.ValueString(),
		plan.FailIfNonzero.ValueBool(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Command execution failed", err.Error())
		return
	}
	plan.Output = types.StringValue(output)
	plan.ExitCode = types.Int64Value(exitCode)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SSHExecResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SSHExecResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If there's an on_destroy command, execute it
	if !data.OnDestroy.IsNull() {
		// Get SSH client
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

		_, _, err = executeCommand(
			client,
			data.OnDestroy.ValueString(),
			data.FailIfNonzero.ValueBool(),
		)
		if err != nil {
			resp.Diagnostics.AddError("Failed to execute destroy command", err.Error())
			return
		}
	}
}
