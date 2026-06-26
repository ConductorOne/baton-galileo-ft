package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

const (
	testAPILogin    = "test-login"
	testAPITransKey = "test-trans-key"
	testProviderID  = "test-provider-id"
)

// parseAndValidate parses the form body and validates Galileo-FT credentials.
// Returns the parsed form and true on success; writes an error response and returns false otherwise.
func (s *Server) parseAndValidate(w http.ResponseWriter, r *http.Request) (bool, *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	if err := r.ParseForm(); err != nil {
		writeGalileoError(w, http.StatusBadRequest, 400, "malformed request body")
		return false, r
	}

	login := r.PostForm.Get("apiLogin")
	transKey := r.PostForm.Get("apiTransKey")
	providerID := r.PostForm.Get("providerId")

	if login == "" || transKey == "" || providerID == "" {
		writeGalileoError(w, http.StatusBadRequest, 400, "missing required credentials: apiLogin, apiTransKey, providerId")
		return false, r
	}

	if login != testAPILogin || transKey != testAPITransKey || providerID != testProviderID {
		writeGalileoError(w, http.StatusUnauthorized, 401, "invalid credentials")
		return false, r
	}

	return true, r
}

// writeGalileoError writes the Galileo-FT error envelope.
// Mirror of galileo.ErrorResponse: {"status_code": N, "status": "..."}.
func writeGalileoError(w http.ResponseWriter, httpStatus, code int, status string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"status_code": code,
		"status":      status,
	})
}

// writeEmpty writes HTTP 200 with Content-Type: application/json and no body.
// Used for write endpoints (add/remove) where the connector passes nil as the
// response target and uhttp.WithJSONResponse(nil) only accepts empty bodies.
func writeEmpty(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// writeSuccess writes {"response_data": data}.
// Mirror of galileo.BaseResponse.
func writeSuccess(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"response_data": data,
	})
}

// writeListSuccess writes the paginated list response.
// Mirror of galileo.ListResponse.
func writeListSuccess(w http.ResponseWriter, data any, page, numPages int) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"response_data":   data,
		"page":            page,
		"number_of_pages": numPages,
	})
}

// handlePing validates credentials and returns success.
// The connector calls Ping with a nil response target (uhttp.WithJSONResponse(nil)).
// uhttp only returns nil for nil targets when the body is also empty, so we write
// an empty body here. Any non-empty body causes an InvalidUnmarshalError.
// Doc URL: https://docs.galileo-ft.com/pro/reference/ping
func (s *Server) handlePing(w http.ResponseWriter, r *http.Request) {
	ok, _ := s.parseAndValidate(w, r)
	if !ok {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// handleGetRootGroups returns a page of root groups.
// Pagination: page (1-indexed) + recordCnt.
// Doc URL: https://docs.galileo-ft.com/pro/reference/getrootgroups
func (s *Server) handleGetRootGroups(w http.ResponseWriter, r *http.Request) {
	ok, r := s.parseAndValidate(w, r)
	if !ok {
		return
	}

	page := atoiOr(r.PostForm.Get("page"), 1)
	count := atoiOr(r.PostForm.Get("recordCnt"), 50)

	if page < 1 {
		page = 1
	}
	if count < 1 {
		count = 50
	}

	groups, numPages := s.state.GetRootGroupPage(page, count)
	writeListSuccess(w, groups, page, numPages)
}

// handleGetGroupHierarchy returns the child hierarchy for a group.
// Doc URL: https://docs.galileo-ft.com/pro/reference/getgrouphierarchy
func (s *Server) handleGetGroupHierarchy(w http.ResponseWriter, r *http.Request) {
	ok, r := s.parseAndValidate(w, r)
	if !ok {
		return
	}

	groupID := r.PostForm.Get("groupId")
	if groupID == "" {
		writeGalileoError(w, http.StatusBadRequest, 400, "groupId is required")
		return
	}

	children := s.state.GetGroupChildren(groupID)
	if children == nil {
		children = []GroupHierarchy{}
	}
	writeSuccess(w, children)
}

// handleGetGroupsInfo returns full group metadata for a list of group IDs.
// Doc URL: https://docs.galileo-ft.com/pro/reference/getgroupsinfo
func (s *Server) handleGetGroupsInfo(w http.ResponseWriter, r *http.Request) {
	ok, r := s.parseAndValidate(w, r)
	if !ok {
		return
	}

	// The connector sends groupIds as multi-value form fields.
	groupIDs := r.PostForm["groupIds"]
	if len(groupIDs) == 0 {
		writeGalileoError(w, http.StatusBadRequest, 400, "at least one groupIds is required")
		return
	}

	groups := s.state.GetGroupsInfo(groupIDs)
	if groups == nil {
		groups = []*Group{}
	}
	writeSuccess(w, groups)
}

// handleGetAccountGroupRelationships returns the accounts that belong to a group.
// Always returns a 1-element array even for empty groups — the connector indexes [0] directly.
// Doc URL: https://docs.galileo-ft.com/pro/reference/getaccountgrouprelationships
func (s *Server) handleGetAccountGroupRelationships(w http.ResponseWriter, r *http.Request) {
	ok, r := s.parseAndValidate(w, r)
	if !ok {
		return
	}

	groupID := r.PostForm.Get("groupId")
	if groupID == "" {
		writeGalileoError(w, http.StatusBadRequest, 400, "groupId is required")
		return
	}

	members := s.state.GetGroupMembers(groupID)
	if members == nil {
		members = []string{}
	}

	// Response must be a 1-element array — the connector calls res.Data[0] unconditionally.
	writeSuccess(w, []map[string]any{
		{
			"group_id":   groupID,
			"pmt_ref_no": members,
		},
	})
}

// handleGetRelatedAccounts returns child accounts for a given account.
// Doc URL: https://docs.galileo-ft.com/pro/reference/getrelatedaccounts
func (s *Server) handleGetRelatedAccounts(w http.ResponseWriter, r *http.Request) {
	ok, r := s.parseAndValidate(w, r)
	if !ok {
		return
	}

	accountNo := r.PostForm.Get("accountNo")
	if accountNo == "" {
		writeGalileoError(w, http.StatusBadRequest, 400, "accountNo is required")
		return
	}

	related := s.state.GetRelatedAccounts(accountNo)
	if related == nil {
		related = []*Account{}
	}
	writeSuccess(w, map[string]any{
		"child_accounts": related,
	})
}

// handleGetAccountOverview returns the customer profile for an account.
// Doc URL: https://docs.galileo-ft.com/pro/reference/getaccountoverview
func (s *Server) handleGetAccountOverview(w http.ResponseWriter, r *http.Request) {
	ok, r := s.parseAndValidate(w, r)
	if !ok {
		return
	}

	accountNo := r.PostForm.Get("accountNo")
	if accountNo == "" {
		writeGalileoError(w, http.StatusBadRequest, 400, "accountNo is required")
		return
	}

	customer, found := s.state.GetCustomer(accountNo)
	if !found {
		writeGalileoError(w, http.StatusNotFound, 404, "account not found")
		return
	}

	writeSuccess(w, map[string]any{
		"profile": customer,
	})
}

// handleSetAccountGroupRelationships adds an account to a group (provisioning Grant).
// Doc URL: https://docs.galileo-ft.com/pro/reference/setaccountgrouprelationships
func (s *Server) handleSetAccountGroupRelationships(w http.ResponseWriter, r *http.Request) {
	ok, r := s.parseAndValidate(w, r)
	if !ok {
		return
	}

	groupID := r.PostForm.Get("groupId")
	if groupID == "" {
		writeGalileoError(w, http.StatusBadRequest, 400, "groupId is required")
		return
	}

	// The connector sends accountNos as a multi-value field.
	accountNos := r.PostForm["accountNos"]
	if len(accountNos) == 0 {
		writeGalileoError(w, http.StatusBadRequest, 400, "at least one accountNos is required")
		return
	}

	for _, accID := range accountNos {
		alreadyMember, groupExists, accountExists := s.state.AddAccountToGroup(groupID, accID)
		if !groupExists {
			writeGalileoError(w, http.StatusNotFound, 404, "group not found: "+groupID)
			return
		}
		if !accountExists {
			writeGalileoError(w, http.StatusNotFound, 404, "account not found: "+accID)
			return
		}
		if alreadyMember {
			writeEmpty(w)
			return
		}
	}

	writeEmpty(w)
}

// handleRemoveAccountGroupRelationship removes an account from its current group (provisioning Revoke).
// Note: the connector does NOT send groupId — it only sends accountNos; the account is removed from
// whatever group it currently belongs to (Galileo accounts can only be in one group at a time).
// Doc URL: https://docs.galileo-ft.com/pro/reference/removeaccountgrouprelationship
func (s *Server) handleRemoveAccountGroupRelationship(w http.ResponseWriter, r *http.Request) {
	ok, r := s.parseAndValidate(w, r)
	if !ok {
		return
	}

	accountNos := r.PostForm["accountNos"]
	if len(accountNos) == 0 {
		writeGalileoError(w, http.StatusBadRequest, 400, "at least one accountNos is required")
		return
	}

	for _, accID := range accountNos {
		notMember, accountExists := s.state.RemoveAccountFromGroup(accID)
		if !accountExists {
			writeGalileoError(w, http.StatusNotFound, 404, "account not found: "+accID)
			return
		}
		if notMember {
			writeEmpty(w)
			return
		}
	}

	writeEmpty(w)
}

func atoiOr(s string, def int) int {
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}
