# Architecture Overview

GoStrategy is architected as an event-driven, decoupled system divided into a Go backend game server and a modern Svelte 5 single-page application.

---

## High-Level Topology

```
[Browser Client]
       │
       ├─ (HTTP REST) ────────► [Gin Engine API Routes] ──► [PostgreSQL]
       │
       └─ (WebSocket Stream) ──► [WSHub Pump] ────────────► [Game Runner]
                                                                │
                                                    ┌───────────┴───────────┐
                                                    ▼                       ▼
                                              [Rule Engine]          [AI Opponents]
                                        (Checks moves & combat)     (Heuristic/MCTS)
```

---

## 1. Backend Core (`code/backend`)

* **Game Engine (`internal/game`):**
  * Models the 10x10 board, lake impassable zones, and piece ranks (Spy through Marshal, Bombs, and Flag).
  * Enforces the Two-Square Rule, scout straight-line dashes, and strike resolution.
* **WebSocket Hub (`internal/api`):**
  * Manages thread-safe bi-directional client connections via non-blocking write channels.
  * Broadcasts turn changes, revealed combat pieces, and game-over state.
* **AI Bot Agents (`internal/ai`):**
  * **Heuristic:** Fast, rule-based tactical moves prioritizing piece capture and flag hunting.
  * **Minimax:** Alpha-beta pruning search tree for depth-bounded tactical analysis.
  * **MCTS:** Monte Carlo Tree Search running stochastic simulations for strategic depth.

---

## 2. Frontend Client (`code/frontend`)

* **Framework:** SvelteKit with Svelte 5 Runes.
* **Design & Styling:** Tailwind CSS v4 with custom theme.
* **State Management:** Reactive state controllers colocated inside feature routes (e.g., `routes/game/[id]/_state/game-session.svelte.ts`).
* **Runtime:** Distributed as pre-compiled static assets served by Nginx.
