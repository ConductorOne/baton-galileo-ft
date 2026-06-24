package main

import "fmt"

func seed(s *State) {
	// Seed 53 root groups. The connector requests 50 per page (ResourcesPageSize=50),
	// so this forces 2 pages and exercises the full pagination loop.
	for i := 1; i <= 53; i++ {
		id := fmt.Sprintf("group-%02d", i)
		g := &Group{
			ID:           id,
			Name:         fmt.Sprintf("Group %02d", i),
			LegalName:    fmt.Sprintf("Group %02d LLC", i),
			ContactEmail: fmt.Sprintf("contact-%02d@example.com", i),
			ContactName:  fmt.Sprintf("Contact %02d", i),
			Level:        1,
		}
		s.groups[id] = g
		s.rootGroups = append(s.rootGroups, g)
	}

	// group-01 has two children (tests ListChildrenGroups + GetGroupsInfo code path).
	childrenOf01 := []*Group{
		{ID: "group-01-child-a", ParentGroupID: "group-01", Name: "Group 01 Child A", Level: 2},
		{ID: "group-01-child-b", ParentGroupID: "group-01", Name: "Group 01 Child B", Level: 2},
	}
	for _, c := range childrenOf01 {
		s.groups[c.ID] = c
		s.groupChildren["group-01"] = append(s.groupChildren["group-01"], c.ID)
	}

	// group-02 has one child.
	child02 := &Group{ID: "group-02-child-a", ParentGroupID: "group-02", Name: "Group 02 Child A", Level: 2}
	s.groups[child02.ID] = child02
	s.groupChildren["group-02"] = append(s.groupChildren["group-02"], child02.ID)

	// Accounts — 5 primary accounts.
	primaryAccounts := []struct {
		id, accNum string
	}{
		{"acc-prn-001", "ACC001"},
		{"acc-prn-002", "ACC002"},
		{"acc-prn-003", "ACC003"},
		{"acc-prn-004", "ACC004"}, // no group membership — tests empty-grants path
		{"acc-prn-005", "ACC005"},
	}
	for _, a := range primaryAccounts {
		s.accounts[a.id] = &Account{
			ID:        a.id,
			AccNumber: a.accNum,
			Active:    "1",
			Status:    "N",
			ProdID:    "prod-001",
		}
	}

	// Customer profiles for primary accounts.
	profiles := []struct {
		id, first, last, email string
	}{
		{"acc-prn-001", "Alice", "Adams", "alice@example.com"},
		{"acc-prn-002", "Bob", "Baker", "bob@example.com"},
		{"acc-prn-003", "Carol", "Clark", "carol@example.com"},
		{"acc-prn-004", "Dave", "Davis", "dave@example.com"},
		{"acc-prn-005", "Eve", "Evans", "eve@example.com"},
	}
	for _, p := range profiles {
		s.customers[p.id] = &Customer{
			FirstName:   p.first,
			LastName:    p.last,
			Email:       p.email,
			City:        "Testville",
			State:       "CA",
			PostalCode:  "90000",
			CountryCode: "US",
		}
	}

	// acc-prn-001 has one related child account (tests ListRelatedAccounts code path).
	relatedAcc := &Account{
		ID:        "acc-prn-001-child",
		AccNumber: "ACC001C",
		Active:    "1",
		Status:    "N",
		ProdID:    "prod-001",
	}
	s.accounts[relatedAcc.ID] = relatedAcc
	s.customers[relatedAcc.ID] = &Customer{
		FirstName:   "Alice",
		LastName:    "Adams Jr",
		Email:       "alice-jr@example.com",
		City:        "Testville",
		State:       "CA",
		PostalCode:  "90000",
		CountryCode: "US",
	}
	s.accountRelated["acc-prn-001"] = []string{"acc-prn-001-child"}

	// Group memberships:
	//   acc-prn-001, acc-prn-002 → group-01
	//   acc-prn-002, acc-prn-003 → group-02  (acc-prn-002 in two groups, tests overlapping membership)
	//   acc-prn-005             → group-03
	//   acc-prn-004             → (none, tests empty-grants path)
	s.groupMembers["group-01"] = []string{"acc-prn-001", "acc-prn-002"}
	s.groupMembers["group-02"] = []string{"acc-prn-002", "acc-prn-003"}
	s.groupMembers["group-03"] = []string{"acc-prn-005"}
}
