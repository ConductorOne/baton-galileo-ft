![Baton Logo](./docs/images/baton-logo.png)
# baton-galileo-ft
Welcome to your new connector!

# `baton-galileo-ft` [![Go Reference](https://pkg.go.dev/badge/github.com/conductorone/baton-galileo-ft.svg)](https://pkg.go.dev/github.com/conductorone/baton-galileo-ft) ![ci](https://github.com/conductorone/baton-galileo-ft/actions/workflows/ci.yaml/badge.svg)

`baton-galileo-ft` is an example connector built using the [Baton SDK](https://github.com/conductorone/baton-sdk). It uses hardcoded data to provide a simple example of how to build your own connector with Baton.

Check out [Baton](https://github.com/conductorone/baton) to learn more about the project in general.

# Getting Started
To start out, you will want to update the dependencies.
Do this by running `make update-deps`.

## brew
```
brew install conductor/baton/baton conductor/baton/baton-galileo-ft

baton-galileo-ft
baton resources
```

## docker
```
docker run --rm -v $(pwd):/out -e BATON_API_LOGIN=api_login -e BATON_API_TRANS_KEY=api_trans_key -e BATON_PROVIDER_ID=provider_id ghcr.io/conductorone/baton-galileo-ft:latest -f "/out/sync.c1z"
docker run --rm -v $(pwd):/out ghcr.io/conductorone/baton:latest -f "/out/sync.c1z" resources
```

## source
```
go install github.com/conductorone/baton/cmd/baton@main
go install github.com/conductorone/baton-galileo-ft/cmd/baton-galileo-ft@main

BATON_API_LOGIN=api_login BATON_API_TRANS_KEY=api_trans_key BATON_PROVIDER_ID=provider_id baton-galileo-ft
baton resources
```

# Contributing, Support and Issues

We started Baton because we were tired of taking screenshots and manually building spreadsheets.  We welcome contributions, and ideas, no matter how small -- our goal is to make identity and permissions sprawl less painful for everyone.  If you have questions, problems, or ideas: Please open a Github Issue!

See [CONTRIBUTING.md](https://github.com/ConductorOne/baton/blob/main/CONTRIBUTING.md) for more details.

# `baton-galileo-ft` Command Line Usage
```
baton-galileo-ft

Usage:
  baton-galileo-ft [flags]
  baton-galileo-ft [command]

Available Commands:
  capabilities       Get connector capabilities
  completion         Generate the autocompletion script for the specified shell
  config             Get the connector config schema
  health-check       Check the health of a running connector
  help               Help about any command

Flags:
      --api-login string                                 required: The username provided by Galileo-FT for API access. ($BATON_API_LOGIN)
      --api-trans-key string                             required: The password provided by Galileo-FT, used alongside the api-login. ($BATON_API_TRANS_KEY)
      --auth-method string                               ($BATON_AUTH_METHOD)
      --client-id string                                 The client ID used to authenticate with ConductorOne ($BATON_CLIENT_ID)
      --client-secret string                             The client secret used to authenticate with ConductorOne ($BATON_CLIENT_SECRET)
      --external-resource-c1z string                     The path to the c1z file to sync external baton resources with ($BATON_EXTERNAL_RESOURCE_C1Z)
      --external-resource-entitlement-id-filter string   The entitlement that external users, groups must have access to sync external baton resources ($BATON_EXTERNAL_RESOURCE_ENTITLEMENT_ID_FILTER)
      --external-resource-traits strings                 Resource type traits (e.g. "user", "group", "app") to sync and match from the external resource c1z. When unset the matcher falls back to user and group; passing this flag replaces the full set rather than adding to it. ($BATON_EXTERNAL_RESOURCE_TRAITS)
  -f, --file string                                      The path to the c1z file to sync with ($BATON_FILE) (default "sync.c1z")
      --health-check                                     Enable the HTTP health check endpoint ($BATON_HEALTH_CHECK)
      --health-check-port int                            Port for the HTTP health check endpoint ($BATON_HEALTH_CHECK_PORT) (default 8081)
  -h, --help                                             help for baton-galileo-ft
      --hostname string                                  URL hostname for production hostname. ($BATON_HOSTNAME)
      --http-timeout-seconds int                         HTTP client timeout in seconds (max 1800) ($BATON_HTTP_TIMEOUT_SECONDS) (default 300)
      --keep-previous-sync-c1z                           Keep the previously synced c1z on disk to enable ETag replay across service-mode syncs (requires a connector that supports ETag replay; costs one c1z of local disk) ($BATON_KEEP_PREVIOUS_SYNC_C1Z)
      --log-format string                                The output format for logs: json, console ($BATON_LOG_FORMAT) (default "json")
      --log-level string                                 The log level: debug, info, warn, error ($BATON_LOG_LEVEL) (default "info")
      --log-level-debug-expires-at string                The timestamp indicating when debug-level logging should expire ($BATON_LOG_LEVEL_DEBUG_EXPIRES_AT)
      --log-path strings                                 The file path to write logs to ($BATON_LOG_PATH)
      --otel-collector-endpoint string                   The endpoint of the OpenTelemetry collector to send observability data to (used for both tracing and logging if specific endpoints are not provided) ($BATON_OTEL_COLLECTOR_ENDPOINT)
      --parallel-sync                                    Deprecated: use --workers instead. ($BATON_PARALLEL_SYNC)
      --provider-id string                               required: A unique identifier from Galileo-FT representing your organization, used for tracking transactions and data. ($BATON_PROVIDER_ID)
  -p, --provisioning                                     This must be set in order for provisioning actions to be enabled ($BATON_PROVISIONING)
      --skip-entitlements-and-grants                     This must be set to skip syncing of entitlements and grants ($BATON_SKIP_ENTITLEMENTS_AND_GRANTS)
      --skip-full-sync                                   This must be set to skip a full sync ($BATON_SKIP_FULL_SYNC)
      --storage-engine string                            The storage engine to use when opening the sync c1z file: sqlite or pebble. Leave unset to use the baton-sdk default. ($BATON_STORAGE_ENGINE)
      --sync-resource-types strings                      The resource type IDs to sync ($BATON_SYNC_RESOURCE_TYPES)
      --sync-resources strings                           The resource IDs to sync ($BATON_SYNC_RESOURCES)
      --task-concurrency int                             The number of Baton tasks to run concurrently in service mode. Tasks may include sync, grant, revoke, and more. Minimum value is 1, maximum value is 100. ($BATON_TASK_CONCURRENCY) (default 3)
      --ticketing                                        This must be set to enable ticketing support ($BATON_TICKETING)
  -v, --version                                          version for baton-galileo-ft
      --workers int                                      The number of sync workers to use. -1 for auto-detect, 0 for sequential, >0 for parallel ($BATON_WORKERS)

Use "baton-galileo-ft [command] --help" for more information about a command.
```
