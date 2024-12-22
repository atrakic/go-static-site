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
		 && git tag $(V) -m "release: $(V)" \
		 && git push origin $(V) -f \
		 && git fetch --tags --force --all -p \
		 && if [ ! -z "$(GITHUB_TOKEN)" ] ; then \
			curl \
			  -H "Authorization: token $(GITHUB_TOKEN)" \
				-X POST	\
				-H "Accept: application/vnd.github.v3+json"	\
				https://api.github.com/repos/atrakic/$(shell basename $$PWD)/releases \
				-d "{\"tag_name\":\"$(V)\",\"generate_release_notes\":true}"; \
			fi;

export_docker:
	docker export go-static-site | tar t
