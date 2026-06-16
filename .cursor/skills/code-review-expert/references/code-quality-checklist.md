# Code Quality Checklist

## Structural Quality & Maintainability

### Non-Negotiable Rules

1. **Do not let a PR push a file from under 1k lines to over 1k lines without a very strong reason.**
   - Treat this as a strong code-quality smell by default.
   - Prefer extracting helpers, subcomponents, modules, or local abstractions.
   - If the diff crosses that threshold, explicitly ask whether the code should be decomposed first.
   - Only waive if there is a compelling structural reason and the file is still clearly organized.

2. **Do not allow random spaghetti growth in existing code.**
   - Be highly suspicious of new ad-hoc conditionals, scattered special cases, or one-off branches.
   - If a change adds "weird if statements in random places", treat that as a design problem, not a stylistic nit.
   - Prefer pushing the logic into a dedicated abstraction, helper, state machine, or separate module.
   - Call out changes that make surrounding code harder to reason about, even if they technically work.

3. **Bias toward cleaning the design, not just accepting working code.**
   - If behavior can stay the same while the structure becomes meaningfully cleaner, push for the cleaner version.
   - Do not rubber-stamp "it works" implementations that leave the codebase messier.
   - Strongly prefer simplifications that remove moving pieces altogether over refactors that merely spread the same complexity around.

4. **Prefer direct, boring, maintainable code over hacky or magical code.**
   - Treat brittle, ad-hoc, or "magic" behavior as a code-quality problem.
   - Be skeptical of generic mechanisms that hide simple data-shape assumptions.
   - Flag thin abstractions, identity wrappers, or pass-through helpers that add indirection without buying clarity.

5. **Push hard on type and boundary cleanliness when they affect maintainability.**
   - Question unnecessary `interface{}`, `any`, or cast-heavy code when a clearer type boundary could exist.
   - Prefer explicit typed models or shared contracts over loosely-shaped ad-hoc objects.
   - If a branch relies on silent fallback to paper over an unclear invariant, ask whether the boundary should be made explicit.

6. **Keep logic in the canonical layer and reuse existing helpers.**
   - Call out feature logic leaking into shared paths or implementation details leaking through APIs.
   - Prefer existing canonical utilities/helpers over bespoke one-offs.
   - Push code toward the right package, service, or module instead of normalizing architectural drift.

7. **Treat unnecessary sequential orchestration and non-atomic updates as design smells.**
   - If independent work is serialized for no good reason, ask whether the flow should run in parallel.
   - If related updates can leave state half-applied, push for a more atomic structure.
   - Flag avoidable orchestration complexity that makes the implementation more brittle.

### What to Flag Aggressively

- A file crossing 1000 lines due to the PR, especially if the new code could be split out.
- New conditionals bolted onto unrelated code paths.
- One-off booleans, nullable modes, or flags that complicate existing control flow.
- Feature-specific logic leaking into general-purpose modules.
- Generic "magic" handling that hides simple structure and makes the code harder to reason about.
- Thin wrappers or identity abstractions that add indirection without simplifying anything.
- Unnecessary casts, `any`, `interface{}`, or optional params that muddy the real contract.
- **Knowledge duplication** (DRY is about decisions, not code lines):
  - Same business decision expressed in multiple places — these copies will drift apart silently
  - Same domain concept named differently in different parts of the codebase (e.g., `user`/`account`/`member`/`customer` all referring to one entity)
  - Configuration values repeated as literals in multiple files instead of a single source of truth
  - Two modules independently implementing the same algorithm or validation rule
  - Copy-pasted logic instead of extracted helpers
- Narrow edge-case handling implemented in the middle of an already busy function.
- Refactors that technically pass tests but make the code less modular or less readable.
- "Temporary" branching that is likely to become permanent debt.
- Bespoke helpers where the codebase already has a canonical utility for the job.
- Logic added in the wrong layer/package when it should live somewhere more central.
- Sequential flow where obviously independent work could be simpler with parallel execution.
- Partial-update logic that leaves state less atomic than necessary.
- A complicated implementation where a cleaner reframing could delete whole categories of complexity.
- Refactors that move code around but fail to reduce the number of concepts a reader must hold in their head.

### Preferred Remedies

When you identify a structural problem, prefer suggestions like:

- Delete a whole layer of indirection rather than polishing it.
- Reframe the state model so conditionals disappear instead of getting centralized.
- Change the ownership boundary so the feature becomes a natural extension of an existing abstraction.
- Turn special-case logic into a simpler default flow with fewer exceptions.
- Extract a helper or pure function.
- Split a large file into smaller focused modules.
- Move feature-specific logic behind a dedicated abstraction.
- Replace condition chains with a typed model or explicit dispatcher.
- Separate orchestration from business logic.
- Collapse duplicate branches into a single clearer flow.
- Delete wrappers that do not meaningfully clarify the API.
- Reuse the existing canonical helper instead of introducing a near-duplicate.
- Make type boundaries more explicit so the control flow gets simpler.
- Move the logic to the package/module/layer that already owns the concept.
- Parallelize independent work when that also simplifies the orchestration.
- Restructure related updates into a more atomic flow.

Do not be satisfied with "maybe rename this" feedback when the real issue is structural.

### Questions to Ask
- "Is there a code-judo move that would make this dramatically simpler?"
- "Can this change be reframed so fewer concepts, branches, or helper layers are needed?"
- "Does this improve or worsen the local architecture?"
- "Is this abstraction actually earning its keep, or is it just a wrapper?"
- "Is this logic living in the canonical layer, or did the diff leak details across a boundary?"
- "Is this orchestration more sequential or less atomic than it needs to be?"

### What Not to Flag

- Linear code with clear names and guard clauses is not automatically high cognitive load — a 40-line function can be perfectly readable
- A long routine may be acceptable when it is linear, well-named, and single-purpose
- Internal implementation detail hidden behind a deep, simple module boundary is not a shallow-module problem
- Temporary duplication during an active extraction or migration is not necessarily debt
- Domain-specific terminology should not be flagged if it matches how experts actually speak
- Thin wrappers that absorb vendor churn or hide instability may be justified
- CRUD-heavy workflows may legitimately use simpler patterns than full DDD
- Repetition across separate bounded contexts is not automatically duplicate knowledge — local ownership may be clearer than a shared dependency
- Shared protocol constants repeated at explicit module boundaries may be acceptable when coupling cost exceeds duplication cost

---

## Error Handling

### Anti-patterns to Flag

- **Ignored errors**: Discarding error return values
  ```go
  result, _ := someFunction()  // Silent failure
  _ = db.Save(&entity)         // Ignoring DB error
  ```
- **Log and forget**: Logging error but not returning/handling it
- **Error information leakage**: Stack traces or internal details exposed to users
- **Missing error wrapping**: Use `utils.AppError` instead of `errors.New`
- **Goroutine errors**: Errors in goroutines not propagated to caller

### Best Practices to Check

- [ ] Errors are caught at appropriate boundaries
- [ ] Error messages are user-friendly (no internal details exposed)
- [ ] Errors are logged with sufficient context for debugging
- [ ] Async errors are properly propagated or handled
- [ ] Fallback behavior is defined for recoverable errors
- [ ] Critical errors trigger alerts/monitoring

### Questions to Ask
- "What happens when this operation fails?"
- "Will the caller know something went wrong?"
- "Is there enough context to debug this error?"

---

## Performance & Caching

### CPU-Intensive Operations

- **Expensive operations in hot paths**: Regex compilation, JSON parsing, crypto in loops
- **Blocking main thread**: Sync I/O, heavy computation without worker/async
- **Unnecessary recomputation**: Same calculation done multiple times
- **Missing memoization**: Pure functions called repeatedly with same inputs

### Database & I/O

- **N+1 queries**: Loop that makes a query per item instead of batch
  ```go
  // Bad: N+1
  for _, id := range ids {
      var user User
      db.First(&user, id)
  }
  // Good: Use Preload or batch query
  db.Preload("Users").Find(&orders)
  db.Where("id IN ?", ids).Find(&users)
  ```
- **Missing indexes**: Queries on unindexed columns
- **Over-fetching**: SELECT * when only few columns needed
- **No pagination**: Loading entire dataset into memory

### Caching Issues

- **Missing cache for expensive operations**: Repeated API calls, DB queries, computations
- **Cache without TTL**: Stale data served indefinitely
- **Cache without invalidation strategy**: Data updated but cache not cleared
- **Cache key collisions**: Insufficient key uniqueness
- **Caching user-specific data globally**: Security/privacy issue

### Memory

- **Unbounded collections**: Arrays/maps that grow without limit
- **Large object retention**: Holding references preventing GC
- **String concatenation in loops**: Use `strings.Builder` instead
- **Loading large files entirely**: Use streaming instead

### Questions to Ask
- "What's the time complexity of this operation?"
- "How does this behave with 10x/100x data?"
- "Is this result cacheable? Should it be?"
- "Can this be batched instead of one-by-one?"

---

## Boundary Conditions

### Nil/Zero Value Handling

- **Missing nil checks**: Accessing fields on potentially nil pointers
- **Zero value confusion**: `if value != ""` when empty string is valid
- **Pointer vs value**: Returning pointer when nil is not a valid state

### Empty Collections

- **Empty array not handled**: Code assumes array has items
- **Empty object edge case**: `for...in` or `Object.keys` on empty object
- **First/last element access**: `arr[0]` or `arr[arr.length-1]` without length check

### Numeric Boundaries

- **Division by zero**: Missing check before division
- **Integer overflow**: Large numbers exceeding safe integer range
- **Floating point comparison**: Using `===` instead of epsilon comparison
- **Negative values**: Index or count that shouldn't be negative
- **Off-by-one errors**: Loop bounds, array slicing, pagination

### String Boundaries

- **Empty string**: Not handled as edge case
- **Whitespace-only string**: Passes truthy check but is effectively empty
- **Very long strings**: No length limits causing memory/display issues
- **Unicode edge cases**: Emoji, RTL text, combining characters

### Common Patterns to Flag

```go
// Dangerous: no nil check
name := user.Profile.Name

// Dangerous: slice access without check
first := items[0]

// Dangerous: division without check
avg := total / count

// Dangerous: zero value check excludes valid values
if value != 0 { ... }  // fails when 0 is valid
```

### Questions to Ask
- "What if this pointer is nil?"
- "What if this slice is empty?"
- "What's the valid range for this number?"
- "What happens at the boundaries (0, -1, MaxInt)?"
