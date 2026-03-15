# Automated Payment Reconciliation Engine for SMBs

## Project Overview

A production-grade reconciliation system designed to automatically match, classify, and reconcile financial transactions across multiple payment sources commonly used by Brazilian SMBs (Pequenas e Médias Empresas). The engine ingests data from bank statements (OFX/CNAB), Pix settlement reports, boleto registries, card acquirer files (CIELO, Rede, Stone, PagSeguro), and ERP/accounting entries — then applies rule-based and fuzzy matching algorithms to identify discrepancies, flag exceptions, and produce audit-ready reconciliation reports.

This is a portfolio project built to production standards. No complexity has been reduced. Every component is designed as if serving real customers in the Brazilian financial ecosystem.

---

## Problem Statement

Brazilian SMBs operate in one of the most fragmented payment ecosystems in the world. A single small retailer might receive payments through Pix (instant payments via BCB's SPI), boletos bancários (bank slips), multiple card acquirers with different settlement schedules, TEDs, and even recurring charges via débito automático. Each of these channels produces settlement files in different formats, with different identification keys, settlement windows, and fee structures.

The reconciliation burden is enormous: a mid-size PME processing R$500k/month might deal with 4–6 different payment sources, each with its own reporting cadence. Manual reconciliation through spreadsheets leads to missed fee overcharges (acquirers charging above the contracted MDR is extremely common), undetected chargebacks, delayed detection of Pix devolutions, and general cash flow opacity.

This engine exists to automate what a team of 2–3 financial analysts would do manually: match every expected receivable against every actual bank credit, surface discrepancies, and produce clean reconciliation outputs.

---

## Deep Dive: Architecture & Design

### Core Reconciliation Pipeline

The system operates as a multi-stage pipeline with clear separation of concerns:

**Stage 1 — Ingestion & Normalization**

Raw files from heterogeneous sources are parsed into a unified internal transaction model. This is the hardest part of the system because Brazilian payment file formats are notoriously inconsistent. CNAB 240 and CNAB 400 (the Febraban standards for bank remittance/return files) have bank-specific variations that break parsers. OFX files from different banks embed metadata differently. Card acquirer APIs each return settlement data in proprietary schemas.

The normalizer produces a canonical `TransactionRecord` with fields like: source identifier, amount (in cents, always integer arithmetic to avoid floating-point), timestamp (UTC-normalized from BRT), counterparty identifiers (CNPJ/CPF, Pix key, NSU for card transactions), fee breakdown, expected settlement date, and a hash fingerprint for deduplication.

**Stage 2 — Matching Engine**

The core algorithm operates in three passes:

1. **Exact match**: Direct key-based matching using deterministic identifiers. For card transactions, this means matching the acquirer's NSU (Número Sequencial Único) against the bank statement entry. For Pix, matching the EndToEndId (E2EID) from the SPI against the bank credit. For boletos, matching the nosso número or linha digitável.

2. **Fuzzy match**: When exact keys are unavailable or corrupted (which happens more often than anyone admits), the engine falls back to a scoring-based approach. It computes similarity across amount (with configurable tolerance for rounding and fee deduction), date proximity (accounting for D+1, D+2, D+30 settlement windows per acquirer), and counterparty name similarity (Levenshtein distance on razão social/nome fantasia).

3. **Aggregate match**: Some payment sources settle in batches. Card acquirers, for example, often deposit a single lump sum representing multiple transactions. The engine attempts N-to-1 matching by finding subsets of expected receivables whose sum matches a single bank deposit (subset-sum with pruning heuristics for performance).

**Stage 3 — Exception Detection & Classification**

Unmatched records after all passes are classified into exception categories: fee divergence (acquirer charged more than contracted MDR/taxa de antecipação), timing mismatch (settlement arrived outside expected window), partial payment, duplicate transaction, chargeback/contestação, Pix devolução, and unknown. Each exception gets a severity score and suggested resolution action.

**Stage 4 — Output & Reporting**

Reconciliation results are persisted and exposed through both an API and generated reports. Reports include: matched transactions with confidence scores, exception details with drill-down, fee analysis (actual vs. contracted rates per acquirer), cash flow projection adjustments, and aging of unreconciled items.

### Data Model

The internal data model is designed around the realities of Brazilian payments:

- **TransactionRecord**: The canonical normalized record. Every ingested transaction becomes one of these regardless of source. Key fields include `source_type` (enum: PIX, BOLETO, CARD_CREDIT, CARD_DEBIT, TED, DOC, DEBITO_AUTOMATICO), `amount_cents` (i64, always positive), `fee_cents` (i64), `net_amount_cents` (i64), `expected_settlement_date`, `actual_settlement_date`, `counterparty_document` (CPF/CNPJ), `external_id` (source-specific unique key), and `fingerprint_hash` (SHA-256 of key fields for dedup).

- **ReconciliationPair**: Links two TransactionRecords (expected vs. actual) with a `match_type` (EXACT, FUZZY, AGGREGATE), `confidence_score` (0.0–1.0), and `discrepancy_details` (optional struct with amount delta, date delta, fee delta).

- **Exception**: An unmatched or problematic record with `exception_type`, `severity` (LOW, MEDIUM, HIGH, CRITICAL), `suggested_action`, and `resolution_status` (OPEN, IN_REVIEW, RESOLVED, IGNORED).

- **AcquirerContract**: Stores contracted commercial conditions per acquirer — MDR per bandeira (Visa, Mastercard, Elo, Hipercard, Amex), antecipação rates, settlement schedule (D+n), monthly fee caps. Used to validate actual fees against contracted rates.

- **ReconciliationRun**: Metadata for each reconciliation execution — timestamp, files ingested, match statistics, exception counts, duration.

### File Parser Architecture

Each payment source gets a dedicated parser module implementing a common `Parser` trait/interface. This is intentionally not abstracted into a generic parser because the formats are different enough that forced abstraction creates more problems than it solves.

Supported formats:

- **OFX (Open Financial Exchange)**: Bank statements. Parsed with a streaming XML/SGML parser (OFX 1.x is SGML, not XML — a common gotcha). Extracts STMTTRN entries.
- **CNAB 240/400**: Febraban standard for bank remittance returns. Fixed-width positional files with header/trailer records. Bank-specific field variations are handled through configuration profiles (one per bank: Itaú, Bradesco, BB, Santander, Caixa, BTG Pactual, Inter, etc.).
- **Pix Settlement (DICT/SPI reports)**: JSON or CSV exports from PSP dashboards. Key fields: E2EID, valor, data liquidação, chave Pix.
- **Card Acquirer Extratos**: CIELO (EEFI/EEVC formats), Rede (CSV), Stone (API JSON), PagSeguro (CSV), Getnet (positional), SafraPay. Each has its own settlement report layout with NSU, valor bruto, valor líquido, taxa, bandeira, parcela, previsão de pagamento.
- **NF-e/NFC-e XML** (optional): Cross-reference against fiscal documents for additional validation.

### Settlement Calendar Awareness

The engine maintains a calendar module that understands Brazilian banking holidays (feriados bancários), which differ from civil holidays. Settlement date calculations account for dias úteis only. The calendar is pre-loaded with Anbima's official holiday list and can be updated annually. This is critical because a D+2 settlement from a Friday would land on Tuesday (skipping the weekend), but if Monday is a feriado, it lands on Wednesday. Getting this wrong cascades into false-positive timing mismatches.

---

## Tech Stack

### Core Application

| Layer | Technology | Rationale |
|---|---|---|
| **Language** | Go (Golang) 1.22+ | High concurrency for parallel file parsing, excellent performance for batch processing, strong standard library for file I/O and HTTP. Aligns with your current learning path. |
| **Secondary Language** | C# (.NET 8) | Used for the API layer and business logic orchestration. Leverages your existing expertise. Acts as the "command center" that coordinates Go-based processing workers. |
| **Database** | PostgreSQL 16 | ACID compliance is non-negotiable for financial data. JSONB columns for flexible metadata storage. Partitioning by reconciliation_run_date for query performance at scale. |
| **Cache** | Redis 7 | Caching parsed file results, deduplication fingerprints, and acquirer contract lookups. Also used as the message broker for async job coordination. |
| **Message Queue** | Redis Streams (or NATS if scaling beyond single-node) | Decouples file ingestion from processing. Enables retry logic for failed parsing jobs without blocking the pipeline. |
| **File Storage** | MinIO (S3-compatible) | Raw file archival with immutable retention. Every ingested file is stored with SHA-256 hash for audit trail. Can be swapped for AWS S3 in cloud deployment. |

### Processing & Matching

| Component | Technology | Rationale |
|---|---|---|
| **CNAB Parser** | Go — custom positional parser | No reliable open-source CNAB parser handles all bank variations. Built from Febraban specs + bank-specific documentation. |
| **OFX Parser** | Go — custom SGML/XML streaming parser | Must handle both OFX 1.x (SGML) and 2.x (XML). Existing Go libraries are incomplete. |
| **Fuzzy Matching** | Go — custom scoring engine | Weighted multi-field similarity with configurable thresholds. Levenshtein for string fields, windowed comparison for dates, tolerance-based for amounts. |
| **Subset-Sum Solver** | Go — branch-and-bound with pruning | For N-to-1 aggregate matching. Pruned search with early termination. Max subset size capped at 50 for practical performance. |
| **Fee Validation** | C# — rule engine | Compares actual fees against AcquirerContract terms. Handles tiered MDR (different rates per bandeira/produto), antecipação calculations, and monthly volume-based discounts. |

### API & Interface

| Component | Technology | Rationale |
|---|---|---|
| **REST API** | C# — ASP.NET Core Minimal APIs | Clean endpoints for triggering reconciliation runs, querying results, managing acquirer contracts, and exception handling workflows. |
| **Background Jobs** | Hangfire (C#) or Go worker goroutines | Scheduled reconciliation runs (daily at 7:00 BRT after bank file availability), retry logic, report generation. |
| **Dashboard** | React + TypeScript + Recharts | Reconciliation overview, exception drill-down, fee analysis charts, aging reports. Tailwind CSS for styling. |
| **Authentication** | JWT + refresh tokens | Role-based access: admin (full config), analyst (view + resolve exceptions), viewer (read-only dashboards). |

### Infrastructure & Observability

| Component | Technology | Rationale |
|---|---|---|
| **Containerization** | Docker + Docker Compose | Local development and single-machine deployment. Dockerfile per service (API, workers, dashboard). |
| **CI/CD** | GitHub Actions | Automated tests, linting, build, and container image push. |
| **Logging** | Structured JSON logs (Serilog for C#, zerolog for Go) | Correlation IDs across the pipeline. Every transaction touch is logged for audit. |
| **Metrics** | Prometheus + Grafana | Reconciliation success rates, processing latency per file type, exception counts by category, queue depth. |
| **Testing** | xUnit (C#), Go testing + testify | Unit tests for every parser against real anonymized file samples. Integration tests for the full pipeline with known-answer test sets. |

---

## Features

### Core Reconciliation

- **Multi-source ingestion**: Parse and normalize OFX, CNAB 240/400, Pix reports, card acquirer files (CIELO, Rede, Stone, PagSeguro, Getnet, SafraPay), and generic CSV/Excel inputs.
- **Three-pass matching**: Exact key match → fuzzy scored match → aggregate subset-sum match. Configurable confidence thresholds per match type.
- **Bank-specific CNAB profiles**: Configuration-driven parser variations for Itaú, Bradesco, Banco do Brasil, Santander, Caixa, BTG Pactual, Banco Inter, Nubank, Sicredi, and Sicoob.
- **Integer arithmetic throughout**: All monetary calculations in cents (i64). Zero floating-point operations on financial values. This prevents the classic R$0.01 rounding drift that plagues spreadsheet-based reconciliation.
- **Idempotent processing**: Re-ingesting the same file produces no duplicates. SHA-256 fingerprinting on file content + individual transaction hashing.

### Fee Intelligence

- **Acquirer contract management**: Store and version contracted MDR rates per bandeira (Visa, Mastercard, Elo, Hipercard, Amex), per produto (crédito à vista, parcelado lojista, débito), and per acquirer.
- **Fee divergence detection**: Automatically flag when the actual fee charged by the acquirer exceeds the contracted rate. Calculates the overpayment in R$ and aggregates monthly for dispute evidence.
- **Antecipação validation**: For merchants using receivables anticipation (antecipação de recebíveis), validates that the antecipação rate applied matches the contracted terms, including checking against the CDI-based pricing if applicable.
- **Split payment tracking**: For marketplaces and franquias operating under split payment arrangements, tracks that the split percentages applied by the subadquirente match the configured rules.

### Exception Management

- **Automated classification**: Unmatched records are classified into: fee divergence, timing mismatch, partial payment, duplicate, chargeback/contestação, Pix devolução, boleto not compensated, unknown.
- **Severity scoring**: Exceptions are scored based on amount (higher value = higher severity), age (older unresolved = escalating severity), and pattern (recurring exceptions from same source get amplified).
- **Resolution workflow**: Exceptions move through OPEN → IN_REVIEW → RESOLVED/IGNORED states. Resolution requires a justification note for audit trail.
- **Aging dashboard**: Unreconciled items tracked by age bucket (0–7 days, 8–15 days, 16–30 days, 30+ days) with drill-down.

### Brazilian Payment Ecosystem Specifics

- **Pix ecosystem awareness**: Handles Pix instant credits, scheduled Pix (Pix Agendado), Pix Cobrança (with txid matching), Pix Devolução (partial and full), and Pix Troco/Saque edge cases. Understands that Pix settlements hit the account in real-time but may appear in bank statements with a lag.
- **Boleto lifecycle tracking**: Tracks boletos from registration (remessa) through compensation (retorno) including partial payments (pagamento parcial), late payments with multa/juros calculations, and DDA (Débito Direto Autorizado) accelerated clearing.
- **Card settlement calendar**: Maintains per-acquirer settlement schedules. Knows that CIELO settles D+1 for débito and D+30 for crédito (or D+2 with antecipação), while Stone might have different terms. Accounts for bandeira-specific settlement rules.
- **Registro de recebíveis (receivables registry)**: Awareness of the CIP/B3 interoperável registry and how receivable unit registration affects settlement flows, particularly when the merchant has used recebíveis as collateral for credit lines.
- **Banking holiday calendar**: Full Anbima feriado bancário calendar with automatic yearly updates. Correctly calculates dias úteis for settlement date expectations.

### Reporting & Analytics

- **Daily reconciliation summary**: Matched count/amount, exception count/amount, reconciliation rate percentage, processing time.
- **Fee analysis report**: Monthly aggregation of actual vs. contracted fees per acquirer, per bandeira. Highlights overpayment opportunities.
- **Cash flow reconciliation**: Expected vs. actual cash inflows by day, with variance explanations linked to specific exceptions.
- **Audit trail**: Immutable log of every reconciliation run, every match decision (with confidence scores), and every exception resolution. Exportable for external audit.
- **Export formats**: PDF reports, Excel spreadsheets (for the accountant who will inevitably ask), CSV for ERP integration, and JSON via API.

---

## Expected Results

### Quantitative Targets

- **Reconciliation rate ≥ 95%**: Of all transactions ingested, at least 95% should be automatically matched without human intervention. The remaining 5% are genuine exceptions requiring analyst review.
- **Processing throughput**: Handle 50,000 transactions per reconciliation run in under 5 minutes on a single 4-core machine. This covers a mid-size PME doing ~R$2M/month across all payment channels.
- **Fee divergence detection accuracy ≥ 99%**: When an acquirer overcharges relative to contracted rates, the system should catch it. False positive rate for fee alerts below 2%.
- **False match rate < 0.1%**: Incorrectly matched transactions should be extremely rare. The fuzzy matching threshold is tuned to prioritize precision over recall — it is better to leave a transaction unmatched (for human review) than to match it incorrectly.
- **Duplicate detection rate ≥ 99.5%**: Re-ingested files or duplicate transactions across sources should be caught by fingerprint hashing.

### Qualitative Outcomes

- A financial analyst who currently spends 3–4 hours daily on manual reconciliation should be able to reduce that to 30 minutes of exception review.
- The fee analysis module should surface at least one actionable overpayment finding per quarter for a typical SMB with 3+ acquirers (this is based on how frequently acquirers silently adjust rates).
- The system should serve as a single source of truth for "did this payment arrive, and did we receive the correct amount" — replacing the spreadsheet jungle that most PMEs operate with.
- The aging dashboard should make it impossible for an unreconciled transaction to "fall through the cracks" beyond 30 days without explicit acknowledgment.

### Portfolio Demonstration Value

- Demonstrates deep understanding of Brazilian financial infrastructure (SPI, SPB, CIP, Febraban standards).
- Shows ability to design and implement a multi-stage data pipeline with complex matching logic.
- Proves capability in handling real-world messy data (inconsistent file formats, bank-specific quirks, encoding issues).
- Exhibits production-grade concerns: audit trails, idempotency, integer arithmetic for money, observability, and proper error handling.
- Showcases polyglot architecture (Go + C#) with clear reasoning for technology choices.

---

## What This Project Will NOT Do

### Explicitly Out of Scope

- **It is not an ERP or accounting system.** It does not manage chart of accounts, produce balance sheets, or handle tax obligations. It reconciles payment data and feeds results to whatever ERP the business uses (TOTVS Protheus, Omie, Bling, Conta Azul, etc.) via exports or API.

- **It does not initiate payments.** The engine is read-only with respect to bank accounts. It never sends PIX, registers boletos, or triggers TEDs. It only ingests and analyzes settlement data after the fact.

- **It does not integrate with Open Finance Brasil APIs directly.** While the system could theoretically pull bank statements via Open Finance (fase 3 — iniciação de transação de pagamento is out of scope anyway), this project uses file-based ingestion. Open Finance integration would require becoming a participante regulado or partnering with one, which is a regulatory undertaking beyond portfolio scope.

- **It does not handle nota fiscal reconciliation end-to-end.** While it can optionally cross-reference NF-e/NFC-e XMLs for additional validation, it is not a fiscal compliance tool. It does not validate ICMS/ISS/PIS/COFINS calculations or interact with SEFAZ.

- **It does not replace the function of a registradora de recebíveis.** It is aware of how receivable registration affects settlement flows, but it does not interact with CIP, B3 (TAG), or CERC registries directly.

- **It does not provide credit scoring or risk analysis.** While the data could theoretically feed a credit model (reconciliation quality as a proxy for business health), this is not a lending or credit product.

- **It does not handle cryptocurrency or digital asset reconciliation.** Only BRL-denominated transactions through regulated Brazilian payment channels.

- **It does not process real-time streaming transactions.** The engine operates in batch mode (typically daily). While Pix is real-time, the reconciliation of Pix transactions happens in batch against the bank statement, not as a real-time event stream.

- **It does not manage user onboarding or KYC.** There is no CNPJ validation against Receita Federal, no document upload for compliance, no anti-fraud screening. It assumes the operator is already authenticated and authorized.

- **It will not automatically dispute fee overcharges with acquirers.** It surfaces the evidence and quantifies the overpayment, but the actual dispute process (opening a ticket with CIELO, calling Rede's backoffice, etc.) remains manual.

- **No mobile application.** Dashboard is web-only (responsive, but not a native app).

- **No multi-tenant SaaS infrastructure.** The project runs as a single-tenant deployment. Adding tenant isolation, usage metering, billing, and tenant provisioning would be a separate project.

---

## Project Structure

The monorepo is organized by **deployable unit**, with clear language boundaries. Each service is self-contained with its own build system, dependency management, and test suite. Shared resources (migrations, docs) live at the root level.

```
reconciliation-engine/
├── services/
│   ├── api/                              # C# ASP.NET Core — API, domain, infrastructure
│   │   ├── ReconciliationEngine.sln
│   │   ├── src/
│   │   │   ├── Reconciliation.Api/       # REST API, auth, job orchestration (Hangfire)
│   │   │   ├── Reconciliation.Core/      # Domain models, value objects, services, specs
│   │   │   └── Reconciliation.Infra/     # Database (EF Core), Redis, MinIO adapters
│   │   └── tests/
│   │       ├── Reconciliation.Core.Tests/ # xUnit — domain, services, specifications
│   │       └── Reconciliation.Infra.Tests/# xUnit — repositories (Testcontainers)
│   │
│   └── worker/                           # Go — File parsing and matching engine
│       ├── go.mod
│       ├── go.sum
│       ├── cmd/
│       │   ├── worker/                   # Redis Stream consumer entry point
│       │   └── cli/                      # CLI for manual runs and debugging
│       ├── internal/
│       │   ├── parsers/                  # All file parsers
│       │   │   ├── ofx/                  # OFX 1.x (SGML) / 2.x (XML)
│       │   │   ├── cnab/
│       │   │   │   ├── cnab240/
│       │   │   │   ├── cnab400/
│       │   │   │   └── profiles/         # Bank-specific CNAB field configs
│       │   │   ├── pix/                  # Pix settlement reports (JSON/CSV)
│       │   │   ├── acquirers/
│       │   │   │   ├── cielo/            # EEFI/EEVC formats
│       │   │   │   ├── rede/             # CSV
│       │   │   │   ├── stone/            # JSON API
│       │   │   │   ├── pagseguro/        # CSV
│       │   │   │   ├── getnet/           # Positional
│       │   │   │   └── safrapay/         # CSV/Positional
│       │   │   └── common/               # Shared types, normalization, fingerprinting
│       │   └── matching/                 # Three-pass matching engine
│       │       ├── exact/                # Key-based (NSU, E2EID, NossoNumero)
│       │       ├── fuzzy/                # Scored similarity (Levenshtein, amount, date)
│       │       └── aggregate/            # N-to-1 subset-sum solver
│       ├── testdata/                     # Anonymized sample files for parser tests
│       │   ├── ofx/
│       │   ├── cnab240/
│       │   ├── cnab400/
│       │   ├── pix/
│       │   ├── cielo/
│       │   ├── rede/
│       │   ├── stone/
│       │   └── pagseguro/
│       └── scripts/
│           ├── seed-calendar.go          # Populate banking holiday calendar
│           └── generate-test-data.go     # Generate synthetic test datasets
│
├── dashboard/                            # React + TypeScript + Tailwind frontend
│   ├── src/
│   │   ├── pages/                        # Route-level page components
│   │   ├── components/                   # Reusable UI components
│   │   ├── api/                          # Axios HTTP client
│   │   └── types/                        # TypeScript domain types
│   ├── package.json
│   └── vite.config.ts
│
├── migrations/                           # PostgreSQL schema migrations (shared)
├── deploy/
│   ├── docker/                           # Dockerfiles per service
│   └── docker-compose.yml
├── infra/                                # Terraform — AWS infrastructure
│   ├── modules/                          # Reusable TF modules
│   └── environments/                     # dev / prod variable overrides
├── docs/
│   ├── architecture.md
│   ├── file-formats.md
│   └── api-spec.yaml                     # OpenAPI 3.0
├── .github/workflows/                    # CI/CD pipelines per service
├── .env.example
├── .gitignore
└── SPEC.md
```

### Why This Structure

**Clear language boundaries.** The Go module (`go.mod`) lives inside `services/worker/`, so Go tooling only sees Go code. The C# solution (`ReconciliationEngine.sln`) lives inside `services/api/`, with no collision. The dashboard is already self-contained. Each service can be opened independently in an IDE without the other language's files cluttering the workspace.

**Tests live next to code.** Go tests use `_test.go` files alongside source (idiomatic Go convention), with anonymized fixture files under `services/worker/testdata/` — the only code that parses raw payment files. C# tests are separate xUnit projects within the solution, seeding domain objects and database state directly — they never touch raw payment files. No ambiguous shared test folders.

**Each service builds independently.** `cd services/worker && go build ./...` works. `cd services/api && dotnet build` works. `cd dashboard && npm run build` works. No cross-language build interference.

**Shared resources stay at the root.** Migrations, deploy configs, docs, infrastructure, and CI pipelines are cross-cutting concerns that belong at the monorepo root, not inside any single service.

---

## Development Phases

### Phase 1 — Foundation (Weeks 1–3)
Set up the monorepo structure, database schema, Docker Compose environment. Implement the core `TransactionRecord` data model and persistence layer. Build the OFX parser and a basic CNAB 240 parser (Itaú profile first, since BTG Pactual uses a similar layout). Create the ingestion API endpoint and file storage in MinIO.

### Phase 2 — Matching Engine (Weeks 4–6)
Implement the three-pass matching algorithm. Start with exact matching only, validate against test fixtures, then layer fuzzy and aggregate matching. Build the reconciliation run orchestration (trigger → ingest → normalize → match → persist results). Add the exception classification logic.

### Phase 3 — Acquirer Parsers & Fee Intelligence (Weeks 7–9)
Build parsers for CIELO, Stone, and Rede (the three largest acquirers). Implement the acquirer contract management module and fee validation engine. This phase produces the fee divergence detection feature.

### Phase 4 — Dashboard & Reporting (Weeks 10–12)
Build the React dashboard with reconciliation overview, exception management UI, and fee analysis charts. Implement PDF and Excel report generation. Add the aging dashboard for unreconciled items.

### Phase 5 — Hardening & Polish (Weeks 13–14)
Integration tests with realistic data volumes. Performance optimization (particularly the subset-sum solver). Prometheus/Grafana observability setup. Documentation, README, and architecture diagrams. Record a demo video for the portfolio.

---

## Key Technical Decisions & Rationale

**Why integer arithmetic for money?** Floating-point representation of decimal values is fundamentally broken for financial calculations. R$19.99 cannot be exactly represented in IEEE 754. Over thousands of transactions, rounding errors accumulate into real discrepancies. Using cents as i64 eliminates this entirely. Every amount in the system is stored and computed as an integer number of cents.

**Why Go for parsers and matching, C# for API?** Parsers are I/O-bound and benefit from Go's goroutine model for parallel file processing. The matching engine is CPU-bound and benefits from Go's low-overhead concurrency. The API layer, on the other hand, benefits from ASP.NET Core's mature middleware ecosystem (auth, validation, OpenAPI generation, Hangfire integration). This is also a deliberate portfolio choice: demonstrating comfort across both ecosystems.

**Why not use an existing reconciliation platform?** Products like Transfeera, Fitbank, or Celcoin offer partial reconciliation features, but they are either locked to their own payment processing ecosystem or prohibitively expensive for SMBs. The open-source alternatives (Ledger, Beancount) are personal finance tools not designed for multi-source commercial reconciliation. Building from scratch is the only way to handle the full complexity of the Brazilian payment ecosystem.

**Why batch instead of real-time?** Bank statement files (OFX/CNAB) are inherently batch artifacts — they are generated and made available at specific times (typically early morning). Card acquirer settlement files follow the same pattern. While Pix notifications could theoretically enable real-time reconciliation, the bank statement (the authoritative source of truth for "money actually arrived") is still batch. Building a real-time architecture would add complexity without meaningful benefit for the reconciliation use case.

---

## Domain-Driven Design (DDD)

This project follows DDD principles as described in Eric Evans' "Domain-Driven Design: Tackling Complexity in the Heart of Software." The reconciliation domain is complex enough to justify full DDD — there are multiple bounded contexts with distinct ubiquitous languages, non-trivial invariants that must be enforced at the domain level, and business rules that should not leak into infrastructure or application layers.

### Strategic Design

**Bounded Contexts**

The system is decomposed into four bounded contexts, each with its own ubiquitous language and internal model:

1. **Ingestion Context**: Concerned with file parsing, normalization, and deduplication. The language here is about "sources," "raw records," "normalization rules," "file fingerprints," and "parser profiles." A `TransactionRecord` in this context is a freshly parsed, normalized data point — it has no knowledge of matching or reconciliation outcomes.

2. **Reconciliation Context**: The core domain. This is where matching, exception detection, and reconciliation runs live. The language shifts to "matching passes," "confidence scores," "reconciliation pairs," "exceptions," "resolution workflows," and "reconciliation runs." This context owns the most complex business logic.

3. **Fee Intelligence Context**: Deals with acquirer contracts, fee validation, and overpayment detection. The language is about "contracted rates," "MDR," "bandeiras," "antecipação terms," "fee divergence," and "overpayment evidence." While closely related to reconciliation, it has its own aggregate roots and rules — a fee divergence can exist independently of a reconciliation pair.

4. **Reporting Context**: Read-model focused. Consumes events/data from the other contexts to produce dashboards, reports, and analytics. The language is about "reconciliation summaries," "aging buckets," "fee analysis periods," and "audit trails." This context is deliberately kept as a thin projection layer.

**Context Mapping**

The Ingestion Context has a **Customer-Supplier** relationship with the Reconciliation Context — Ingestion produces normalized `TransactionRecord`s that Reconciliation consumes. The contract between them is the canonical `TransactionRecord` schema, which Ingestion is responsible for producing correctly.

Fee Intelligence has a **Partnership** relationship with Reconciliation — they share the concept of a matched transaction but interpret it differently (Reconciliation cares about the match; Fee Intelligence cares about the fee delta on that match).

Reporting is a **Conformist** to all other contexts — it consumes their data without influencing their models.

### Tactical Design

**Value Objects**

Value Objects are the backbone of the domain model. They enforce invariants at construction time and are immutable:

- `Money`: Wraps an `int64` representing cents. Constructed via factory methods like `Money.FromCents(1999)` or `Money.FromReais(19.99m)` (the latter converts internally and is only used at system boundaries). Supports arithmetic operations that return new `Money` instances. Prevents negative amounts where domain rules require it. This is the single most important Value Object — it eliminates an entire class of bugs.

- `Cnpj` / `Cpf`: Self-validating document number types. Constructed from string input, validates the check digits (módulo 11 for both), and stores the raw 14/11-digit number without punctuation. Exposes formatted output (`XX.XXX.XXX/XXXX-XX`) as a presentation concern only.

- `PixKey`: Represents a Pix key with its type (CPF, CNPJ, email, phone, EVP/random). Validates format based on type. Two `PixKey` instances with the same normalized value are equal regardless of how they were input.

- `Nsu`: Acquirer-specific transaction identifier. Wraps the NSU string with its source acquirer, since NSU values are only unique within a single acquirer's universe.

- `EndToEndId`: The 32-character Pix E2EID. Validates format (starts with 'E', contains the ISPB, date, and sequence). Immutable.

- `NossoNumero`: Boleto identifier assigned by the bank. Includes the bank code and convenio as context, since nosso número alone is not globally unique.

- `SettlementDate`: A date that is always a dia útil. Construction validates against the banking holiday calendar. Provides methods like `NextBusinessDay()` and `AddBusinessDays(n int)`.

- `ConfidenceScore`: Bounded float between 0.0 and 1.0. Constructed via `ConfidenceScore.Of(0.87)` which panics/throws on out-of-range. Used exclusively in the matching engine.

- `DateRange`: Immutable pair of dates with the invariant that start ≤ end. Used for reconciliation run periods and report date ranges.

- `FeeRate`: Represents a percentage rate (e.g., MDR of 1.89%). Stored internally as basis points (189 bps) to avoid floating-point. Provides `ApplyTo(Money) → Money` to calculate the fee amount from a transaction value.

**Entities**

Entities have identity and lifecycle:

- `ReconciliationRun`: Identified by a unique `RunId` (UUID). Has a lifecycle: CREATED → INGESTING → MATCHING → CLASSIFYING → COMPLETED / FAILED. Aggregates statistics (matched count, exception count, duration). Each run is an immutable record of what happened — once COMPLETED, it cannot be modified.

- `Exception`: Identified by `ExceptionId`. Has a resolution lifecycle: OPEN → IN_REVIEW → RESOLVED / IGNORED. Carries a `resolution_note` (set on resolution) and `resolved_by` (analyst identity). The state transitions enforce invariants — you cannot resolve an exception that is already resolved, you cannot move from IGNORED back to OPEN.

- `AcquirerContract`: Identified by `ContractId`. Versioned — when terms change, a new version is created rather than mutating the existing record. This preserves the ability to validate historical fees against the rates that were in effect at the time of the transaction.

**Aggregates**

Each Aggregate enforces a consistency boundary:

- **ReconciliationRun Aggregate** (root: `ReconciliationRun`): Contains the run metadata, the list of `ReconciliationPair`s produced during this run, and the `Exception`s generated. The invariant is that a run's match count + exception count must equal the total ingested record count (every record is either matched or an exception — nothing falls through). All pairs and exceptions within a run are created and persisted atomically with the run's state transition to COMPLETED.

- **AcquirerContract Aggregate** (root: `AcquirerContract`): Contains the contract terms — a collection of `FeeSchedule` value objects (one per bandeira/produto combination). The invariant is that no two FeeSchedules within the same contract can have overlapping effective date ranges for the same bandeira/produto pair. Versioning is handled by creating a new AcquirerContract entity with an incremented version number and a new effective date.

- **Exception Aggregate** (root: `Exception`): Owns its resolution lifecycle. The invariant is the state machine: only valid transitions are allowed. The resolution note is mandatory when transitioning to RESOLVED. The `resolved_by` field must be a valid analyst identity.

**Repositories**

Repositories provide collection-like interfaces for Aggregate persistence. They operate on Aggregate roots only — you never persist a child entity directly:

- `IReconciliationRunRepository`: Methods like `Save(ReconciliationRun)`, `GetById(RunId)`, `GetByDateRange(DateRange)`, `GetLatestByStatus(RunStatus)`. The implementation uses PostgreSQL with the run, its pairs, and its exceptions persisted in a single transaction.

- `IAcquirerContractRepository`: Methods like `Save(AcquirerContract)`, `GetById(ContractId)`, `GetActiveByAcquirer(AcquirerId, Date)` (returns the contract version effective on that date), `GetAll()`.

- `IExceptionRepository`: Methods like `Save(Exception)`, `GetById(ExceptionId)`, `GetOpenByAging(AgingBucket)`, `GetByRunId(RunId)`. While exceptions live within the ReconciliationRun aggregate conceptually, they have their own repository because exception resolution happens independently of run lifecycle (an analyst resolves exceptions days after the run completed).

- `ITransactionRecordRepository`: For the Ingestion context. Methods like `Save(TransactionRecord)`, `GetByFingerprint(Hash)` (for dedup), `GetUnmatchedByDateRange(DateRange)`.

**Domain Services**

Operations that don't naturally belong to a single Entity or Value Object live in Domain Services:

- `MatchingService`: Orchestrates the three-pass matching algorithm. Takes a set of expected `TransactionRecord`s and a set of actual `TransactionRecord`s, runs them through exact → fuzzy → aggregate passes, and returns a collection of `ReconciliationPair`s and unmatched remainders. This is a stateless service — all state comes from the inputs.

- `FeeValidationService`: Takes a `ReconciliationPair` and the relevant `AcquirerContract`, calculates the expected fee based on contracted terms, and returns a `FeeValidationResult` (which may contain a `FeeDivergence` if the actual fee exceeds the contracted rate).

- `ExceptionClassificationService`: Takes an unmatched `TransactionRecord` and classifies it into an `ExceptionType` with a `Severity`. Uses a chain of classification rules (fee divergence check, timing window check, duplicate check, chargeback pattern check, etc.).

**Specifications**

The Specification pattern is used for complex query predicates and business rule validation that need to be composable and reusable:

- `UnreconciledTransactionSpecification`: Defines what it means for a transaction to be "unreconciled" — not part of any `ReconciliationPair` in any COMPLETED run, not already classified as an exception, and within the active reconciliation window. Used by both the matching engine (to find candidates) and the aging dashboard (to query unreconciled items).

- `FeeToleranceSpecification`: Defines when a fee difference is considered a "divergence" vs. acceptable rounding. Configurable tolerance (default: R$0.02 per transaction). Used by the `FeeValidationService` and by the exception classifier.

- `SettlementWindowSpecification`: Given a transaction's source type, acquirer, and bandeira, defines the acceptable settlement date window. A card transaction from CIELO with bandeira Visa on crédito à vista should settle within D+30 ± 1 dia útil. Used by the fuzzy matcher for date scoring and by the exception classifier for timing mismatch detection.

- `DuplicateTransactionSpecification`: Defines when two `TransactionRecord`s are considered duplicates — same fingerprint hash, or same (source, amount, date, counterparty) tuple within a configurable time window. Used during ingestion for dedup and during matching to flag potential double-processing.

- `AggregateMatchCandidateSpecification`: Defines eligibility for aggregate (N-to-1) matching — only transactions from sources known to batch-settle (card acquirers), within the same settlement date, from the same acquirer. Prevents the subset-sum solver from wasting cycles on impossible candidates.

Specifications are composable using `And`, `Or`, and `Not` combinators. They can be evaluated in-memory (for domain logic) or translated to SQL predicates (for repository queries) via a visitor pattern, avoiding the need to load entire datasets into memory for filtering.

**Domain Events**

Domain Events capture meaningful state changes for cross-context communication:

- `ReconciliationRunCompleted`: Emitted when a run transitions to COMPLETED. Consumed by the Reporting Context to update dashboards and by the Fee Intelligence Context to trigger fee validation on newly matched pairs.
- `ExceptionResolved`: Emitted when an analyst resolves an exception. Consumed by Reporting to update aging counts.
- `FeeDivergenceDetected`: Emitted by Fee Intelligence when an overcharge is found. Consumed by Reporting to update the fee analysis view and could trigger notifications.

Events are persisted to an outbox table in the same transaction as the aggregate change (transactional outbox pattern), then dispatched asynchronously. For the scope of this project, an in-process event bus is sufficient — no need for Kafka or RabbitMQ.

---

## Test-Driven Development (TDD)

Every line of domain logic in this project is written test-first. TDD is not optional or aspirational here — it is the development methodology. The reconciliation domain has enough edge cases and invariants that writing code without tests first would be reckless.

### TDD Workflow

The strict Red-Green-Refactor cycle applies:

1. **Red**: Write a failing test that describes the next increment of behavior. The test name should read like a specification: `MatchingService_ExactMatch_ByNSU_ReturnsHighConfidencePair`, `Money_FromReais_ConvertsCorrectlyToCents`, `Exception_Resolve_WithoutNote_ThrowsInvariantViolation`.

2. **Green**: Write the minimum code to make the test pass. No more.

3. **Refactor**: Clean up the implementation without changing behavior. All tests must remain green.

### Test Strategy by Layer

**Value Object Tests (Unit)**

Every Value Object has exhaustive unit tests covering: valid construction, invalid construction (rejected inputs), equality semantics, and arithmetic/behavior. These tests are fast, isolated, and form the foundation of the test suite.

Examples:
- `Money`: Construction from cents and reais, addition, subtraction, comparison, zero handling, negative rejection where applicable.
- `Cnpj`: Valid CNPJ passes, invalid check digits rejected, formatted vs. raw output, equality between formatted and unformatted input.
- `SettlementDate`: Construction on a business day succeeds, construction on a weekend or feriado is rejected (or auto-adjusted, depending on design choice), `AddBusinessDays` correctly skips holidays.

**Entity and Aggregate Tests (Unit)**

Test the invariants and lifecycle transitions:
- `ReconciliationRun`: Valid state transitions succeed, invalid transitions throw, completed run statistics are consistent.
- `Exception`: Resolution requires a note, state machine enforced, double-resolution rejected.

**Domain Service Tests (Unit)**

The matching engine is tested with known-answer fixtures:
- Exact match: Given transaction A with NSU X in expected set and transaction B with NSU X in actual set, `MatchingService` returns a pair with EXACT type and 1.0 confidence.
- Fuzzy match: Given transaction A (R$100.00 on 2026-03-05) and transaction B (R$99.98 on 2026-03-06), verify the confidence score and match type.
- Aggregate match: Given three expected transactions (R$50 + R$30 + R$20) and one actual deposit of R$100, verify the N-to-1 pair.
- No match: Given transactions with no plausible counterpart, verify they remain unmatched.

**Specification Tests (Unit)**

Each specification is tested in isolation and in combination:
- `FeeToleranceSpecification`: Fee delta of R$0.01 is within tolerance, R$0.05 is not.
- Composed: `UnreconciledTransactionSpecification.And(SettlementWindowSpecification)` correctly filters.

**Parser Tests (Unit/Integration)**

Every parser is tested against real anonymized file samples (stored in `test/fixtures/`). For each supported format and bank/acquirer variation:
- Parse a known file, assert the count of extracted records.
- Assert specific field values on known records (amount, date, NSU, etc.).
- Verify deduplication: parse the same file twice, assert no new records on second pass.
- Verify error handling: malformed files produce descriptive errors, not panics.

**Repository Tests (Integration)**

Tested against a real PostgreSQL instance (via Docker in CI):
- Save and retrieve aggregates, verify all children are persisted.
- Test specification-to-SQL translation: verify that in-memory and database evaluation produce the same results for the same inputs.
- Test concurrent save of the same fingerprint (dedup race condition).

**End-to-End Pipeline Tests (Integration)**

Full pipeline tests that ingest a set of fixture files, run the reconciliation, and assert the final output:
- A curated test dataset with known matches, known exceptions, and known fee divergences.
- Assert reconciliation rate, exception count, and fee divergence amounts match expected values exactly.
- These tests are slower and run in CI, not on every local save.

### Test Infrastructure

- **C# tests**: xUnit with FluentAssertions for readable assertions. Moq for interface mocking where needed (repositories in service tests). Testcontainers for PostgreSQL in integration tests.
- **Go tests**: Standard `testing` package with `testify/assert` and `testify/suite`. `testcontainers-go` for database tests. Table-driven tests for parser variations.
- **Test data**: The `test/fixtures/` directory contains anonymized but structurally realistic files for every supported format. A `generate-test-data.go` script produces synthetic datasets of configurable size for load testing.
- **Coverage target**: 90%+ line coverage on domain and service layers. Coverage on infrastructure/adapter layers is lower (repositories are tested via integration, not mocked). Coverage is tracked in CI but not used as a gate — a test that exercises a meaningful behavior is worth more than a test that bumps coverage numbers.

---

## Services Architecture

### Design Philosophy: Right-Sized for Scope

This project deliberately avoids microservices. A reconciliation engine for SMBs does not need Kubernetes, service meshes, API gateways, or distributed tracing across 12 services. That would be over-engineering that kills productivity and inflates the infrastructure budget for zero benefit at this scale.

The architecture follows a **modular monolith with a separated processing worker** — two deployable units that communicate through the database and Redis, not HTTP:

1. **API Service (C# / ASP.NET Core)**: Handles HTTP requests, orchestrates reconciliation runs, serves the dashboard API, manages acquirer contracts, and processes exception resolution. This is the "brain" — it owns the domain logic, the repositories, and the job scheduling (Hangfire).

2. **Processing Worker (Go)**: A long-running process that consumes file parsing and matching jobs from Redis Streams. When the API Service triggers a reconciliation run, it enqueues jobs. The Go worker picks them up, parses files, runs the matching engine, and writes results directly to PostgreSQL. The worker is stateless — it can be restarted or scaled horizontally without coordination.

3. **Dashboard (React + TypeScript)**: A static SPA served from the API Service (or a CDN in production). Communicates exclusively through the API Service's REST endpoints.

This gives us the performance benefits of Go where it matters (parsing, matching) and the developer productivity of C# where it matters (API, business rules, job scheduling), without the operational overhead of inter-service communication, distributed transactions, or independent deployment pipelines that would slow down a solo developer.

### Why Not Microservices

At the scale of this project — a single-tenant reconciliation engine processing tens of thousands of transactions daily — the operational cost of microservices outweighs every benefit. Specific reasons:

- **No independent scaling need**: The parser/matcher workload spikes during the daily batch run (early morning) and is idle otherwise. A single Go worker process handles this fine. If it ever needed more throughput, horizontal scaling of stateless workers is trivial without a microservices architecture.
- **No independent deployment need**: The domain model is tightly coupled by design (a reconciliation run touches ingestion, matching, fee validation, and exception classification). Deploying these independently would require versioned contracts and backward compatibility — pure overhead for a single-developer project.
- **No polyglot persistence need**: Everything lives in PostgreSQL. Adding a separate database per service would fragment the data model and require eventual consistency patterns (sagas, compensating transactions) for what is fundamentally a transactional workflow.
- **Debugging is local**: With two processes and a shared database, the entire request lifecycle is traceable with structured logs and correlation IDs. No distributed tracing infrastructure needed.

### Communication Patterns

- **API → Worker**: Via Redis Streams. The API publishes a job message (e.g., `{run_id, file_paths, job_type: "parse_and_match"}`), the worker consumes it. Simple, reliable, and Redis is already in the stack for caching.
- **Worker → API/DB**: The worker writes results directly to PostgreSQL. No callback HTTP requests. The API polls run status from the database (or uses PostgreSQL LISTEN/NOTIFY for near-real-time UI updates).
- **Dashboard → API**: Standard REST over HTTPS. No WebSockets needed — reconciliation is a batch process and the dashboard doesn't need sub-second updates. Polling every 5 seconds during an active run is sufficient.
- **Domain Events**: In-process event bus within the API Service. Events are dispatched synchronously within the same transaction boundary (transactional outbox for reliability). No external message broker needed.

---

## Infrastructure & Deployment

### Development Environment

For local development and early-stage work, the infrastructure is intentionally simple:

- **Docker Compose** orchestrates all services: PostgreSQL, Redis, MinIO, the API Service, the Go Worker, and the React dashboard (via dev server or nginx).
- **Single `docker-compose.yml`** brings up the entire stack with `docker compose up`. No external dependencies, no cloud accounts needed.
- **Local volumes** for PostgreSQL data and MinIO storage. Persistent across restarts, disposable with `docker compose down -v`.
- **Environment variables** for configuration (database connection strings, Redis URL, MinIO credentials). A `.env.example` file documents all required variables.
- **Hot reload**: The C# API runs via `dotnet watch`, the Go worker uses `air` for live reload, and the React dashboard uses Vite's dev server. Changes are reflected without manual restarts.

This is the development baseline. It runs on a developer laptop and in CI.

### Production Deployment (AWS) — Next Step

Production deployment targets AWS, provisioned entirely via Terraform. This is scoped as the next phase after the core engine is functional and tested. The infrastructure design prioritizes cost-efficiency for an SMB-targeted product while maintaining production-grade reliability.

**Compute**

- **API Service**: ECS Fargate (serverless containers). A single task with autoscaling based on CPU/memory. Fargate eliminates EC2 instance management and scales to zero cost when idle. The ASP.NET Core container runs behind an ALB (Application Load Balancer) with HTTPS termination via ACM (AWS Certificate Manager).
- **Go Worker**: ECS Fargate task triggered by EventBridge Scheduler (daily at 06:00 BRT) or on-demand via the API. Runs as a "batch job" — spins up, processes the reconciliation run, and terminates. This avoids paying for an always-on worker when reconciliation is a once-daily operation.
- **Dashboard**: Static assets deployed to S3 + CloudFront. No compute cost beyond CDN egress. Invalidation on deploy.

**Data**

- **PostgreSQL**: Amazon RDS (PostgreSQL 16) with a `db.t4g.medium` instance (2 vCPU, 4GB RAM) as the starting tier. Multi-AZ disabled initially (cost saving), enabled when the product has paying customers. Automated backups with 7-day retention. Parameter group tuned for the reconciliation workload (higher `work_mem`, `shared_buffers` sized to instance).
- **Redis**: Amazon ElastiCache (Redis 7) with a `cache.t4g.micro` single-node. Used for job queuing and caching. No cluster mode — the workload doesn't justify it.
- **MinIO → S3**: In production, MinIO is replaced by S3 directly. The application code uses the S3-compatible SDK in both environments (MinIO implements the S3 API). Bucket lifecycle policies archive files older than 90 days to S3 Glacier for cost reduction. Bucket versioning enabled for audit compliance.

**Networking**

- VPC with public and private subnets across 2 AZs. API Service and RDS in private subnets. ALB in public subnet. NAT Gateway for outbound internet from private subnets (needed for external API calls if acquirer integrations are added later).
- Security groups restrict access: ALB → API Service on port 8080, API Service → RDS on port 5432, API Service → ElastiCache on port 6379. No direct public access to any backend service.

**Observability**

- **CloudWatch Logs**: Structured JSON logs from ECS tasks shipped to CloudWatch. Log groups per service with 30-day retention.
- **CloudWatch Metrics + Alarms**: CPU/memory on Fargate tasks, RDS connection count, ElastiCache memory usage, ALB 5xx rate. Alarms notify via SNS → email (or Slack webhook).
- **Prometheus + Grafana**: Deployed as a sidecar or separate ECS task for application-level metrics (reconciliation rates, processing latency, exception counts). CloudWatch covers infrastructure; Prometheus covers business metrics.

**Terraform Structure**

```
infra/
├── environments/
│   ├── dev/
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   └── terraform.tfvars
│   └── prod/
│       ├── main.tf
│       ├── variables.tf
│       └── terraform.tfvars
├── modules/
│   ├── networking/          # VPC, subnets, NAT, security groups
│   ├── database/            # RDS PostgreSQL, parameter groups
│   ├── cache/               # ElastiCache Redis
│   ├── storage/             # S3 buckets, lifecycle policies
│   ├── compute/             # ECS cluster, task definitions, ALB
│   ├── cdn/                 # CloudFront + S3 for dashboard
│   └── observability/       # CloudWatch log groups, alarms, SNS
├── backend.tf               # S3 + DynamoDB remote state
└── versions.tf              # Provider version constraints
```

State is stored remotely in S3 with DynamoDB locking. Each environment (`dev`, `prod`) has its own state file and variable overrides. Modules are reusable across environments with different sizing parameters (e.g., `db_instance_class = "db.t4g.micro"` in dev, `"db.t4g.medium"` in prod).

**Estimated Monthly AWS Cost (Production, Low Traffic)**

| Service | Spec | Estimated Cost |
|---|---|---|
| ECS Fargate (API) | 0.5 vCPU, 1GB, always-on | ~$15 |
| ECS Fargate (Worker) | 1 vCPU, 2GB, ~30 min/day | ~$2 |
| RDS PostgreSQL | db.t4g.micro, 20GB, single-AZ | ~$15 |
| ElastiCache Redis | cache.t4g.micro | ~$12 |
| S3 | 10GB stored + lifecycle to Glacier | ~$1 |
| CloudFront | Low traffic CDN | ~$1 |
| ALB | Single ALB | ~$18 |
| NAT Gateway | Single AZ | ~$33 |
| CloudWatch | Logs + basic metrics | ~$5 |
| **Total** | | **~$102/month** |

The NAT Gateway is the single largest cost item. For a portfolio/early-stage deployment, it can be eliminated by placing the API in a public subnet with a security group (less secure but saves $33/month). This tradeoff is documented in the Terraform variables.

---

## CI/CD Pipeline

Each deployable unit has its own independent pipeline. The pipelines are implemented in GitHub Actions and follow the principle of fast feedback — unit tests run first, integration tests after, and deployment only on the main branch.

### Pipeline: API Service (C#)

```
Trigger: push to main OR pull request affecting src/Reconciliation.*/**

Steps:
1. Checkout
2. Setup .NET 8 SDK
3. Restore dependencies
4. Build (Release configuration)
5. Run unit tests (xUnit) — domain, services, specifications
6. Start PostgreSQL + Redis via Docker (services in GH Actions)
7. Run integration tests (Testcontainers) — repositories, end-to-end pipeline
8. [main only] Build Docker image, push to ECR
9. [main only] Update ECS task definition, deploy via rolling update
```

### Pipeline: Processing Worker (Go)

```
Trigger: push to main OR pull request affecting parsers/** OR matching/**

Steps:
1. Checkout
2. Setup Go 1.22
3. Run go vet + staticcheck (linting)
4. Run unit tests — parsers (against fixture files), matching engine
5. Start PostgreSQL + Redis via Docker
6. Run integration tests — full parse-and-match pipeline with fixture datasets
7. [main only] Build Docker image (multi-stage, scratch base), push to ECR
8. [main only] Update ECS task definition, deploy
```

### Pipeline: Dashboard (React)

```
Trigger: push to main OR pull request affecting dashboard/**

Steps:
1. Checkout
2. Setup Node 20
3. Install dependencies (npm ci)
4. Lint (ESLint) + type check (tsc --noEmit)
5. Run unit tests (Vitest) — component tests, utility functions
6. Build production bundle (vite build)
7. [main only] Sync build output to S3
8. [main only] Invalidate CloudFront distribution
```

### Pipeline: Infrastructure (Terraform)

```
Trigger: push to main OR pull request affecting infra/**

Steps:
1. Checkout
2. Setup Terraform
3. terraform fmt -check (formatting)
4. terraform init
5. terraform validate
6. terraform plan (output saved as artifact)
7. [main only, manual approval] terraform apply -auto-approve
```

The Terraform pipeline requires manual approval for `apply` on production. The plan output is posted as a PR comment for review. This prevents accidental infrastructure changes.

### Cross-Cutting CI Concerns

- **Branch protection**: `main` requires passing CI + at least one review. Direct pushes disabled.
- **Dependency scanning**: Dependabot (GitHub native) for C#, Go, and Node.js dependencies. Security alerts auto-create PRs.
- **Secret management**: AWS credentials stored as GitHub Actions secrets. Terraform uses an IAM role with least-privilege (only the permissions needed for the resources it manages). No secrets in code, ever.
- **Artifact versioning**: Docker images tagged with git SHA and `latest`. S3 dashboard deployments tagged with git SHA in metadata. Terraform state versioned via S3 bucket versioning.

---

*Last updated: March 2026*
*Author: Software Engineer — BTG Pactual Empresas / Fee Team*
*Status: Pre-development — Architecture & Design Phase*