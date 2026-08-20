# Anti False Positives

Cases that look like smells but are usually fine. Do not flag these by default.

## SOLID / DI

- Composition roots and dependency-injection configuration wiring concrete dependencies — that is their job, not a DIP violation.
- Thin adapters importing domain and infrastructure as explicit boundary glue.
- A switch over an external protocol, wire format, or closed enum is not automatically missing polymorphism.
- High fan-out in an orchestration / use-case layer is not automatically dependency disorder.

## Coupling / change propagation

- Coordinated edits inside one bounded context are normal, not shotgun surgery.
- Stable public API behavior that is intentionally supported is not Hyrum's Law debt.
- Repetition across separate bounded contexts may be clearer than a forced shared dependency.
- Shared protocol constants at explicit module boundaries can beat a shared package when coupling cost is higher.

## Domain / DDD

- CRUD / transaction-script flows need not be rich domain models.
- DTOs, persistence records, and API payloads may be data-only.
- Thin entities are fine when the domain itself is simple.
- Shared infra names (`Repository`, `Handler`) are not ubiquitous-language drift.

## Structure / size

- Linear, well-named, single-purpose functions are fine even at 30–50 lines.
- Temporary duplication during an active extraction/migration is not automatically debt.
- Thin wrappers that absorb vendor churn or hide instability may be justified.
- A 40-line function with clear names and guard clauses is not automatically high cognitive load.
- Internal implementation detail hidden behind a deep, simple module boundary is not a shallow-module problem.

## Error handling

- A function that intentionally wraps and re-throws errors at a boundary is not "swallowing" them.
- Logging an error AND returning it is acceptable at module boundaries — the caller still sees it.
- Error messages that include internal identifiers for debugging are acceptable as long as they don't expose secrets or raw stack traces.

## Concurrency

- A `sync.Mutex` used briefly in a request-scoped struct is not automatically a bottleneck.
- Intentionally fire-and-forget background work with its own error logging can be acceptable for non-critical side effects (metrics, audit logs).
- Sequential operations where each depends on the previous result are not "pointless serialization."
