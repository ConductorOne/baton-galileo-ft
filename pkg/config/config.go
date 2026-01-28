package config

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	ApiLogin = field.StringField(
		"api-login",
		field.WithRequired(true),
		field.WithDescription("The username provided by Galileo-FT for API access."),
		field.WithDisplayName("API Login"),
	)
	ApiTransKey = field.StringField(
		"api-trans-key",
		field.WithRequired(true),
		field.WithDescription("The password provided by Galileo-FT, used alongside the api-login."),
		field.WithIsSecret(true),
		field.WithDisplayName("API Trans Key"),
	)
	Hostname = field.StringField(
		"hostname",
		field.WithDescription("URL hostname for production hostname."),
		field.WithDisplayName("Hostname"),
	)
	ProviderID = field.StringField(
		"provider-id",
		field.WithRequired(true),
		field.WithDescription("A unique identifier from Galileo-FT representing your organization, used for tracking transactions and data."),
		field.WithDisplayName("Provider ID"),
	)

	// FieldRelationships defines relationships between the fields.
	FieldRelationships = []field.SchemaFieldRelationship{}
)

//go:generate go run ./gen
var Configuration = field.NewConfiguration([]field.SchemaField{
	ApiLogin,
	ApiTransKey,
	ProviderID,
	Hostname,
}, field.WithConstraints(FieldRelationships...))
