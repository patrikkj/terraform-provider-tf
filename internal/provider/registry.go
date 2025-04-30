package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type Subprovider interface {
	Name() string
	Schema() schema.Schema
	Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse)
	DataSources() []func() datasource.DataSource
	Resources() []func() resource.Resource
}

var registry = []Subprovider{}

func RegisterProvider(subprovider Subprovider) {
	registry = append(registry, subprovider)
}
