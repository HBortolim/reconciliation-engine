using ReconciliationEngine.Core.Domain.Enums;
using ReconciliationEngine.Core.Domain.ValueObjects;

namespace ReconciliationEngine.Core.Domain.Entities;

/// <summary>
/// Entity representing a transaction record from a source system.
/// </summary>
public class TransactionRecord
{
    public Guid Id { get; private set; }
    public SourceType SourceType { get; private set; }
    public Money Amount { get; private set; }
    public Money Fee { get; private set; }
    public Money NetAmount { get; private set; }
    public DateOnly? ExpectedSettlementDate { get; private set; }
    public DateOnly? ActualSettlementDate { get; private set; }
    public string? CounterpartyDocument { get; private set; }
    public string ExternalId { get; private set; }
    public string FingerprintHash { get; private set; }
    public string SourceFile { get; private set; }
    public DateTime ParsedAt { get; private set; }

    private TransactionRecord() { }

    /// <summary>
    /// Creates a new transaction record.
    /// </summary>
    public TransactionRecord(
        SourceType sourceType,
        Money amount,
        Money fee,
        DateOnly? expectedSettlementDate,
        DateOnly? actualSettlementDate,
        string externalId,
        string fingerprintHash,
        string sourceFile,
        string? counterpartyDocument = null)
    {
        if (amount == null) throw new ArgumentNullException(nameof(amount));
        if (fee == null) throw new ArgumentNullException(nameof(fee));
        if (string.IsNullOrWhiteSpace(externalId))
            throw new ArgumentException("External ID cannot be null or empty.", nameof(externalId));
        if (string.IsNullOrWhiteSpace(fingerprintHash))
            throw new ArgumentException("Fingerprint hash cannot be null or empty.", nameof(fingerprintHash));
        if (string.IsNullOrWhiteSpace(sourceFile))
            throw new ArgumentException("Source file cannot be null or empty.", nameof(sourceFile));

        Id = Guid.NewGuid();
        SourceType = sourceType;
        Amount = amount;
        Fee = fee;
        NetAmount = amount - fee;
        ExpectedSettlementDate = expectedSettlementDate;
        ActualSettlementDate = actualSettlementDate;
        CounterpartyDocument = counterpartyDocument;
        ExternalId = externalId;
        FingerprintHash = fingerprintHash;
        SourceFile = sourceFile;
        ParsedAt = DateTime.UtcNow;
    }

    public override bool Equals(object? obj)
    {
        return obj is TransactionRecord record && record.Id == Id;
    }

    public override int GetHashCode()
    {
        return Id.GetHashCode();
    }
}
