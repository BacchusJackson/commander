default:
	@just --choose

bin := "./bin"
app := bin / "cmdr"
src := "./cmd/commander"

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
# example: GITHUB_TOKEN just release
release:
	goreleaser release --clean

serve:
	cd docs && mdbook serve
	
clean:
	rm -rf {{bin}}
