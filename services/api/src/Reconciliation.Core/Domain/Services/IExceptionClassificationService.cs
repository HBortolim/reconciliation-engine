using ReconciliationEngine.Core.Domain.Entities;

namespace ReconciliationEngine.Core.Domain.Services;

/// <summary>
/// Domain service interface for classifying unmatched transactions as exceptions.
/// </summary>
public interface IExceptionClassificationService
{
    /// <summary>
    /// Classifies unmatched transactions and returns exceptions.
    /// </summary>
    Task<IList<ReconciliationException>> ClassifyAsync(
        IEnumerable<TransactionRecord> unmatchedTransactions,
        CancellationToken cancellationToken = default);
}
