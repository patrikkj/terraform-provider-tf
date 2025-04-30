package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type ProviderModel struct {
	// No configuration needed for local provider
}

var ProviderSchema = schema.Schema{
	Description: "Provider for managing local files and executing local commands",
}

var _ provider.Provider = &Provider{}

type Provider struct {
	version string
}

func (p *Provider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "tf"
	resp.Version = p.version
}

func (p *Provider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = ProviderSchema
}

func (p *Provider) Configure(_ context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
}

func (p *Provider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return dataSourceRegistry
}

func (p *Provider) Resources(ctx context.Context) []func() resource.Resource {
	return resourceRegistry
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &Provider{version: version}
	}
}
