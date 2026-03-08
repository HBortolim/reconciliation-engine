namespace ReconciliationEngine.Core.Domain.ValueObjects;

/// <summary>
/// Self-validating value object representing a Brazilian CPF (Cadastro de Pessoa Física).
/// Validates modulo 11 check digits and stores the raw 11-digit string.
/// </summary>
public record Cpf
{
    private const int CpfLength = 11;

    public string Value { get; }

    private Cpf(string value)
    {
        Value = value;
    }

    /// <summary>
    /// Creates a CPF instance after validation.
    /// </summary>
    /// <param name="rawCpf">11-digit CPF string (can contain formatting characters)</param>
    /// <returns>CPF instance</returns>
    public static Cpf Create(string rawCpf)
    {
        if (string.IsNullOrWhiteSpace(rawCpf))
            throw new ArgumentException("CPF cannot be null or empty.", nameof(rawCpf));

        // Remove formatting
        var cleaned = System.Text.RegularExpressions.Regex.Replace(rawCpf, @"\D", "");

        if (cleaned.Length != CpfLength)
            throw new ArgumentException($"CPF must have exactly {CpfLength} digits.", nameof(rawCpf));

        if (!ValidateCpf(cleaned))
            throw new ArgumentException("Invalid CPF check digit.", nameof(rawCpf));

        return new Cpf(cleaned);
    }

    /// <summary>
    /// Gets the formatted CPF as XXX.XXX.XXX-XX.
    /// </summary>
    public string Formatted => $"{Value[0..3]}.{Value[3..6]}.{Value[6..9]}-{Value[9..11]}";

    private static bool ValidateCpf(string cpf)
    {
        if (string.IsNullOrEmpty(cpf) || cpf.Length != CpfLength)
            return false;

        // All same digits is invalid
        if (cpf.All(c => c == cpf[0]))
            return false;

        // Calculate first check digit
        int sum = 0;
        for (int i = 0; i < 9; i++)
            sum += int.Parse(cpf[i].ToString()) * (10 - i);

        int remainder = sum % 11;
        int firstDigit = remainder < 2 ? 0 : 11 - remainder;

        if (int.Parse(cpf[9].ToString()) != firstDigit)
            return false;

        // Calculate second check digit
        sum = 0;
        for (int i = 0; i < 10; i++)
            sum += int.Parse(cpf[i].ToString()) * (11 - i);

        remainder = sum % 11;
        int secondDigit = remainder < 2 ? 0 : 11 - remainder;

        return int.Parse(cpf[10].ToString()) == secondDigit;
    }

    public override string ToString() => Formatted;

    public bool Equals(Cpf? other) => other != null && Value == other.Value;
    public override int GetHashCode() => Value.GetHashCode();
}
