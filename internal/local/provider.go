package local

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	_provider "github.com/patrikkj/terraform-provider-tf/internal/provider"
)

type LocalProvider struct{}

func (p *LocalProvider) Name() string {
	return "local"
}

func (p *LocalProvider) Schema() schema.Schema {
	return schema.Schema{
		Description: "Provider for managing local files and executing local commands",
	}
}

func (p *LocalProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// No configuration needed
}

func (p *LocalProvider) DataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewLocalExecDataSource,
		NewLocalFileDataSource,
	}
}

func (p *LocalProvider) Resources() []func() resource.Resource {
	return []func() resource.Resource{
		NewLocalExecResource,
		NewLocalFileResource,
	}
}

func init() {
	_provider.RegisterSubprovider(&LocalProvider{})
}
