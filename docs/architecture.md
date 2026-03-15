# Architecture Overview

## System Architecture

The Reconciliation Engine follows a **modular monolith with a separated processing worker** pattern — two deployable units communicating through PostgreSQL and Redis.

### Deployable Units

1. **API Service (C# / ASP.NET Core 8)** — HTTP API, domain logic, job orchestration, Hangfire scheduling
2. **Processing Worker (Go 1.22)** — File parsing, matching engine, Redis Stream consumer
3. **Dashboard (React + TypeScript)** — Static SPA served via nginx/CDN

### Core Pipeline

```
[Files] → Stage 1: Ingestion & Normalization (Go)
       → Stage 2: Matching Engine (Go) — Exact → Fuzzy → Aggregate
       → Stage 3: Exception Detection & Classification (C#)
       → Stage 4: Output & Reporting (C#/React)
```

### Bounded Contexts (DDD)

| Context | Responsibility | Language |
|---------|---------------|----------|
| **Ingestion** | File parsing, normalization, deduplication | Go |
| **Reconciliation** | Matching, exception detection, run lifecycle | Go + C# |
| **Fee Intelligence** | Contract management, fee validation, overpayment detection | C# |
| **Reporting** | Dashboards, reports, analytics (read-model) | C# + React |

### Communication Patterns

- **API → Worker**: Redis Streams (job messages)
- **Worker → DB**: Direct PostgreSQL writes
- **Dashboard → API**: REST over HTTPS
- **Domain Events**: In-process event bus (transactional outbox)

### Data Flow

```
Bank/Acquirer Files (OFX, CNAB, CSV, JSON)
        │
        ▼
   [MinIO Storage] ← Raw file archival (SHA-256 for audit)
        │
        ▼
   [Go Parsers] → Normalize to canonical TransactionRecord
        │
        ▼
   [PostgreSQL] ← Persist normalized records (dedup by fingerprint)
        │
        ▼
   [Matching Engine]
   ├── Pass 1: Exact (NSU, E2EID, NossoNumero)
   ├── Pass 2: Fuzzy (amount tolerance, date window, name similarity)
   └── Pass 3: Aggregate (N-to-1 subset-sum)
        │
        ▼
   [Exception Classifier] → Categorize unmatched records
        │
        ▼
   [Fee Validator] → Check actual vs contracted rates
        │
        ▼
   [API/Dashboard] → Results, reports, exception resolution
```

### Infrastructure

| Component | Technology | Purpose |
|-----------|-----------|---------|
| Database | PostgreSQL 16 | ACID-compliant financial data storage |
| Cache | Redis 7 | Job queuing, caching, dedup fingerprints |
| Object Storage | MinIO (dev) / S3 (prod) | Raw file archival |
| Containers | Docker + Compose | Local development |
| CI/CD | GitHub Actions | Automated testing and deployment |
| IaC | Terraform | AWS infrastructure provisioning |
| Monitoring | Prometheus + Grafana | Application metrics |
| Logging | Serilog (C#) + zerolog (Go) | Structured JSON logs |

### Key Design Decisions

- **Integer arithmetic for money** — All amounts in cents (int64). Zero floating-point on financial values.
- **Batch processing** — Daily reconciliation runs, not real-time streaming.
- **No microservices** — Modular monolith is right-sized for single-tenant SMB use case.
- **Polyglot by purpose** — Go for I/O-bound parsing and CPU-bound matching; C# for API and business rules.
- **Idempotent ingestion** — SHA-256 fingerprinting prevents duplicate processing.
