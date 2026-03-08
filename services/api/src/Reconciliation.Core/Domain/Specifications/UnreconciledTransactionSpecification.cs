using ReconciliationEngine.Core.Domain.Entities;

namespace ReconciliationEngine.Core.Domain.Specifications;

/// <summary>
/// Specification for identifying unreconciled transactions.
/// </summary>
public class UnreconciledTransactionSpecification : ISpecification<TransactionRecord>
{
    private readonly ISet<Guid> _reconciledTransactionIds;

    public UnreconciledTransactionSpecification(IEnumerable<Guid> reconciledTransactionIds)
    {
        _reconciledTransactionIds = new HashSet<Guid>(reconciledTransactionIds);
    }

    public bool IsSatisfiedBy(TransactionRecord candidate)
    {
        return candidate != null && !_reconciledTransactionIds.Contains(candidate.Id);
    }

    public ISpecification<TransactionRecord> And(ISpecification<TransactionRecord> other)
    {
        return new CompositeSpecification(this, other, CompositeType.And);
    }

    public ISpecification<TransactionRecord> Or(ISpecification<TransactionRecord> other)
    {
        return new CompositeSpecification(this, other, CompositeType.Or);
    }

    public ISpecification<TransactionRecord> Not()
    {
        return new NotSpecification(this);
    }

    private class CompositeSpecification : ISpecification<TransactionRecord>
    {
        private readonly ISpecification<TransactionRecord> _left;
        private readonly ISpecification<TransactionRecord> _right;
        private readonly CompositeType _type;

        public CompositeSpecification(ISpecification<TransactionRecord> left, ISpecification<TransactionRecord> right, CompositeType type)
        {
            _left = left;
            _right = right;
            _type = type;
        }

        public bool IsSatisfiedBy(TransactionRecord candidate)
        {
            return _type == CompositeType.And
                ? _left.IsSatisfiedBy(candidate) && _right.IsSatisfiedBy(candidate)
                : _left.IsSatisfiedBy(candidate) || _right.IsSatisfiedBy(candidate);
        }

        public ISpecification<TransactionRecord> And(ISpecification<TransactionRecord> other) =>
            new CompositeSpecification(this, other, CompositeType.And);

        public ISpecification<TransactionRecord> Or(ISpecification<TransactionRecord> other) =>
            new CompositeSpecification(this, other, CompositeType.Or);

        public ISpecification<TransactionRecord> Not() =>
            new NotSpecification(this);
    }

    private class NotSpecification : ISpecification<TransactionRecord>
    {
        private readonly ISpecification<TransactionRecord> _specification;

        public NotSpecification(ISpecification<TransactionRecord> specification)
        {
            _specification = specification;
        }

        public bool IsSatisfiedBy(TransactionRecord candidate)
        {
            return !_specification.IsSatisfiedBy(candidate);
        }

        public ISpecification<TransactionRecord> And(ISpecification<TransactionRecord> other) =>
            new CompositeSpecification(this, other, CompositeType.And);

        public ISpecification<TransactionRecord> Or(ISpecification<TransactionRecord> other) =>
            new CompositeSpecification(this, other, CompositeType.Or);

        public ISpecification<TransactionRecord> Not() =>
            _specification;
    }

    private enum CompositeType { And, Or }
}
