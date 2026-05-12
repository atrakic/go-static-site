SHELL := /bin/sh

APP := $(notdir $(CURDIR))
CHART ?= $(APP)
IMAGE ?= go-static-site

.PHONY: all run deploy release export_docker help

all: run

run: ## Run the app locally
	go run .

deploy: ## Refresh generated site output
	rm -rf _site
	cp -a tmp/makesite/_site .
	git add _site

release: ## Release (usage: make release V=0.0.1)
	@test -n "$(V)" || (echo "Error: V is required, e.g. make release V=0.0.1" && exit 1)
	@printf "About to create and push tag %s. Press Enter to continue or Ctrl+C to abort..." "$(V)"; \
	read dummy
	./scripts/bump_chart.sh $(CHART) $(V)
	git tag "$(V)" -m "release: $(V)"
	git push origin "$(V)"
	git fetch --tags --force --all -p
	@if [ -n "$(GITHUB_TOKEN)" ]; then \
		curl -fsSL \
			-H "Authorization: token $(GITHUB_TOKEN)" \
			-H "Accept: application/vnd.github.v3+json" \
			-X POST \
			https://api.github.com/repos/atrakic/$(APP)/releases \
			-d "{\"tag_name\":\"$(V)\",\"generate_release_notes\":true}"; \
	else \
		echo "GITHUB_TOKEN not set; skipping GitHub release creation"; \
	fi

export_docker: ## List files in exported container filesystem
	docker export $(IMAGE) | tar t

help: ## Show targets
	@awk 'BEGIN {FS = ":.*## "}; /^[a-zA-Z0-9_-]+:.*## / {printf "%-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)
