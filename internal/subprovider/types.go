package subprovider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	pschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rschma "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type SubproviderMeta struct {
	Name    string
	Version string
	Schema  pschema.Schema
}

type DataSourceMeta struct {
	Name    string
	Version string
	Schema  dschema.Schema
}

type ResourceMeta struct {
	Name    string
	Version string
	Schema  rschma.Schema
}

type Subprovider[C any] interface {
	Metadata() SubproviderMeta
	Configure(ctx context.Context, config C, resp *provider.ConfigureResponse)
}

type DataSource[T any, C any, S interface{ Subprovider[C] }] interface {
	Metadata() DataSourceMeta
	Read(ctx context.Context, subprovider S, config T, resp *datasource.ReadResponse)
}

type Resource[T any, C any, S interface{ Subprovider[C] }] interface {
	Metadata() ResourceMeta
	Create(ctx context.Context, subprovider S, config T, resp *resource.CreateResponse)
	Read(ctx context.Context, subprovider S, config T, resp *resource.ReadResponse)
	Update(ctx context.Context, subprovider S, config T, resp *resource.UpdateResponse)
	Delete(ctx context.Context, subprovider S, config T, resp *resource.DeleteResponse)
}
