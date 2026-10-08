# GoStrategy

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go&logoColor=white)](https://go.dev/)
[![SvelteKit Version](https://img.shields.io/badge/SvelteKit-5-FF3E00?style=flat&logo=svelte&logoColor=white)](https://kit.svelte.dev/)
[![Project Version](https://img.shields.io/github/v/release/Thomas-More-Digital-Innovation/2526-DI-004-GoStrategy?color=blue&style=flat)](CHANGELOG.md)
[![License](https://img.shields.io/badge/license-MIT-green?style=flat)](LICENSE)

<img src="documents/pictures/logo.png" width=200>

[![Production](https://github.com/Thomas-More-Digital-Innovation/2526-DI-004-GoStrategy/actions/workflows/cd.yml/badge.svg)](https://github.com/Thomas-More-Digital-Innovation/2526-DI-004-GoStrategy/actions/workflows/cd.yml)
[![Backend CI](https://github.com/Thomas-More-Digital-Innovation/2526-DI-004-GoStrategy/actions/workflows/backend-ci.yml/badge.svg)](https://github.com/Thomas-More-Digital-Innovation/2526-DI-004-GoStrategy/actions/workflows/backend-ci.yml)
[![Frontend CI](https://github.com/Thomas-More-Digital-Innovation/2526-DI-004-GoStrategy/actions/workflows/frontend-ci.yml/badge.svg)](https://github.com/Thomas-More-Digital-Innovation/2526-DI-004-GoStrategy/actions/workflows/frontend-ci.yml)
[![Codecov](https://codecov.io/gh/Thomas-More-Digital-Innovation/2526-DI-004-GoStrategy/branch/main/graph/badge.svg)](https://codecov.io/gh/Thomas-More-Digital-Innovation/2526-DI-004-GoStrategy)
[![CodeQL](https://github.com/Thomas-More-Digital-Innovation/2526-DI-004-GoStrategy/actions/workflows/codeql.yml/badge.svg)](https://github.com/Thomas-More-Digital-Innovation/2526-DI-004-GoStrategy/actions/workflows/codeql.yml)
[![golangci-lint](https://img.shields.io/badge/linting-golangci--lint-blue?style=flat&logo=golangci-lint)](https://github.com/Thomas-More-Digital-Innovation/2526-DI-004-GoStrategy/actions/workflows/backend-ci.yml)

Open-source implementation of the classic GoStrategy board game. Built with a high-performance Go backend and a modern SvelteKit frontend, designed for both Human-vs-AI and AI-vs-AI experimentation.

Play the game here: https://gostrategy.dotsem.be

Or run it yourself locally: [see "Running with Docker (Recommended)"](#running-with-docker-recommended).

## Project Overview

This project aims to provide a robust platform for playing GoStrategy while serving as a testing ground for various AI strategies. It features real-time gameplay via WebSockets, secure user authentication, and a declarative infrastructure managed via NixOS.

### Tech Stack
- **Backend**: Go (Gin, Gorilla WebSocket)
- **Frontend**: SvelteKit 5 (Runes, Tailwind CSS)
- **Database**: PostgreSQL 15
- **Infrastructure**: NixOS, Docker, Cloudflare Tunnels
- **Monitoring**: Loki, Promtail

## Getting Started

### Prerequisites
- Docker and Docker Compose
- Bun (for local frontend development)
- Go 1.25+ (for local backend development)
- [just](https://github.com/casey/just) command runner
- [lefthook](https://github.com/evilmartians/lefthook) Git hooks manager

### Running with Docker (Recommended)
The easiest way to get the full stack running is using docker compose:
```bash
docker compose up
```
This will spin up the backend, frontend, and database containers. The app will be available at `http://localhost:5000`.

### Local Development
If you prefer running services outside of Docker:

**Git Hooks (Lefthook & Just):**
Install the pre-commit hooks to automatically format and lint code before committing:
```bash
lefthook install
```
Lefthook runs on `git commit` to auto-format staged Go files and execute `just go lint` (which runs `golangci-lint`). You can also run the hooks manually at any time:
```bash
lefthook run pre-commit
```

**Backend:**
```bash
cd code/backend
go run cmd/server/main.go
```

**Frontend:**
```bash
cd code/frontend
bun install
bun run dev
```
or use pnpm:
```bash
cd code/frontend
pnpm install
pnpm dev
```

### Simulate AI vs AI Games
From the `code/backend` directory:

```bash
cd code/backend
go run cmd/simulation/main.go --ai={ai1}:{ai2} --format md --logging=false --matches {n}
```

Supported AI identifiers for `{ai1}` and `{ai2}` include `fafo`, `fato`, `heuristic`, `minimax`, and `mcts`. Optional flags include `--setup` to load predetermined piece setups:

```bash
# Run 100 matches between MCTS and FATO
go run cmd/simulation/main.go --ai=mcts:fato --format md --matches 100

# Run matches with a specific board configuration
go run cmd/simulation/main.go --ai=mcts:minimax --matches 10 --setup ../../documents/files/board_setups/best_setup_ever.json
```

## Architecture & Performance

### High-Performance Concurrency
The backend uses an event-driven architecture with Go channels to manage game states. This ensures minimal mutex contention and high throughput for concurrent matches.

## AI Architecture & Agents

GoStrategy includes an AI experimentation suite with multiple agent implementations tailored for imperfect-information gameplay, benchmark tournaments, and parameter optimization:

| Agent Identifier | Strategy / Algorithm | Description | Key Parameters |
| :--- | :--- | :--- | :--- |
| `fafo` | Pure Random | Baseline agent selecting uniform random legal moves without state tracking or evaluation. | — |
| `fato` | Rule-based Heuristic | Greedy tactical agent with piece memory, loop evasion, and aggression-scaled attack thresholds. | `aggression` |
| `heuristic` | Static Evaluation Engine | Evaluates board states using strategic piece values, forward territory exploration, combat risk, and aggression weightings. | `aggression`, feature weights |
| `minimax` | Adversarial Tree Search | Minimax with alpha-beta pruning over determinized hidden-information boards, evaluated via the static evaluation engine. | `depth`, `aggression`, feature weights |
| `mcts` | Monte Carlo Tree Search | Information-Set MCTS using root-level determinization sampling, concurrent UCB1 tree selection across multi-worker goroutines, and fast rollout simulations. | `iterations`, `total_rollouts`, `exploration_constant`, `aggression` |

Detailed match histories and benchmark analytics are archived in the [AI Data folder](documents/files/ai-data/).

> [!NOTE]
> The AI's are still in development, and are not yet optimized to their full potential. 

## Documentation
- [Changelog](CHANGELOG.md) - Project history and versioning.
- [Roadmap](ROADMAP.md) - Future plans and hardening steps.
- [API Documentation](http://localhost:8080/swagger/index.html) - Swagger UI (available when running locally).

## Contributing
- [Contributing Guide](CONTRIBUTING.md) - Guidelines for contributing to the project.

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
