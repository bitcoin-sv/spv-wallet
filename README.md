# BUX: Server

> Get started using [BUX](https://getbux.io) in five minutes

[![Release](https://img.shields.io/github/release-pre/BuxOrg/bux-server.svg?logo=github&style=flat&v=3)](https://github.com/BuxOrg/bux-server/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/BuxOrg/bux-server/run-tests.yml?branch=master&v=3)](https://github.com/BuxOrg/bux-server/actions)
[![Report](https://goreportcard.com/badge/github.com/BuxOrg/bux-server?style=flat&v=3)](https://goreportcard.com/report/github.com/BuxOrg/bux-server)
[![codecov](https://codecov.io/gh/BuxOrg/bux-server/branch/master/graph/badge.svg?v=3)](https://codecov.io/gh/BuxOrg/bux-server)
[![Mergify Status](https://img.shields.io/endpoint.svg?url=https://api.mergify.com/v1/badges/BuxOrg/bux-server&style=flat&v=3)](https://mergify.io)
[![Go](https://img.shields.io/github/go-mod/go-version/BuxOrg/bux-server?v=3)](https://golang.org/)
<br>
[![Gitpod Ready-to-Code](https://img.shields.io/badge/Gitpod-ready--to--code-blue?logo=gitpod&v=3)](https://gitpod.io/#https://github.com/BuxOrg/bux-server)
[![standard-readme compliant](https://img.shields.io/badge/readme%20style-standard-brightgreen.svg?style=flat&v=3)](https://github.com/RichardLitt/standard-readme)
[![Makefile Included](https://img.shields.io/badge/Makefile-Supported%20-brightgreen?=flat&logo=probot&v=3)](Makefile)
[![Sponsor](https://img.shields.io/badge/sponsor-mrz1836-181717.svg?logo=github&style=flat&v=3)](https://github.com/sponsors/BuxOrg)
[![Donate](https://img.shields.io/badge/donate-bitcoin-ff9900.svg?logo=bitcoin&style=flat&v=3)](https://gobitcoinsv.com/#sponsor?utm_source=github&utm_medium=sponsor-link&utm_campaign=bux-server&utm_term=bux-server&utm_content=bux-server)

<br/>

## Table of Contents

-   [What is BUX?](#what-is-bux)
-   [Config Variables](#config-variables)
-   [Installation](#installation)
-   [Documentation](#documentation)
-   [Examples & Tests](#examples--tests)
-   [Benchmarks](#benchmarks)
-   [Code Standards](#code-standards)
-   [Usage](#usage)
-   [Contributing](#contributing)
-   [License](#license)

<br/>

## What is BUX?

[Read more about BUX](https://getbux.io)

<br/>

## Config Variables

Default config variables can be overridden by (in this order):
1. Flags
2. ENV variables
3. Config file

Available flags:

```bash
  -C, --config_file string                       custom config file path
  -h, --help                                     show help
  -v, --version                                  show version
  -d, --dump_config                              dump config to file, specified by config_file flag
      --debug                                    enable debug logging (default true)
      --debug_profiling                          enable debug profiling (default true)
      --disable_itc                              disable ITC - Incoming Transaction Checking (default true)
      --import_block_headers string              path or URL to blockheaders file
      --request_logging                          request logging from api routers (rest and graphql) (default true)
      --auth.admin_key string                    key that is used for administrative requests (default "xpub661MyMwAqRbcFaYeQLxmExXvTCjw9jjBRpifkoGggkAitXNNjva4TStLJuYjjEmU4AzXRPGwoECjXo3Rgqg8zQqW6UPVfkKtsrogGBw8xz7")
      --auth.require_signing                     require signing
      --auth.scheme string                       authentication scheme to use (default "xpub")
      --auth.signing_disabled                    NOTE: Only for development, turns off signing (default true)
      --beef.use_beef                            enables BEEF transaction format, requires Pulse settings (default true)
      --beef.pulse.url string                    pulse url for validating merkle roots (default "http://localhost:8000/api/v1/chain/merkleroot/verify")
      --beef.pulse.auth_token string             authentication token for pulse (default "asd")
      --cache.engine string                      cache engine: redis, freecache or empty (default "freecache")
      --cluster.coordinator string               redis or memory (default "redis")
      --cluster.prefix string                    prefix string to use for all cluster keys (default "bux_cluser_")
      --cluster.redis.url string                 Redis URL for cluster coordinator, if redis is chosen (default "localhost:6379")
      --cluster.redis.max_idle_timeout string    max idle timeout for redis for cluster, if redis is chosen (default "10s")
      --cluster.redis.use_tls                    should redis cluster coordinator use tls, if redis is chosen
      --db.datastore.auto_migrate                loads a blank database (default true)
      --db.datastore.debug                       show sql statements
      --db.datastore.engine string               mysql, sqlite, postgresql, mongodb, empty (default "sqlite")
      --db.datastore.table_prefix string         prefix for all tables in db (default "xapi")
      --db.mongodb.db_name string                database name for MongoDB (default "xapi")
      --db.mongodb.transactions                  has transactions
      --db.mongodb.uri string                    connection uri to MongoDB (default "mongodb://localhost:27017/xapi")
      --db.sql.driver string                     mysql, postgresql (default "postgresql")
      --db.sql.host string                       db host (default "localhost")
      --db.sql.user string                       db user (default "postgres")
      --db.sql.name string                       db name (default "xapi")
      --db.sql.password string                   db password
      --db.sql.port string                       db port (default "5432")
      --db.sql.replica                           true if it's a replica (Read-Only)
      --db.sql.skip_initialize_with_version      skip using MySQL in test mode (default true)
      --db.sql.time_zone string                  time zone for db (default "UTC")
      --db.sql.tx_timeout string                 timeout for transactions (default "10s")
      --db.sqlite.database_path string           db path for sqlite (default "./test-json.db")
      --db.sqlite.shared                         adds a shared param to the connection string (default true)
      --graphql.enabled                          enable graphql (default true)
      --graphql.playground_path string           playground path for graphql (default "/graphql")
      --graphql.server_path string               server path (default "/graphql")
      --monitor.auth_token string                token to connect to the server with
      --monitor.bux_agent_url string             the bux agent server url address (default "ws://localhost:8000/websocket")
      --monitor.debug                            enable debug
      --monitor.enabled                          enable monitor
      --monitor.false_positive_rate float        percentage of false positives to expect (default 0.01)
      --monitor.load_monitored_destinations      load monitored destinations
      --monitor.max_number_of_destinations int   number of destinations that the filter can hold (default 100000)
      --monitor.monitor_days int                 number of days in the past that an address should be monitored for (default 7)
      --monitor.processor_type string            type of processor to start monitor with (default "bloom")
      --monitor.save_transaction_destinations    save destinations on monitored transactions (default true)
      --new_relic.domain_name string             used for hostname display (default "domain.com")
      --new_relic.enabled                        enable NewRelic
      --new_relic.license_key string             license key (default "BOGUS-LICENSE-KEY-1234567890987654321234")
      --nodes.use_mapi_fee_quotes                use mAPI fee quotes (default true)
      --nodes.minercraft_api string              type of api to use by minercraft, arc of mapi (default "mAPI")
      --nodes.broadcast_client_apis strings      go-broadcastClient api keys in fromat 'api_url|token' (default [url|token])
      --notifications.enabled                    enable notifications
      --notifications.webhook_endpoint string    webhook endpoint for notifications
      --paymail.default_from_paymail string      default 'from:@domain.com' paymail (default "from@domain.com")
      --paymail.default_note string              default paymail note, IE: message needed for address resolution (default "bux Address Resolution")
      --paymail.domains strings                  list of allowed paymail domains (default [localhost])
      --paymail.domain_validation_enabled        enable paymail domain validation, turn off if hosted domain is not paymail related
      --paymail.enabled                          enable paymail (default true)
      --paymail.sender_validation_enabled        enable paymail sender validation - extra security (default true)
      --redis.dependency_mode                    only in Redis with script enabled (default true)
      --redis.max_active_connections int         max active redis connections
      --redis.max_connection_lifetime string     max redis connection lifetime (default "60s")
      --redis.max_idle_connections int           max idle redis connections (default 10)
      --redis.max_idle_timeout string            max idle redis timeout (default "10s")
      --redis.url string                         redis url connections string (default "redis://localhost:6379")
      --redis.use_tls                            enable redis TLS
      --task_manager.engine string               tasq, machinery, empty (default "taskq")
      --task_manager.factory string              memory, redis, empty (default "memory")
      --task_manager.queue_name string           name of task manager queue (default "development_queue")
      --server.idle_timeout string               server idle timeout (default "60s")
      --server.read_timeout string               server read timeout (default "15s")
      --server.write_timeout string              server write timout (default "15s")
      --server.port string                       server port (default "3003")
```

To override with ENV variables, take the value and use the formula:
- `BUX_` prefix
- replace dot "." with underscore "\_"
- all uppercase

So for example, to override authentication key for admin, for which a path is `admin.auth_key`, do:
```BUX_ADMIN_AUTH_KEY="auth123```

To generate config file with defaults, use the --dump flag, or:
```bash
go run ./cmd/server/main.go -d
```

The default config file path is **project root**, and the default file name is **config.json**. This can be overridden by -C flag.
```bash
go run ./cmd/server/main.go -C /my/config.json
```

## Installation

**bux-server** requires a [supported release of Go](https://golang.org/doc/devel/release.html#policy).

```shell script
go get -u github.com/BuxOrg/bux-server
```

<br/>

#### build

```shell script
go build -o bux-server cmd/server/*
```

#### run

```shell script
./bux-server
```

## Documentation

View the generated [documentation](https://pkg.go.dev/github.com/BuxOrg/bux-server)

[![GoDoc](https://godoc.org/github.com/BuxOrg/bux-server?status.svg&style=flat&v=3)](https://pkg.go.dev/github.com/BuxOrg/bux-server)

<br/>

<details>
<summary><strong><code>Repository Features</code></strong></summary>
<br/>

This repository was created using [MrZ's `go-template`](https://github.com/mrz1836/go-template#about)

#### Built-in Features

-   Continuous integration via [GitHub Actions](https://github.com/features/actions)
-   Build automation via [Make](https://www.gnu.org/software/make)
-   Dependency management using [Go Modules](https://github.com/golang/go/wiki/Modules)
-   Code formatting using [gofumpt](https://github.com/mvdan/gofumpt) and linting with [golangci-lint](https://github.com/golangci/golangci-lint) and [yamllint](https://yamllint.readthedocs.io/en/stable/index.html)
-   Unit testing with [testify](https://github.com/stretchr/testify), [race detector](https://blog.golang.org/race-detector), code coverage [HTML report](https://blog.golang.org/cover) and [Codecov report](https://codecov.io/)
-   Releasing using [GoReleaser](https://github.com/goreleaser/goreleaser) on [new Tag](https://git-scm.com/book/en/v2/Git-Basics-Tagging)
-   Dependency scanning and updating thanks to [Dependabot](https://dependabot.com) and [Nancy](https://github.com/sonatype-nexus-community/nancy)
-   Security code analysis using [CodeQL Action](https://docs.github.com/en/github/finding-security-vulnerabilities-and-errors-in-your-code/about-code-scanning)
-   Automatic syndication to [pkg.go.dev](https://pkg.go.dev/) on every release
-   Generic templates for [Issues and Pull Requests](https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/configuring-issue-templates-for-your-repository) in GitHub
-   All standard GitHub files such as `LICENSE`, `CONTRIBUTING.md`, `CODE_OF_CONDUCT.md`, and `SECURITY.md`
-   Code [ownership configuration](.github/CODEOWNERS) for GitHub
-   All your ignore files for [vs-code](.editorconfig), [docker](.dockerignore) and [git](.gitignore)
-   Automatic sync for [labels](.github/labels.yml) into GitHub using a pre-defined [configuration](.github/labels.yml)
-   Built-in powerful merging rules using [Mergify](https://mergify.io/)
-   Welcome [new contributors](.github/mergify.yml) on their first Pull-Request
-   Follows the [standard-readme](https://github.com/RichardLitt/standard-readme/blob/master/spec.md) specification
-   [Visual Studio Code](https://code.visualstudio.com) configuration with [Go](https://code.visualstudio.com/docs/languages/go)
-   (Optional) [Slack](https://slack.com), [Discord](https://discord.com) or [Twitter](https://twitter.com) announcements on new GitHub Releases
-   (Optional) Easily add [contributors](https://allcontributors.org/docs/en/bot/installation) in any Issue or Pull-Request

</details>

<details>
<summary><strong><code>Package Dependencies</code></strong></summary>
<br/>

-   [BitcoinSchema/go-bitcoin](https://github.com/BitcoinSchema/go-bitcoin)
-   [BuxOrg/bux](https://github.com/BuxOrg/bux)
-   [mrz1836/go-api-router](https://github.com/mrz1836/go-api-router)
-   [mrz1836/go-sanitize](https://github.com/mrz1836/go-sanitize)
-   [stretchr/testify](https://github.com/stretchr/testify)
-   [tonicpow/go-paymail](https://github.com/tonicpow/go-paymail)
-   [See all dependencies](go.mod)
</details>

<details>
<summary><strong><code>Library Deployment</code></strong></summary>
<br/>

Releases are automatically created when you create a new [git tag](https://git-scm.com/book/en/v2/Git-Basics-Tagging)!

If you want to manually make releases, please install GoReleaser:

[goreleaser](https://github.com/goreleaser/goreleaser) for easy binary or library deployment to GitHub and can be installed:

-   **using make:** `make install-releaser`
-   **using brew:** `brew install goreleaser`

The [.goreleaser.yml](.goreleaser.yml) file is used to configure [goreleaser](https://github.com/goreleaser/goreleaser).

<br/>

### Automatic Releases on Tag Creation (recommended)

Automatic releases via [GitHub Actions](.github/workflows/release.yml) from creating a new tag:

```shell
make tag version=1.2.3
```

<br/>

### Manual Releases (optional)

Use `make release-snap` to create a snapshot version of the release, and finally `make release` to ship to production (manually).

<br/>

</details>

<details>
<summary><strong><code>Makefile Commands</code></strong></summary>
<br/>

View all `makefile` commands

```shell script
make help
```

List of all current commands:

```text
all                           Runs multiple commands
clean                         Remove previous builds and any cached data
clean-mods                    Remove all the Go mod cache
coverage                      Shows the test coverage
diff                          Show the git diff
generate                      Runs the go generate command in the base of the repo
godocs                        Sync the latest tag with GoDocs
graphql                       Generates the graphql schemas
help                          Show this help message
install                       Install the application
install-all-contributors      Installs all contributors locally
install-go                    Install the application (Using Native Go)
install-releaser              Install the GoReleaser application
lint                          Run the golangci-lint application (install if not found)
release                       Full production release (creates release in GitHub)
release                       Runs common.release then runs godocs
release-snap                  Test the full release (build binaries)
release-test                  Full production test release (everything except deploy)
replace-version               Replaces the version in HTML/JS (pre-deploy)
tag                           Generate a new tag and push (tag version=0.0.0)
tag-remove                    Remove a tag if found (tag-remove version=0.0.0)
tag-update                    Update an existing tag to current commit (tag-update version=0.0.0)
test                          Runs lint and ALL tests
test-ci                       Runs all tests via CI (exports coverage)
test-ci-no-race               Runs all tests via CI (no race) (exports coverage)
test-ci-short                 Runs unit tests via CI (exports coverage)
test-no-lint                  Runs just tests
test-short                    Runs vet, lint and tests (excludes integration tests)
test-unit                     Runs tests and outputs coverage
uninstall                     Uninstall the application (and remove files)
update-contributors           Regenerates the contributors html/list
update-linter                 Update the golangci-lint package (macOS only)
vet                           Run the Go vet application
```

</details>

<br/>

## Examples & Tests

All unit tests run via [GitHub Actions](https://github.com/BuxOrg/bux-server/actions) and
uses [Go version 1.19.x](https://golang.org/doc/go1.19). View the [configuration file](.github/workflows/run-tests.yml).

<br/>

Run all tests (including integration tests)

```shell script
make test
```

<br/>

Run tests (excluding integration tests)

```shell script
make test-short
```

<br/>

## Benchmarks

Run the Go benchmarks:

```shell script
make bench
```

<br/>

## Code Standards

Read more about this Go project's [code standards](.github/CODE_STANDARDS.md).

<br/>

## Usage

Checkout the docker compose quickstart below for a quick way to get started.

<br/>

## Docker Compose Quickstart

To get started with development, `bux-server` provides a `start-bux-server.sh` script
which is using `docker-compose.yml` file to starts up Bux Server with selected database
and cache storage. To start, we need to fill the config json which we want to use,
for example: `config/envs/development.json`.

Main configuration is done when running the script.

There are two way of running this script:
1. with manual configuration - Every option is displayed in terminal and user can choose
   which database/cache storage use and configure how to run bux-server.
  ```bash
  ./start-bux-server.sh
  ```
2. with flags which define how to set up docker services. Ever option is displayed when
   you ran the script with flag `-h` or `--help`. Possible options:

  ```bash
  ./start-bux-server.sh --help

  Welcome in Bux Server!
  Usage: ./start-bux-server.sh [OPTIONS]

  This script helps you to run Bux server with your preferred database and cache storage.

  Options:

  -db,  --database       Define database - postgresql, mongodb, sqlite
  -c,   --cache          Define cache storage - freecache(in-memory), redis
  -bs,  --bux-server     Whether the bux-server should be run - true/false
  -env, --environment    Define bux-server environment - development/staging/production
  -b,   --background     Whether the bux-server should be run in background - true/false
  -x,   --xpub           Define admin xPub
  -l,   --load           Load .env.config file and run bux-server with its settings
  ```

  ```bash
  ./start-bux-server.sh -db postgresql -c redis -bs true -env development -b false 
  ```

`-l/--load` option add possibility to use previously created `.env.config` file and run bux-server with simple command:
  ```bash
  ./start-bux-server.sh -l
  ```

## Contributing

View the [contributing guidelines](.github/CONTRIBUTING.md) and follow the [code of conduct](.github/CODE_OF_CONDUCT.md).

<br/>

### How can I help?

All kinds of contributions are welcome :raised_hands:!
The most basic way to show your support is to star :star2: the project, or to raise issues :speech_balloon:.
You can also support this project by [becoming a sponsor on GitHub](https://github.com/sponsors/BuxOrg) :clap:
or by making a [**bitcoin donation**](https://gobitcoinsv.com/#sponsor?utm_source=github&utm_medium=sponsor-link&utm_campaign=bux-server&utm_term=bux-server&utm_content=bux-server) to ensure this journey continues indefinitely! :rocket:

[![Stars](https://img.shields.io/github/stars/BuxOrg/bux-server?label=Please%20like%20us&style=social&v=3)](https://github.com/BuxOrg/bux-server/stargazers)

<br/>

### Contributors ✨

Thank you to these wonderful people ([emoji key](https://allcontributors.org/docs/en/emoji-key)):

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tr>
    <td align="center"><a href="https://mrz1818.com"><img src="https://avatars.githubusercontent.com/u/3743002?v=4?s=100" width="100px;" alt=""/><br /><sub><b>Mr. Z</b></sub></a><br /><a href="#infra-mrz1836" title="Infrastructure (Hosting, Build-Tools, etc)">🚇</a> <a href="https://github.com/BuxOrg/bux-server/commits?author=mrz1836" title="Code">💻</a> <a href="#maintenance-mrz1836" title="Maintenance">🚧</a> <a href="#security-mrz1836" title="Security">🛡️</a></td>
    <td align="center"><a href="https://github.com/icellan"><img src="https://avatars.githubusercontent.com/u/4411176?v=4?s=100" width="100px;" alt=""/><br /><sub><b>Siggi</b></sub></a><br /><a href="#infra-icellan" title="Infrastructure (Hosting, Build-Tools, etc)">🚇</a> <a href="https://github.com/BuxOrg/bux-server/commits?author=icellan" title="Code">💻</a> <a href="#security-icellan" title="Security">🛡️</a></td>
    <td align="center"><a href="https://github.com/galt-tr"><img src="https://avatars.githubusercontent.com/u/64976002?v=4?s=100" width="100px;" alt=""/><br /><sub><b>Dylan</b></sub></a><br /><a href="#infra-galt-tr" title="Infrastructure (Hosting, Build-Tools, etc)">🚇</a> <a href="https://github.com/BuxOrg/bux-server/commits?author=galt-tr" title="Code">💻</a></td>
  </tr>
</table>

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->

> This project follows the [all-contributors](https://github.com/all-contributors/all-contributors) specification.

<br/>

## License

[![License](https://img.shields.io/github/license/BuxOrg/bux-server.svg?style=flat&v=3)](LICENSE)
