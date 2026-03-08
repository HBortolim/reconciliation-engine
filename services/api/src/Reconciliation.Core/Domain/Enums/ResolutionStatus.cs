namespace ReconciliationEngine.Core.Domain.Enums;

/// <summary>
/// Represents the resolution status of a reconciliation exception.
/// </summary>
public enum ResolutionStatus
{
    Open = 1,
    InReview = 2,
    Resolved = 3,
    Ignored = 4
}
