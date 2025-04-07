package provider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type LocalExecDataSourceModel struct {
	Command       types.String `tfsdk:"command"`
	Output        types.String `tfsdk:"output"`
	ExitCode      types.Int64  `tfsdk:"exit_code"`
	FailIfNonzero types.Bool   `tfsdk:"fail_if_nonzero"`
	Id            types.String `tfsdk:"id"`
}

var LocalExecDataSourceSchema = schema.Schema{
	Description: "Execute local commands",
	Attributes: map[string]schema.Attribute{
		"command":         schema.StringAttribute{Required: true, Description: "Command to execute"},
		"output":          schema.StringAttribute{Computed: true, Description: "Output of the command"},
		"exit_code":       schema.Int64Attribute{Computed: true, Description: "Exit code of the command"},
		"fail_if_nonzero": schema.BoolAttribute{Optional: true, Description: "Whether to fail if the command returns a non-zero exit code"},
		"id":              schema.StringAttribute{Computed: true, Description: "Unique identifier for this execution"},
	},
}

var _ datasource.DataSource = &LocalExecDataSource{}

func NewLocalExecDataSource() datasource.DataSource {
	return &LocalExecDataSource{}
}

type LocalExecDataSource struct{}

func (d *LocalExecDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_local_exec"
}

func (d *LocalExecDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = LocalExecDataSourceSchema
}

func (d *LocalExecDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// No configuration needed
}

func (d *LocalExecDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data LocalExecDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set default value for fail_if_nonzero if not specified
	if data.FailIfNonzero.IsNull() {
		data.FailIfNonzero = types.BoolValue(true)
	}

	// Generate ID early, based on the command
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
