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
