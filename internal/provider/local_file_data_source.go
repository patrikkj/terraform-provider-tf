package provider

import (
	"context"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type LocalFileDataSourceModel struct {
	Path         types.String `tfsdk:"path"`
	Content      types.String `tfsdk:"content"`
	Permissions  types.String `tfsdk:"permissions"`
	FailIfAbsent types.Bool   `tfsdk:"fail_if_absent"`
	Id           types.String `tfsdk:"id"`
}

var LocalFileDataSourceSchema = schema.Schema{
	Description: "Read local files",
	Attributes: map[string]schema.Attribute{
		"path":           schema.StringAttribute{Required: true, Description: "Path to the file"},
		"content":        schema.StringAttribute{Computed: true, Description: "Content of the file"},
		"permissions":    schema.StringAttribute{Computed: true, Optional: true, Description: "File permissions (e.g., '0644')"},
		"fail_if_absent": schema.BoolAttribute{Optional: true, Description: "Whether to fail if the file does not exist"},
		"id":             schema.StringAttribute{Computed: true, Description: "Unique identifier for this file"},
	},
}

var _ datasource.DataSource = &LocalFileDataSource{}

func NewLocalFileDataSource() datasource.DataSource {
	return &LocalFileDataSource{}
}

type LocalFileDataSource struct{}

func (d *LocalFileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_local_file"
}

func (d *LocalFileDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = LocalFileDataSourceSchema
}

func (d *LocalFileDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// No configuration needed
}

func (d *LocalFileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data LocalFileDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate a unique ID early, based on the path
	data.Id = types.StringValue(generateFileID(data.Path.ValueString(), time.Now()))

	content, err := os.ReadFile(data.Path.ValueString())
	if err != nil {
		if data.FailIfAbsent.ValueBool() {
			resp.Diagnostics.AddError("Failed to read file", err.Error())
			return
		}
		// If fail_if_absent is false, return empty content
		data.Content = types.StringValue("")
	} else {
		data.Content = types.StringValue(string(content))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
