using ReconciliationEngine.Core.Domain.Entities;

namespace ReconciliationEngine.Core.Domain.Specifications;

/// <summary>
/// Specification for determining eligibility for N-to-1 (aggregate) matching.
/// </summary>
public class AggregateMatchCandidateSpecification : ISpecification<IEnumerable<TransactionRecord>>
{
    private readonly Money _minAmount;
    private readonly int _minTransactionCount;

    public AggregateMatchCandidateSpecification(Money? minAmount = null, int minTransactionCount = 2)
    {
        _minAmount = minAmount ?? Money.FromCentavos(0);
        _minTransactionCount = minTransactionCount;
    }

    public bool IsSatisfiedBy(IEnumerable<TransactionRecord> candidate)
    {
        if (candidate == null)
            return false;

        var transactions = candidate.ToList();

        // Must have minimum number of transactions
        if (transactions.Count < _minTransactionCount)
            return false;

        // All must be from same source type
        var firstSource = transactions.First().SourceType;
        if (!transactions.All(t => t.SourceType == firstSource))
            return false;

        // Aggregate amount must meet minimum
        var totalAmount = transactions.Aggregate(
            Money.FromCentavos(0),
            (acc, t) => acc + t.Amount);

        return totalAmount.Centavos >= _minAmount.Centavos;
    }

    public ISpecification<IEnumerable<TransactionRecord>> And(ISpecification<IEnumerable<TransactionRecord>> other)
    {
        throw new NotImplementedException();
    }

    public ISpecification<IEnumerable<TransactionRecord>> Or(ISpecification<IEnumerable<TransactionRecord>> other)
    {
        throw new NotImplementedException();
    }

    public ISpecification<IEnumerable<TransactionRecord>> Not()
    {
        throw new NotImplementedException();
    }
}
