mod go 'just/go.just'
mod web 'just/web.just'

# list available recipes
default:
    @just --list

# start full-stack environment with Docker Compose
dev:
    docker compose up --build

# run AI engine simulation
simulation *args:
    just go simulation {{args}}

# run Playwright e2e against the production Docker images (same stack as CI)
e2e *args:
    #!/usr/bin/env bash
    set -u
    compose="docker compose --env-file .env.e2e -f docker-compose.yml -f docker-compose.e2e.yml"
    $compose up -d --build --wait || exit 1
    (cd code/frontend && CI=1 pnpm exec playwright test {{args}})
    status=$?
    $compose down -v
    exit $status
