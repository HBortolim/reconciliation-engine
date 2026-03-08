namespace ReconciliationEngine.Core.Domain.ValueObjects;

/// <summary>
/// Value object representing a fee rate stored as basis points (1 bp = 0.01%).
/// </summary>
public record FeeRate
{
    public int BasisPoints { get; }

    private FeeRate(int basisPoints)
    {
        if (basisPoints < 0)
            throw new ArgumentException("Fee rate cannot be negative.", nameof(basisPoints));

        BasisPoints = basisPoints;
    }

    /// <summary>
    /// Creates a fee rate from basis points.
    /// </summary>
    /// <param name="basisPoints">Rate in basis points (1 bp = 0.01%)</param>
    /// <returns>FeeRate instance</returns>
    public static FeeRate FromBasisPoints(int basisPoints) => new(basisPoints);

    /// <summary>
    /// Creates a fee rate from a percentage.
    /// </summary>
    /// <param name="percentage">Rate as percentage (e.g., 1.5 for 1.5%)</param>
    /// <returns>FeeRate instance</returns>
    public static FeeRate FromPercentage(decimal percentage)
    {
        var basisPoints = (int)(percentage * 100);
        return new(basisPoints);
    }

    /// <summary>
    /// Applies this fee rate to a Money amount.
    /// </summary>
    /// <param name="amount">Amount to apply fee to</param>
    /// <returns>Fee amount</returns>
    public Money ApplyTo(Money amount)
    {
        if (amount == null) throw new ArgumentNullException(nameof(amount));

        // (amount * basisPoints) / 10000
        var feeCentavos = (amount.Centavos * BasisPoints) / 10000;
        return Money.FromCentavos(feeCentavos);
    }

    /// <summary>
    /// Gets the rate as a percentage.
    /// </summary>
    public decimal ToPercentage() => BasisPoints / 100m;

    public override string ToString() => $"{ToPercentage():F2}%";

    public bool Equals(FeeRate? other) => other != null && BasisPoints == other.BasisPoints;
    public override int GetHashCode() => BasisPoints.GetHashCode();
}
