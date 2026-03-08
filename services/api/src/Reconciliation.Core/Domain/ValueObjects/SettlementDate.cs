namespace ReconciliationEngine.Core.Domain.ValueObjects;

/// <summary>
/// Value object representing a settlement date that must be a business day (dia útil).
/// Validates against holiday calendar.
/// </summary>
public record SettlementDate
{
    // Simple hardcoded holidays for 2025-2026 (can be extended)
    private static readonly HashSet<DateOnly> BrazilianHolidays = new()
    {
        // 2025
        new(2025, 1, 1),   // New Year
        new(2025, 2, 17),  // Carnival Monday
        new(2025, 2, 18),  // Carnival Tuesday
        new(2025, 3, 28),  // Good Friday
        new(2025, 4, 21),  // Tiradentes' Day
        new(2025, 5, 1),   // Labor Day
        new(2025, 9, 7),   // Independence Day
        new(2025, 10, 12), // Our Lady Aparecida
        new(2025, 11, 2),  // All Souls' Day
        new(2025, 11, 20), // Black Consciousness Day
        new(2025, 12, 25), // Christmas
        // 2026
        new(2026, 1, 1),   // New Year
        new(2026, 2, 9),   // Carnival Monday
        new(2026, 2, 10),  // Carnival Tuesday
        new(2026, 4, 3),   // Good Friday
        new(2026, 4, 21),  // Tiradentes' Day
        new(2026, 5, 1),   // Labor Day
        new(2026, 9, 7),   // Independence Day
        new(2026, 10, 12), // Our Lady Aparecida
        new(2026, 11, 2),  // All Souls' Day
        new(2026, 11, 20), // Black Consciousness Day
        new(2026, 12, 25), // Christmas
    };

    public DateOnly Value { get; }

    private SettlementDate(DateOnly value)
    {
        Value = value;
    }

    /// <summary>
    /// Creates a settlement date, validating it's a business day.
    /// </summary>
    public static SettlementDate Create(DateOnly date)
    {
        if (!IsBusinessDay(date))
            throw new ArgumentException("Settlement date must be a business day (não pode ser fim de semana ou feriado).", nameof(date));

        return new SettlementDate(date);
    }

    /// <summary>
    /// Gets the next business day from this settlement date.
    /// </summary>
    public SettlementDate NextBusinessDay()
    {
        var next = Value.AddDays(1);
        while (!IsBusinessDay(next))
        {
            next = next.AddDays(1);
        }
        return new SettlementDate(next);
    }

    /// <summary>
    /// Adds business days to this settlement date.
    /// </summary>
    public SettlementDate AddBusinessDays(int days)
    {
        if (days < 0)
            throw new ArgumentException("Days must be non-negative.", nameof(days));

        var current = Value;
        for (int i = 0; i < days; i++)
        {
            current = current.AddDays(1);
            while (!IsBusinessDay(current))
            {
                current = current.AddDays(1);
            }
        }
        return new SettlementDate(current);
    }

    /// <summary>
    /// Checks if a date is a business day.
    /// </summary>
    private static bool IsBusinessDay(DateOnly date)
    {
        // Check if weekend
        if (date.DayOfWeek == DayOfWeek.Saturday || date.DayOfWeek == DayOfWeek.Sunday)
            return false;

        // Check if holiday
        return !BrazilianHolidays.Contains(date);
    }

    public override string ToString() => Value.ToString("yyyy-MM-dd");

    public bool Equals(SettlementDate? other) => other != null && Value == other.Value;
    public override int GetHashCode() => Value.GetHashCode();
}
