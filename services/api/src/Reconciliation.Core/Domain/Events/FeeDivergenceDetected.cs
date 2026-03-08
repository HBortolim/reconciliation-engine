using ReconciliationEngine.Core.Domain.ValueObjects;

namespace ReconciliationEngine.Core.Domain.Events;

/// <summary>
/// Domain event raised when a fee divergence is detected.
/// </summary>
public record FeeDivergenceDetected : IDomainEvent
{
    public Guid PairId { get; init; }
    public Money ExpectedFee { get; init; }
    public Money ActualFee { get; init; }
    public Money DeltaCentavos { get; init; }
    public DateTime OccurredAt { get; init; }

    public FeeDivergenceDetected(Guid pairId, Money expectedFee, Money actualFee, Money deltaCentavos)
    {
        if (pairId == Guid.Empty)
            throw new ArgumentException("Pair ID cannot be empty.", nameof(pairId));
        if (expectedFee == null)
            throw new ArgumentNullException(nameof(expectedFee));
        if (actualFee == null)
            throw new ArgumentNullException(nameof(actualFee));
        if (deltaCentavos == null)
            throw new ArgumentNullException(nameof(deltaCentavos));

        PairId = pairId;
        ExpectedFee = expectedFee;
        ActualFee = actualFee;
        DeltaCentavos = deltaCentavos;
        OccurredAt = DateTime.UtcNow;
    }
}
