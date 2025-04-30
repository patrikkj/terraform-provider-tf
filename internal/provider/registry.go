package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
	resourceRegistry   []func() resource.Resource
	dataSourceRegistry []func() datasource.DataSource
)

func RegisterResource(factory func() resource.Resource) {
	resourceRegistry = append(resourceRegistry, factory)
}

func RegisterDataSource(factory func() datasource.DataSource) {
	dataSourceRegistry = append(dataSourceRegistry, factory)
}
