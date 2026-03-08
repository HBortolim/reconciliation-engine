using ReconciliationEngine.Core.Domain.Entities;
using ReconciliationEngine.Core.Domain.ValueObjects;

namespace ReconciliationEngine.Core.Domain.Specifications;

/// <summary>
/// Specification for evaluating if fees are within acceptable tolerance.
/// Default tolerance is R$0.02 (2 centavos).
/// </summary>
public class FeeToleranceSpecification : ISpecification<(TransactionRecord transaction, Money expectedFee)>
{
    private readonly Money _tolerance;

    public FeeToleranceSpecification(Money? tolerance = null)
    {
        _tolerance = tolerance ?? Money.FromCentavos(2); // Default 2 centavos
    }

    public bool IsSatisfiedBy((TransactionRecord transaction, Money expectedFee) candidate)
    {
        if (candidate.transaction?.Fee == null || candidate.expectedFee == null)
            return false;

        var delta = candidate.transaction.Fee.Centavos >= candidate.expectedFee.Centavos
            ? candidate.transaction.Fee.Centavos - candidate.expectedFee.Centavos
            : candidate.expectedFee.Centavos - candidate.transaction.Fee.Centavos;

        return delta <= _tolerance.Centavos;
    }

    public ISpecification<(TransactionRecord, Money)> And(ISpecification<(TransactionRecord, Money)> other)
    {
        throw new NotImplementedException();
    }

    public ISpecification<(TransactionRecord, Money)> Or(ISpecification<(TransactionRecord, Money)> other)
    {
        throw new NotImplementedException();
    }

    public ISpecification<(TransactionRecord, Money)> Not()
    {
        throw new NotImplementedException();
    }
}
