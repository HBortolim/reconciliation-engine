using ReconciliationEngine.Core.Domain.Aggregates;
using ReconciliationEngine.Core.Domain.Enums;
using ReconciliationEngine.Core.Domain.ValueObjects;

namespace ReconciliationEngine.Core.Interfaces;

/// <summary>
/// Repository interface for ReconciliationRun aggregate root.
/// </summary>
public interface IReconciliationRunRepository
{
    /// <summary>
    /// Saves a reconciliation run.
    /// </summary>
    Task SaveAsync(ReconciliationRun run, CancellationToken cancellationToken = default);

    /// <summary>
    /// Retrieves a run by ID.
    /// </summary>
    Task<ReconciliationRun?> GetByIdAsync(Guid runId, CancellationToken cancellationToken = default);

    /// <summary>
    /// Retrieves runs within a date range.
    /// </summary>
    Task<IList<ReconciliationRun>> GetByDateRangeAsync(DateRange range, CancellationToken cancellationToken = default);

    /// <summary>
    /// Retrieves the latest run with a specific status.
    /// </summary>
    Task<ReconciliationRun?> GetLatestByStatusAsync(RunStatus status, CancellationToken cancellationToken = default);
}
