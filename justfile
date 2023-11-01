default:
	@just --choose

bin := "./bin"
app := bin / "guardian"
src := "./cmd/guardian"

build:
	mkdir -p {{bin}}
	go build -o {{app}} {{src}}

test:
	go test -v ./...

coverage:
	go test -coverprofile cover.cov ./...

open-coverage:
	go tool cover -html=cover.cov

fmt:
	gofmt -x .

release-snapshot:
	goreleaser release --snapshot --clean

# release with goreleaser
# Needs GITHUB_TOKEN defined
release:
	goreleaser release --clean

serve_docs:
	cd docs && mdbook serve
	
clean:
	rm -rf {{bin}}

# ex:
# 	echo {{APP_MAIN}}
# ex:
# 	#!/usr/bin/env nu
# 	ls
