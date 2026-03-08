using ReconciliationEngine.Core.Domain.Aggregates;
using ReconciliationEngine.Core.Domain.Enums;
using ReconciliationEngine.Core.Domain.ValueObjects;
using ReconciliationEngine.Core.Interfaces;

namespace ReconciliationEngine.Infra.Persistence.Repositories;

/// <summary>
/// Repository implementation for ReconciliationRun aggregate root.
/// </summary>
public class ReconciliationRunRepository : IReconciliationRunRepository
{
    private readonly ReconciliationDbContext _dbContext;

    public ReconciliationRunRepository(ReconciliationDbContext dbContext)
    {
        _dbContext = dbContext ?? throw new ArgumentNullException(nameof(dbContext));
    }

    public async Task SaveAsync(ReconciliationRun run, CancellationToken cancellationToken = default)
    {
        if (run == null) throw new ArgumentNullException(nameof(run));

        // TODO: Implement save logic using EF Core
        await Task.CompletedTask;
    }

    public async Task<ReconciliationRun?> GetByIdAsync(Guid runId, CancellationToken cancellationToken = default)
    {
        // TODO: Implement retrieval by ID
        return await Task.FromResult<ReconciliationRun?>(null);
    }

    public async Task<IList<ReconciliationRun>> GetByDateRangeAsync(DateRange range, CancellationToken cancellationToken = default)
    {
        if (range == null) throw new ArgumentNullException(nameof(range));

        // TODO: Implement date range query
        return await Task.FromResult<IList<ReconciliationRun>>(new List<ReconciliationRun>());
    }

    public async Task<ReconciliationRun?> GetLatestByStatusAsync(RunStatus status, CancellationToken cancellationToken = default)
    {
        // TODO: Implement latest by status query
        return await Task.FromResult<ReconciliationRun?>(null);
    }
}
