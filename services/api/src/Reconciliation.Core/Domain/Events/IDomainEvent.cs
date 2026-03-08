namespace ReconciliationEngine.Core.Domain.Events;

/// <summary>
/// Marker interface for domain events.
/// </summary>
public interface IDomainEvent
{
    /// <summary>
    /// Gets the timestamp when this event occurred.
    /// </summary>
    DateTime OccurredAt { get; }
}
