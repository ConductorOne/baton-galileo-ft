package config

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	apiLoginField = field.StringField(
		"api-login",
		field.WithDisplayName("Galileo-FT API Login"),
		field.WithRequired(true),
		field.WithDescription("The username provided by Galileo-FT for API access."),
		field.WithIsSecret(true),
	)
	apiTransKeyField = field.StringField(
		"api-trans-key",
		field.WithDisplayName("Galileo-FT API Transaction Key"),
		field.WithRequired(true),
		field.WithDescription("The password provided by Galileo-FT, used alongside the api-login."),
		field.WithIsSecret(true),
	)
	hostnameField = field.StringField(
		"hostname",
		field.WithDisplayName("Hostname"),
		field.WithDescription("URL hostname for production hostname."),
	)
	providerIDField = field.StringField(
		"provider-id",
		field.WithDisplayName("Provider ID"),
		field.WithRequired(true),
		field.WithDescription("A unique identifier from Galileo-FT representing your organization, used for tracking transactions and data."),
	)
	configurationFields = []field.SchemaField{apiLoginField, apiTransKeyField, providerIDField, hostnameField}
	fieldRelationships  = []field.SchemaFieldRelationship{}
)

//go:generate go run ./gen
var Config = field.NewConfiguration(
	configurationFields,
	field.WithConstraints(fieldRelationships...),
	field.WithConnectorDisplayName("Galileo-FT"),
	field.WithHelpUrl("/docs/baton/galileo-ft"),
	field.WithIconUrl("/static/app-icons/galileo-ft.svg"),
)
