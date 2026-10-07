set working-directory := '../code/backend'

# format Go code
fmt:
    gofmt -w -s .

# check Go formatting
fmt-check:
    #!/usr/bin/env bash
    set -eo pipefail
    UNFORMATTED=$(gofmt -s -l .)
    if [ -n "$UNFORMATTED" ]; then
        echo "The following Go files are not formatted:"
        echo "$UNFORMATTED"
        exit 1
    fi

# lint Go code
lint:
    golangci-lint run --config=../../.golangci.yaml

# run Go tests
test:
    go test ./...

# run Go tests for CI with gotestsum and coverage
test-ci:
    #!/usr/bin/env bash
    set -uo pipefail
    gotestsum --format pkgname-and-test-fails --junitfile unit-tests.xml -- -race ./... -coverprofile=coverage.out -covermode=atomic

# run Go tests with coverage report
coverage:
    #!/usr/bin/env bash
    set -eo pipefail
    go test -coverprofile=coverage.out ./...
    go tool cover -func=coverage.out

# run Go backend development server
dev:
    go run ./cmd/server

# run database migrations
migrate:
    go run scripts/migrate.go

# update Swagger documentation
swagger:
    go run github.com/swaggo/swag/cmd/swag init -g cmd/server/main.go

# update AsyncAPI documentation
asyncapi:
    npx -y @asyncapi/cli@latest generate fromTemplate docs/asyncapi.yaml @asyncapi/html-template -o docs/asyncapi-docs --force-write

# update Swagger and AsyncAPI documentation
docs: swagger asyncapi

# run simulation with custom args
simulation *args:
    go run ./cmd/simulation {{args}}