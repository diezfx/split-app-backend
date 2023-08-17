GOLANG_CI_VERSION=v1.54.0

generate:
	go generate ./...
	go run -mod=mod entgo.io/ent/cmd/ent generate --feature sql/versioned-migration ./internal/ent/schema --target gen/ent

fmt:
	gofumpt -l -w .

test:
	go test ./... -v

lint/download:
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./.bin $(GOLANG_CI_VERSION) 
lint: lint/download
	./.bin/golangci-lint run -c ./.golangci.yml
lint/fix: lint/download
	./.bin/golangci-lint run -c ./.golangci.yml --fix


docker/build:
	docker build . -t ghcr.io/diezfx/split-app-backend:latest -f "deployment/Dockerfile" --build-arg="APP_NAME=split-app-backend"


