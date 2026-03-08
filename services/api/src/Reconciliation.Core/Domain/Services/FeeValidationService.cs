using ReconciliationEngine.Core.Domain.Entities;
using ReconciliationEngine.Core.Domain.ValueObjects;

namespace ReconciliationEngine.Core.Domain.Services;

/// <summary>
/// Stub implementation of fee validation service.
/// Validates actual fees against contracted rates.
/// </summary>
public class FeeValidationService : IFeeValidationService
{
    public Task<(bool IsValid, Money? DeltaCentavos)> ValidateFeeAsync(
        TransactionRecord transaction,
        AcquirerContract contract,
        CancellationToken cancellationToken = default)
    {
        // TODO: Implement fee validation logic
        // Compare transaction.Fee against contract.FeeSchedules
        // Calculate delta if divergence exists
        return Task.FromResult((true, (Money?)null));
    }
}
