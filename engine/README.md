<div align="center">

# SPV Wallet Engine

[![Release](https://img.shields.io/github/release-pre/bitcoin-sv/spv-wallet/engine.svg?logo=github&style=flat&v=2)](https://github.com/bitcoin-sv/spv-wallet/engine/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/bitcoin-sv/spv-wallet/engine/run-tests.yml?branch=master&v=2)](https://github.com/bitcoin-sv/spv-wallet/engine/actions)
[![Report](https://goreportcard.com/badge/github.com/bitcoin-sv/spv-wallet/engine?style=flat&v=2)](https://goreportcard.com/report/github.com/bitcoin-sv/spv-wallet/engine)
[![codecov](https://codecov.io/gh/bitcoin-sv/spv-wallet/engine/branch/master/graph/badge.svg?v=2)](https://codecov.io/gh/bitcoin-sv/spv-wallet/engine)
[![Mergify Status](https://img.shields.io/endpoint.svg?url=https://api.mergify.com/v1/badges/bitcoin-sv/spv-wallet/engine&style=flat&v=2)](https://mergify.com)
<br>

[![Go](https://img.shields.io/github/go-mod/go-version/bitcoin-sv/spv-wallet/engine?v=2)](https://golang.org/)
[![Gitpod Ready-to-Code](https://img.shields.io/badge/Gitpod-ready--to--code-blue?logo=gitpod&v=2)](https://gitpod.io/#https://github.com/bitcoin-sv/spv-wallet/engine)
[![standard-readme compliant](https://img.shields.io/badge/readme%20style-standard-brightgreen.svg?style=flat&v=2)](https://github.com/RichardLitt/standard-readme)
[![Makefile Included](https://img.shields.io/badge/Makefile-Supported%20-brightgreen?=flat&logo=probot&v=2)](Makefile)
<br/>

</div>

> Bitcoin UTXO & xPub Management Engine

## Table of Contents

-   [SPV Wallet Engine](#spv-wallet-engine)
    -   [Table of Contents](#table-of-contents)
    -   [About](#about)
        -   [DISCLAIMER](#disclaimer)
        -   [SPV Wallet Engine: Out-of-the-box Features:](#spv-wallet-engine-out-of-the-box-features)
        -   [**Project Assumptions: MVP**](#project-assumptions-mvp)
    -   [Installation](#installation)
    -   [Documentation](#documentation)
        -   [Built-in Features](#built-in-features)
    -   [Usage](#usage)
        -   [Examples \& Tests](#examples--tests)
        -   [Benchmarks](#benchmarks)
    -   [Code Standards](#code-standards)
    -   [Usage](#usage-1)
    -   [Contributing](#contributing)
    -   [License](#license)

<br/>

## About

> **TLDR;**
>
> Application developers should focus on their applications and should not be bogged down with managing UTXOs or XPubs. Developers should be able to use an open-source, easy to install solution to rapidly build full-featured Bitcoin applications.

<br/>

---

#### DISCLAIMER

> SPV Wallet Engine is still considered _"ALPHA"_ and should not be used in production until a major v1.0.0 is released.

---

<br/>

#### SPV Wallet Engine: Out-of-the-box Features:

-   xPub & UTXO State Management (state, balance, utxos, destinations)
-   Bring your own Database ([PostgreSQL](https://www.postgresql.org/), [SQLite](https://www.sqlite.org), [MySQL](https://www.mysql.com/), [Mongo](https://www.mongodb.com/) or [interface](./datastore/interface.go) your own)
-   Caching ([FreeCache](https://github.com/github.com/coocood/freecache), [Redis](https://redis.io/) or [interface](https://github.com/mrz1836/go-cachestore/blob/master/interface.go) your own)
-   Task Management ([TaskQ](https://github.com/vmihailenco/taskq) or [interface](taskmanager/interface.go) your own)
-   Transaction Syncing (queue, broadcast, push to mempool or on-chain, or [interface](chainstate/interface.go) your own)
-   Future plugins using [BRFC standards](http://bsvalias.org/01-brfc-specifications.html)

#### **Project Assumptions: MVP**

-   _No private keys are used_, only the xPub (or access key) is given to SPV Wallet Engine
-   (BYOX) `Bring your own xPub`
-   Signing a transaction is outside this application (IE: [spv-wallet](https://github.com/bitcoin-sv/spv-wallet) or [spv-wallet-client](https://github.com/bitcoin-sv/spv-wallet-go-client))
-   All transactions need to be submitted to the SPV Wallet service to effectively track utxo states
-   Database can be backed up, but not regenerated from chain
    -   Certain data is not on chain, plus re-scanning an xPub is expensive and not easily possible with 3rd party limitations

<br/>

## Installation

**spv-wallet/engine** requires a [supported release of Go](https://golang.org/doc/devel/release.html#policy).

```shell script
go get -u github.com/bitcoin-sv/spv-wallet/engine
```

<br/>

## Documentation

View the generated [documentation](https://pkg.go.dev/github.com/bitcoin-sv/spv-wallet/engine)

[![GoDoc](https://godoc.org/github.com/bitcoin-sv/spv-wallet/engine?status.svg&style=flat&v=2)](https://pkg.go.dev/github.com/bitcoin-sv/spv-wallet/engine)

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

-   [bitcoinschema/go-bitcoin](https://github.com/bitcoinschema/go-bitcoin)
-   [bitcoinschema/go-map](https://github.com/bitcoinschema/go-map)
-   [coocood/freecache](https://github.com/coocood/freecache)
-   [gorm.io/gorm](https://gorm.io/gorm)
-   [libsv/go-bk](https://github.com/libsv/go-bk)
-   [libsv/go-bt](https://github.com/libsv/go-bt)
-   [mrz1836/go-cache](https://github.com/mrz1836/go-cache)
-   [mrz1836/go-cachestore](https://github.com/mrz1836/go-cachestore)
-   [mrz1836/go-logger](https://github.com/mrz1836/go-logger)
-   [newrelic/go-agent](https://github.com/newrelic/go-agent)
-   [robfig/cron](https://github.com/robfig/cron)
-   [stretchr/testify](https://github.com/stretchr/testify)
-   [tonicpow/go-minercraft](https://github.com/tonicpow/go-minercraft)
-   [bitcoin-sv/go-paymail](https://github.com/bitcoin-sv/go-paymail)
-   [vmihailenco/taskq](https://github.com/vmihailenco/taskq)
</details>

## Usage

### Examples & Tests

Checkout all the [examples](examples)!

All unit tests and [examples](examples) run via [GitHub Actions](https://github.com/bitcoin-sv/spv-wallet/engine/actions) and
uses [Go version 1.19.x](https://golang.org/doc/go1.19). View the [configuration file](.github/workflows/run-tests.yml).

<br/>

Run all unit tests (excluding database tests)

```shell script
make test
```

<br/>

Run database integration tests

```shell script
make test-all-db
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

## Code Standards

Read more about this Go project's [code standards](.github/CODE_STANDARDS.md).

<br/>

## Usage

```
func main() {
	client, err := engine.NewClient(
		context.Background(), // Set context
	)
	if err != nil {
		log.Fatalln("error: " + err.Error())
	}

	defer func() {
		_ = client.Close(context.Background())
	}()

	log.Println("client loaded!", client.UserAgent())
}
```

Checkout all the [examples](examples)!

<br/>

## Contributing

All kinds of contributions are welcome!
<br/>
To get started, take a look at [code standards](.github/CODE_STANDARDS.md).
<br/>
View the [contributing guidelines](.github/CODE_STANDARDS.md#3-contributing) and follow the [code of conduct](.github/CODE_OF_CONDUCT.md).

<br/>

## License

[![License](https://img.shields.io/github/license/bitcoin-sv/spv-wallet/engine.svg?style=flat&v=2)](LICENSE)
