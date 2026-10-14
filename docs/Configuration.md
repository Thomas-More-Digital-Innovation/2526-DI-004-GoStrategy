# Configuration Reference

All GoStrategy services are configured via environment variables. Copy `.env.template` to `.env` before starting services.

---

## Environment Variables

### Global Settings

| Variable | Default | Description |
| :--- | :--- | :--- |
| `APP_ENV` | `development` | Runtime mode: `development` or `production`. Enables secure cookie headers in production. |

### Backend Configuration

| Variable | Default | Description |
| :--- | :--- | :--- |
| `BACKEND_PORT` | `8080` | Port on which the Go HTTP and WebSocket server listens. |
| `ALLOWED_ORIGINS` | `*` | Comma-separated list of allowed CORS origins (e.g., `https://stratego.example.com`). |
| `JWT_SECRET` | *required* | High-entropy secret used to sign and verify user session JWT tokens. Generate with `openssl rand -base64 32`. |
| `TELEMETRY_ALLOWED_IPS` | `""` | Comma-separated whitelist of client IP addresses permitted to access Prometheus telemetry metrics. |

### Database Settings

| Variable | Default | Description |
| :--- | :--- | :--- |
| `DB_HOST` | `gostrategy` | Hostname or IP of the PostgreSQL database server. |
| `DB_PORT` | `5432` | Port for the PostgreSQL service. |
| `DB_USER` | `gostrategy` | Database username. |
| `DB_PASSWORD` | `password` | Database password. |
| `DB_NAME` | `gostrategy` | Target PostgreSQL database name. |

### Frontend Build Variables

These variables are consumed during frontend compilation (`bun run build` or Vite dev):

| Variable | Default | Description |
| :--- | :--- | :--- |
| `FRONTEND_PORT` | `5000` | Host port mapped to the web interface in development. |
| `VITE_API_BASE` | `http://localhost:8080` | Base URL used by the browser to execute REST API requests against the Go backend. |
| `VITE_WS_BASE` | `ws://localhost:8080` | WebSocket endpoint URL used by the browser to establish game socket connections. |

---

## Database Migrations

Database schema migrations are embedded inside the Go binary (`internal/db/migrations/`) and execute automatically upon application startup. No manual migration scripts are required.
