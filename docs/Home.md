# GoStrategy Documentation

Welcome to the official documentation for **GoStrategy**, an open-source, full-stack implementation of the classic Stratego board game with modern competitive online play and artificial intelligence research capabilities.

---

## Quick Navigation

* **[Docker Deployment](Docker-Deployment):** Run GoStrategy locally or on your own VPS with Docker Compose.
* **[Configuration](Configuration):** Complete reference of all backend and frontend environment variables.
* **[Architecture](Architecture):** Deep dive into the Go game engine, WebSocket protocol, AI agents, and SvelteKit frontend.

---

## Core Capabilities

* **Game Engine:** Built in Go with strict game rule validation (including all generic stratego rules).
* **AI Engine & Benchmarking:** Play against Heuristic, Minimax (alpha-beta pruning), and Monte Carlo Tree Search (MCTS) algorithms.
* **Real-Time Multiplayer:** Concurrency-safe WebSocket hubs supporting low-latency state synchronization and spectator modes.
* **Custom Board Setups:** Drag-and-drop piece deployment with support for exporting setups to JSON.
* **Self-Hostable & Air-Gapped:** Distributed as lightweight multi-stage Docker containers with zero external cloud dependencies.

---

## Contributing & Issues

Documentation source files live directly in the main repository inside the `/docs` directory. If you find a typo or want to expand a guide, submit a Pull Request. Editing via the GitHub interface is disabled and will be overwritten by the CI/CD pipeline.
