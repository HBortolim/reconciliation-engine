namespace ReconciliationEngine.Core.Domain.ValueObjects;

/// <summary>
/// Value object representing a Network Service Unit (NSU) which is unique only within an acquirer.
/// </summary>
public record Nsu
{
    public string Value { get; }
    public string SourceAcquirer { get; }

    private Nsu(string value, string sourceAcquirer)
    {
        Value = value;
        SourceAcquirer = sourceAcquirer;
    }

    /// <summary>
    /// Creates an NSU instance.
    /// </summary>
    /// <param name="value">NSU value</param>
    /// <param name="sourceAcquirer">Source acquirer identifier</param>
    /// <returns>NSU instance</returns>
    public static Nsu Create(string value, string sourceAcquirer)
    {
        if (string.IsNullOrWhiteSpace(value))
            throw new ArgumentException("NSU value cannot be null or empty.", nameof(value));

        if (string.IsNullOrWhiteSpace(sourceAcquirer))
            throw new ArgumentException("Source acquirer cannot be null or empty.", nameof(sourceAcquirer));

        return new Nsu(value.Trim(), sourceAcquirer.Trim());
    }

    public override string ToString() => $"{SourceAcquirer}:{Value}";

    public bool Equals(Nsu? other) => other != null && Value == other.Value && SourceAcquirer == other.SourceAcquirer;
    public override int GetHashCode() => HashCode.Combine(Value, SourceAcquirer);
}
