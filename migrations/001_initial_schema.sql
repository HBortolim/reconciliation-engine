-- Extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Enums
CREATE TYPE source_type AS ENUM ('PIX', 'BOLETO', 'CARD_CREDIT', 'CARD_DEBIT', 'TED', 'DOC', 'DEBITO_AUTOMATICO');
CREATE TYPE match_type AS ENUM ('EXACT', 'FUZZY', 'AGGREGATE');
CREATE TYPE exception_type AS ENUM ('FEE_DIVERGENCE', 'TIMING_MISMATCH', 'PARTIAL_PAYMENT', 'DUPLICATE', 'CHARGEBACK', 'PIX_DEVOLUCAO', 'BOLETO_NOT_COMPENSATED', 'UNKNOWN');
CREATE TYPE severity AS ENUM ('LOW', 'MEDIUM', 'HIGH', 'CRITICAL');
CREATE TYPE run_status AS ENUM ('CREATED', 'INGESTING', 'MATCHING', 'CLASSIFYING', 'COMPLETED', 'FAILED');
CREATE TYPE resolution_status AS ENUM ('OPEN', 'IN_REVIEW', 'RESOLVED', 'IGNORED');

-- Transaction Records table
CREATE TABLE transaction_records (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    source_type source_type NOT NULL,
    amount_centavos BIGINT NOT NULL CHECK (amount_centavos >= 0),
    fee_centavos BIGINT NOT NULL DEFAULT 0,
    net_amount_centavos BIGINT NOT NULL,
    expected_settlement_date DATE,
    actual_settlement_date DATE,
    counterparty_document VARCHAR(14),
    external_id VARCHAR(255),
    fingerprint_hash VARCHAR(64) NOT NULL,
    source_file VARCHAR(500) NOT NULL,
    raw_data JSONB,
    parsed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_fingerprint UNIQUE (fingerprint_hash)
);

-- Indexes on transaction_records
CREATE INDEX idx_tr_source_type ON transaction_records(source_type);
CREATE INDEX idx_tr_external_id ON transaction_records(external_id);
CREATE INDEX idx_tr_counterparty ON transaction_records(counterparty_document);
CREATE INDEX idx_tr_expected_settlement ON transaction_records(expected_settlement_date);
CREATE INDEX idx_tr_actual_settlement ON transaction_records(actual_settlement_date);
CREATE INDEX idx_tr_amount ON transaction_records(amount_centavos);
CREATE INDEX idx_tr_parsed_at ON transaction_records(parsed_at);

-- Reconciliation Runs table
CREATE TABLE reconciliation_runs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    status run_status NOT NULL DEFAULT 'CREATED',
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    failed_at TIMESTAMPTZ,
    failure_reason TEXT,
    files_ingested TEXT[] NOT NULL DEFAULT '{}',
    total_records_ingested INTEGER NOT NULL DEFAULT 0,
    matched_count INTEGER NOT NULL DEFAULT 0,
    exception_count INTEGER NOT NULL DEFAULT 0,
    processing_duration_ms BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_rr_status ON reconciliation_runs(status);
CREATE INDEX idx_rr_started_at ON reconciliation_runs(started_at);

-- Reconciliation Pairs table
CREATE TABLE reconciliation_pairs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    run_id UUID NOT NULL REFERENCES reconciliation_runs(id) ON DELETE CASCADE,
    expected_record_id UUID NOT NULL REFERENCES transaction_records(id),
    actual_record_id UUID NOT NULL REFERENCES transaction_records(id),
    match_type match_type NOT NULL,
    confidence_score DOUBLE PRECISION NOT NULL CHECK (confidence_score >= 0 AND confidence_score <= 1),
    amount_delta_centavos BIGINT DEFAULT 0,
    date_delta_days INTEGER DEFAULT 0,
    fee_delta_centavos BIGINT DEFAULT 0,
    matched_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_rp_run_id ON reconciliation_pairs(run_id);
CREATE INDEX idx_rp_expected ON reconciliation_pairs(expected_record_id);
CREATE INDEX idx_rp_actual ON reconciliation_pairs(actual_record_id);
CREATE INDEX idx_rp_match_type ON reconciliation_pairs(match_type);

-- Exceptions table
CREATE TABLE exceptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    run_id UUID NOT NULL REFERENCES reconciliation_runs(id) ON DELETE CASCADE,
    transaction_record_id UUID NOT NULL REFERENCES transaction_records(id),
    exception_type exception_type NOT NULL,
    severity severity NOT NULL,
    suggested_action TEXT,
    resolution_status resolution_status NOT NULL DEFAULT 'OPEN',
    resolution_note TEXT,
    resolved_by VARCHAR(255),
    resolved_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_ex_run_id ON exceptions(run_id);
CREATE INDEX idx_ex_status ON exceptions(resolution_status);
CREATE INDEX idx_ex_type ON exceptions(exception_type);
CREATE INDEX idx_ex_severity ON exceptions(severity);
CREATE INDEX idx_ex_created_at ON exceptions(created_at);

-- Acquirer Contracts table
CREATE TABLE acquirer_contracts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    acquirer_id VARCHAR(100) NOT NULL,
    acquirer_name VARCHAR(255) NOT NULL,
    version INTEGER NOT NULL DEFAULT 1,
    effective_from DATE NOT NULL,
    effective_to DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_acquirer_version UNIQUE (acquirer_id, version)
);

CREATE INDEX idx_ac_acquirer ON acquirer_contracts(acquirer_id);
CREATE INDEX idx_ac_effective ON acquirer_contracts(effective_from, effective_to);

-- Fee Schedules table (child of acquirer_contracts)
CREATE TABLE fee_schedules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    contract_id UUID NOT NULL REFERENCES acquirer_contracts(id) ON DELETE CASCADE,
    bandeira VARCHAR(50) NOT NULL,
    produto VARCHAR(50) NOT NULL,
    mdr_basis_points INTEGER NOT NULL CHECK (mdr_basis_points >= 0),
    antecipacao_basis_points INTEGER DEFAULT 0,
    settlement_days INTEGER NOT NULL DEFAULT 30,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_fee_schedule UNIQUE (contract_id, bandeira, produto)
);

CREATE INDEX idx_fs_contract ON fee_schedules(contract_id);

-- Banking Holidays calendar
CREATE TABLE banking_holidays (
    id SERIAL PRIMARY KEY,
    holiday_date DATE NOT NULL UNIQUE,
    description VARCHAR(255) NOT NULL,
    year INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_bh_date ON banking_holidays(holiday_date);
CREATE INDEX idx_bh_year ON banking_holidays(year);

-- Domain Events outbox (transactional outbox pattern)
CREATE TABLE domain_events_outbox (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_type VARCHAR(255) NOT NULL,
    aggregate_id UUID NOT NULL,
    aggregate_type VARCHAR(255) NOT NULL,
    payload JSONB NOT NULL,
    occurred_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    published_at TIMESTAMPTZ,
    published BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX idx_deo_unpublished ON domain_events_outbox(published) WHERE published = FALSE;
CREATE INDEX idx_deo_occurred ON domain_events_outbox(occurred_at);

-- Ingested Files tracking (for idempotency)
CREATE TABLE ingested_files (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    file_name VARCHAR(500) NOT NULL,
    file_hash VARCHAR(64) NOT NULL,
    file_size_bytes BIGINT NOT NULL,
    storage_path VARCHAR(1000) NOT NULL,
    record_count INTEGER NOT NULL DEFAULT 0,
    ingested_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    run_id UUID REFERENCES reconciliation_runs(id),
    CONSTRAINT uq_file_hash UNIQUE (file_hash)
);

CREATE INDEX idx_if_hash ON ingested_files(file_hash);
CREATE INDEX idx_if_run ON ingested_files(run_id);

-- Seed: 2026 Brazilian banking holidays (Anbima)
INSERT INTO banking_holidays (holiday_date, description, year) VALUES
('2026-01-01', 'Confraternização Universal', 2026),
('2026-02-16', 'Carnaval', 2026),
('2026-02-17', 'Carnaval', 2026),
('2026-04-03', 'Sexta-feira Santa', 2026),
('2026-04-21', 'Tiradentes', 2026),
('2026-05-01', 'Dia do Trabalho', 2026),
('2026-06-04', 'Corpus Christi', 2026),
('2026-09-07', 'Independência do Brasil', 2026),
('2026-10-12', 'Nossa Senhora Aparecida', 2026),
('2026-11-02', 'Finados', 2026),
('2026-11-15', 'Proclamação da República', 2026),
('2026-12-25', 'Natal', 2026);
