package main

import (
	"github.com/apparentlymart/terraform-provider-javascript/internal/provider"
	tfsdk "github.com/apparentlymart/terraform-sdk"
)

func main() {
	tfsdk.ServeProviderPlugin(provider.Provider())
}
