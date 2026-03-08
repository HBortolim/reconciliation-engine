using ReconciliationEngine.Core.Domain.Entities;

namespace ReconciliationEngine.Core.Domain.Services;

/// <summary>
/// Stub implementation of matching service.
/// Orchestrates 3-pass matching strategy:
/// 1. Exact matching (same amount, date, fee)
/// 2. Fuzzy matching (amount tolerance, date variance)
/// 3. Aggregate matching (multiple small transactions to one large)
/// </summary>
public class MatchingService : IMatchingService
{
    public Task<IList<ReconciliationPair>> MatchTransactionsAsync(
        IEnumerable<TransactionRecord> sourceTransactions,
        IEnumerable<TransactionRecord> destinationTransactions,
        CancellationToken cancellationToken = default)
    {
        // TODO: Implement 3-pass matching algorithm
        return Task.FromResult<IList<ReconciliationPair>>(new List<ReconciliationPair>());
    }
}
