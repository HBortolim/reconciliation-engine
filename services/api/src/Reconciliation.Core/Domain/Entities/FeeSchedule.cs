using ReconciliationEngine.Core.Domain.ValueObjects;

namespace ReconciliationEngine.Core.Domain.Entities;

/// <summary>
/// Value object within AcquirerContract representing a fee schedule.
/// </summary>
public record FeeSchedule
{
    public string Bandeira { get; init; }
    public string Produto { get; init; }
    public FeeRate MdrRate { get; init; }
    public FeeRate? AntecipacaoRate { get; init; }
    public int SettlementDays { get; init; }

    private FeeSchedule() { }

    /// <summary>
    /// Creates a fee schedule.
    /// </summary>
    public FeeSchedule(
        string bandeira,
        string produto,
        FeeRate mdrRate,
        int settlementDays,
        FeeRate? antecipacaoRate = null)
    {
        if (string.IsNullOrWhiteSpace(bandeira))
            throw new ArgumentException("Bandeira cannot be null or empty.", nameof(bandeira));
        if (string.IsNullOrWhiteSpace(produto))
            throw new ArgumentException("Produto cannot be null or empty.", nameof(produto));
        if (mdrRate == null)
            throw new ArgumentNullException(nameof(mdrRate));
        if (settlementDays < 0)
            throw new ArgumentException("Settlement days cannot be negative.", nameof(settlementDays));

        Bandeira = bandeira;
        Produto = produto;
        MdrRate = mdrRate;
        AntecipacaoRate = antecipacaoRate;
        SettlementDays = settlementDays;
    }
}
