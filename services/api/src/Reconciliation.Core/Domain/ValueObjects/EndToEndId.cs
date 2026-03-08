namespace ReconciliationEngine.Core.Domain.ValueObjects;

/// <summary>
/// Value object representing a Pix End-to-End ID (E2EID).
/// Must be 32 characters, start with 'E', and contain an ISPB code.
/// </summary>
public record EndToEndId
{
    private const int E2EIdLength = 32;

    public string Value { get; }

    private EndToEndId(string value)
    {
        Value = value;
    }

    /// <summary>
    /// Creates an End-to-End ID instance with format validation.
    /// </summary>
    /// <param name="value">32-character E2EID</param>
    /// <returns>EndToEndId instance</returns>
    public static EndToEndId Create(string value)
    {
        if (string.IsNullOrWhiteSpace(value))
            throw new ArgumentException("End-to-End ID cannot be null or empty.", nameof(value));

        if (value.Length != E2EIdLength)
            throw new ArgumentException($"End-to-End ID must be exactly {E2EIdLength} characters.", nameof(value));

        if (!value.StartsWith('E') && !value.StartsWith('e'))
            throw new ArgumentException("End-to-End ID must start with 'E'.", nameof(value));

        // ISPB check: positions 1-8 should be numeric (ISPB code)
        var ispbPart = value.Substring(1, 8);
        if (!ispbPart.All(char.IsDigit))
            throw new ArgumentException("End-to-End ID must contain numeric ISPB code in positions 2-9.", nameof(value));

        return new EndToEndId(value.ToUpper());
    }

    public override string ToString() => Value;

    public bool Equals(EndToEndId? other) => other != null && Value == other.Value;
    public override int GetHashCode() => Value.GetHashCode();
}
