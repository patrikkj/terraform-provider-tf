package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type LocalProviderModel struct {
	// No configuration needed for local provider
}

var LocalProviderSchema = schema.Schema{
	Description: "Provider for managing local files and executing local commands",
}

var _ provider.Provider = &LocalProvider{}

type LocalProvider struct {
	version string
}

func (p *LocalProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "tf"
	resp.Version = p.version
}

func (p *LocalProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = LocalProviderSchema
}

func (p *LocalProvider) Configure(_ context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// No configuration needed
}

func (p *LocalProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewLocalExecResource,
		NewLocalFileResource,
	}
}

func (p *LocalProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewLocalExecDataSource,
		NewLocalFileDataSource,
	}
}

func (p *LocalProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &LocalProvider{
			version: version,
		}
	}
}
