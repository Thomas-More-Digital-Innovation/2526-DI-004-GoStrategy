set working-directory := '../code/frontend'

# install web dependencies
install:
    pnpm install

# update lockfiles (pnpm and bun)
update-locks:
    pnpm install && bun install

# run Svelte type checking
check:
    pnpm check

# lint web code (alias to check)
lint: check

# build Svelte static assets
build:
    pnpm build

# preview production build locally
preview:
    pnpm preview

# audit frontend dependencies
audit:
    pnpm audit --audit-level high

# run frontend development server
dev:
    pnpm dev
