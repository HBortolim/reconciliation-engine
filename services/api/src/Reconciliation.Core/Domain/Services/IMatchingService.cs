using ReconciliationEngine.Core.Domain.Entities;

namespace ReconciliationEngine.Core.Domain.Services;

/// <summary>
/// Domain service interface for transaction matching orchestration.
/// Implements 3-pass matching strategy.
/// </summary>
public interface IMatchingService
{
    /// <summary>
    /// Performs matching on provided transactions.
    /// Returns list of matched pairs.
    /// </summary>
    Task<IList<ReconciliationPair>> MatchTransactionsAsync(
        IEnumerable<TransactionRecord> sourceTransactions,
        IEnumerable<TransactionRecord> destinationTransactions,
        CancellationToken cancellationToken = default);
}
