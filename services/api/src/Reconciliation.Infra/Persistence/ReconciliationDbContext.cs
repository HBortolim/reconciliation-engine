using Microsoft.EntityFrameworkCore;
using ReconciliationEngine.Core.Domain.Entities;
using ReconciliationEngine.Core.Domain.Aggregates;

namespace ReconciliationEngine.Infra.Persistence;

/// <summary>
/// EF Core DbContext for reconciliation domain.
/// </summary>
public class ReconciliationDbContext : DbContext
{
    public ReconciliationDbContext(DbContextOptions<ReconciliationDbContext> options) : base(options)
    {
    }

    public DbSet<TransactionRecord> TransactionRecords { get; set; } = null!;
    public DbSet<ReconciliationRun> ReconciliationRuns { get; set; } = null!;
    public DbSet<ReconciliationPair> ReconciliationPairs { get; set; } = null!;
    public DbSet<ReconciliationException> ReconciliationExceptions { get; set; } = null!;
    public DbSet<AcquirerContract> AcquirerContracts { get; set; } = null!;

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        base.OnModelCreating(modelBuilder);

        // TODO: Configure entity mappings, value object conversions, indexes, etc.
    }
}
