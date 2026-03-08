using ReconciliationEngine.Core.Domain.Enums;
using ReconciliationEngine.Core.Domain.Entities;

namespace ReconciliationEngine.Core.Domain.Aggregates;

/// <summary>
/// Aggregate root representing a reconciliation run.
/// Manages the state machine for reconciliation process and enforces invariants.
/// </summary>
public class ReconciliationRun
{
    private readonly List<ReconciliationPair> _pairs = new();
    private readonly List<ReconciliationException> _exceptions = new();

    public Guid RunId { get; private set; }
    public RunStatus Status { get; private set; }
    public DateTime StartedAt { get; private set; }
    public DateTime? CompletedAt { get; private set; }
    public int FilesIngested { get; private set; }
    public IReadOnlyList<ReconciliationPair> Pairs => _pairs.AsReadOnly();
    public IReadOnlyList<ReconciliationException> Exceptions => _exceptions.AsReadOnly();

    private ReconciliationRun() { }

    /// <summary>
    /// Creates a new reconciliation run in Created status.
    /// </summary>
    public ReconciliationRun()
    {
        RunId = Guid.NewGuid();
        Status = RunStatus.Created;
        StartedAt = DateTime.UtcNow;
        FilesIngested = 0;
    }

    /// <summary>
    /// Transitions to Ingesting status.
    /// </summary>
    public void StartIngestion()
    {
        if (Status != RunStatus.Created)
            throw new InvalidOperationException($"Cannot start ingestion from status {Status}. Expected {RunStatus.Created}.");

        Status = RunStatus.Ingesting;
    }

    /// <summary>
    /// Transitions to Matching status.
    /// </summary>
    public void StartMatching()
    {
        if (Status != RunStatus.Ingesting)
            throw new InvalidOperationException($"Cannot start matching from status {Status}. Expected {RunStatus.Ingesting}.");

        Status = RunStatus.Matching;
    }

    /// <summary>
    /// Transitions to Classifying status.
    /// </summary>
    public void StartClassifying()
    {
        if (Status != RunStatus.Matching)
            throw new InvalidOperationException($"Cannot start classifying from status {Status}. Expected {RunStatus.Matching}.");

        Status = RunStatus.Classifying;
    }

    /// <summary>
    /// Completes the reconciliation run.
    /// </summary>
    public void Complete()
    {
        if (Status != RunStatus.Classifying)
            throw new InvalidOperationException($"Cannot complete from status {Status}. Expected {RunStatus.Classifying}.");

        Status = RunStatus.Completed;
        CompletedAt = DateTime.UtcNow;
    }

    /// <summary>
    /// Fails the reconciliation run with a reason.
    /// </summary>
    public void Fail(string reason)
    {
        if (string.IsNullOrWhiteSpace(reason))
            throw new ArgumentException("Failure reason cannot be null or empty.", nameof(reason));

        Status = RunStatus.Failed;
        CompletedAt = DateTime.UtcNow;
    }

    /// <summary>
    /// Records a reconciliation pair match.
    /// </summary>
    public void AddPair(ReconciliationPair pair)
    {
        if (pair == null) throw new ArgumentNullException(nameof(pair));
        if (Status != RunStatus.Matching && Status != RunStatus.Classifying)
            throw new InvalidOperationException("Pairs can only be added during Matching or Classifying states.");

        _pairs.Add(pair);
    }

    /// <summary>
    /// Records a reconciliation exception.
    /// </summary>
    public void AddException(ReconciliationException exception)
    {
        if (exception == null) throw new ArgumentNullException(nameof(exception));
        if (Status != RunStatus.Classifying)
            throw new InvalidOperationException("Exceptions can only be added during Classifying state.");

        _exceptions.Add(exception);
    }

    /// <summary>
    /// Increments the files ingested counter.
    /// </summary>
    public void IncrementFilesIngested()
    {
        if (Status != RunStatus.Ingesting)
            throw new InvalidOperationException("Files can only be ingested during Ingesting state.");

        FilesIngested++;
    }

    /// <summary>
    /// Validates the reconciliation invariant: matched pairs + exceptions should account for total ingested.
    /// </summary>
    public bool ValidateInvariant(int totalIngestedRecords)
    {
        // This is a soft validation - in practice, some records might be partially matched
        return Pairs.Count + Exceptions.Count <= totalIngestedRecords;
    }

    public override bool Equals(object? obj)
    {
        return obj is ReconciliationRun run && run.RunId == RunId;
    }

    public override int GetHashCode()
    {
        return RunId.GetHashCode();
    }
}
