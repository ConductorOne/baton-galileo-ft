package galileo

import "encoding/json"

type BaseResponse[T any] struct {
	Data T `json:"response_data"`
}

// StatusEnvelope captures the status_code/status fields Galileo-FT includes on every response
// (https://docs.galileo-ft.com/pro/reference/api-reference-global-response-statuses). Used for
// calls whose response_data carries no useful payload, so a non-zero status_code on an HTTP 2xx
// is still caught rather than silently treated as success.
type StatusEnvelope struct {
	StatusCode json.Number     `json:"status_code"`
	Status     string          `json:"status"`
	Data       json.RawMessage `json:"response_data"`
}

type ListResponse[T any] struct {
	Data       []T  `json:"response_data"`
	Page       uint `json:"page"`
	NumOfPages uint `json:"number_of_pages"`
}

type RelatedAccountsResponse struct {
	Children []Account `json:"child_accounts"`
}

type Account struct {
	ID        string `json:"prn"`
	Active    string `json:"active"`
	Status    string `json:"status"`
	AccNumber string `json:"galileo_account_number"`
	ProdID    string `json:"product_id"`
}

type Customer struct {
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`

	Address1    string `json:"address_1"`
	Address2    string `json:"address_2"`
	City        string `json:"city"`
	State       string `json:"state"`
	PostalCode  string `json:"postal_code"`
	CountryCode string `json:"country_code"`
	HomePhone   string `json:"home_phone"`
	MobilePhone string `json:"mobile_phone"`
}

type Group struct {
	ID            string `json:"group_id"`
	ExternalID    string `json:"external_id"`
	ParentGroupID string `json:"parent_group_id"`

	Name      string `json:"group_name"`
	LegalName string `json:"business_legal_name"`
	Business  string `json:"doing_business_as"`

	Level        int    `json:"max_level"`
	ContactEmail string `json:"primary_contact_email"`
	ContactName  string `json:"primary_contact_name"`
}

type GroupHierarchy struct {
	ID       string           `json:"group_id"`
	Name     string           `json:"group_name"`
	Children []GroupHierarchy `json:"children"`
}

type GroupToAccounts struct {
	GroupID    string   `json:"group_id"`
	AccountIDs []string `json:"pmt_ref_no"`
}
