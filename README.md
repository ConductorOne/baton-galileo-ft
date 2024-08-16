![Baton Logo](./docs/images/baton-logo.png)
# baton-galileo-ft
Welcome to your new connector! 

# `baton-galileo-ft` [![Go Reference](https://pkg.go.dev/badge/github.com/conductorone/baton-galileo-ft.svg)](https://pkg.go.dev/github.com/conductorone/baton-galileo-ft) ![main ci](https://github.com/conductorone/baton-galileo-ft/actions/workflows/main.yaml/badge.svg)

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
  help               Help about any command

Flags:
      --api-login string       required: The username provided by Galileo-FT for API access. ($BATON_API_LOGIN)
      --api-trans-key string   required: The password provided by Galileo-FT, used alongside the api-login. ($BATON_API_TRANS_KEY)
      --client-id string       The client ID used to authenticate with ConductorOne ($BATON_CLIENT_ID)
      --client-secret string   The client secret used to authenticate with ConductorOne ($BATON_CLIENT_SECRET)
  -f, --file string            The path to the c1z file to sync with ($BATON_FILE) (default "sync.c1z")
  -h, --help                   help for baton-galileo-ft
      --hostname string        URL hostname for production hostname. ($BATON_HOSTNAME)
      --log-format string      The output format for logs: json, console ($BATON_LOG_FORMAT) (default "json")
      --log-level string       The log level: debug, info, warn, error ($BATON_LOG_LEVEL) (default "info")
      --provider-id string     required: A unique identifier from Galileo-FT representing your organization, used for tracking transactions and data. ($BATON_PROVIDER_ID)
  -p, --provisioning           This must be set in order for provisioning actions to be enabled ($BATON_PROVISIONING)
      --skip-full-sync         This must be set to skip a full sync ($BATON_SKIP_FULL_SYNC)
      --ticketing              This must be set to enable ticketing support ($BATON_TICKETING)
  -v, --version                version for baton-galileo-ft

Use "baton-galileo-ft [command] --help" for more information about a command.
```
