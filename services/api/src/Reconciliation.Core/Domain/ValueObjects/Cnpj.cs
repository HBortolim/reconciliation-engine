namespace ReconciliationEngine.Core.Domain.ValueObjects;

/// <summary>
/// Self-validating value object representing a Brazilian CNPJ (Cadastro Nacional da Pessoa Jurídica).
/// Validates modulo 11 check digits and stores the raw 14-digit string.
/// </summary>
public record Cnpj
{
    private const int CnpjLength = 14;

    public string Value { get; }

    private Cnpj(string value)
    {
        Value = value;
    }

    /// <summary>
    /// Creates a CNPJ instance after validation.
    /// </summary>
    /// <param name="rawCnpj">14-digit CNPJ string (can contain formatting characters)</param>
    /// <returns>CNPJ instance</returns>
    public static Cnpj Create(string rawCnpj)
    {
        if (string.IsNullOrWhiteSpace(rawCnpj))
            throw new ArgumentException("CNPJ cannot be null or empty.", nameof(rawCnpj));

        // Remove formatting
        var cleaned = System.Text.RegularExpressions.Regex.Replace(rawCnpj, @"\D", "");

        if (cleaned.Length != CnpjLength)
            throw new ArgumentException($"CNPJ must have exactly {CnpjLength} digits.", nameof(rawCnpj));

        if (!ValidateCnpj(cleaned))
            throw new ArgumentException("Invalid CNPJ check digit.", nameof(rawCnpj));

        return new Cnpj(cleaned);
    }

    /// <summary>
    /// Gets the formatted CNPJ as XX.XXX.XXX/XXXX-XX.
    /// </summary>
    public string Formatted => $"{Value[0..2]}.{Value[2..5]}.{Value[5..8]}/{Value[8..12]}-{Value[12..14]}";

    private static bool ValidateCnpj(string cnpj)
    {
        if (string.IsNullOrEmpty(cnpj) || cnpj.Length != CnpjLength)
            return false;

        // Calculate first check digit
        int sum = 0;
        int multiplier = 2;
        for (int i = 11; i >= 0; i--)
        {
            sum += int.Parse(cnpj[i].ToString()) * multiplier;
            multiplier++;
            if (multiplier > 9) multiplier = 2;
        }

        int remainder = sum % 11;
        int firstDigit = remainder < 2 ? 0 : 11 - remainder;

        if (int.Parse(cnpj[12].ToString()) != firstDigit)
            return false;

        // Calculate second check digit
        sum = 0;
        multiplier = 2;
        for (int i = 12; i >= 1; i--)
        {
            sum += int.Parse(cnpj[i].ToString()) * multiplier;
            multiplier++;
            if (multiplier > 9) multiplier = 2;
        }

        remainder = sum % 11;
        int secondDigit = remainder < 2 ? 0 : 11 - remainder;

        return int.Parse(cnpj[13].ToString()) == secondDigit;
    }

    public override string ToString() => Formatted;

    public bool Equals(Cnpj? other) => other != null && Value == other.Value;
    public override int GetHashCode() => Value.GetHashCode();
}
