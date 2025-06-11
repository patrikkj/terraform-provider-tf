package subprovider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

type DataSourceWrapper[T any, C any, S interface{ Subprovider[C] }] struct {
	Subprovider *S
	DataSource  DataSource[T, C, S]
}

func (d *DataSourceWrapper[T, C, S]) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + d.DataSource.Metadata().Name
}

func (d *DataSourceWrapper[T, C, S]) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = d.DataSource.Metadata().Schema
}

func (d *DataSourceWrapper[T, C, S]) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
}

func (d *DataSourceWrapper[T, C, S]) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data T
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	d.DataSource.Read(ctx, *d.Subprovider, data, resp)
}
