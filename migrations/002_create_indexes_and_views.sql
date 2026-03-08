-- Views for analytics and operational querying

-- v_reconciliation_summary: Run summary with key metrics and percentages
CREATE VIEW v_reconciliation_summary AS
SELECT
    rr.id,
    rr.status,
    rr.started_at,
    rr.completed_at,
    EXTRACT(EPOCH FROM (COALESCE(rr.completed_at, NOW()) - rr.started_at)) / 1000.0 AS duration_ms,
    rr.total_records_ingested,
    rr.matched_count,
    rr.exception_count,
    CASE 
        WHEN rr.total_records_ingested > 0 
        THEN ROUND(100.0 * rr.matched_count / rr.total_records_ingested, 2)
        ELSE 0
    END AS match_percentage,
    CASE 
        WHEN rr.total_records_ingested > 0 
        THEN ROUND(100.0 * rr.exception_count / rr.total_records_ingested, 2)
        ELSE 0
    END AS exception_percentage,
    rr.created_at
FROM reconciliation_runs rr;

-- v_exception_aging: Exceptions grouped by age buckets
CREATE VIEW v_exception_aging AS
SELECT
    CASE
        WHEN ex.resolution_status = 'RESOLVED' THEN 'RESOLVED'
        WHEN ex.resolution_status = 'IGNORED' THEN 'IGNORED'
        WHEN NOW() - ex.created_at < INTERVAL '1 day' THEN '0-1d (Critical)'
        WHEN NOW() - ex.created_at < INTERVAL '7 days' THEN '1-7d (High)'
        WHEN NOW() - ex.created_at < INTERVAL '30 days' THEN '7-30d (Medium)'
        ELSE '30d+ (Aged)'
    END AS aging_bucket,
    ex.exception_type,
    ex.severity,
    COUNT(*) AS count,
    ROUND(100.0 * COUNT(*) / (SELECT COUNT(*) FROM exceptions WHERE resolution_status IN ('OPEN', 'IN_REVIEW')), 2) AS percentage_of_open,
    MIN(ex.created_at) AS oldest_exception,
    MAX(ex.created_at) AS newest_exception
FROM exceptions ex
WHERE ex.resolution_status IN ('OPEN', 'IN_REVIEW')
GROUP BY aging_bucket, ex.exception_type, ex.severity
ORDER BY aging_bucket, COUNT(*) DESC;

-- v_fee_analysis: Fee divergences aggregated by acquirer and bandeira
CREATE VIEW v_fee_analysis AS
SELECT
    ac.acquirer_id,
    ac.acquirer_name,
    fs.bandeira,
    fs.produto,
    COUNT(rp.id) AS pair_count,
    ROUND(AVG(ABS(rp.fee_delta_centavos))::NUMERIC / 100, 2) AS avg_fee_delta_reais,
    ROUND(MAX(ABS(rp.fee_delta_centavos))::NUMERIC / 100, 2) AS max_fee_delta_reais,
    ROUND(STDDEV(ABS(rp.fee_delta_centavos))::NUMERIC / 100, 2) AS stddev_fee_delta_reais,
    ROUND(SUM(ABS(rp.fee_delta_centavos))::NUMERIC / 100, 2) AS total_fee_variance_reais,
    COUNT(CASE WHEN rp.fee_delta_centavos > 0 THEN 1 END) AS overpayment_count,
    COUNT(CASE WHEN rp.fee_delta_centavos < 0 THEN 1 END) AS underpayment_count
FROM reconciliation_pairs rp
JOIN transaction_records tr ON rp.actual_record_id = tr.id
JOIN acquirer_contracts ac ON tr.source_type::TEXT = ac.acquirer_id
JOIN fee_schedules fs ON fs.contract_id = ac.id
WHERE rp.fee_delta_centavos IS NOT NULL
GROUP BY ac.acquirer_id, ac.acquirer_name, fs.bandeira, fs.produto
ORDER BY total_fee_variance_reais DESC NULLS LAST;

-- v_unreconciled_transactions: Transactions not yet in any reconciliation pair
CREATE VIEW v_unreconciled_transactions AS
SELECT
    tr.id,
    tr.source_type,
    tr.amount_centavos,
    ROUND(tr.amount_centavos::NUMERIC / 100, 2) AS amount_reais,
    tr.fee_centavos,
    ROUND(tr.fee_centavos::NUMERIC / 100, 2) AS fee_reais,
    tr.counterparty_document,
    tr.external_id,
    tr.expected_settlement_date,
    tr.actual_settlement_date,
    tr.source_file,
    tr.parsed_at,
    tr.created_at,
    NOW() - tr.created_at AS age_since_ingestion
FROM transaction_records tr
WHERE tr.id NOT IN (
    SELECT expected_record_id FROM reconciliation_pairs
    UNION
    SELECT actual_record_id FROM reconciliation_pairs
)
ORDER BY tr.created_at DESC;

-- v_recent_runs: Latest reconciliation runs with aggregate stats
CREATE VIEW v_recent_runs AS
SELECT
    rr.id,
    rr.status,
    rr.started_at,
    rr.completed_at,
    rr.total_records_ingested,
    rr.matched_count,
    rr.exception_count,
    CASE 
        WHEN rr.total_records_ingested > 0 
        THEN ROUND(100.0 * rr.matched_count / rr.total_records_ingested, 2)
        ELSE 0
    END AS match_percentage,
    CASE
        WHEN rr.status = 'FAILED' THEN 'ERROR'
        WHEN rr.status = 'COMPLETED' THEN 'OK'
        ELSE 'IN_PROGRESS'
    END AS health_status,
    ARRAY_LENGTH(rr.files_ingested, 1) AS file_count,
    rr.processing_duration_ms
FROM reconciliation_runs rr
ORDER BY rr.started_at DESC
LIMIT 100;

-- v_exception_summary: Aggregate statistics on exceptions by type and severity
CREATE VIEW v_exception_summary AS
SELECT
    ex.exception_type,
    ex.severity,
    COUNT(*) AS total_count,
    COUNT(CASE WHEN ex.resolution_status = 'OPEN' THEN 1 END) AS open_count,
    COUNT(CASE WHEN ex.resolution_status = 'IN_REVIEW' THEN 1 END) AS in_review_count,
    COUNT(CASE WHEN ex.resolution_status = 'RESOLVED' THEN 1 END) AS resolved_count,
    COUNT(CASE WHEN ex.resolution_status = 'IGNORED' THEN 1 END) AS ignored_count,
    ROUND(100.0 * COUNT(CASE WHEN ex.resolution_status = 'RESOLVED' THEN 1 END) / COUNT(*), 2) AS resolution_rate_percent
FROM exceptions ex
GROUP BY ex.exception_type, ex.severity
ORDER BY total_count DESC, ex.severity DESC;

-- v_transaction_source_distribution: Transaction counts and amounts by source type
CREATE VIEW v_transaction_source_distribution AS
SELECT
    tr.source_type,
    COUNT(*) AS transaction_count,
    SUM(tr.amount_centavos) AS total_amount_centavos,
    ROUND(SUM(tr.amount_centavos)::NUMERIC / 100, 2) AS total_amount_reais,
    ROUND(AVG(tr.amount_centavos)::NUMERIC / 100, 2) AS avg_amount_reais,
    SUM(tr.fee_centavos) AS total_fees_centavos,
    ROUND(SUM(tr.fee_centavos)::NUMERIC / 100, 2) AS total_fees_reais,
    MIN(tr.parsed_at) AS first_transaction_date,
    MAX(tr.parsed_at) AS last_transaction_date
FROM transaction_records tr
GROUP BY tr.source_type
ORDER BY total_amount_centavos DESC;
