package main

import (
	"github.com/aarani/terraform-provider-vultr/vultr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: vultr.Provider,
	})
}
