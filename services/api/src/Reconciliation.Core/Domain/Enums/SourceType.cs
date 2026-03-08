namespace ReconciliationEngine.Core.Domain.Enums;

/// <summary>
/// Represents the source type of a transaction.
/// </summary>
public enum SourceType
{
    Pix = 1,
    Boleto = 2,
    CardCredit = 3,
    CardDebit = 4,
    Ted = 5,
    Doc = 6,
    DebitoAutomatico = 7
}
