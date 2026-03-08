# Reconciliation Engine

A production-grade automated payment reconciliation system for Brazilian SMBs. Built to handle the full complexity of Brazil's fragmented payment ecosystem — OFX/CNAB bank statements, Pix settlement reports, boleto registries, and card acquirer files from CIELO, Rede, Stone, and others — without cutting corners on correctness, observability, or domain accuracy.

This is a portfolio project built to real production standards. No complexity has been reduced.

---

## The Problem

A mid-size Brazilian retailer processing R$500k/month can receive payments through 5+ channels simultaneously: Pix, card acquirers with different settlement schedules, boletos, TEDs, and recurring débito automático. Each channel produces settlement files in different formats, with different identification keys and fee structures.

Manual reconciliation through spreadsheets leads to missed fee overcharges (acquirers silently billing above contracted MDR is extremely common), undetected chargebacks, delayed detection of Pix devoluções, and general cash flow opacity. A team of 2–3 financial analysts can spend 3–4 hours daily on this alone.

This engine automates that entirely: ingests every settlement file, matches every expected receivable against every actual bank credit, surfaces discrepancies, and produces clean audit-ready outputs.

---

## What It Does

**Multi-source ingestion** — Parses and normalizes OFX 1.x/2.x (SGML and XML), CNAB 240/400 with bank-specific profiles (Itaú, Bradesco, BB, Santander, Caixa, BTG Pactual, Inter, Nubank, Sicredi, Sicoob), Pix settlement reports (JSON/CSV), and card acquirer extratos from CIELO (EEFI/EEVC), Rede, Stone, PagSeguro, Getnet, and SafraPay.

**Three-pass matching engine** — Exact key-based matching (NSU, E2EID, NossoNumero) → fuzzy scored matching (amount tolerance, settlement date windows, Levenshtein name similarity) → aggregate N-to-1 matching for batch-settled acquirer deposits (branch-and-bound subset-sum solver).

**Fee intelligence** — Stores versioned acquirer contracts with MDR rates per bandeira (Visa, Mastercard, Elo, Hipercard, Amex) and produto (crédito à vista, parcelado, débito). Automatically flags when actual fees exceed contracted rates and quantifies the overpayment for dispute evidence.

**Exception management** — Classifies unmatched records into typed exceptions (fee divergence, timing mismatch, partial payment, chargeback, Pix devolução, boleto not compensated) with severity scoring and a resolution workflow with mandatory audit notes.

**Settlement calendar awareness** — Full Anbima banking holiday calendar for correct dias úteis calculations. A D+2 settlement from a Friday before a feriado lands on Wednesday, not Tuesday. Getting this wrong cascades into false-positive timing exceptions.

**Reporting** — Daily reconciliation summaries, fee analysis vs. contracted rates, aging dashboard (0–7d / 8–15d / 16–30d / 30+d buckets), audit trail with immutable match decisions and resolution history. Exports to PDF, Excel, CSV, and JSON.

---

## Stack

| Layer | Technology | Why |
|-------|-----------|-----|
| File parsing & matching | **Go 1.22** | Goroutine-based parallel file processing; low-overhead concurrency for the CPU-bound subset-sum solver |
| API & domain logic | **C# / ASP.NET Core 8** | Mature middleware ecosystem (auth, OpenAPI, Hangfire job scheduling); DDD pattern support |
| Frontend | **React + TypeScript + Tailwind** | Type-safe dashboard with Recharts for fee analysis visualization |
| Database | **PostgreSQL 16** | ACID compliance for financial data; BIGINT amounts (centavos), JSONB metadata, partitioning |
| Cache & queue | **Redis 7** | Redis Streams for API→Worker job dispatch; caching for dedup fingerprints and contract lookups |
| File storage | **MinIO / S3** | Immutable raw file archival with SHA-256 hash for audit trail |
| Infrastructure | **Docker Compose + Terraform** | Local dev via Compose; AWS (ECS Fargate + RDS + ElastiCache) via Terraform |
| CI/CD | **GitHub Actions** | Per-service pipelines: C# (xUnit + Testcontainers), Go (race detector + staticcheck), React (Vitest) |

---

## Architecture

The system runs as two deployable units communicating through PostgreSQL and Redis — no microservices overhead, but clean language separation:

```
Payment Files (OFX, CNAB, Pix, Acquirers)
        │
        ▼
  [Go Worker] — MinIO archival → Parsers → Normalize → 3-pass Matching
        │
        ▼ (direct PostgreSQL writes)
  [PostgreSQL] ←→ [C# API] — Fee Validation, Exception Mgmt, Hangfire Jobs
                       │
                       ▼ (REST)
                 [React Dashboard]
```

The Go worker consumes file parsing and matching jobs from a Redis Stream. The C# API owns the domain logic (DDD aggregates, domain events, specifications), schedules daily reconciliation runs via Hangfire, and serves the dashboard's REST API. The dashboard is a static SPA.

### Bounded Contexts (Domain-Driven Design)

| Context | Responsibility | Language |
|---------|---------------|----------|
| **Ingestion** | Parsing, normalization, SHA-256 deduplication | Go |
| **Reconciliation** | 3-pass matching, run lifecycle, exception detection | Go + C# |
| **Fee Intelligence** | Contract management, MDR validation, overpayment evidence | C# |
| **Reporting** | Read-model projections, aging, audit trail | C# + React |

### Key Design Decisions

**Integer arithmetic for money.** All amounts are stored and computed in centavos as `int64`/`long`. Zero floating-point operations on financial values — a `Money` value object enforces this at the domain boundary. This eliminates the R$0.01 rounding drift that plagues spreadsheet-based reconciliation.

**Idempotent ingestion.** Re-processing the same file produces zero duplicates. SHA-256 fingerprinting on file content and on individual transaction field tuples enforces this at the database level with a unique constraint.

**TDD throughout.** Every line of domain logic is written test-first. Value objects have exhaustive construction and invariant tests. The matching engine is tested against known-answer fixtures. Parser tests run against real anonymized payment files.

**Batch over real-time.** Bank statement files are inherently batch artifacts, generated at specific times each morning. Building a real-time architecture would add complexity without meaningful reconciliation benefit — the authoritative source of truth (the bank statement) is still batch.

---

## Repository Structure

```
reconciliation-engine/
├── services/
│   ├── api/                    # C# solution (Core, Infra, API projects)
│   └── worker/                 # Go module (parsers, matching engine, CLI)
│       └── testdata/           # Anonymized sample files per format
├── dashboard/                  # React + TypeScript + Vite
├── migrations/                 # PostgreSQL schema (shared)
├── deploy/                     # Dockerfiles + docker-compose.yml
├── infra/                      # Terraform modules (AWS)
├── docs/                       # Architecture diagram, OpenAPI spec, file format reference
└── .github/workflows/          # CI pipelines per service
```

---

## Development Phases

| Phase | Scope | Status |
|-------|-------|--------|
| **1 — Foundation** | Monorepo setup, DB schema, Docker Compose, `TransactionRecord` model, OFX parser, CNAB 240 (Itaú) | 🔧 In progress |
| **2 — Matching Engine** | 3-pass matching, reconciliation run orchestration, exception classification | ⏳ Planned |
| **3 — Acquirer Parsers & Fees** | CIELO/Stone/Rede parsers, contract CRUD, fee validation engine | ⏳ Planned |
| **4 — Dashboard & Reports** | React UI, exception management, fee charts, PDF/Excel export | ⏳ Planned |
| **5 — Hardening** | Integration tests at volume, performance tuning, Prometheus/Grafana | ⏳ Planned |

---

## Getting Started

```bash
# Prerequisites: Docker, Go 1.22+, .NET 8 SDK, Node 20+

# 1. Clone and configure environment
git clone https://github.com/hbortolim/reconciliation-engine
cp .env.example .env

# 2. Start infrastructure
docker compose up -d postgres redis minio

# 3. Apply database schema
psql -h localhost -U recon_user -d reconciliation -f migrations/001_initial_schema.sql
psql -h localhost -U recon_user -d reconciliation -f migrations/002_create_indexes_and_views.sql

# 4. Run the services
cd services/api  && dotnet run --project src/Reconciliation.Api
cd services/worker && go run ./cmd/worker
cd dashboard && npm install && npm run dev
```

---

## What This Is Not

This is a reconciliation and financial analysis tool, not an all-in-one finance platform. It does not initiate payments, replace an ERP or accounting system, integrate with Open Finance Brasil APIs, handle nota fiscal compliance, provide credit scoring, or support multi-tenant SaaS deployment. It is deliberately scoped to do one thing — reconcile payment data — and do it correctly.

---

## Domain References

- [Febraban CNAB 240/400 standards](https://www.febraban.org.br/)
- [Banco Central do Brasil — Pix / SPI](https://www.bcb.gov.br/estabilidadefinanceira/pix)
- [Anbima banking holiday calendar](https://www.anbima.com.br/feriados)
- [CIP — Câmara Interbancária de Pagamentos](https://www.cip-bancos.org.br/)
