namespace ReconciliationEngine.Core.Domain.ValueObjects;

/// <summary>
/// Immutable value object representing a date range with start <= end invariant.
/// </summary>
public record DateRange
{
    public DateOnly Start { get; }
    public DateOnly End { get; }

    private DateRange(DateOnly start, DateOnly end)
    {
        Start = start;
        End = end;
    }

    /// <summary>
    /// Creates a date range, ensuring start <= end.
    /// </summary>
    /// <param name="start">Start date</param>
    /// <param name="end">End date</param>
    /// <returns>DateRange instance</returns>
    public static DateRange Create(DateOnly start, DateOnly end)
    {
        if (start > end)
            throw new ArgumentException("Start date must be less than or equal to end date.", nameof(start));

        return new DateRange(start, end);
    }

    /// <summary>
    /// Checks if a date is within this range (inclusive).
    /// </summary>
    public bool Contains(DateOnly date) => date >= Start && date <= End;

    /// <summary>
    /// Gets the number of days in this range.
    /// </summary>
    public int DayCount => (End.DayNumber - Start.DayNumber) + 1;

    /// <summary>
    /// Checks if two ranges overlap.
    /// </summary>
    public bool Overlaps(DateRange other)
    {
        if (other == null) throw new ArgumentNullException(nameof(other));
        return Start <= other.End && End >= other.Start;
    }

    public override string ToString() => $"{Start:yyyy-MM-dd} to {End:yyyy-MM-dd}";

    public bool Equals(DateRange? other) => other != null && Start == other.Start && End == other.End;
    public override int GetHashCode() => HashCode.Combine(Start, End);
}
