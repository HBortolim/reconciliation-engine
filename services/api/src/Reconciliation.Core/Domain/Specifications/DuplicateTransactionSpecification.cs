using ReconciliationEngine.Core.Domain.Entities;

namespace ReconciliationEngine.Core.Domain.Specifications;

/// <summary>
/// Specification for detecting duplicate transactions by fingerprint or tuple.
/// </summary>
public class DuplicateTransactionSpecification : ISpecification<TransactionRecord>
{
    private readonly ISet<string> _knownFingerprints;

    public DuplicateTransactionSpecification(IEnumerable<string> knownFingerprints)
    {
        _knownFingerprints = new HashSet<string>(knownFingerprints);
    }

    public bool IsSatisfiedBy(TransactionRecord candidate)
    {
        return candidate != null && _knownFingerprints.Contains(candidate.FingerprintHash);
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
