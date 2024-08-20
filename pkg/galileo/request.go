package galileo

import (
	"bytes"
	"fmt"
	"io"
	"net/url"

	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/google/uuid"
)

type FormData struct {
	APILogin    string
	APITransKey string
	ProviderID  string
	AccountNo   string
	GroupID     string
	GroupIDs    []string
	AccountIDs  []string
}

type PaginationVars struct {
	Page  uint
	Count uint
}

func NewPaginationVars(page, count uint) *PaginationVars {
	return &PaginationVars{
		Page:  page,
		Count: count,
	}
}

func (pg *PaginationVars) PrepareVars(form *url.Values) {
	form.Set("page", fmt.Sprintf("%d", pg.Page))
	form.Set("recordCnt", fmt.Sprintf("%d", pg.Count))
}

func generateTransactionID() string {
	id := uuid.New()

	return id.String()
}

func prepareForm(data *FormData) *url.Values {
	form := &url.Values{}

	// set username and password
	form.Set("apiLogin", data.APILogin)
	form.Set("apiTransKey", data.APITransKey)

	// set providerId and transactionId
	form.Set("providerId", data.ProviderID)
	form.Set("transactionId", generateTransactionID())

	// set account, if provided
	if data.AccountNo != "" {
		form.Set("accountNo", data.AccountNo)
	}

	// set group id, if provided
	if data.GroupID != "" {
		form.Set("groupId", data.GroupID)
	}

	// In Go, if `data.GroupIDs` is nil, this is a noop.
	for _, id := range data.GroupIDs {
		form.Add("groupIds", id)
	}

	// In Go, if `data.AccountIDs` is nil, this is a noop.
	for _, id := range data.AccountIDs {
		form.Add("accountNos", id)
	}

	return form
}

func WithResponseContentTypeJSON() uhttp.RequestOption {
	return func() (io.ReadWriter, map[string]string, error) {
		return nil, map[string]string{
			"Response-Content-Type": "json",
		}, nil
	}
}

func WithContentTypeFormHeader() uhttp.RequestOption {
	return func() (io.ReadWriter, map[string]string, error) {
		return nil, map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		}, nil
	}
}

func WithFormBody(values *url.Values) uhttp.RequestOption {
	return func() (io.ReadWriter, map[string]string, error) {
		b := bytes.NewBufferString(values.Encode())

		return b, nil, nil
	}
}
