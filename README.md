<div align="center">

# SPV Wallet


[![Release](https://img.shields.io/github/release-pre/BuxOrg/spv-wallet.svg?logo=github&style=flat&v=3)](https://github.com/BuxOrg/spv-wallet/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/BuxOrg/spv-wallet/run-tests.yml?branch=master&v=3)](https://github.com/BuxOrg/spv-wallet/actions)
[![Report](https://goreportcard.com/badge/github.com/BuxOrg/spv-wallet?style=flat&v=3)](https://goreportcard.com/report/github.com/BuxOrg/spv-wallet)
[![codecov](https://codecov.io/gh/BuxOrg/spv-wallet/branch/master/graph/badge.svg?v=3)](https://codecov.io/gh/BuxOrg/spv-wallet)
[![Mergify Status](https://img.shields.io/endpoint.svg?url=https://api.mergify.com/v1/badges/BuxOrg/spv-wallet&style=flat&v=3)](https://mergify.io)
<br>

[![Go](https://img.shields.io/github/go-mod/go-version/BuxOrg/spv-wallet?v=3)](https://golang.org/)
[![Gitpod Ready-to-Code](https://img.shields.io/badge/Gitpod-ready--to--code-blue?logo=gitpod&v=3)](https://gitpod.io/#https://github.com/BuxOrg/spv-wallet)
[![standard-readme compliant](https://img.shields.io/badge/readme%20style-standard-brightgreen.svg?style=flat&v=3)](https://github.com/RichardLitt/standard-readme)
[![Makefile Included](https://img.shields.io/badge/Makefile-Supported%20-brightgreen?=flat&logo=probot&v=3)](Makefile)
<br/>
</div>

## Table of Contents

- [About](#about)
- [Installation](#installation)
- [Documentation](#documentation)
- [Usage](#usage)
  - [Config Variables](#config-variables)
  - [Examples & Tests](#examples--tests)
  - [Benchmarks](#benchmarks)
- [Code Standards](#code-standards)
- [Contributing](#contributing)
- [License](#license)

<br/>

## About

Complete stand-alone server using the SPV Wallet engine (UTXOs, xPubs, Paymail & More!)

<br/>

## Installation

**spv-wallet** requires a [supported release of Go](https://golang.org/doc/devel/release.html#policy).

```shell script
go get -u github.com/BuxOrg/spv-wallet
```

#### build

```shell script
go build -o spv-wallet cmd/server/*
```

#### run

```shell script
./spv-wallet
```
<br/>

## Documentation

View the generated [documentation](https://pkg.go.dev/github.com/BuxOrg/spv-wallet)

[![GoDoc](https://godoc.org/github.com/BuxOrg/spv-wallet?status.svg&style=flat&v=3)](https://pkg.go.dev/github.com/BuxOrg/spv-wallet)

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

## Usage

> Every variable which is used and can be configured is described in [config.example.yaml](config.example.yaml)


### Defaults

If you run spv-wallet without editing anything, it will use the default configuration from file [defaults.go](/config/defaults.go). It is set up to use _freecache_, _sqlite_ with enabled _paymail_ with _signing disabled_ and with _beef_.


### Config Variables

Default config variables can be overridden by (in this order of importance):
1. Flags (only the ones below)
2. ENV variables
3. Config file

#### Flags

Available flags:

```bash
  -C, --config_file string                       custom config file path
  -h, --help                                     show help
  -v, --version                                  show version
  -d, --dump_config                              dump config to file, specified by config_file (-C) flag
```

To generate config file with defaults, use the --dump flag, or:
```bash
go run ./cmd/server/main.go -d
```

The default config file path is **project root**, and the default file name is **config.yaml**. This can be overridden by -C flag.
```bash
go run ./cmd/server/main.go -C /my/config.json
```

#### Environment variables

To override any config variable with ENV, use the "SPV\_" prefix with mapstructure annotation path with "_" as a delimiter in all uppercase. Example:

Let's take this fragment of AppConfig from `config.example.yaml`:

```yaml
auth:
    admin_key: xpub661MyMwAqRbcFgfmdkPgE2m5UjHXu9dj124DbaGLSjaqVESTWfCD4VuNmEbVPkbYLCkykwVZvmA8Pbf8884TQr1FgdG2nPoHR8aB36YdDQh
    require_signing: false
    scheme: xpub
    signing_disabled: true
```

To override admin_key in auth config, use the path with "_" as a path delimiter and SPV\_ as prefix. So:
```bash
SPV_AUTH_ADMIN_KEY="admin_key"
```

To be able to use TAAL API Key is needed. 

To get and API Key:


1. Enter the URL https://platform.taal.com/ in your browser.
2. Register or login on to TAAL PLATFORM.
3. Your mainnet and testnet API keys will be displayed on dashboard tab.

https://docs.taal.com/introduction/get-an-api-key

To use your API key put key in ``token`` field in ```config.example.yaml```

``nodes`` -> ``apis`` -> ``token``


<br/>

### Examples & Tests

All unit tests run via [GitHub Actions](https://github.com/BuxOrg/spv-wallet/actions) and
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

### Benchmarks

Run the Go benchmarks:

```shell script
make bench
```

<br/>

### Docker Compose Quickstart

To get started with development, `spv-wallet` provides a `start.sh` script
which is using `docker-compose.yml` file to starts up SPV Wallet serer with selected database
and cache storage. To start, we need to fill the config json which we want to use,
for example: `config/envs/development.json`.

Main configuration is done when running the script.

There are two way of running this script:
1. with manual configuration - Every option is displayed in terminal and user can choose
   which database/cache storage use and configure how to run spv-wallet.
  ```bash
  ./start.sh
  ```
2. with flags which define how to set up docker services. Ever option is displayed when
   you ran the script with flag `-h` or `--help`. Possible options:
  ```bash
  ./start.sh -db postgresql -c redis -bs true -env development -b false 
  ```

`-l/--load` option add possibility to use previously created `.env.config` file and run spv-wallet with simple command:
  ```bash
  ./start.sh -l
  ```
<br/>

## Code Standards
Read more about this Go project's [code standards](.github/CODE_STANDARDS.md).

<br/>

## Contributing
All kinds of contributions are welcome!
<br/>
To get started, take a look at [code standards](.github/CODE_STANDARDS.md).
<br/>
View the [contributing guidelines](.github/CODE_STANDARDS.md#3-contributing) and follow the [code of conduct](.github/CODE_OF_CONDUCT.md).

<br/>

## License

[![License](https://img.shields.io/github/license/BuxOrg/spv-wallet.svg?style=flat&v=3)](LICENSE)
