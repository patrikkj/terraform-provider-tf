package subprovider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var ExampleProviderSchema = schema.Schema{
	Description: "Example provider",
	Attributes: map[string]schema.Attribute{
		"greeting": schema.StringAttribute{Required: true, Description: "Greeting to the world"},
	},
}

type ExampleProviderModel struct {
	Greeting types.String `tfsdk:"greeting"`
}

type ExampleProvider struct {
}

var _ Subprovider[*ExampleProviderModel] = &ExampleProvider{}

func (p *ExampleProvider) Metadata() SubproviderMeta {
	return SubproviderMeta{
		Name:    "example",
		Version: "1.0.0",
		Schema:  ExampleProviderSchema,
	}
}

func (p *ExampleProvider) Configure(ctx context.Context, config *ExampleProviderModel, resp *provider.ConfigureResponse) {
	p.greeting = config.Greeting.ValueString()
}
