using ReconciliationEngine.Core.Domain.Entities;
using ReconciliationEngine.Core.Domain.ValueObjects;

namespace ReconciliationEngine.Core.Domain.Services;

/// <summary>
/// Domain service interface for validating actual vs. contracted fees.
/// </summary>
public interface IFeeValidationService
{
    /// <summary>
    /// Validates fee against contract for a transaction.
    /// Returns fee delta if divergence exists.
    /// </summary>
    Task<(bool IsValid, Money? DeltaCentavos)> ValidateFeeAsync(
        TransactionRecord transaction,
        AcquirerContract contract,
        CancellationToken cancellationToken = default);
}
