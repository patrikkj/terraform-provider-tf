package subprovider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"go.opencensus.io/resource"
)

type Provider struct {
	subprovider Subprovider
}

var _ provider.Provider = &Provider{}

type ProviderState map[string]interface{}

func (p *Provider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "tf"
	resp.Version = p.version
}

func (p *Provider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	// Build schema from subproviders
	attributes := map[string]schema.Attribute{}
	blocks := map[string]schema.Block{}
	for _, prov := range registry {
		subschema := prov.Schema()
		attributes[prov.Name()] = schema.SingleNestedAttribute{
			Description: subschema.Description,
			Attributes:  subschema.Attributes,
			Optional:    true,
		}
		for name, block := range subschema.Blocks {
			blocks[name] = block
		}
	}
	resp.Schema = schema.Schema{
		Attributes:  attributes,
		Blocks:      blocks,
		Description: "Internal provider for local commands, SSH and docker.",
	}
}

func (p *Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Initialize the provider state as a map
	state := make(ProviderState)
	resp.DataSourceData = state
	resp.ResourceData = state

	// Configure each provider
	for _, prov := range registry {
		prov.Configure(ctx, req, resp)
	}
}

func (p *Provider) DataSources(ctx context.Context) []func() datasource.DataSource {
	datasources := []func() datasource.DataSource{}
	for _, prov := range registry {
		datasources = append(datasources, prov.DataSources()...)
	}
	return datasources
}

func (p *Provider) Resources(ctx context.Context) []func() resource.Resource {
	resources := []func() resource.Resource{}
	for _, prov := range registry {
		resources = append(resources, prov.Resources()...)
	}
	return resources
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &Provider{
			version: version,
		}
	}
}

type S[T any] struct {
	name   string
	schema schema.Schema
	state  T
}
