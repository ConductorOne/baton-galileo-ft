package main

import (
	"context"
	"fmt"
	"os"

	"github.com/conductorone/baton-galileo-ft/pkg/connector"
	"github.com/conductorone/baton-galileo-ft/pkg/galileo"
	configSchema "github.com/conductorone/baton-sdk/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/field"
	"github.com/conductorone/baton-sdk/pkg/types"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	version       = "dev"
	connectorName = "baton-galileo-ft"
	apiLogin      = "api-login"
	apiTransKey   = "api-trans-key"
	providerID    = "provider-id"
	hostname      = "hostname"
)

var (
	apiLoginField = field.StringField(
		apiLogin,
		field.WithRequired(true),
		field.WithDescription("The username provided by Galileo-FT for API access."),
	)
	apiTransKeyField = field.StringField(
		apiTransKey,
		field.WithRequired(true),
		field.WithDescription("The password provided by Galileo-FT, used alongside the api-login."),
	)
	hostnameField = field.StringField(
		hostname,
		field.WithDescription("URL hostname for production hostname."),
	)
	providerIDField = field.StringField(
		providerID,
		field.WithRequired(true),
		field.WithDescription("A unique identifier from Galileo-FT representing your organization, used for tracking transactions and data."),
	)
	configurationFields = []field.SchemaField{apiLoginField, apiTransKeyField, providerIDField, hostnameField}
)

func main() {
	ctx := context.Background()
	_, cmd, err := configSchema.DefineConfiguration(ctx,
		connectorName,
		getConnector,
		field.NewConfiguration(configurationFields),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	cmd.Version = version
	err = cmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func getConnector(ctx context.Context, cfg *viper.Viper) (types.ConnectorServer, error) {
	l := ctxzap.Extract(ctx)
	cb, err := connector.New(ctx, &galileo.Config{
		Hostname:    cfg.GetString(hostname),
		APILogin:    cfg.GetString(apiLogin),
		APITransKey: cfg.GetString(apiTransKey),
		ProviderID:  cfg.GetString(providerID),
	})
	if err != nil {
		l.Error("error creating connector", zap.Error(err))
		return nil, err
	}

	c, err := connectorbuilder.NewConnector(ctx, cb)
	if err != nil {
		l.Error("error creating connector", zap.Error(err))
		return nil, err
	}

	return c, nil
}
