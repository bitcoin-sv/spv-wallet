# Common makefile commands & variables between projects
include .make/common.mk

# Common Golang makefile commands & variables between projects
include .make/go.mk

## Not defined? Use default repo name which is the application
ifeq ($(REPO_NAME),)
	REPO_NAME="spv-wallet-models"
endif

## Not defined? Use default repo owner
ifeq ($(REPO_OWNER),)
	REPO_OWNER="bitcoin-sv"
endif

.PHONY: all
all: ## Runs multiple commands
	@$(MAKE) test

.PHONY: clean
clean: ## Remove previous builds and any cached data
	@echo "cleaning local cache..."
	@go clean -cache -testcache -i -r
	@$(MAKE) clean-mods
	@test $(DISTRIBUTIONS_DIR)
	@if [ -d $(DISTRIBUTIONS_DIR) ]; then rm -r $(DISTRIBUTIONS_DIR); fi

.PHONY: install-all-contributors
install-all-contributors: ## Installs all contributors locally
	@echo "installing all-contributors cli tool..."
	@yarn global add all-contributors-cli

.PHONY: release
release:: ## Runs common.release then runs godocs
	@$(MAKE) godocs

.PHONY: update-contributors
update-contributors: ## Regenerates the contributors html/list
	@echo "generating contributor html..."
	@all-contributors generate
