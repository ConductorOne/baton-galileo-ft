# Test server

Mock API for baton-galileo-ft, used by CI when no real-tenant credentials are available.
Replicates the Galileo-FT API's auth flow, endpoints, and error envelopes.

## Auth

Galileo-FT embeds credentials in every POST form body — there is no token endpoint.

| Real API | Test server |
|---|---|
| `apiLogin` / `apiTransKey` / `providerId` in every POST body | Same fields; hardcoded values below |
| Credentials from the customer | `apiLogin=test-login`, `apiTransKey=test-trans-key`, `providerId=test-provider-id` |

## Endpoints

| Path | Method | Doc URL |
|---|---|---|
| `/intserv/4.0/ping` | POST | https://docs.galileo-ft.com/pro/reference/ping |
| `/intserv/4.0/getRootGroups` | POST | https://docs.galileo-ft.com/pro/reference/getrootgroups |
| `/intserv/4.0/getGroupHierarchy` | POST | https://docs.galileo-ft.com/pro/reference/getgrouphierarchy |
| `/intserv/4.0/getGroupsInfo` | POST | https://docs.galileo-ft.com/pro/reference/getgroupsinfo |
| `/intserv/4.0/getAccountGroupRelationships` | POST | https://docs.galileo-ft.com/pro/reference/getaccountgrouprelationships |
| `/intserv/4.0/getRelatedAccounts` | POST | https://docs.galileo-ft.com/pro/reference/getrelatedaccounts |
| `/intserv/4.0/getAccountOverview` | POST | https://docs.galileo-ft.com/pro/reference/getaccountoverview |
| `/intserv/4.0/setAccountGroupRelationships` | POST | https://docs.galileo-ft.com/pro/reference/setaccountgrouprelationships |
| `/intserv/4.0/removeAccountGroupRelationship` | POST | https://docs.galileo-ft.com/pro/reference/removeaccountgrouprelationship |

## Seed data

53 root groups (forces 2-page pagination since the connector requests 50 per page):

- `group-01` — 2 child groups (`group-01-child-a`, `group-01-child-b`); members: `acc-prn-001`, `acc-prn-002`
- `group-02` — 1 child group (`group-02-child-a`); members: `acc-prn-002`, `acc-prn-003` (overlapping with group-01)
- `group-03` — no children; member: `acc-prn-005`
- `group-04` through `group-53` — no children, no members

5 primary accounts:

| Account (prn) | Name | Group |
|---|---|---|
| `acc-prn-001` | Alice Adams | group-01 |
| `acc-prn-002` | Bob Baker | group-01, group-02 (overlapping) |
| `acc-prn-003` | Carol Clark | group-02 |
| `acc-prn-004` | Dave Davis | (none — tests empty-grants path) |
| `acc-prn-005` | Eve Evans | group-03 |

`acc-prn-001` has one related child account `acc-prn-001-child` (tests `ListRelatedAccounts`).

## Running locally

```bash
# Start the test server (from project root)
go run ./cmd/test-server/

# In a separate terminal, point the connector at it
./baton-galileo-ft \
  --base-url http://localhost:8765 \
  --api-login test-login \
  --api-trans-key test-trans-key \
  --provider-id test-provider-id
```

## Curl examples

```bash
# Ping (validates credentials)
curl -s -X POST http://localhost:8765/intserv/4.0/ping \
  -d 'apiLogin=test-login&apiTransKey=test-trans-key&providerId=test-provider-id&transactionId=abc'

# List root groups (page 1)
curl -s -X POST http://localhost:8765/intserv/4.0/getRootGroups \
  -d 'apiLogin=test-login&apiTransKey=test-trans-key&providerId=test-provider-id&transactionId=abc&page=1&recordCnt=50'

# Get group members
curl -s -X POST http://localhost:8765/intserv/4.0/getAccountGroupRelationships \
  -d 'apiLogin=test-login&apiTransKey=test-trans-key&providerId=test-provider-id&transactionId=abc&groupId=group-01'

# Add account to group (Grant)
curl -s -X POST http://localhost:8765/intserv/4.0/setAccountGroupRelationships \
  -d 'apiLogin=test-login&apiTransKey=test-trans-key&providerId=test-provider-id&transactionId=abc&groupId=group-03&accountNos=acc-prn-004'

# Remove account from group (Revoke)
curl -s -X POST http://localhost:8765/intserv/4.0/removeAccountGroupRelationship \
  -d 'apiLogin=test-login&apiTransKey=test-trans-key&providerId=test-provider-id&transactionId=abc&accountNos=acc-prn-004'
```
