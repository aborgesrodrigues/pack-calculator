# Pack Calculator

HTTP API to configure pack sizes and calculate how many packs to ship for an order.

## Prerequisites

- [Docker](https://www.docker.com/) and Docker Compose
- [Make](https://www.gnu.org/software/make/)
- A `.env` file in the project root (not committed; see example below)

## Environment

Create `.env` at the repository root:

```env
DB_HOST=db
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=pack_calculator
POSTGRES_PORT=5432
```

`DB_HOST=db` matches the Docker Compose service name when using `make start`.

Local development assets live under `dev/` (Docker Compose, app Dockerfile, database init script, and the web UI in `dev/ui/`).

## Run locally (Makefile)

| Command | Description |
|---------|-------------|
| `make support` | Start PostgreSQL only (port `5432`) |
| `make start` | Build and start the app, UI, and database (API on `8080`, UI on `8081`) |
| `make stop` | Stop and remove containers |
| `make test` | Run tests with the race detector |
| `make swagger` | Regenerate OpenAPI docs under `docs/` |
| `make build` | Build the app Docker image only |

Typical flow:

```bash
make start
```

API base URL: `http://localhost:8080`

Web UI: `http://localhost:8081` (`dev/ui/` — forms for `/pack_size/batch` and `/calculate`)

Swagger UI: `http://localhost:8080/swagger`

Stop when finished:

```bash
make stop
```

## API examples (curl)

List configured pack sizes:

```bash
curl -s http://localhost:8080/pack_size/batch
```

Configure pack sizes (replaces any existing sizes):

```bash
curl -s -X POST http://localhost:8080/pack_size/batch \
  -H 'Content-Type: application/json' \
  -d '{"sizes":[5000,2000,1000,500,250]}'
```

Calculate fulfilment for an order (returns the order with a generated `event_id` and the pack mix in `packs`):

```bash
curl -s -X POST http://localhost:8080/calculate \
  -H 'Content-Type: application/json' \
  -d '{"items":12001}'
```

Example response (`packs` maps each pack size to how many of that pack to ship):

```json
{
  "event_id": "550e8400-e29b-41d4-a716-446655440000",
  "items": 12001,
  "packs": {
    "5000": 2,
    "2000": 1,
    "250": 1
  }
}
```

Invalid request body (400):

```bash
curl -s -X POST http://localhost:8080/calculate \
  -H 'Content-Type: application/json' \
  -d 'not-json'
```
