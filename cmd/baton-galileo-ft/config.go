package main

import (
	"context"

	"github.com/conductorone/baton-galileo-ft/pkg/galileo"
	"github.com/conductorone/baton-sdk/pkg/cli"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// config defines the external configuration required for the connector to run.
type config struct {
	cli.BaseConfig `mapstructure:",squash"` // Puts the base config options in the same place as the connector options
	galileo.Config `mapstructure:",squash"`
}

// validateConfig is run after the configuration is loaded, and should return an error if it isn't valid.
func validateConfig(ctx context.Context, cfg *config) error {
	if cfg.APILogin == "" || cfg.APITransKey == "" {
		return status.Error(codes.InvalidArgument, "api-login and api-trans-key must be provided, use --help for more information")
	}

	if cfg.ProviderID == "" {
		return status.Error(codes.InvalidArgument, "provider-id must be provided, use --help for more information")
	}

	return nil
}

func cmdFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String("api-login", "", "The username provided by Galileo-FT for API access. ($BATON_API_LOGIN)")
	cmd.PersistentFlags().String("api-trans-key", "", "The password provided by Galileo-FT, used alongside the api-login. ($BATON_API_TRANS_KEY)")
	cmd.PersistentFlags().String("provider-id", "", "A unique identifier from Galileo-FT representing your organization, used for tracking transactions and data. ($BATON_PROVIDER_ID)")
	cmd.PersistentFlags().String("hostname", "", "URL hostname for production hostname. ($BATON_HOSTNAME)")
}
