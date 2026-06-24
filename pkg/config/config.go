package config

//go:generate go run ./gen

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	apiLoginField = field.StringField(
		"api-login",
		field.WithRequired(true),
		field.WithDisplayName("API Login"),
		field.WithDescription("The username provided by Galileo-FT for API access."),
		field.WithPlaceholder("Your Galileo-FT API login"),
		field.WithIsSecret(true),
	)
	apiTransKeyField = field.StringField(
		"api-trans-key",
		field.WithRequired(true),
		field.WithDisplayName("API Trans Key"),
		field.WithDescription("The password provided by Galileo-FT, used alongside the api-login."),
		field.WithPlaceholder("Your Galileo-FT API transaction key"),
		field.WithIsSecret(true),
	)
	providerIDField = field.StringField(
		"provider-id",
		field.WithRequired(true),
		field.WithDisplayName("Provider ID"),
		field.WithDescription("A unique identifier from Galileo-FT representing your organization, used for tracking transactions and data."),
		field.WithPlaceholder("Your Galileo-FT provider ID"),
	)
	hostnameField = field.StringField(
		"hostname",
		field.WithDisplayName("Hostname"),
		field.WithDescription("URL hostname for production hostname."),
		field.WithPlaceholder("Your Galileo-FT hostname"),
	)
	baseURLField = field.StringField(
		"base-url",
		field.WithDescription("Override the Galileo FT API URL (for testing)."),
		field.WithExportTarget(field.ExportTargetCLIOnly),
		field.WithHidden(true),
	)
)

var Config = field.NewConfiguration(
	[]field.SchemaField{
		apiLoginField,
		apiTransKeyField,
		providerIDField,
		hostnameField,
		baseURLField,
	},
	field.WithConnectorDisplayName("Galileo FT"),
	field.WithIconUrl("/static/app-icons/galileo-ft.svg"),
	field.WithHelpUrl("/docs/baton/galileo-ft"),
)
