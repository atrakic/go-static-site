all:
	go run main.go

deploy:
	rm -rf _site
	cp -a tmp/makesite/_site .
	git add _site


CHART ?= $(shell basename $$PWD)
release: ## Release (eg. V=0.0.1)
	 @[ "$(V)" ] \
		 && read -p "Press enter to confirm and push tag v$(V) to origin, <Ctrl+C> to abort ..." \
		 && ./scripts/bump_chart.sh $(CHART) $(V) \
		 && git tag v$(V) -m "chore: v$(V)" \
		 && git push origin v$(V) -f \
		 && git fetch --tags --force --all -p \
		 && git describe --tags $(shell git rev-list --tags --max-count=1)

export_docker:
	docker export go-static-site | tar t
