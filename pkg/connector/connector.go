package connector

import (
	"context"
	"io"

	"github.com/conductorone/baton-galileo-ft/pkg/galileo"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Galileo struct {
	client *galileo.Client
}

// ResourceSyncers returns a ResourceSyncer for each resource type that should be synced from the upstream service.
func (g *Galileo) ResourceSyncers(ctx context.Context) []connectorbuilder.ResourceSyncer {
	return []connectorbuilder.ResourceSyncer{
		newUserBuilder(g.client),
		newGroupBuilder(g.client),
	}
}

// Asset takes an input AssetRef and attempts to fetch it using the connector's authenticated http client
// It streams a response, always starting with a metadata object, following by chunked payloads for the asset.
func (g *Galileo) Asset(ctx context.Context, asset *v2.AssetRef) (string, io.ReadCloser, error) {
	return "", nil, nil
}

// Metadata returns metadata about the connector.
func (g *Galileo) Metadata(ctx context.Context) (*v2.ConnectorMetadata, error) {
	return &v2.ConnectorMetadata{
		DisplayName: "Galileo-FT",
		Description: "Connector syncing Galileo-FT accounts and groups to Baton",
	}, nil
}

// Validate is called to ensure that the connector is properly configured. It should exercise any API credentials
// to be sure that they are valid.
func (g *Galileo) Validate(ctx context.Context) (annotations.Annotations, error) {
	err := g.client.Ping(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "galileo-ft-connector: failed to validate credentials")
	}

	return nil, nil
}

// New returns a new instance of the connector.
func New(ctx context.Context, cfg *galileo.Config) (*Galileo, error) {
	httpClient, err := uhttp.NewClient(ctx, uhttp.WithLogger(true, nil))
	if err != nil {
		return nil, err
	}

	client, err := galileo.NewClient(httpClient, cfg)
	if err != nil {
		return nil, err
	}

	return &Galileo{
		client: client,
	}, nil
}
