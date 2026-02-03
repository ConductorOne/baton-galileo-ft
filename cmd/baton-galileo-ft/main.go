package main

import (
	"context"

	cfg "github.com/conductorone/baton-galileo-ft/pkg/config"
	"github.com/conductorone/baton-galileo-ft/pkg/connector"
	"github.com/conductorone/baton-sdk/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
)

var version = "dev"

func main() {
	ctx := context.Background()
	config.RunConnector(ctx,
		config.WithDefaultCapabilitiesConnectorBuilder(
			"baton-galileo-ft",
			version,
			cfg.Config,
			connector.New,
		),
		connectorbuilder.WithSkipFullSync(),
	)
}
