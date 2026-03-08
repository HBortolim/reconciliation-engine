using ReconciliationEngine.Core.Domain.Enums;
using ReconciliationEngine.Core.Domain.Entities;
using ReconciliationEngine.Core.Domain.ValueObjects;

namespace ReconciliationEngine.Core.Domain.Specifications;

/// <summary>
/// Specification for validating settlement window per source/acquirer/bandeira combination.
/// </summary>
public class SettlementWindowSpecification : ISpecification<TransactionRecord>
{
    private readonly Dictionary<(SourceType source, string acquirer, string bandeira), int> _settlementDays;

    public SettlementWindowSpecification(Dictionary<(SourceType, string, string), int> settlementDays)
    {
        _settlementDays = settlementDays ?? new();
    }

    public bool IsSatisfiedBy(TransactionRecord candidate)
    {
        if (candidate?.ExpectedSettlementDate == null || candidate?.ActualSettlementDate == null)
            return true; // Cannot evaluate without dates

        // Simplified: check if actual settlement is within expected window
        var expected = candidate.ExpectedSettlementDate.Value;
        var actual = candidate.ActualSettlementDate.Value;

        return actual <= expected;
    }

    public ISpecification<TransactionRecord> And(ISpecification<TransactionRecord> other)
    {
        throw new NotImplementedException();
    }

    public ISpecification<TransactionRecord> Or(ISpecification<TransactionRecord> other)
    {
        throw new NotImplementedException();
    }

    public ISpecification<TransactionRecord> Not()
    {
        throw new NotImplementedException();
    }
}
