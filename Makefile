GIT_HASH := $(shell git rev-parse --short HEAD)
DOCKER_USER := marie20767
IMAGE_NAME := keto-granola-server

dep:
	go mod download

run:
	go run main.go

lint: lint/install lint/run

lint/install:
	curl -sSfL https://golangci-lint.run/install.sh | sh -s -- v2.12.2

lint/run:
	bin/golangci-lint run --config .golangci.yml

lint/fix:
	bin/golangci-lint run --config .golangci.yml --fix

test: # Runs unit and e2e tests
	go test -shuffle=on -tags=e2e ./...

test/unit: # unit tests only
	go test -shuffle=on ./...

test/e2e: # e2e tests only
	go test -shuffle=on -tags=e2e ./internal/tests

sqlc:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.31.1 generate

mocks:
	go install github.com/matryer/moq@latest && \
	go generate ./...

migrate/create:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migrate/create name=<migration_name>"; \
		exit 1; \
	fi
	migrate create -ext sql -dir internal/store/migrations -seq $(name)

build:
	CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64 \
	go build -o keto-granola-server .

docker/local-build:
	DOCKER_BUILDKIT=1 docker buildx build \
	-t $(DOCKER_USER)/$(IMAGE_NAME):local .

# For CI caching
CACHE_FROM ?=
CACHE_TO ?=

docker/ci-build:
	DOCKER_BUILDKIT=1 docker buildx build \
	$(CACHE_FROM) \
	$(CACHE_TO) \
	--platform linux/amd64 \
	-t $(DOCKER_USER)/$(IMAGE_NAME):latest \
	-t $(DOCKER_USER)/$(IMAGE_NAME):$(GIT_HASH) .

docker/push:
	docker push --all-tags $(DOCKER_USER)/$(IMAGE_NAME)

docker/build-and-push:
	DOCKER_BUILDKIT=1 docker buildx build \
	--platform linux/amd64 \
	--push \
	-t $(DOCKER_USER)/$(IMAGE_NAME):latest \
	-t $(DOCKER_USER)/$(IMAGE_NAME):$(GIT_HASH) .

	