using ReconciliationEngine.Core.Domain.Entities;
using ReconciliationEngine.Core.Domain.ValueObjects;

namespace ReconciliationEngine.Core.Interfaces;

/// <summary>
/// Repository interface for TransactionRecord.
/// </summary>
public interface ITransactionRecordRepository
{
    /// <summary>
    /// Saves a transaction record.
    /// </summary>
    Task SaveAsync(TransactionRecord transaction, CancellationToken cancellationToken = default);

    /// <summary>
    /// Retrieves a transaction by fingerprint hash.
    /// </summary>
    Task<TransactionRecord?> GetByFingerprintAsync(string fingerprintHash, CancellationToken cancellationToken = default);

    /// <summary>
    /// Retrieves unmatched transactions within a date range.
    /// </summary>
    Task<IList<TransactionRecord>> GetUnmatchedByDateRangeAsync(DateRange range, CancellationToken cancellationToken = default);
}
