# GoStrategy — Coding Guidelines & Agent Directives

This document defines the coding standards, architectural decisions, and workflows for **GoStrategy**. Every AI assistant and developer working on this project must adhere strictly to these rules.

---

## 1. Universal Principles (All Stacks)

### 1.1 The Ponytail Workflow (Pragmatic Senior Dev Mode)
Before writing any code, stop at the first rung that holds:
1. **YAGNI:** Does this need to be built at all?
2. **Reuse:** Does it already exist in this codebase? Reuse existing helpers/patterns.
3. **Stdlib First:** Does the Go or Web standard library already do this? Use it.
4. **Native Platform:** Does a native browser or platform feature cover it? Use it.
5. **Existing Deps:** Does an already-installed dependency solve it? Use it.
6. **Simplicity:** Can this be one clean line? Make it one line.
7. **Minimal Diff:** Only then write the minimum code that works.

#### Ponytail Execution Rules:
- **No Unrequested Abstractions:** No boilerplate, speculative layers, or interfaces nobody asked for.
- **Deletion Over Addition:** Boring over clever. Fewest files possible.
- **Shortcuts & Ceilings:** Mark intentional simplifications with a `// ponytail:` comment naming the ceiling and upgrade path (e.g., `// ponytail: in-memory session map; upgrade to redis if multi-instance clustering needed`).
- **Non-Negotiables:** Never skip input validation at trust boundaries, WebSocket frame validation, security checks, or error handling preventing game state corruption.

### 1.2 Comments & Documentation
- **No Redundant Comments:** Absolutely no comments echoing code syntax (e.g., `turn++ // increment turn`).
- **Omission Rule:** If a comment does not answer "How?", "Why?", or "WTF?", omit it entirely.
- **Allowed Comments:**
  - Package-level documentation.
  - Standard doc comments on exported symbols (Godoc, TSDoc).
  - High-level `// why:` context explaining non-obvious game engine or network constraints.
  - Actionable `// todo:` comments explaining future improvements or known rough edges.
  - `// ponytail:` comments documenting intentional MVP simplifications.
- **Style:** Short comments in lowercase. Capital letters only for multi-line contextual explanations.
- **Preservation:** Always preserve Swagger (`@Summary`, `@Tags`, `@Router`), AsyncAPI, and Godoc annotations.

### 1.3 Function & File Limits
- **Single Responsibility (SOLID):** Functions must execute one unambiguous action.
- **DRY:** Never duplicate logic; generalize shared domain logic.
- **File Length Ceiling:** Strict maximum of **300 lines per file**. If a file exceeds this:
  - Decouple, isolate state, or split sub-components/modules into colocated files.
  - Notify the developer before performing major structural refactoring.
- **Clean Code:** Zero dead code, unused imports, or magic numbers (use named constants or enums).

### 1.4 Development & Tooling Guardrails
- **Package Managers:** Frontend strictly uses `pnpm` or `bun`. Never run `npm` or `yarn`.
- **Pre-commit Automation:** Enforced via `lefthook` (runs `gofmt` and `just go lint`).
- **AI Tooling Constraints:**
  - AI assistants must **never** run `git commit` or `git push`.
  - AI assistants must **never** execute destructive migration or refactoring plans without approval.
  - Test suite guard: do not invoke full test suites natively unless explicitly instructed.
- **Commit Format:** `action(scope): description` (e.g., `feat(game): add piece movement validation`).

---

## 2. Go (Backend Game Engine & API)

### 2.1 Architecture & Layout
- **Module:** `digital-innovation/gostrategy`
- **Standard Go Layout:**
  - `/code/backend/cmd/server` -> HTTP & WebSocket server entry point (`main.go`).
  - `/code/backend/cmd/simulation` -> AI engine simulation & benchmarking binaries.
  - `/code/backend/internal` -> Private core logic. Absolute ban on external package exposure:
    - `/internal/game` -> Game board, turn lifecycle, movement rules, validation.
    - `/internal/ai` -> AI algorithms, heuristics, bot agents.
    - `/internal/api` -> Gin route handlers, WebSocket hubs, and middleware.
    - `/internal/db` -> GORM models, connection pool, migrations.
    - `/internal/auth` -> JWT issuance, verification, user sessions.
    - `/internal/logging` & `/internal/telemetry` -> slog-based logging and Prometheus metrics.
  - `/code/backend/pkg` -> Explicitly exportable, standalone utility libraries (keep minimal).

### 2.2 Structural Rules & Best Practices
- **Implicit Interfaces:** Rely on Go's implicit interface satisfaction. Define interfaces at the consumer layer, never speculative producer interfaces.
- **Error Handling:** Return early on errors (`if err != nil { return err }`). Avoid deep nesting. Wrap errors with context (`fmt.Errorf("failed to process move: %w", err)`).
- **Concurrency & WebSockets:** Keep mutex critical sections minimal around game session state. Ensure thread-safe WebSocket reads and writes using decoupled channel pumps.
- **Logging:** Structured logging only via `internal/logging` (`slog`). Never log sensitive credentials, raw secrets, or PII.

### 2.3 Documentation & Contracts
- **Swagger / OpenAPI:** Decorate all HTTP endpoints with `swaggo/swag` annotations (`@Summary`, `@Tags`, `@Produce`, `@Success`, `@Failure`, `@Router`).
- **AsyncAPI:** Maintain WebSocket event schemas via AsyncAPI specifications in `./scripts/update-asyncapi.sh`.

### 2.4 Testing
- Group subtests using `t.Run`.
- Target files match target naming: `*_test.go`.
- Assertions strictly via `stretchr/testify/assert` and `stretchr/testify/require`.

---

## 3. SvelteKit & Svelte 5 (Frontend)

### 3.1 Architectural Paradigm: Feature-First Routing
The frontend follows a **Feature-First / Page-Colocated Architecture** (similar to Flutter feature-first packaging). 

#### Core Routing Boundary:
- **`src/lib/` is strictly for cross-page shared items:**
  - Put code in `src/lib/` **only** if it is consumed by two or more distinct routes.
  - Items in `src/lib/`: global API clients (`$lib/api`), shared TypeScript definitions (`$lib/types`), base UI primitives (`$lib/components`), shared math/game utilities (`$lib/utils`), and global user session state (`$lib/state`).
- **Everything page-specific stays inside the route directory:**
  - Any component, state machine, store, or helper used only by a specific route must reside in that route’s directory.
  - Use leading underscore (`_`) on folders to prevent SvelteKit from treating internal helpers as route endpoints:
    - `_components/` -> Page-local sub-components.
    - `_state/` -> Page-local runes classes or state controllers (e.g., `game-session.svelte.ts`).
    - `_utils/` -> Helpers local to the feature.
    - `_types/` -> Feature-specific local types.

#### Directory Structure Reference:
```text
code/frontend/src/
├── lib/                             # SHARED ACROSS MULTIPLE PAGES ONLY
│   ├── api/                         # Backend HTTP & WebSocket base clients
│   ├── components/                  # Atomic global UI primitives (buttons, modals, icons)
│   ├── state/                       # Global auth / theme state
│   ├── types/                       # Shared API & game domain interfaces
│   └── utils/                       # Universal math, formatting, matrix helpers
└── routes/                          # FEATURE / PAGE DIRECTORIES
    ├── +layout.svelte               # Global app shell
    ├── layout.css                   # Global Tailwind imports
    ├── +page.svelte                 # Landing / home page
    └── game/
        └── [id]/                    # "Game Session" Feature
            ├── +page.svelte         # Game view orchestrator
            ├── +page.ts             # Route loader / params validation
            ├── _components/         # Page-specific components only
            │   ├── board/
            │   │   ├── BoardGrid.svelte
            │   │   └── Tile.svelte
            │   └── right-bar/
            │       ├── GameControls.svelte
            │       └── GameHistory.svelte
            └── _state/              # Page-specific state & controllers
                ├── context.ts
                └── game-session.svelte.ts
```

### 3.2 Svelte 5 Runes & Reactivity
- **Runes Standard:** Exclusively use Svelte 5 runes (`$state`, `$derived`, `$props`, `$effect`, `$bindable`).
- **Legacy Syntax Ban:** Absolute ban on legacy Svelte v4 reactive statements (`$:`) and `export let` syntax.
- **Universal State Modules:** Encapsulate state using `.svelte.ts` files with reactive classes or functions using runes.
- **Styling:** Exclusively Tailwind CSS utility classes (Tailwind v4 via `@tailwindcss/vite`). No inline `style="..."` or unstructured raw CSS blocks.

---

## 4. Verification Commands Reference

### Root Workflows

| Command | Description |
| :--- | :--- |
| `just default` | List all available recipes and submodules (`just`) |
| `just dev` | Start full local development stack with Docker Compose |

### Backend Submodule (`just go <cmd>`)

| Command | Description |
| :--- | :--- |
| `just go fmt` | Format Go code (`gofmt -w -s .`) |
| `just go fmt-check` | Verify Go code formatting |
| `just go lint` | Run `golangci-lint` with root `.golangci.yaml` |
| `just go test` | Run Go unit/integration tests |
| `just go test-ci` | Run tests with `gotestsum` and atomic coverage profiling |
| `just go coverage` | Run test suite and print function coverage summary |
| `just go dev` | Run backend standalone (`go run ./cmd/server`) |
| `just go migrate` | Run database migrations tool (`go run scripts/migrate.go`) |
| `just go swagger` | Update Swag OpenAPI specifications |
| `just go asyncapi` | Update AsyncAPI HTML documentation |
| `just go docs` | Update both Swagger and AsyncAPI documentation |

### Frontend Submodule (`just web <cmd>`)

| Command | Description |
| :--- | :--- |
| `just web install` | Install frontend dependencies (`pnpm install`) |
| `just web update-locks` | Update and sync lockfiles (`pnpm install && bun install`) |
| `just web check` | Run `svelte-check` type checking against `tsconfig.json` |
| `just web lint` | Alias to `check` |
| `just web dev` | Run Vite development server (`pnpm dev`) |
| `just web build` | Build static production assets (`pnpm build`) |
| `just web preview` | Preview production build locally (`pnpm preview`) |
| `just web audit` | Audit frontend dependencies for security advisories |
