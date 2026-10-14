# Docker Deployment Guide

GoStrategy can be self-hosted on any system supporting Docker and Docker Compose. The official deployment stack bundles the Go backend, PostgreSQL database, and Nginx-powered SvelteKit web client.

---

## Prerequisites

* **Docker Engine:** 24.0+
* **Docker Compose:** v2.20+
* **System Resources:** 1 vCPU and 1 GB RAM minimum (2 vCPUs recommended for MCTS AI simulation).

---

## Quickstart (Production Stack)

Clone the repository and spin up the production stack:

```bash
git clone https://github.com/Thomas-More-Digital-Innovation/2526-DI-004-Stratego.git gostrategy
cd gostrategy

# Copy environment template
cp .env.template .env

# Start all services
docker compose -f docker-compose.prod.yml up -d
```

Once running, access the web client at:
* **Web UI:** `http://localhost:8080` (or your server's public IP / reverse proxy)
* **Backend Healthcheck:** `http://localhost:8081/health`

---

## Service Architecture

| Service | Port | Description |
| :--- | :--- | :--- |
| `frontend` | `8080` (mapped to internal `80`) | Nginx serving static SvelteKit SPA assets and documentation |
| `backend` | `8081` (mapped to internal `8080`) | Go HTTP/WebSocket game engine and API server |
| `database` | `5432` | PostgreSQL database storing user profiles and match history |

---

## Reverse Proxy Setup (Traefik / Nginx / Caddy)

If placing GoStrategy behind a reverse proxy with TLS/SSL:

1. Forward HTTP traffic to port `8080` (`frontend`).
2. Ensure WebSocket upgrades are enabled:
   ```nginx
   proxy_http_version 1.1;
   proxy_set_header Upgrade $http_upgrade;
   proxy_set_header Connection "upgrade";
   ```
3. Set `VITE_API_BASE` and `VITE_WS_BASE` in `.env` to match your domain (e.g., `https://stratego.example.com/api` and `wss://stratego.example.com/api/ws`).

---

## Upgrading

To pull the latest release:

```bash
git pull origin main
docker compose -f docker-compose.prod.yml build
docker compose -f docker-compose.prod.yml up -d
```
Database schema migrations apply automatically on backend startup.
