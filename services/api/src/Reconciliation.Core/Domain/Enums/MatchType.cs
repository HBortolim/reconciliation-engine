namespace ReconciliationEngine.Core.Domain.Enums;

/// <summary>
/// Represents the type of matching performed on transaction pairs.
/// </summary>
public enum MatchType
{
    Exact = 1,
    Fuzzy = 2,
    Aggregate = 3
}
