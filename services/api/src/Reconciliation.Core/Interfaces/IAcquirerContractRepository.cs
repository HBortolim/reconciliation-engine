using ReconciliationEngine.Core.Domain.Entities;

namespace ReconciliationEngine.Core.Interfaces;

/// <summary>
/// Repository interface for AcquirerContract.
/// </summary>
public interface IAcquirerContractRepository
{
    /// <summary>
    /// Saves an acquirer contract.
    /// </summary>
    Task SaveAsync(AcquirerContract contract, CancellationToken cancellationToken = default);

    /// <summary>
    /// Retrieves a contract by ID.
    /// </summary>
    Task<AcquirerContract?> GetByIdAsync(Guid contractId, CancellationToken cancellationToken = default);

    /// <summary>
    /// Retrieves active contracts for an acquirer.
    /// </summary>
    Task<IList<AcquirerContract>> GetActiveByAcquirerAsync(string acquirerId, CancellationToken cancellationToken = default);

    /// <summary>
    /// Retrieves all contracts for an acquirer.
    /// </summary>
    Task<IList<AcquirerContract>> GetAllAsync(string acquirerId, CancellationToken cancellationToken = default);
}
