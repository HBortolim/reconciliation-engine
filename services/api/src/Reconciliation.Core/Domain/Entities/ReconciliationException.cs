using ReconciliationEngine.Core.Domain.Enums;

namespace ReconciliationEngine.Core.Domain.Entities;

/// <summary>
/// Entity representing a reconciliation exception.
/// Implements a state machine for resolution status transitions.
/// </summary>
public class ReconciliationException
{
    public Guid Id { get; private set; }
    public ExceptionType ExceptionType { get; private set; }
    public Severity Severity { get; private set; }
    public string SuggestedAction { get; private set; }
    public ResolutionStatus ResolutionStatus { get; private set; }
    public string? ResolutionNote { get; private set; }
    public string? ResolvedBy { get; private set; }
    public DateTime? ResolvedAt { get; private set; }
    public DateTime CreatedAt { get; private set; }

    private ReconciliationException() { }

    /// <summary>
    /// Creates a new reconciliation exception in Open status.
    /// </summary>
    public ReconciliationException(
        ExceptionType exceptionType,
        Severity severity,
        string suggestedAction)
    {
        if (string.IsNullOrWhiteSpace(suggestedAction))
            throw new ArgumentException("Suggested action cannot be null or empty.", nameof(suggestedAction));

        Id = Guid.NewGuid();
        ExceptionType = exceptionType;
        Severity = severity;
        SuggestedAction = suggestedAction;
        ResolutionStatus = ResolutionStatus.Open;
        CreatedAt = DateTime.UtcNow;
    }

    /// <summary>
    /// Transitions the exception to InReview status.
    /// </summary>
    public void StartReview()
    {
        if (ResolutionStatus != ResolutionStatus.Open)
            throw new InvalidOperationException($"Cannot start review: status is {ResolutionStatus}, expected {ResolutionStatus.Open}.");

        ResolutionStatus = ResolutionStatus.InReview;
    }

    /// <summary>
    /// Resolves the exception with optional analyst note.
    /// </summary>
    public void Resolve(string? note, string analyst)
    {
        if (ResolutionStatus == ResolutionStatus.Resolved)
            throw new InvalidOperationException("Exception is already resolved.");

        if (ResolutionStatus != ResolutionStatus.InReview && ResolutionStatus != ResolutionStatus.Open)
            throw new InvalidOperationException($"Cannot resolve from status {ResolutionStatus}.");

        if (string.IsNullOrWhiteSpace(analyst))
            throw new ArgumentException("Analyst name cannot be null or empty.", nameof(analyst));

        ResolutionStatus = ResolutionStatus.Resolved;
        ResolutionNote = note;
        ResolvedBy = analyst;
        ResolvedAt = DateTime.UtcNow;
    }

    /// <summary>
    /// Marks the exception as ignored.
    /// </summary>
    public void Ignore(string? note, string analyst)
    {
        if (ResolutionStatus == ResolutionStatus.Ignored)
            throw new InvalidOperationException("Exception is already ignored.");

        if (string.IsNullOrWhiteSpace(analyst))
            throw new ArgumentException("Analyst name cannot be null or empty.", nameof(analyst));

        ResolutionStatus = ResolutionStatus.Ignored;
        ResolutionNote = note;
        ResolvedBy = analyst;
        ResolvedAt = DateTime.UtcNow;
    }

    public override bool Equals(object? obj)
    {
        return obj is ReconciliationException exception && exception.Id == Id;
    }

    public override int GetHashCode()
    {
        return Id.GetHashCode();
    }
}
