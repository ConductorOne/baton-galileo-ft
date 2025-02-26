package galileo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/conductorone/baton-sdk/pkg/uhttp"
)

const (
	BaseHost = "api-sandbox.cv.gpsrv.com"

	RelatedAccountsEndpoint        = "/intserv/4.0/getRelatedAccounts"
	AccountOverviewEndpoint        = "/intserv/4.0/getAccountOverview"
	RootGroupsEndpoint             = "/intserv/4.0/getRootGroups"
	GroupHierarchyEndpoint         = "/intserv/4.0/getGroupHierarchy"
	GroupInfoEndpoint              = "/intserv/4.0/getGroupsInfo"
	GroupsToAccountsEndpoint       = "/intserv/4.0/getAccountGroupRelationships"
	AddAccountToGroupEndpoint      = "/intserv/4.0/setAccountGroupRelationships"
	RemoveAccountFromGroupEndpoint = "/intserv/4.0/removeAccountGroupRelationship"

	PingEndpoint = "/intserv/4.0/ping"
)

type Config struct {
	Hostname    string `mapstructure:"hostname"`
	APILogin    string `mapstructure:"api-login"`
	APITransKey string `mapstructure:"api-trans-key"`
	ProviderID  string `mapstructure:"provider-id"`
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

func (c *Client) Ping(ctx context.Context) error {
	data := &FormData{
		APILogin:    c.config.APILogin,
		APITransKey: c.config.APITransKey,
		ProviderID:  c.config.ProviderID,
	}

	err := c.post(ctx, PingEndpoint, prepareForm(data), nil)
	if err != nil {
		return err
	}

	return nil
}

// https://docs.galileo-ft.com/pro/reference/post_getrelatedaccounts
func (c *Client) ListRelatedAccounts(ctx context.Context, accountID string) ([]Account, error) {
	var accounts BaseResponse[RelatedAccountsResponse]

	data := &FormData{
		APILogin:    c.config.APILogin,
		APITransKey: c.config.APITransKey,
		ProviderID:  c.config.ProviderID,
		AccountNo:   accountID,
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

// https://docs.galileo-ft.com/pro/reference/post_getaccountoverview
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

// https://docs.galileo-ft.com/pro/reference/post_getrootgroups
func (c *Client) ListRootGroups(ctx context.Context, pgVars *PaginationVars) ([]Group, uint, error) {
	var res ListResponse[Group]

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

// https://docs.galileo-ft.com/pro/reference/post_getgrouphierarchy
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

// https://docs.galileo-ft.com/pro/reference/post_getgroupsinfo
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

// https://docs.galileo-ft.com/pro/reference/post_getaccountgrouprelationships
func (c *Client) ListGroupMembers(ctx context.Context, groupID string) (*GroupToAccounts, error) {
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

	if len(res.Data) > 1 {
		return nil, fmt.Errorf("unexpected number of group to accounts responses: %d", len(res.Data))
	}

	return &res.Data[0], nil
}

// https://docs.galileo-ft.com/pro/reference/post_setaccountgrouprelationships
func (c *Client) AddAccountToGroup(ctx context.Context, groupID, accountID string) error {
	data := &FormData{
		APILogin:    c.config.APILogin,
		APITransKey: c.config.APITransKey,
		ProviderID:  c.config.ProviderID,
		GroupID:     groupID,
		AccountIDs:  []string{accountID},
	}

	err := c.post(ctx, AddAccountToGroupEndpoint, prepareForm(data), nil)
	if err != nil {
		return err
	}

	return nil
}

// https://docs.galileo-ft.com/pro/reference/post_removeaccountgrouprelationship
func (c *Client) RemoveAccountFromGroup(ctx context.Context, groupID, accountID string) error {
	data := &FormData{
		APILogin:    c.config.APILogin,
		APITransKey: c.config.APITransKey,
		ProviderID:  c.config.ProviderID,
		AccountIDs:  []string{accountID},
	}

	err := c.post(ctx, RemoveAccountFromGroupEndpoint, prepareForm(data), nil)
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

	resp, err := c.httpClient.Do(req, uhttp.WithJSONResponse(response), WithErrorResponse())
	if err != nil {
		return fmt.Errorf("failed to do request: %w", err)
	}

	defer resp.Body.Close()

	return nil
}

func checkContentType(contentType string) error {
	if contentType != "application/json" {
		return fmt.Errorf("unexpected content type %s", contentType)
	}

	return nil
}

// More information about error status codes can be found here: https://docs.galileo-ft.com/pro/reference/api-reference-global-response-statuses
type ErrorResponse struct {
	Code   uint   `json:"status_code"`
	Status string `json:"status"`
}

func WithErrorResponse() uhttp.DoOption {
	return func(resp *uhttp.WrapperResponse) error {
		if err := checkContentType(resp.Header.Get("Content-Type")); err != nil {
			return fmt.Errorf("%w - %v", err, string(resp.Body))
		}

		var response ErrorResponse
		if err := json.Unmarshal(resp.Body, &response); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}

		return fmt.Errorf("%s (%d)", response.Status, response.Code)
	}
}
