namespace ReconciliationEngine.Core.Domain.Specifications;

/// <summary>
/// Generic specification interface for encapsulating business rules.
/// Supports AND, OR, NOT combinators for rule composition.
/// </summary>
/// <typeparam name="T">The type being evaluated</typeparam>
public interface ISpecification<T>
{
    /// <summary>
    /// Evaluates if the given object satisfies this specification.
    /// </summary>
    bool IsSatisfiedBy(T candidate);

    /// <summary>
    /// Combines this specification with another using AND logic.
    /// </summary>
    ISpecification<T> And(ISpecification<T> other);

    /// <summary>
    /// Combines this specification with another using OR logic.
    /// </summary>
    ISpecification<T> Or(ISpecification<T> other);

    /// <summary>
    /// Negates this specification.
    /// </summary>
    ISpecification<T> Not();
}
