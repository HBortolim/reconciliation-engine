using ReconciliationEngine.Core.Domain.Enums;
using ReconciliationEngine.Core.Domain.ValueObjects;

namespace ReconciliationEngine.Core.Domain.Entities;

/// <summary>
/// Entity representing a matched pair of transactions.
/// </summary>
public class ReconciliationPair
{
    public Guid Id { get; private set; }
    public Guid SourceTransactionId { get; private set; }
    public Guid DestinationTransactionId { get; private set; }
    public MatchType MatchType { get; private set; }
    public ConfidenceScore ConfidenceScore { get; private set; }
    public Money AmountDelta { get; private set; }
    public int? DateDeltaDays { get; private set; }
    public Money FeeDelta { get; private set; }
    public DateTime CreatedAt { get; private set; }

    private ReconciliationPair() { }

    /// <summary>
    /// Creates a new reconciliation pair.
    /// </summary>
    public ReconciliationPair(
        Guid sourceTransactionId,
        Guid destinationTransactionId,
        MatchType matchType,
        ConfidenceScore confidenceScore,
        Money amountDelta,
        Money feeDelta,
        int? dateDeltaDays = null)
    {
        if (sourceTransactionId == Guid.Empty)
            throw new ArgumentException("Source transaction ID cannot be empty.", nameof(sourceTransactionId));
        if (destinationTransactionId == Guid.Empty)
            throw new ArgumentException("Destination transaction ID cannot be empty.", nameof(destinationTransactionId));
        if (amountDelta == null) throw new ArgumentNullException(nameof(amountDelta));
        if (feeDelta == null) throw new ArgumentNullException(nameof(feeDelta));
        if (confidenceScore == null) throw new ArgumentNullException(nameof(confidenceScore));

        Id = Guid.NewGuid();
        SourceTransactionId = sourceTransactionId;
        DestinationTransactionId = destinationTransactionId;
        MatchType = matchType;
        ConfidenceScore = confidenceScore;
        AmountDelta = amountDelta;
        DateDeltaDays = dateDeltaDays;
        FeeDelta = feeDelta;
        CreatedAt = DateTime.UtcNow;
    }

    public override bool Equals(object? obj)
    {
        return obj is ReconciliationPair pair && pair.Id == Id;
    }

    public override int GetHashCode()
    {
        return Id.GetHashCode();
    }
}
