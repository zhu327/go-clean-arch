# Anti False Positives

Cases that look like smells but are usually fine. Do not flag these by default.

## SOLID / DI

- Composition roots (`main`, DI containers, Wire) wiring concrete deps — that is their job, not a DIP violation.
- Thin adapters importing domain + infrastructure as explicit boundary glue.
- Switch over external protocol, wire format, or closed enum — not automatically missing polymorphism.
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
