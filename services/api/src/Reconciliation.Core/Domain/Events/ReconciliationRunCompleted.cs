namespace ReconciliationEngine.Core.Domain.Events;

/// <summary>
/// Domain event raised when a reconciliation run completes.
/// </summary>
public record ReconciliationRunCompleted : IDomainEvent
{
    public Guid RunId { get; init; }
    public int MatchedCount { get; init; }
    public int ExceptionCount { get; init; }
    public TimeSpan Duration { get; init; }
    public DateTime OccurredAt { get; init; }

    public ReconciliationRunCompleted(Guid runId, int matchedCount, int exceptionCount, TimeSpan duration)
    {
        if (runId == Guid.Empty)
            throw new ArgumentException("Run ID cannot be empty.", nameof(runId));
        if (matchedCount < 0)
            throw new ArgumentException("Matched count cannot be negative.", nameof(matchedCount));
        if (exceptionCount < 0)
            throw new ArgumentException("Exception count cannot be negative.", nameof(exceptionCount));

        RunId = runId;
        MatchedCount = matchedCount;
        ExceptionCount = exceptionCount;
        Duration = duration;
        OccurredAt = DateTime.UtcNow;
    }
}
