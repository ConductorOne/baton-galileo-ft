package galileo

import (
	"encoding/json"
	"fmt"
)

// BaseResponse mirrors Galileo-FT's response envelope: every response — success or failure —
// carries status_code/status (https://docs.galileo-ft.com/pro/reference/api-reference-global-response-statuses),
// so a real tenant can answer HTTP 200 with a non-zero status_code. checkStatus lets post() catch
// that uniformly instead of relying solely on the HTTP status.
type BaseResponse[T any] struct {
	StatusCode json.Number `json:"status_code"`
	Status     string      `json:"status"`
	Data       T           `json:"response_data"`
}

func (r BaseResponse[T]) checkStatus() error {
	return checkGalileoStatus(r.StatusCode, r.Status)
}

type ListResponse[T any] struct {
	StatusCode json.Number `json:"status_code"`
	Status     string      `json:"status"`
	Data       []T         `json:"response_data"`
	Page       uint        `json:"page"`
	NumOfPages uint        `json:"number_of_pages"`
}

func (r ListResponse[T]) checkStatus() error {
	return checkGalileoStatus(r.StatusCode, r.Status)
}

// statusChecker is implemented by every response envelope type so post() can validate status_code
// generically, regardless of the concrete response_data shape the caller decoded into.
type statusChecker interface {
	checkStatus() error
}

// checkGalileoStatus treats a missing status_code as success (some mocks/older responses may omit
// it) and otherwise requires it to parse as the numeric zero value.
func checkGalileoStatus(code json.Number, status string) error {
	if code == "" {
		return nil
	}
	n, err := code.Int64()
	if err != nil || n != 0 {
		return fmt.Errorf("galileo-ft: request failed: %s (%s)", status, code)
	}
	return nil
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
