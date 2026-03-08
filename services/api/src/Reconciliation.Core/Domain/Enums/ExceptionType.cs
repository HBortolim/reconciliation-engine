namespace ReconciliationEngine.Core.Domain.Enums;

/// <summary>
/// Represents the type of reconciliation exception.
/// </summary>
public enum ExceptionType
{
    FeeDivergence = 1,
    TimingMismatch = 2,
    PartialPayment = 3,
    Duplicate = 4,
    Chargeback = 5,
    PixDevolucao = 6,
    BoletoNotCompensated = 7,
    Unknown = 8
}
