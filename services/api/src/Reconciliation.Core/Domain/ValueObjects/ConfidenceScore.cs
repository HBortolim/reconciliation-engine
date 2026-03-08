namespace ReconciliationEngine.Core.Domain.ValueObjects;

/// <summary>
/// Value object representing a confidence score bounded between 0.0 and 1.0.
/// </summary>
public record ConfidenceScore
{
    private const double MinScore = 0.0;
    private const double MaxScore = 1.0;

    public double Value { get; }

    private ConfidenceScore(double value)
    {
        Value = value;
    }

    /// <summary>
    /// Creates a confidence score instance.
    /// </summary>
    /// <param name="value">Score between 0.0 and 1.0</param>
    /// <returns>ConfidenceScore instance</returns>
    public static ConfidenceScore Create(double value)
    {
        if (value < MinScore || value > MaxScore)
            throw new ArgumentException($"Confidence score must be between {MinScore} and {MaxScore}.", nameof(value));

        return new ConfidenceScore(value);
    }

    /// <summary>
    /// Determines if this score is above a threshold.
    /// </summary>
    public bool IsAboveThreshold(double threshold) => Value >= threshold;

    /// <summary>
    /// Determines if this score is below a threshold.
    /// </summary>
    public bool IsBelowThreshold(double threshold) => Value <= threshold;

    public override string ToString() => Value.ToString("P2"); // Format as percentage

    public bool Equals(ConfidenceScore? other) => other != null && Value.Equals(other.Value);
    public override int GetHashCode() => Value.GetHashCode();
}
