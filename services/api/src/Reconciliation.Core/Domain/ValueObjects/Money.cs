namespace ReconciliationEngine.Core.Domain.ValueObjects;

/// <summary>
/// Immutable value object representing monetary amounts in Brazilian Real (cents).
/// Stored internally as long (int64) cents for precision.
/// </summary>
public record Money
{
    private const long MinCents = 0;
    private const long MaxCents = long.MaxValue;

    public long Cents { get; }

    private Money(long cents)
    {
        if (cents < MinCents)
            throw new ArgumentException("Money cannot be negative.", nameof(cents));
        
        if (cents > MaxCents)
            throw new ArgumentException("Money value exceeds maximum allowed.", nameof(cents));

        Cents = cents;
    }

    /// <summary>
    /// Creates a Money instance from cents.
    /// </summary>
    /// <param name="cents">Amount in cents (cents)</param>
    /// <returns>Money instance</returns>
    public static Money FromCents(long cents) => new(cents);

    /// <summary>
    /// Creates a Money instance from Brazilian Real (Reais).
    /// </summary>
    /// <param name="reais">Amount in reais</param>
    /// <returns>Money instance</returns>
    public static Money FromReais(decimal reais)
    {
        if (reais < 0)
            throw new ArgumentException("Money cannot be negative.", nameof(reais));

        var cents = (long)(reais * 100);
        return new Money(cents);
    }

    /// <summary>
    /// Adds two Money instances.
    /// </summary>
    public static Money operator +(Money left, Money right)
    {
        if (left == null) throw new ArgumentNullException(nameof(left));
        if (right == null) throw new ArgumentNullException(nameof(right));
        
        return new Money(left.Cents + right.Cents);
    }

    /// <summary>
    /// Subtracts one Money instance from another.
    /// </summary>
    public static Money operator -(Money left, Money right)
    {
        if (left == null) throw new ArgumentNullException(nameof(left));
        if (right == null) throw new ArgumentNullException(nameof(right));
        
        var result = left.Cents - right.Cents;
        if (result < 0)
            throw new InvalidOperationException("Subtraction would result in negative money.");
        
        return new Money(result);
    }

    /// <summary>
    /// Gets the Money as decimal Reais.
    /// </summary>
    public decimal ToReais() => Cents / 100m;

    public override string ToString() => $"R$ {ToReais():F2}";

    public bool Equals(Money? other) => other != null && Cents == other.Cents;
    public override int GetHashCode() => Cents.GetHashCode();
}
