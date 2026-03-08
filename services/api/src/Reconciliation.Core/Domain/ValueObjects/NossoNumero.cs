namespace ReconciliationEngine.Core.Domain.ValueObjects;

/// <summary>
/// Value object representing a Boleto "Nosso Número" (Our Number).
/// Includes bank code and convenio information.
/// </summary>
public record NossoNumero
{
    public string Value { get; }
    public string BankCode { get; }
    public string Convenio { get; }

    private NossoNumero(string value, string bankCode, string convenio)
    {
        Value = value;
        BankCode = bankCode;
        Convenio = convenio;
    }

    /// <summary>
    /// Creates a Nosso Número instance.
    /// </summary>
    /// <param name="value">Nosso Número value</param>
    /// <param name="bankCode">Bank code (typically 3 digits)</param>
    /// <param name="convenio">Convenio (agreement) code</param>
    /// <returns>NossoNumero instance</returns>
    public static NossoNumero Create(string value, string bankCode, string convenio)
    {
        if (string.IsNullOrWhiteSpace(value))
            throw new ArgumentException("Nosso Número cannot be null or empty.", nameof(value));

        if (string.IsNullOrWhiteSpace(bankCode))
            throw new ArgumentException("Bank code cannot be null or empty.", nameof(bankCode));

        if (string.IsNullOrWhiteSpace(convenio))
            throw new ArgumentException("Convenio cannot be null or empty.", nameof(convenio));

        return new NossoNumero(value.Trim(), bankCode.Trim(), convenio.Trim());
    }

    public override string ToString() => $"{BankCode}/{Convenio}/{Value}";

    public bool Equals(NossoNumero? other) => 
        other != null && 
        Value == other.Value && 
        BankCode == other.BankCode && 
        Convenio == other.Convenio;

    public override int GetHashCode() => HashCode.Combine(Value, BankCode, Convenio);
}
