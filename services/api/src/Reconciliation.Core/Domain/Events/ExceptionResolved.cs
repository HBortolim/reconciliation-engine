using ReconciliationEngine.Core.Domain.Enums;

namespace ReconciliationEngine.Core.Domain.Events;

/// <summary>
/// Domain event raised when a reconciliation exception is resolved.
/// </summary>
public record ExceptionResolved : IDomainEvent
{
    public Guid ExceptionId { get; init; }
    public ResolutionStatus ResolutionStatus { get; init; }
    public string ResolvedBy { get; init; }
    public DateTime OccurredAt { get; init; }

    public ExceptionResolved(Guid exceptionId, ResolutionStatus resolutionStatus, string resolvedBy)
    {
        if (exceptionId == Guid.Empty)
            throw new ArgumentException("Exception ID cannot be empty.", nameof(exceptionId));
        if (string.IsNullOrWhiteSpace(resolvedBy))
            throw new ArgumentException("Resolved by cannot be null or empty.", nameof(resolvedBy));

        ExceptionId = exceptionId;
        ResolutionStatus = resolutionStatus;
        ResolvedBy = resolvedBy;
        OccurredAt = DateTime.UtcNow;
    }
}
