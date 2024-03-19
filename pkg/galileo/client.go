package galileo

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/conductorone/baton-sdk/pkg/uhttp"
)

const (
	BaseHost = "api-sandbox.cv.gpsrv.com"

	RelatedAccountsEndpoint  = "/intserv/4.0/getRelatedAccounts"
	AccountOverviewEndpoint  = "/intserv/4.0/getAccountOverview"
	RootGroupsEndpoint       = "/intserv/4.0/getRootGroups"
	GroupHierarchyEndpoint   = "/intserv/4.0/getGroupHierarchy"
	GroupInfoEndpoint        = "/intserv/4.0/getGroupsInfo"
	GroupsToAccountsEndpoint = "/intserv/4.0/getAccountGroupRelationships"
)

type Config struct {
	Hostname       string `mapstructure:"hostname"`
	APILogin       string `mapstructure:"api-login"`
	APITransKey    string `mapstructure:"api-trans-key"`
	ProviderID     string `mapstructure:"provider-id"`
	PrimaryAccount string `mapstructure:"primary-account"`
}

type Client struct {
	httpClient *uhttp.BaseHttpClient
	config     *Config
	baseUrl    *url.URL
}

func NewClient(httpClient *http.Client, config *Config) *Client {
	b := &url.URL{
		Scheme: "https",
		Host:   BaseHost,
	}

	// Override the default host if a hostname is provided
	if config.Hostname != "" {
		b.Host = config.Hostname
	}

	return &Client{
		httpClient: uhttp.NewBaseHttpClient(httpClient),
		config:     config,
		baseUrl:    b,
	}
}

func (c *Client) GetPrimaryAccountNumber() string {
	return c.config.PrimaryAccount
}

func (c *Client) ListRelatedAccounts(ctx context.Context) ([]Account, error) {
	var accounts BaseResponse[RelatedAccountsResponse]

	data := &FormData{
		APILogin:    c.config.APILogin,
		APITransKey: c.config.APITransKey,
		ProviderID:  c.config.ProviderID,
		AccountNo:   c.config.PrimaryAccount,
	}

	err := c.post(ctx, RelatedAccountsEndpoint, prepareForm(data), &accounts)
	if err != nil {
		return nil, err
	}

	return accounts.Data.Children, nil
}

type AccountOverviewResponse struct {
	Profile *Customer `json:"profile"`
}

func (c *Client) GetCustomer(ctx context.Context, accountID string) (*Customer, error) {
	var res BaseResponse[AccountOverviewResponse]

	data := &FormData{
		APILogin:    c.config.APILogin,
		APITransKey: c.config.APITransKey,
		ProviderID:  c.config.ProviderID,
		AccountNo:   accountID,
	}

	err := c.post(ctx, AccountOverviewEndpoint, prepareForm(data), &res)
	if err != nil {
		return nil, err
	}

	return res.Data.Profile, nil
}

func (c *Client) ListRootGroups(ctx context.Context, pgVars *PaginationVars) ([]Group, uint, error) {
	var res ListReponse[Group]

	data := &FormData{
		APILogin:    c.config.APILogin,
		APITransKey: c.config.APITransKey,
		ProviderID:  c.config.ProviderID,
	}

	// prepare form arguments
	form := prepareForm(data)
	pgVars.PrepareVars(form)

	err := c.post(ctx, RootGroupsEndpoint, form, &res)
	if err != nil {
		return nil, 0, err
	}

	return res.Data, res.NumOfPages, nil
}

func mapGroupIDs(groups []GroupHierarchy) []string {
	var ids []string
	for _, g := range groups {
		ids = append(ids, g.ID)

		if len(g.Children) > 0 {
			ids = append(ids, mapGroupIDs(g.Children)...)
		}
	}

	return ids
}

func (c *Client) ListChildrenGroups(ctx context.Context, parentGroupID string) ([]string, error) {
	var res BaseResponse[[]GroupHierarchy]

	data := &FormData{
		APILogin:    c.config.APILogin,
		APITransKey: c.config.APITransKey,
		ProviderID:  c.config.ProviderID,
		GroupID:     parentGroupID,
	}

	err := c.post(ctx, GroupHierarchyEndpoint, prepareForm(data), &res)
	if err != nil {
		return nil, err
	}

	// The maximum number of levels below a root group is five, making six levels total.
	ids := mapGroupIDs(res.Data)

	return ids, nil
}

func (c *Client) GetGroupsInfo(ctx context.Context, groupIDs []string) ([]Group, error) {
	var res BaseResponse[[]Group]

	data := &FormData{
		APILogin:    c.config.APILogin,
		APITransKey: c.config.APITransKey,
		ProviderID:  c.config.ProviderID,
		GroupIDs:    groupIDs,
	}

	err := c.post(ctx, GroupInfoEndpoint, prepareForm(data), &res)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (c *Client) ListGroupMembers(ctx context.Context, groupID string) ([]GroupToAccounts, error) {
	var res BaseResponse[[]GroupToAccounts]

	data := &FormData{
		APILogin:    c.config.APILogin,
		APITransKey: c.config.APITransKey,
		ProviderID:  c.config.ProviderID,
		GroupID:     groupID,
	}

	err := c.post(ctx, GroupsToAccountsEndpoint, prepareForm(data), &res)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (c *Client) AddAccountToGroup(ctx context.Context, groupID, accountID string) error {
	data := &FormData{
		APILogin:    c.config.APILogin,
		APITransKey: c.config.APITransKey,
		ProviderID:  c.config.ProviderID,
		GroupID:     groupID,
		AccountIDs:  []string{accountID},
	}

	err := c.post(ctx, GroupsToAccountsEndpoint, prepareForm(data), nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RemoveAccountFromGroup(ctx context.Context, groupID, accountID string) error {
	data := &FormData{
		APILogin:    c.config.APILogin,
		APITransKey: c.config.APITransKey,
		ProviderID:  c.config.ProviderID,
		AccountIDs:  []string{accountID},
	}

	err := c.post(ctx, GroupsToAccountsEndpoint, prepareForm(data), nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) createRequest(ctx context.Context, path string, form *url.Values) (*http.Request, error) {
	u := *c.baseUrl
	u.Path = path

	req, err := c.httpClient.NewRequest(
		ctx,
		http.MethodPost,
		&u,
		WithContentTypeFormHeader(),
		WithResponseContentTypeJSON(),
		WithFormBody(form),
	)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) post(ctx context.Context, path string, form *url.Values, response interface{}) error {
	req, err := c.createRequest(ctx, path, form)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req, uhttp.WithJSONResponse(response))
	if err != nil {
		return fmt.Errorf("failed to do request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
