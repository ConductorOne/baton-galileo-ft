package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

const port = ":8765"

type Server struct {
	state *State
}

func main() {
	state := NewState()
	server := &Server{state: state}
	mux := http.NewServeMux()

	// All Galileo-FT endpoints are POST with form-encoded bodies.
	// Credentials (apiLogin, apiTransKey, providerId) are validated per-request — there is no token endpoint.
	// Doc: https://docs.galileo-ft.com/pro/reference/

	// Doc URL: https://docs.galileo-ft.com/pro/reference/ping
	mux.HandleFunc("POST /intserv/4.0/ping", server.handlePing)

	// Doc URL: https://docs.galileo-ft.com/pro/reference/getrootgroups
	mux.HandleFunc("POST /intserv/4.0/getRootGroups", server.handleGetRootGroups)

	// Doc URL: https://docs.galileo-ft.com/pro/reference/getgrouphierarchy
	mux.HandleFunc("POST /intserv/4.0/getGroupHierarchy", server.handleGetGroupHierarchy)

	// Doc URL: https://docs.galileo-ft.com/pro/reference/getgroupsinfo
	mux.HandleFunc("POST /intserv/4.0/getGroupsInfo", server.handleGetGroupsInfo)

	// Doc URL: https://docs.galileo-ft.com/pro/reference/getaccountgrouprelationships
	mux.HandleFunc("POST /intserv/4.0/getAccountGroupRelationships", server.handleGetAccountGroupRelationships)

	// Doc URL: https://docs.galileo-ft.com/pro/reference/getrelatedaccounts
	mux.HandleFunc("POST /intserv/4.0/getRelatedAccounts", server.handleGetRelatedAccounts)

	// Doc URL: https://docs.galileo-ft.com/pro/reference/getaccountoverview
	mux.HandleFunc("POST /intserv/4.0/getAccountOverview", server.handleGetAccountOverview)

	// Doc URL: https://docs.galileo-ft.com/pro/reference/setaccountgrouprelationships
	mux.HandleFunc("POST /intserv/4.0/setAccountGroupRelationships", server.handleSetAccountGroupRelationships)

	// Doc URL: https://docs.galileo-ft.com/pro/reference/removeaccountgrouprelationship
	mux.HandleFunc("POST /intserv/4.0/removeAccountGroupRelationship", server.handleRemoveAccountGroupRelationship)

	log.Printf("Galileo-FT test server listening on %s\n", port)
	srv := &http.Server{
		Addr:         port,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Printf("server error: %v", err)
		os.Exit(1)
	}
}
