namespace ReconciliationEngine.Core.Domain.Enums;

/// <summary>
/// Represents the status of a reconciliation run.
/// </summary>
public enum RunStatus
{
    Created = 1,
    Ingesting = 2,
    Matching = 3,
    Classifying = 4,
    Completed = 5,
    Failed = 6
}
