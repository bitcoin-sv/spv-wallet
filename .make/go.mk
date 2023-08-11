## Default to the repo name if empty
ifndef BINARY_NAME
	override BINARY_NAME=app
endif

## Define the binary name
ifdef CUSTOM_BINARY_NAME
	override BINARY_NAME=$(CUSTOM_BINARY_NAME)
endif

## Set the binary release names
DARWIN=$(BINARY_NAME)-darwin
LINUX=$(BINARY_NAME)-linux
WINDOWS=$(BINARY_NAME)-windows.exe

## Define the binary name
TAGS=
ifdef GO_BUILD_TAGS
	override TAGS=-tags $(GO_BUILD_TAGS)
endif

.PHONY: bench
bench:  ## Run all benchmarks in the Go application
	@echo "running benchmarks..."
	@go test -bench=. -benchmem $(TAGS)

.PHONY: build-go
build-go:  ## Build the Go application (locally)
	@echo "building go app..."
	@go build -o bin/$(BINARY_NAME) $(TAGS)

.PHONY: clean-mods
clean-mods: ## Remove all the Go mod cache
	@echo "cleaning mods..."
	@go clean -modcache

.PHONY: coverage
coverage: ## Shows the test coverage
	@echo "creating coverage report..."
	@go test -coverprofile=coverage.out ./... $(TAGS) && go tool cover -func=coverage.out $(TAGS)

.PHONY: generate
generate: ## Runs the go generate command in the base of the repo
	@echo "generating files..."
	@go generate -v $(TAGS)

.PHONY: godocs
godocs: ## Sync the latest tag with GoDocs
	@echo "syndicating to GoDocs..."
	@test $(GIT_DOMAIN)
	@test $(REPO_OWNER)
	@test $(REPO_NAME)
	@test $(VERSION_SHORT)
	@curl https://proxy.golang.org/$(GIT_DOMAIN)/$(REPO_OWNER)/$(REPO_NAME)/@v/$(VERSION_SHORT).info

.PHONY: install
install: ## Install the application
	@echo "installing binary..."
	@go build -o $$GOPATH/bin/$(BINARY_NAME) $(TAGS)

.PHONY: install-go
install-go: ## Install the application (Using Native Go)
	@echo "installing package..."
	@go install $(GIT_DOMAIN)/$(REPO_OWNER)/$(REPO_NAME) $(TAGS)

.PHONY: lint
lint: ## Run the golangci-lint application (install if not found)
	@echo "installing golangci-lint..."
	@#Travis (has sudo)
	@if [ "$(shell command -v golangci-lint)" = "" ] && [ $(TRAVIS) ]; then curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.53.3 && sudo cp ./bin/golangci-lint $(go env GOPATH)/bin/; fi;
	@#AWS CodePipeline
	@if [ "$(shell command -v golangci-lint)" = "" ] && [ "$(CODEBUILD_BUILD_ID)" != "" ]; then curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.53.3; fi;
	@#GitHub Actions
	@if [ "$(shell command -v golangci-lint)" = "" ] && [ "$(GITHUB_WORKFLOW)" != "" ]; then curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sudo sh -s -- -b $(go env GOPATH)/bin v1.53.3; fi;
	@#Brew - MacOS
	@if [ "$(shell command -v golangci-lint)" = "" ] && [ "$(shell command -v brew)" != "" ]; then brew install golangci-lint; fi;
	@#MacOS Vanilla
	@if [ "$(shell command -v golangci-lint)" = "" ] && [ "$(shell command -v brew)" != "" ]; then curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- v1.53.3; fi;
	@echo "running golangci-lint..."
	@golangci-lint run --verbose

.PHONY: test
test: ## Runs lint and ALL tests
	@$(MAKE) lint
	@echo "running tests..."
	@go test ./... -v $(TAGS)

.PHONY: test-unit
test-unit: ## Runs tests and outputs coverage
	@echo "running unit tests..."
	@go test ./... -race -coverprofile=coverage.txt -covermode=atomic $(TAGS)

.PHONY: test-short
test-short: ## Runs vet, lint and tests (excludes integration tests)
	@$(MAKE) lint
	@echo "running tests (short)..."
	@go test ./... -v -test.short $(TAGS)

.PHONY: test-ci
test-ci: ## Runs all tests via CI (exports coverage)
	@$(MAKE) lint
	@echo "running tests (CI)..."
	@go test ./... -race -coverprofile=coverage.txt -covermode=atomic $(TAGS)

.PHONY: test-ci-no-race
test-ci-no-race: ## Runs all tests via CI (no race) (exports coverage)
	@$(MAKE) lint
	@echo "running tests (CI - no race)..."
	@go test ./... -coverprofile=coverage.txt -covermode=atomic $(TAGS)

.PHONY: test-ci-short
test-ci-short: ## Runs unit tests via CI (exports coverage)
	@$(MAKE) lint
	@echo "running tests (CI - unit tests only)..."
	@go test ./... -test.short -race -coverprofile=coverage.txt -covermode=atomic $(TAGS)

.PHONY: test-no-lint
test-no-lint: ## Runs just tests
	@echo "running tests..."
	@go test ./... -v $(TAGS)

.PHONY: uninstall
uninstall: ## Uninstall the application (and remove files)
	@echo "uninstalling go application..."
	@test $(BINARY_NAME)
	@test $(GIT_DOMAIN)
	@test $(REPO_OWNER)
	@test $(REPO_NAME)
	@go clean -i $(GIT_DOMAIN)/$(REPO_OWNER)/$(REPO_NAME)
	@rm -rf $$GOPATH/src/$(GIT_DOMAIN)/$(REPO_OWNER)/$(REPO_NAME)
	@rm -rf $$GOPATH/bin/$(BINARY_NAME)

.PHONY: update
update:  ## Update all project dependencies
	@echo "updating dependencies..."
	@go get -u ./... && go mod tidy

.PHONY: update-linter
update-linter: ## Update the golangci-lint package (macOS only)
	@echo "upgrading golangci-lint..."
	@brew upgrade golangci-lint

.PHONY: vet
vet: ## Run the Go vet application
	@echo "running go vet..."
	@go vet -v ./... $(TAGS)
