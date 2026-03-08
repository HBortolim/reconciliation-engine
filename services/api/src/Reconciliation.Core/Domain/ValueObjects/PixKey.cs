namespace ReconciliationEngine.Core.Domain.ValueObjects;

/// <summary>
/// Value object representing a Pix key with validation per type.
/// </summary>
public record PixKey
{
    public PixKeyType Type { get; }
    public string Value { get; }

    private PixKey(PixKeyType type, string value)
    {
        Type = type;
        Value = value;
    }

    /// <summary>
    /// Creates a Pix key instance with format validation.
    /// </summary>
    public static PixKey Create(PixKeyType type, string value)
    {
        if (string.IsNullOrWhiteSpace(value))
            throw new ArgumentException("Pix key value cannot be null or empty.", nameof(value));

        var normalized = value.Trim().ToLower();

        return type switch
        {
            PixKeyType.Cpf => ValidateCpfKey(normalized),
            PixKeyType.Cnpj => ValidateCnpjKey(normalized),
            PixKeyType.Email => ValidateEmailKey(normalized),
            PixKeyType.Phone => ValidatePhoneKey(normalized),
            PixKeyType.Evp => ValidateEvpKey(normalized),
            _ => throw new ArgumentException($"Unknown Pix key type: {type}", nameof(type))
        };
    }

    private static PixKey ValidateCpfKey(string value)
    {
        var cleaned = System.Text.RegularExpressions.Regex.Replace(value, @"\D", "");
        if (cleaned.Length != 11)
            throw new ArgumentException("CPF key must have 11 digits.", nameof(value));
        return new PixKey(PixKeyType.Cpf, cleaned);
    }

    private static PixKey ValidateCnpjKey(string value)
    {
        var cleaned = System.Text.RegularExpressions.Regex.Replace(value, @"\D", "");
        if (cleaned.Length != 14)
            throw new ArgumentException("CNPJ key must have 14 digits.", nameof(value));
        return new PixKey(PixKeyType.Cnpj, cleaned);
    }

    private static PixKey ValidateEmailKey(string value)
    {
        if (!value.Contains('@') || !value.Contains('.'))
            throw new ArgumentException("Invalid email format for Pix key.", nameof(value));
        return new PixKey(PixKeyType.Email, value);
    }

    private static PixKey ValidatePhoneKey(string value)
    {
        var cleaned = System.Text.RegularExpressions.Regex.Replace(value, @"\D", "");
        if (cleaned.Length < 10 || cleaned.Length > 11)
            throw new ArgumentException("Phone key must have 10 or 11 digits.", nameof(value));
        return new PixKey(PixKeyType.Phone, cleaned);
    }

    private static PixKey ValidateEvpKey(string value)
    {
        if (value.Length != 32)
            throw new ArgumentException("EVP key must have 32 characters.", nameof(value));
        return new PixKey(PixKeyType.Evp, value);
    }

    public override string ToString() => Value;

    public bool Equals(PixKey? other) => other != null && Type == other.Type && Value == other.Value;
    public override int GetHashCode() => HashCode.Combine(Type, Value);
}

/// <summary>
/// Enum representing Pix key types.
/// </summary>
public enum PixKeyType
{
    Cpf = 1,
    Cnpj = 2,
    Email = 3,
    Phone = 4,
    Evp = 5
}
