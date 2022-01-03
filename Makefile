.PHONY: lint tidy test check build pack deploy clean

help: ## Show this help
	@echo "Help"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "    \033[36m%-20s\033[93m %s\n", $$1, $$2}'

lint: ## Lint code
	@bash scripts/lint.sh

tidy: ## Check go.mod file is up to date
	@bash scripts/tidy.sh

test: ## Run unit tests
	@bash scripts/test.sh

check: lint tidy test ## Run all checks (lint, tidy, test)

build: ## Compile sources to get a binary
	@bash scripts/build.sh almendruco

pack: ## Pack a binary built previously so that it can be deployed
	@bash scripts/pack.sh almendruco almendruco.zip

deploy: ## Deploy code on AWS by updating almendruco function with a previously packaged binary
	@bash scripts/update.sh -r almendruco almendruco.zip

clean: ## Remove output binary and package
	rm almendruco almendruco.zip
