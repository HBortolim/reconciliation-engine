# Reconciliation Engine — Claude Code Instructions

## Spec-Driven Development (SDD) Protocol

This project follows **Spec-Driven Development at Level 2 (Spec-Anchored)**. The specification is a living document that guides all implementation and is maintained as the codebase evolves.

### The Golden Rule

**SPEC.md is the source of truth.** Every implementation decision must trace back to it. If you're unsure about requirements, architecture, data model, or behavior — read `SPEC.md` first. If the spec doesn't cover it, ask the developer before guessing.

### SDD Workflow

#### Before Writing Any Code

1. **Read the relevant section of SPEC.md** — Identify which bounded context, feature, or component the task relates to.
2. **Check the development phase** — The spec defines 5 phases. Work within the current phase's scope. Don't jump ahead.
3. **Write or update the spec first if needed** — If you discover a gap in the spec (missing edge case, undefined behavior, ambiguous requirement), propose a spec update before implementing.

#### Implementation Cycle

1. **Spec → Test → Code → Validate → Update Spec (if needed)**
2. Write the failing test first (TDD is mandatory per SPEC.md).
3. Implement the minimum code to pass the test.
4. Refactor.
5. If the implementation revealed something the spec didn't anticipate, propose a spec update.

#### After Implementation

- Verify the implementation matches the spec exactly.
- If you deviated from the spec (with good reason), update the spec to reflect reality.
- Never leave the spec and code in disagreement.

---

## Project Architecture

### Monorepo Structure

```
services/api/       → C# ASP.NET Core 8 (API, domain logic, Hangfire)
services/worker/    → Go 1.22 (parsers, matching engine, Redis consumer)
dashboard/          → React + TypeScript + Tailwind (Vite)
migrations/         → PostgreSQL schema (shared)
services/worker/testdata/
deploy/             → Docker, docker-compose
infra/              → Terraform (AWS)
docs/               → Architecture, file formats, OpenAPI spec
```

### Bounded Contexts (DDD)

| Context | Location | Language | Responsibility |
|---------|----------|----------|---------------|
| Ingestion | `services/worker/internal/parsers/` | Go | File parsing, normalization, dedup |
| Reconciliation | `services/worker/internal/matching/` + `services/api/src/Reconciliation.Core/` | Go + C# | 3-pass matching, run lifecycle |
| Fee Intelligence | `services/api/src/Reconciliation.Core/` | C# | Contract mgmt, fee validation |
| Reporting | `services/api/` + `dashboard/` | C# + React | Dashboards, reports, exports |

### Communication

- **API → Worker**: Redis Streams (job messages)
- **Worker → DB**: Direct PostgreSQL writes
- **Domain Events**: In-process event bus (transactional outbox pattern)
- **Dashboard → API**: REST over HTTPS

---

## Development Rules

### Non-Negotiable Constraints

- **Integer arithmetic for money.** All monetary values in centavos (`int64`/`long`). Zero floating-point operations on financial data. Use the `Money` value object in C# and `CentavosFromReais()`/`int64` in Go.
- **Idempotent ingestion.** SHA-256 fingerprinting on file content + transaction-level hashing. Re-ingesting the same file must produce zero duplicates.
- **TDD is mandatory.** Every line of domain logic is written test-first. Red → Green → Refactor. No exceptions.
- **Structured JSON logging.** Serilog for C#, zerolog for Go. Correlation IDs across the pipeline.
- **No floating-point for fees.** Fee rates stored as basis points (integer). `FeeRate` value object handles conversion.

### Coding Standards

#### Go (services/worker/)

- Go module: `github.com/hbortolim/reconciliation-engine`
- Use `internal/` package convention (already set up).
- Test files: `*_test.go` next to source (idiomatic Go).
- Table-driven tests for parser variations.
- `testify/assert` and `testify/suite` for test utilities.
- `zerolog` for structured logging.
- Error handling: explicit returns, no panics in library code.

#### C# (services/api/)

- .NET 8, C# 12.
- Solution: `services/api/ReconciliationEngine.sln`
- Namespaces follow folder structure: `ReconciliationEngine.Core.Domain.ValueObjects`, etc.
- xUnit + FluentAssertions for tests.
- Moq for interface mocking in service tests.
- Testcontainers for PostgreSQL integration tests.
- Value Objects are immutable records. Entities have identity. Aggregates enforce consistency boundaries.
- Repository interfaces in `Reconciliation.Core/Interfaces/`. Implementations in `Reconciliation.Infra/`.

#### React/TypeScript (dashboard/)

- Vite + React + TypeScript strict mode.
- Tailwind CSS for styling.
- Vitest for testing.
- All monetary values displayed via `MoneyDisplay` component (centavos → formatted BRL).

### Database

- PostgreSQL 16. Schema in `migrations/`.
- All amounts: `BIGINT` (centavos), never `DECIMAL` or `FLOAT`.
- UUIDs for primary keys (`uuid_generate_v4()`).
- Domain events: transactional outbox table (`domain_events_outbox`).

### Brazilian Payment Domain Knowledge

When implementing parsers or matching logic, keep in mind:

- **OFX 1.x is SGML, not XML.** Standard XML parsers fail on it.
- **CNAB 240/400 have bank-specific variations** despite being "standards." Always check bank profiles in `services/worker/internal/parsers/cnab/profiles/`.
- **Settlement dates use dias úteis (business days only).** The banking holiday calendar (`banking_holidays` table) must be consulted for any date calculation. Anbima publishes the official list.
- **Card acquirer settlement windows vary.** CIELO: D+1 débito, D+30 crédito. Stone may differ. Always check `AcquirerContract` for the merchant's specific terms.
- **Pix E2EID format:** `E{ISPB_8chars}{YYYYMMDD}{SEQUENCE_11chars}` — validate this.
- **Encoding chaos:** Bank files mix UTF-8, ISO-8859-1, Windows-1252. Always handle encoding detection.
- **Amounts:** Some files use centavos (integer), others reais (decimal with comma separator). Normalize to centavos immediately on parse.

---

## Development Phases

Follow these phases sequentially. Each phase builds on the previous one.

| Phase | Focus | Weeks |
|-------|-------|-------|
| **1 — Foundation** | Monorepo setup, DB schema, Docker Compose, TransactionRecord model, OFX parser, basic CNAB 240 (Itaú) | 1–3 |
| **2 — Matching Engine** | 3-pass matching (exact → fuzzy → aggregate), reconciliation run orchestration, exception classification | 4–6 |
| **3 — Acquirer Parsers & Fees** | CIELO/Stone/Rede parsers, acquirer contract CRUD, fee validation engine | 7–9 |
| **4 — Dashboard & Reports** | React dashboard, exception mgmt UI, fee analysis charts, PDF/Excel export | 10–12 |
| **5 — Hardening** | Integration tests, performance optimization, Prometheus/Grafana, docs | 13–14 |

---

## Commands Reference

### Local Development

```bash
# Start all infrastructure
docker compose up -d postgres redis minio

# Run C# API
cd services/api && dotnet run --project src/Reconciliation.Api

# Run Go worker
cd services/worker && go run ./cmd/worker

# Run dashboard
cd dashboard && npm run dev

# Run all tests
cd services/api && dotnet test
cd services/worker && go test ./...
cd dashboard && npx vitest run
```

### Database

```bash
# Apply migrations (via psql)
psql -h localhost -U recon_user -d reconciliation -f migrations/001_initial_schema.sql
psql -h localhost -U recon_user -d reconciliation -f migrations/002_create_indexes_and_views.sql
```

---

## Key Files Quick Reference

| What | Where |
|------|-------|
| Full specification | `SPEC.md` |
| Architecture diagram | `docs/architecture-overview.excalidraw.png` |
| OpenAPI spec | `docs/api-spec.yaml` |
| DB schema | `migrations/001_initial_schema.sql` |
| DB views | `migrations/002_create_indexes_and_views.sql` |
| Domain value objects | `services/api/src/Reconciliation.Core/Domain/ValueObjects/` |
| Domain entities | `services/api/src/Reconciliation.Core/Domain/Entities/` |
| Aggregate roots | `services/api/src/Reconciliation.Core/Domain/Aggregates/` |
| Repository interfaces | `services/api/src/Reconciliation.Core/Interfaces/` |
| Go parsers | `services/worker/internal/parsers/` |
| Parser test fixtures | `services/worker/testdata/` |
| Matching engine | `services/worker/internal/matching/` |
| Bank profiles (CNAB) | `services/worker/internal/parsers/cnab/profiles/` |
| Docker setup | `docker-compose.yml` + `deploy/docker/` |
| CI pipelines | `.github/workflows/` |
| Terraform | `infra/` |
