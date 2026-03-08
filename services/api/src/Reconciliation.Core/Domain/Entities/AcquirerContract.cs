namespace ReconciliationEngine.Core.Domain.Entities;

/// <summary>
/// Entity representing a versioned acquirer contract.
/// Immutable once created to maintain audit trail.
/// </summary>
public class AcquirerContract
{
    public Guid Id { get; private set; }
    public string AcquirerId { get; private set; }
    public int Version { get; private set; }
    public DateTime EffectiveFrom { get; private set; }
    public DateTime? EffectiveTo { get; private set; }
    public IReadOnlyList<FeeSchedule> FeeSchedules { get; private set; }
    public DateTime CreatedAt { get; private set; }

    private AcquirerContract() { }

    /// <summary>
    /// Creates a new acquirer contract (versioned, immutable).
    /// </summary>
    public AcquirerContract(
        string acquirerId,
        int version,
        DateTime effectiveFrom,
        DateTime? effectiveTo,
        IEnumerable<FeeSchedule> feeSchedules)
    {
        if (string.IsNullOrWhiteSpace(acquirerId))
            throw new ArgumentException("Acquirer ID cannot be null or empty.", nameof(acquirerId));
        if (version < 1)
            throw new ArgumentException("Version must be >= 1.", nameof(version));
        if (feeSchedules == null || !feeSchedules.Any())
            throw new ArgumentException("Fee schedules cannot be null or empty.", nameof(feeSchedules));
        if (effectiveTo != null && effectiveFrom > effectiveTo)
            throw new ArgumentException("Effective from date must be <= effective to date.", nameof(effectiveFrom));

        Id = Guid.NewGuid();
        AcquirerId = acquirerId;
        Version = version;
        EffectiveFrom = effectiveFrom;
        EffectiveTo = effectiveTo;
        FeeSchedules = feeSchedules.ToList().AsReadOnly();
        CreatedAt = DateTime.UtcNow;
    }

    /// <summary>
    /// Checks if this contract is currently active.
    /// </summary>
    public bool IsActive()
    {
        var now = DateTime.UtcNow;
        return now >= EffectiveFrom && (EffectiveTo == null || now <= EffectiveTo);
    }

    public override bool Equals(object? obj)
    {
        return obj is AcquirerContract contract && contract.Id == Id;
    }

    public override int GetHashCode()
    {
        return Id.GetHashCode();
    }
}
