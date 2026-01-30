package connector

import (
	"context"
	"fmt"
	"strconv"

	"github.com/conductorone/baton-galileo-ft/pkg/galileo"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	ent "github.com/conductorone/baton-sdk/pkg/types/entitlement"
	"github.com/conductorone/baton-sdk/pkg/types/grant"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

const (
	RootGroupsType  = "root"
	GroupMembership = "member"
)

type groupBuilder struct {
	client       *galileo.Client
	resourceType *v2.ResourceType
}

func (g *groupBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return groupResourceType
}

func groupResource(group *galileo.Group) (*v2.Resource, error) {
	groupProfile := map[string]interface{}{
		"group-id":      group.ID,
		"legal-name":    group.LegalName,
		"business":      group.Business,
		"contact-email": group.ContactEmail,
		"contact-name":  group.ContactName,
	}

	options := []rs.ResourceOption{
		rs.WithAnnotation(&v2.ChildResourceType{ResourceTypeId: userResourceType.Id}),
	}

	if group.ParentGroupID != "" {
		parentID, err := rs.NewResourceID(groupResourceType, group.ParentGroupID)
		if err != nil {
			return nil, err
		}

		options = append(options, rs.WithParentResourceID(parentID))
	}

	resource, err := rs.NewGroupResource(
		group.Name,
		groupResourceType,
		group.ID,
		[]rs.GroupTraitOption{
			rs.WithGroupProfile(groupProfile),
		},
		options...,
	)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

// General characteristics of Group:
// - You can create multiple root groups in your core.
// - A group can have only one parent group.
// - An account can belong to only one group at a time.
// - The maximum number of levels below a root group is five, making six levels total.
// More information about groups and their hierarchy: https://docs.galileo-ft.com/pro/docs/creating-a-corporate-hierarchy
func (g *groupBuilder) List(ctx context.Context, _ *v2.ResourceId, attrs rs.SyncOpAttrs) ([]*v2.Resource, *rs.SyncOpResults, error) {
	var page int
	var err error
	pToken := attrs.PageToken
	if pToken != "" {
		page, err = strconv.Atoi(pToken)
		if err != nil {
			return nil, nil, fmt.Errorf("galileo-ft-connector: failed to parse page token: %w", err)
		}
	}

	pgVars := galileo.NewPaginationVars(page, ResourcesPageSize)
	groups, totalNumOfPages, err := g.client.ListRootGroups(ctx, pgVars)
	if err != nil {
		return nil, nil, fmt.Errorf("galileo-ft-connector: failed to list root groups: %w", err)
	}

	var rv []*v2.Resource
	for _, rootGroup := range groups {
		// Create a resource for each root group.
		gr, err := groupResource(&rootGroup) // #nosec G601
		if err != nil {
			return nil, nil, fmt.Errorf("galileo-ft-connector: failed to create group resource: %w", err)
		}

		rv = append(rv, gr)

		// Check if have any children.
		childrenGroupIDs, err := g.client.ListChildrenGroups(ctx, rootGroup.ID)
		if err != nil {
			return nil, nil, fmt.Errorf("galileo-ft-connector: failed to list children groups: %w", err)
		}

		if len(childrenGroupIDs) == 0 {
			continue
		}

		// Fetch information about children groups
		children, err := g.client.GetGroupsInfo(ctx, childrenGroupIDs)
		if err != nil {
			return nil, nil, fmt.Errorf("galileo-ft-connector: failed to get children groups info: %w", err)
		}

		// Create a resource for each child group.
		for _, group := range children {
			cgr, err := groupResource(&group) // #nosec G601
			if err != nil {
				return nil, nil, fmt.Errorf("galileo-ft-connector: failed to create group resource: %w", err)
			}

			rv = append(rv, cgr)
		}
	}

	results := &rs.SyncOpResults{}
	if page+1 < int(totalNumOfPages) {
		results.NextPage = strconv.Itoa(page + 1)
	}

	return rv, results, nil
}

func (g *groupBuilder) Entitlements(_ context.Context, resource *v2.Resource, attrs rs.SyncOpAttrs) ([]*v2.Entitlement, *rs.SyncOpResults, error) {
	var rv []*v2.Entitlement

	assignmentOptions := []ent.EntitlementOption{
		ent.WithGrantableTo(userResourceType),
		ent.WithDisplayName(fmt.Sprintf("Group %s %s", resource.DisplayName, GroupMembership)),
		ent.WithDescription(fmt.Sprintf("Group %s membership", resource.DisplayName)),
	}

	rv = append(rv, ent.NewAssignmentEntitlement(resource, GroupMembership, assignmentOptions...))

	return rv, nil, nil
}

func (g *groupBuilder) Grants(ctx context.Context, resource *v2.Resource, attrs rs.SyncOpAttrs) ([]*v2.Grant, *rs.SyncOpResults, error) {
	var rv []*v2.Grant

	group, err := g.client.ListGroupMembers(ctx, resource.Id.Resource)
	if err != nil {
		return nil, nil, fmt.Errorf("galileo-ft-connector: failed to list group members: %w", err)
	}

	for _, accID := range group.AccountIDs {
		accID, err := rs.NewResourceID(userResourceType, accID)
		if err != nil {
			return nil, nil, fmt.Errorf("galileo-ft-connector: failed to create user resource ID: %w", err)
		}

		rv = append(rv, grant.NewGrant(resource, GroupMembership, accID))
	}

	return rv, nil, nil
}

func (g *groupBuilder) Grant(ctx context.Context, principal *v2.Resource, entitlement *v2.Entitlement) (annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)

	if principal.Id.ResourceType != userResourceType.Id {
		l.Warn(
			"galileo-ft-connector: only users can be granted group membership",
			zap.String("principal_id", principal.Id.String()),
			zap.String("principal_type", principal.Id.ResourceType),
		)

		return nil, fmt.Errorf("galileo-ft-connector: only users can be granted group membership")
	}

	err := g.client.AddAccountToGroup(ctx, principal.Id.Resource, entitlement.Resource.Id.Resource)
	if err != nil {
		return nil, fmt.Errorf("galileo-ft-connector: failed to grant group membership: %w", err)
	}

	return nil, nil
}

func (g *groupBuilder) Revoke(ctx context.Context, grant *v2.Grant) (annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)

	principal := grant.Principal
	entitlement := grant.Entitlement

	if principal.Id.ResourceType != userResourceType.Id {
		l.Warn(
			"galileo-ft-connector: only users can have group membership revoked",
			zap.String("principal_id", principal.Id.String()),
			zap.String("principal_type", principal.Id.ResourceType),
		)

		return nil, fmt.Errorf("galileo-ft-connector: only users can have group membership revoked")
	}

	err := g.client.RemoveAccountFromGroup(ctx, principal.Id.Resource, entitlement.Resource.Id.Resource)
	if err != nil {
		return nil, fmt.Errorf("galileo-ft-connector: failed to revoke group membership: %w", err)
	}

	return nil, nil
}

func newGroupBuilder(client *galileo.Client) *groupBuilder {
	return &groupBuilder{
		client:       client,
		resourceType: groupResourceType,
	}
}
