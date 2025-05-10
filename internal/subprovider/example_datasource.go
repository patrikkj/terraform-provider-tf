package subprovider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var ExampleDataSourceSchema = schema.Schema{
	Description: "Example data source",
	Attributes: map[string]schema.Attribute{
		"path":    schema.StringAttribute{Required: true, Description: "Path to the file"},
		"content": schema.StringAttribute{Computed: true, Description: "Content of the file"},
		"id":      schema.StringAttribute{Computed: true, Description: "Unique identifier for this file"},
	},
}

type ExampleDataSourceModel struct {
	Path    types.String `tfsdk:"path"`
	Content types.String `tfsdk:"content"`
	Id      types.String `tfsdk:"id"`
}

type ExampleDataSource struct{}

var _ DataSource[ExampleDataSourceModel, *ExampleProviderModel, *ExampleProvider] = &ExampleDataSource{}

func (d *ExampleDataSource) Metadata() DataSourceMeta {
	return DataSourceMeta{
		Name:    "example",
		Version: "1.0.0",
		Schema:  ExampleDataSourceSchema,
	}
}

func (d *ExampleDataSource) Read(ctx context.Context, p *ExampleProvider, data ExampleDataSourceModel, resp *datasource.ReadResponse) {
	data.Id = types.StringValue(p.greeting)
}
