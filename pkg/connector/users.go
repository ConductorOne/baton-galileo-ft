package connector

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-galileo-ft/pkg/galileo"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

type userBuilder struct {
	client       *galileo.Client
	resourceType *v2.ResourceType
}

func (u *userBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return userResourceType
}

func userResource(accID string, user *galileo.Customer, parentResource *v2.ResourceId) (*v2.Resource, error) {
	userProfile := map[string]interface{}{
		"first_name":   user.FirstName,
		"middle_name":  user.MiddleName,
		"last_name":    user.LastName,
		"address_1":    user.Address1,
		"address_2":    user.Address2,
		"city":         user.City,
		"state":        user.State,
		"postal_code":  user.PostalCode,
		"country":      user.CountryCode,
		"home_phone":   user.HomePhone,
		"mobile_phone": user.MobilePhone,
	}

	fullName := fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	resource, err := rs.NewUserResource(
		fullName,
		userResourceType,
		accID,
		[]rs.UserTraitOption{
			rs.WithUserProfile(userProfile),
			rs.WithEmail(user.Email, true),
			rs.WithStatus(v2.UserTrait_Status_STATUS_ENABLED),
			rs.WithAccountType(v2.UserTrait_ACCOUNT_TYPE_HUMAN),
		},
		rs.WithParentResourceID(parentResource),
	)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

func (u *userBuilder) GetAccountCustomer(ctx context.Context, accID string, parentResourceID *v2.ResourceId) (*v2.Resource, error) {
	customer, err := u.client.GetCustomer(ctx, accID)
	if err != nil {
		return nil, fmt.Errorf("galileo-ft-connector: failed to get customer: %w", err)
	}

	return userResource(accID, customer, parentResourceID)
}

func (u *userBuilder) ListRelatedCustomers(ctx context.Context, accID string, parentResourceID *v2.ResourceId) ([]*v2.Resource, error) {
	accounts, err := u.client.ListRelatedAccounts(ctx, accID)
	if err != nil {
		return nil, fmt.Errorf("galileo-ft-connector: failed to list related accounts: %w", err)
	}

	var rv []*v2.Resource
	for _, acc := range accounts {
		customer, err := u.client.GetCustomer(ctx, acc.ID)
		if err != nil {
			return nil, fmt.Errorf("galileo-ft-connector: failed to get customer: %w", err)
		}

		ur, err := userResource(acc.ID, customer, parentResourceID)
		if err != nil {
			return nil, fmt.Errorf("galileo-ft-connector: failed to create user resource: %w", err)
		}

		rv = append(rv, ur)
	}

	return rv, nil
}

// List returns all the users from the database as resource objects.
// Users include a UserTrait because they are the 'shape' of a standard user.
func (u *userBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, attrs rs.SyncOpAttrs) ([]*v2.Resource, *rs.SyncOpResults, error) {
	if parentResourceID == nil {
		return nil, nil, nil
	}

	group, err := u.client.ListGroupMembers(ctx, parentResourceID.Resource)
	if err != nil {
		return nil, nil, fmt.Errorf("galileo-ft-connector: failed to list accounts under group %s: %w", parentResourceID.Resource, err)
	}

	var rv []*v2.Resource
	for _, accID := range group.AccountIDs {
		// first get the customer of the parent user
		parent, err := u.GetAccountCustomer(ctx, accID, parentResourceID)
		if err != nil {
			return nil, nil, err
		}

		rv = append(rv, parent)

		// then get the related children accounts of the parent user
		accounts, err := u.ListRelatedCustomers(ctx, accID, parentResourceID)
		if err != nil {
			return nil, nil, err
		}

		rv = append(rv, accounts...)
	}

	return rv, nil, nil
}

// Entitlements always returns an empty slice for users.
func (u *userBuilder) Entitlements(_ context.Context, resource *v2.Resource, attrs rs.SyncOpAttrs) ([]*v2.Entitlement, *rs.SyncOpResults, error) {
	return nil, nil, nil
}

// Grants always returns an empty slice for users since they don't have any entitlements.
func (u *userBuilder) Grants(ctx context.Context, resource *v2.Resource, attrs rs.SyncOpAttrs) ([]*v2.Grant, *rs.SyncOpResults, error) {
	return nil, nil, nil
}

func newUserBuilder(client *galileo.Client) *userBuilder {
	return &userBuilder{
		client:       client,
		resourceType: userResourceType,
	}
}
