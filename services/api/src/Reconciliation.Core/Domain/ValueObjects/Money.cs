namespace ReconciliationEngine.Core.Domain.ValueObjects;

/// <summary>
/// Immutable value object representing monetary amounts in Brazilian Real (centavos).
/// Stored internally as long (int64) centavos for precision.
/// </summary>
public record Money
{
    private const long MinCentavos = 0;
    private const long MaxCentavos = long.MaxValue;

    public long Centavos { get; }

    private Money(long centavos)
    {
        if (centavos < MinCentavos)
            throw new ArgumentException("Money cannot be negative.", nameof(centavos));
        
        if (centavos > MaxCentavos)
            throw new ArgumentException("Money value exceeds maximum allowed.", nameof(centavos));

        Centavos = centavos;
    }

    /// <summary>
    /// Creates a Money instance from centavos.
    /// </summary>
    /// <param name="centavos">Amount in centavos (cents)</param>
    /// <returns>Money instance</returns>
    public static Money FromCentavos(long centavos) => new(centavos);

    /// <summary>
    /// Creates a Money instance from Brazilian Real (Reais).
    /// </summary>
    /// <param name="reais">Amount in reais</param>
    /// <returns>Money instance</returns>
    public static Money FromReais(decimal reais)
    {
        if (reais < 0)
            throw new ArgumentException("Money cannot be negative.", nameof(reais));

        var centavos = (long)(reais * 100);
        return new Money(centavos);
    }

    /// <summary>
    /// Adds two Money instances.
    /// </summary>
    public static Money operator +(Money left, Money right)
    {
        if (left == null) throw new ArgumentNullException(nameof(left));
        if (right == null) throw new ArgumentNullException(nameof(right));
        
        return new Money(left.Centavos + right.Centavos);
    }

    /// <summary>
    /// Subtracts one Money instance from another.
    /// </summary>
    public static Money operator -(Money left, Money right)
    {
        if (left == null) throw new ArgumentNullException(nameof(left));
        if (right == null) throw new ArgumentNullException(nameof(right));
        
        var result = left.Centavos - right.Centavos;
        if (result < 0)
            throw new InvalidOperationException("Subtraction would result in negative money.");
        
        return new Money(result);
    }

    /// <summary>
    /// Gets the Money as decimal Reais.
    /// </summary>
    public decimal ToReais() => Centavos / 100m;

    public override string ToString() => $"R$ {ToReais():F2}";

    public bool Equals(Money? other) => other != null && Centavos == other.Centavos;
    public override int GetHashCode() => Centavos.GetHashCode();
}
