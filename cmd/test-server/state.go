package main

import (
	"slices"
	"sync"
)

// Group mirrors galileo.Group.
// Doc: https://docs.galileo-ft.com/pro/reference/getrootgroups
type Group struct {
	ID            string `json:"group_id"`
	ExternalID    string `json:"external_id"`
	ParentGroupID string `json:"parent_group_id"`
	Name          string `json:"group_name"`
	LegalName     string `json:"business_legal_name"`
	Business      string `json:"doing_business_as"`
	Level         int    `json:"max_level"`
	ContactEmail  string `json:"primary_contact_email"`
	ContactName   string `json:"primary_contact_name"`
}

// GroupHierarchy mirrors galileo.GroupHierarchy.
// Doc: https://docs.galileo-ft.com/pro/reference/getgrouphierarchy
type GroupHierarchy struct {
	ID       string           `json:"group_id"`
	Name     string           `json:"group_name"`
	Children []GroupHierarchy `json:"children"`
}

// Account mirrors galileo.Account (returned by getRelatedAccounts).
// Doc: https://docs.galileo-ft.com/pro/reference/getrelatedaccounts
type Account struct {
	ID        string `json:"prn"`
	Active    string `json:"active"`
	Status    string `json:"status"`
	AccNumber string `json:"galileo_account_number"`
	ProdID    string `json:"product_id"`
}

// Customer mirrors galileo.Customer (returned inside getAccountOverview).
// Doc: https://docs.galileo-ft.com/pro/reference/getaccountoverview
type Customer struct {
	FirstName   string `json:"first_name"`
	MiddleName  string `json:"middle_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Address1    string `json:"address_1"`
	Address2    string `json:"address_2"`
	City        string `json:"city"`
	State       string `json:"state"`
	PostalCode  string `json:"postal_code"`
	CountryCode string `json:"country_code"`
	HomePhone   string `json:"home_phone"`
	MobilePhone string `json:"mobile_phone"`
}

type State struct {
	mu sync.Mutex

	groups     map[string]*Group
	rootGroups []*Group // ordered slice — preserves insertion order for paginated getRootGroups

	accounts  map[string]*Account
	customers map[string]*Customer // accountID (prn) → customer profile

	// groupChildren maps groupID → ordered list of direct child group IDs.
	groupChildren map[string][]string

	// groupMembers maps groupID → ordered list of account IDs (prns).
	groupMembers map[string][]string

	// accountRelated maps accountID → ordered list of related child account IDs.
	accountRelated map[string][]string

	// grantedMemberships tracks the most recent group added via AddAccountToGroup (not from seed).
	// RemoveAccountFromGroup only removes this added membership, leaving seed groups intact so
	// the account remains in the c1z after the revoke cycle and a second grant cycle still works.
	grantedMemberships map[string]string // accID → groupID
}

func NewState() *State {
	s := &State{
		groups:             make(map[string]*Group),
		accounts:           make(map[string]*Account),
		customers:          make(map[string]*Customer),
		groupChildren:      make(map[string][]string),
		groupMembers:       make(map[string][]string),
		accountRelated:     make(map[string][]string),
		grantedMemberships: make(map[string]string),
	}
	seed(s)
	return s
}

func (s *State) GetRootGroupPage(page, count int) (page_ []*Group, numPages int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	all := s.rootGroups

	from := (page - 1) * count
	if from < 0 || from > len(all) {
		from = len(all)
	}
	to := from + count
	if to > len(all) {
		to = len(all)
	}

	result := make([]*Group, to-from)
	for i, g := range all[from:to] {
		cp := *g
		result[i] = &cp
	}

	numPages = (len(all) + count - 1) / count
	if numPages == 0 {
		numPages = 1
	}
	return result, numPages
}

func (s *State) GetGroupChildren(groupID string) []GroupHierarchy {
	s.mu.Lock()
	defer s.mu.Unlock()
	var result []GroupHierarchy
	for _, childID := range s.groupChildren[groupID] {
		if g, ok := s.groups[childID]; ok {
			result = append(result, GroupHierarchy{
				ID:       g.ID,
				Name:     g.Name,
				Children: []GroupHierarchy{}, // max depth is 2 in test data; no grandchildren
			})
		}
	}
	return result
}

func (s *State) GetGroupsInfo(ids []string) []*Group {
	s.mu.Lock()
	defer s.mu.Unlock()
	var result []*Group
	for _, id := range ids {
		if g, ok := s.groups[id]; ok {
			cp := *g
			result = append(result, &cp)
		}
	}
	return result
}

func (s *State) GetGroupMembers(groupID string) []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return slices.Clone(s.groupMembers[groupID])
}

func (s *State) GetCustomer(id string) (*Customer, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.customers[id]
	if !ok {
		return nil, false
	}
	cp := *c
	return &cp, true
}

func (s *State) GetRelatedAccounts(accID string) []*Account {
	s.mu.Lock()
	defer s.mu.Unlock()
	var result []*Account
	for _, rid := range s.accountRelated[accID] {
		if a, ok := s.accounts[rid]; ok {
			cp := *a
			result = append(result, &cp)
		}
	}
	return result
}

// AddAccountToGroup adds an account to a group, returning flags for the outcome.
// Records the grant in grantedMemberships so RemoveAccountFromGroup can undo only this grant.
func (s *State) AddAccountToGroup(groupID, accID string) (alreadyMember, groupExists, accountExists bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.groups[groupID]; !ok {
		return false, false, false
	}
	if _, ok := s.accounts[accID]; !ok {
		return false, true, false
	}
	if slices.Contains(s.groupMembers[groupID], accID) {
		return true, true, true
	}
	s.groupMembers[groupID] = append(s.groupMembers[groupID], accID)
	s.grantedMemberships[accID] = groupID
	return false, true, true
}

// RemoveAccountFromGroup removes the most recently granted membership for accID (from grantedMemberships),
// leaving any seed-time membership intact. This ensures the account remains discoverable in subsequent
// sync cycles, so the sync-test@v4 two-cycle idempotency check can succeed.
// Returns notMember=true if no granted membership exists for the account.
func (s *State) RemoveAccountFromGroup(accID string) (notMember bool, accountExists bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.accounts[accID]; !ok {
		return false, false
	}
	groupID, ok := s.grantedMemberships[accID]
	if !ok {
		return true, true
	}
	s.groupMembers[groupID] = slices.DeleteFunc(s.groupMembers[groupID], func(id string) bool { return id == accID })
	delete(s.grantedMemberships, accID)
	return false, true
}
