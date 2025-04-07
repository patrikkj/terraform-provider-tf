package provider

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type LocalExecResourceModel struct {
	Command       types.String `tfsdk:"command"`
	Output        types.String `tfsdk:"output"`
	ExitCode      types.Int64  `tfsdk:"exit_code"`
	FailIfNonzero types.Bool   `tfsdk:"fail_if_nonzero"`
	OnDestroy     types.String `tfsdk:"on_destroy"`
	Id            types.String `tfsdk:"id"`
}

var LocalExecResourceSchema = schema.Schema{
	Description: "Execute local commands with potential side effects",
	Attributes: map[string]schema.Attribute{
		"command":         schema.StringAttribute{Required: true, Description: "Command to execute"},
		"output":          schema.StringAttribute{Computed: true, Description: "Output of the command"},
		"exit_code":       schema.Int64Attribute{Computed: true, Description: "Exit code of the command"},
		"fail_if_nonzero": schema.BoolAttribute{Optional: true, Computed: true, Default: booldefault.StaticBool(true), Description: "Whether to fail if the command returns a non-zero exit code. Defaults to true if not specified."},
		"on_destroy":      schema.StringAttribute{Optional: true, Description: "Command to execute when the resource is destroyed"},
		"id":              schema.StringAttribute{Computed: true, Description: "Unique identifier for this execution"},
	},
}

var _ resource.Resource = &LocalExecResource{}

func NewLocalExecResource() resource.Resource {
	return &LocalExecResource{}
}

type LocalExecResource struct{}

func (r *LocalExecResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_local_exec"
}

func (r *LocalExecResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = LocalExecResourceSchema
}

func (r *LocalExecResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// No configuration needed
}

func (r *LocalExecResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data LocalExecResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set default values for computed fields
	if data.Output.IsNull() {
		data.Output = types.StringValue("")
	}
	if data.ExitCode.IsNull() {
		data.ExitCode = types.Int64Value(0)
	}

	// Generate a unique, stable ID before executing the command
	data.Id = types.StringValue(generateExecID(data.Command.ValueString(), time.Now()))

	// Execute the command
	output, exitCode, err := executeLocalCommand(data.Command.ValueString(), data.FailIfNonzero.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError("Command execution failed", err.Error())
		return
	}
	data.Output = types.StringValue(output)
	data.ExitCode = types.Int64Value(exitCode)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *LocalExecResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data LocalExecResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// No need to re-run the command during read
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *LocalExecResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data LocalExecResourceModel

	// Get the current state
	var state LocalExecResourceModel
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

	// Execute the command
	output, exitCode, err := executeLocalCommand(data.Command.ValueString(), data.FailIfNonzero.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError("Command execution failed", err.Error())
		return
	}
	data.Output = types.StringValue(output)
	data.ExitCode = types.Int64Value(exitCode)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *LocalExecResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data LocalExecResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If there's an on_destroy command, execute it
	if !data.OnDestroy.IsNull() {
		_, _, err := executeLocalCommand(data.OnDestroy.ValueString(), data.FailIfNonzero.ValueBool())
		if err != nil {
			resp.Diagnostics.AddError("Failed to execute destroy command", err.Error())
			return
		}
	}
}

func executeLocalCommand(command string, failIfNonzero bool) (string, int64, error) {
	if command == "" {
		return "", 0, fmt.Errorf("empty command")
	}

	// Use the shell to execute the command
	cmd := exec.Command("sh", "-c", command)

	// Execute the command and capture output
	output, err := cmd.CombinedOutput()
	exitCode := int64(0)
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = int64(exitErr.ExitCode())
			if failIfNonzero {
				return string(output), exitCode, fmt.Errorf("command exited with code %d: %s", exitCode, string(output))
			}
		} else {
			return string(output), 0, fmt.Errorf("failed to execute command: %v", err)
		}
	}

	return string(output), exitCode, nil
}
