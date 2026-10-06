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

# run prettier
fmt: 
    pnpm fmt

# run eslint
lint:
    pnpm lint

# run unit tests (vitest)
test:
    pnpm test:unit

# run e2e tests (playwright)
test-e2e:
    pnpm test:e2e

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
