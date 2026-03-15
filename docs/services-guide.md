# Services Guide

A living reference for running, configuring, and understanding each service in the Reconciliation Engine. Update this document as new binaries, flags, or environment variables are added.

---

## Go Worker (`services/worker/`)

The worker is a Go module with two separate entry points (binaries) under `cmd/`. They share the same internal packages but serve different purposes.

### `cmd/worker` — Production Service

Long-running process that consumes reconciliation jobs from a Redis Stream, parses payment files, runs the matching engine, and writes results to PostgreSQL.

**Run locally:**

```bash
cd services/worker
go run ./cmd/worker
```

**Run via Docker Compose:**

```bash
docker compose up worker
# Requires postgres, redis, and minio to be healthy first
docker compose up -d postgres redis minio
docker compose up worker
```

**Environment variables:**

| Variable | Default | Description |
|----------|---------|-------------|
| `POSTGRES_HOST` | — | PostgreSQL host |
| `POSTGRES_PORT` | `5432` | PostgreSQL port |
| `POSTGRES_DB` | `reconciliation` | Database name |
| `POSTGRES_USER` | `recon_user` | Database user |
| `POSTGRES_PASSWORD` | — | Database password |
| `REDIS_HOST` | — | Redis host |
| `REDIS_PORT` | `6379` | Redis port |
| `REDIS_PASSWORD` | — | Redis password (optional) |
| `MINIO_ENDPOINT` | — | MinIO/S3 endpoint |
| `MINIO_ACCESS_KEY` | `minioadmin` | MinIO access key |
| `MINIO_SECRET_KEY` | `minioadmin` | MinIO secret key |
| `MINIO_BUCKET` | `reconciliation-files` | Target bucket for raw file archival |
| `MINIO_USE_SSL` | `false` | Use TLS for MinIO connection |
| `WORKER_CONCURRENCY` | `4` | Number of concurrent job processors |
| `WORKER_REDIS_STREAM` | `reconciliation-jobs` | Redis Stream key to consume from |
| `LOG_LEVEL` | `debug` | Logging level (debug, info, warn, error) |

**Current status:** Redis Stream polling and PostgreSQL persistence are scaffolded but stubbed. The goroutine structure and graceful shutdown (SIGINT/SIGTERM) are in place.

---

### `cmd/cli` — Local Dev Tool

Command-line tool for running reconciliation manually against local files. Useful for testing parsers and the matching engine without spinning up Redis, PostgreSQL, or MinIO.

**Run:**

```bash
cd services/worker
go run ./cmd/cli --file <path> --source-type <type> [--run-id <id>] [--debug]
```

**Flags:**

| Flag | Required | Description |
|------|----------|-------------|
| `--file` | Yes | Path to the input payment file |
| `--source-type` | Yes | Transaction source type (see valid values below) |
| `--run-id` | No | Unique identifier for this reconciliation run |
| `--debug` | No | Enable verbose debug logging |

**Valid `--source-type` values:**

| Value | Description |
|-------|-------------|
| `PIX` | Pix settlement reports (JSON/CSV) |
| `BOLETO` | Boleto registry files |
| `CARD_CREDIT` | Card acquirer credit files (CIELO, Rede, Stone, etc.) |
| `CARD_DEBIT` | Card acquirer debit files |
| `OFX` | OFX 1.x/2.x bank statements |
| `CNAB240` | CNAB 240 bank files |
| `CNAB400` | CNAB 400 bank files |

**Example:**

```bash
go run ./cmd/cli \
  --file testdata/ofx/itau_sample.ofx \
  --source-type OFX \
  --run-id debug-run-001 \
  --debug
```

**Output:**

```
=== Reconciliation Results ===
Matched pairs: 0
Unmatched expected: 0
Unmatched actual: 0
```

**Note:** This binary has no Docker entry — it is local-only. Parser selection is currently stubbed; the CLI structure and exact matching output are functional.

---

## C# API (`services/api/`)

The ASP.NET Core 8 API owns domain logic (DDD aggregates, domain events), Hangfire job scheduling, fee validation, and the REST API consumed by the dashboard.

**Run locally:**

```bash
cd services/api
dotnet run --project src/Reconciliation.Api
```

**Run tests:**

```bash
cd services/api
dotnet test
```

**Key references:**
- OpenAPI spec: [`docs/api-spec.yaml`](api-spec.yaml)
- Domain value objects: `services/api/src/Reconciliation.Core/Domain/ValueObjects/`
- Repository interfaces: `services/api/src/Reconciliation.Core/Interfaces/`

---

## React Dashboard (`dashboard/`)

Static SPA that consumes the C# API. Built with Vite, React, TypeScript, and Tailwind CSS.

**Run locally:**

```bash
cd dashboard
npm install
npm run dev
```

**Run tests:**

```bash
cd dashboard
npx vitest run
```

---

## Infrastructure

### Start only the infrastructure dependencies

```bash
docker compose up -d postgres redis minio
```

### Start all services

```bash
docker compose up
```

### Apply database migrations

```bash
psql -h localhost -U recon_user -d reconciliation -f migrations/001_initial_schema.sql
psql -h localhost -U recon_user -d reconciliation -f migrations/002_create_indexes_and_views.sql
```

### Prerequisites

- Docker
- Go 1.22+
- .NET 8 SDK
- Node 20+
