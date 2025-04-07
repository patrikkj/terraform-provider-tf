package provider

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type LocalFileResourceModel struct {
	Path            types.String `tfsdk:"path"`
	Content         types.String `tfsdk:"content"`
	Permissions     types.String `tfsdk:"permissions"`
	FailIfAbsent    types.Bool   `tfsdk:"fail_if_absent"`
	DeleteOnDestroy types.Bool   `tfsdk:"delete_on_destroy"`
	Id              types.String `tfsdk:"id"`
}

var LocalFileResourceSchema = schema.Schema{
	Description: "Manage local files with potential side effects",
	Attributes: map[string]schema.Attribute{
		"path":              schema.StringAttribute{Required: true, Description: "Path to the file"},
		"content":           schema.StringAttribute{Required: true, Description: "Content of the file"},
		"permissions":       schema.StringAttribute{Optional: true, Computed: true, Default: stringdefault.StaticString("0644"), Description: "File permissions (e.g., '0644')"},
		"fail_if_absent":    schema.BoolAttribute{Optional: true, Description: "Whether to fail if the file does not exist"},
		"delete_on_destroy": schema.BoolAttribute{Optional: true, Computed: true, Default: booldefault.StaticBool(true), Description: "Whether to delete the file when the resource is destroyed. Defaults to true."},
		"id":                schema.StringAttribute{Computed: true, Description: "Unique identifier for this file"},
	},
}

var _ resource.Resource = &LocalFileResource{}

func NewLocalFileResource() resource.Resource {
	return &LocalFileResource{}
}

type LocalFileResource struct{}

func (r *LocalFileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_local_file"
}

func (r *LocalFileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = LocalFileResourceSchema
}

func (r *LocalFileResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// No configuration needed
}

func (r *LocalFileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data LocalFileResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate a unique, stable ID before writing the file
	data.Id = types.StringValue(generateFileID(data.Path.ValueString(), time.Now()))

	// Create parent directories if they don't exist
	dir := filepath.Dir(data.Path.ValueString())
	if err := os.MkdirAll(dir, 0755); err != nil {
		resp.Diagnostics.AddError("Failed to create directory", err.Error())
		return
	}

	// Write the file
	if err := os.WriteFile(data.Path.ValueString(), []byte(data.Content.ValueString()), 0644); err != nil {
		resp.Diagnostics.AddError("Failed to write file", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *LocalFileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data LocalFileResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	content, err := os.ReadFile(data.Path.ValueString())
	if err != nil {
		if os.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Failed to read file", err.Error())
		return
	}

	data.Content = types.StringValue(string(content))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *LocalFileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data LocalFileResourceModel

	// Get the current state
	var state LocalFileResourceModel
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

	// Create parent directories if they don't exist
	dir := filepath.Dir(data.Path.ValueString())
	if err := os.MkdirAll(dir, 0755); err != nil {
		resp.Diagnostics.AddError("Failed to create directory", err.Error())
		return
	}

	// Write the file
	if err := os.WriteFile(data.Path.ValueString(), []byte(data.Content.ValueString()), 0644); err != nil {
		resp.Diagnostics.AddError("Failed to update file", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *LocalFileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data LocalFileResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if delete_on_destroy is set to false
	if !data.DeleteOnDestroy.IsNull() && !data.DeleteOnDestroy.ValueBool() {
		// Skip deletion if delete_on_destroy is false
		return
	}

	if err := os.Remove(data.Path.ValueString()); err != nil {
		if !os.IsNotExist(err) {
			resp.Diagnostics.AddError("Failed to delete file", err.Error())
		}
	}
}
