using ReconciliationEngine.Core.Domain.Entities;
using ReconciliationEngine.Core.Domain.ValueObjects;

namespace ReconciliationEngine.Core.Interfaces;

/// <summary>
/// Repository interface for ReconciliationException.
/// </summary>
public interface IExceptionRepository
{
    /// <summary>
    /// Saves an exception.
    /// </summary>
    Task SaveAsync(ReconciliationException exception, CancellationToken cancellationToken = default);

    /// <summary>
    /// Retrieves an exception by ID.
    /// </summary>
    Task<ReconciliationException?> GetByIdAsync(Guid exceptionId, CancellationToken cancellationToken = default);

    /// <summary>
    /// Retrieves open exceptions by aging (oldest first).
    /// </summary>
    Task<IList<ReconciliationException>> GetOpenByAgingAsync(CancellationToken cancellationToken = default);

    /// <summary>
    /// Retrieves exceptions for a specific run.
    /// </summary>
    Task<IList<ReconciliationException>> GetByRunIdAsync(Guid runId, CancellationToken cancellationToken = default);
}
