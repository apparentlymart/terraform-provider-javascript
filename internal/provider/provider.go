package provider

import (
	"context"

	tfsdk "github.com/apparentlymart/terraform-sdk"
	"github.com/apparentlymart/terraform-sdk/tfschema"
)

func Provider() *tfsdk.Provider {
	return &tfsdk.Provider{
		ConfigSchema: &tfschema.BlockType{
			Attributes: map[string]*tfschema.Attribute{},
		},
		ConfigureFn: func(ctx context.Context, config *config) (*client, tfsdk.Diagnostics) {
			return &client{}, nil
		},

		DataResourceTypes: map[string]tfsdk.DataResourceType{
			"javascript": javascriptDataResourceType(),
		},
	}
}

type config struct {
}

type client struct {
}
