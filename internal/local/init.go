package local

import (
	"github.com/patrikkj/terraform-provider-tf/internal/provider"
)

func init() {
	// Register resources
	provider.RegisterResource(NewLocalExecResource)
	provider.RegisterResource(NewLocalFileResource)

	// Register data sources
	provider.RegisterDataSource(NewLocalExecDataSource)
	provider.RegisterDataSource(NewLocalFileDataSource)
}

// TODO: Define provider subschema
