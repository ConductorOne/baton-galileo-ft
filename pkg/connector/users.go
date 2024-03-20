package connector

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-galileo-ft/pkg/galileo"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

type userBuilder struct {
	client       *galileo.Client
	resourceType *v2.ResourceType
}

func (u *userBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return userResourceType
}

func userResource(accID string, user *galileo.Customer) (*v2.Resource, error) {
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
	)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

func (u *userBuilder) GetPrimaryAccount(ctx context.Context) (*v2.Resource, error) {
	pID := u.client.GetPrimaryAccountNumber()
	customer, err := u.client.GetCustomer(ctx, pID)
	if err != nil {
		return nil, fmt.Errorf("galileo-ft-connector: failed to get customer: %w", err)
	}

	return userResource(pID, customer)
}

func (u *userBuilder) ListRelatedAccounts(ctx context.Context) ([]*v2.Resource, error) {
	accounts, err := u.client.ListRelatedAccounts(ctx)
	if err != nil {
		return nil, fmt.Errorf("galileo-ft-connector: failed to list related accounts: %w", err)
	}

	var rv []*v2.Resource
	for _, acc := range accounts {
		customer, err := u.client.GetCustomer(ctx, acc.ID)
		if err != nil {
			return nil, fmt.Errorf("galileo-ft-connector: failed to get customer: %w", err)
		}

		ur, err := userResource(acc.ID, customer)
		if err != nil {
			return nil, fmt.Errorf("galileo-ft-connector: failed to create user resource: %w", err)
		}

		rv = append(rv, ur)
	}

	return rv, nil
}

// List returns all the users from the database as resource objects.
// Users include a UserTrait because they are the 'shape' of a standard user.
func (u *userBuilder) List(ctx context.Context, _ *v2.ResourceId, _ *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var rv []*v2.Resource

	// first add the primary account customer as a user
	primary, err := u.GetPrimaryAccount(ctx)
	if err != nil {
		return nil, "", nil, err
	}

	rv = append(rv, primary)

	// then add all related accounts and their customers as users
	accounts, err := u.ListRelatedAccounts(ctx)
	if err != nil {
		return nil, "", nil, err
	}

	rv = append(rv, accounts...)

	return rv, "", nil, nil
}

// Entitlements always returns an empty slice for users.
func (u *userBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

// Grants always returns an empty slice for users since they don't have any entitlements.
func (u *userBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func newUserBuilder(client *galileo.Client) *userBuilder {
	return &userBuilder{
		client:       client,
		resourceType: userResourceType,
	}
}
