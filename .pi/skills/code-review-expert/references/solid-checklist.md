# SOLID Smell Prompts

## SRP (Single Responsibility)

- File owns unrelated concerns (e.g., HTTP + DB + domain rules in one file)
- Large class/module with low cohesion or multiple reasons to change
- Functions that orchestrate many unrelated steps
- God objects that know too much about the system
- **Ask**: "What is the single reason this module would change?"

## OCP (Open/Closed)

- Adding a new behavior requires editing many switch/if blocks
- Feature growth requires modifying core logic rather than extending
- No plugin/strategy/hook points for variation
- **Ask**: "Can I add a new variant without touching existing code?"

**What Not to Flag**: A switch over an external protocol, wire format, or closed enum is not automatically missing polymorphism. Not every variation point needs a strategy pattern.

## LSP (Liskov Substitution)

- Subclass checks for concrete type or throws for base method
- Overridden methods weaken preconditions or strengthen postconditions
- Subclass ignores or no-ops parent behavior
- **Ask**: "Can I substitute any subclass without the caller knowing?"

## ISP (Interface Segregation)

- Interfaces with many methods, most unused by implementers
- Callers depend on broad interfaces for narrow needs
- Empty/stub implementations of interface methods
- **Ask**: "Do all implementers use all methods?"

## DIP (Dependency Inversion)

- High-level logic depends on concrete IO, storage, or network types
- Hard-coded implementations instead of abstractions or injection
- Import chains that couple business logic to infrastructure
- **Ask**: "Can I swap the implementation without changing business logic?"

**What Not to Flag**: Composition roots (Wire injectors, main.go) wiring concrete dependencies is not a DIP violation — that is their job. Thin adapter layers may import both directions when they are explicitly boundary glue. High fan-out in an orchestration layer is not automatically dependency disorder.

## Change Propagation & Coupling

Beyond SOLID, check whether the diff increases the blast radius of future changes.

- **Information Leakage**: a design decision (e.g., file format, protocol detail, data shape) is encoded in more than one module — changing it requires coordinated edits across places even though only one module should "own" the concept
- **Hyrum's Law**: the diff exposes observable behavior (error message text, call ordering, undocumented side effects) that callers will depend on as an implicit contract, even though it was never guaranteed by the API
- **Orthogonality violation**: adding a new variant of one dimension (e.g., a new payment type) forces edits in unrelated dimensions (logging, caching, notification) — the dimensions are not independent
- **High change propagation radius**: modifying one feature requires touching > 3 files in unrelated modules (Shotgun Surgery at the architecture level)

- **Ask**: "If I change this one decision, how many unrelated files need to change?"
- **Ask**: "Is this exposing implementation details that callers will accidentally depend on?"
- **Ask**: "Can I add a new variant of X without touching Y?"

### What Not to Flag

- A composition root wiring concrete dependencies is not a coupling problem by itself
- Similar edits inside one bounded context may be normal coordinated change, not shotgun surgery
- A stable public API with intentionally supported behavior is not automatically Hyrum's Law debt
- Adapter modules may depend on both domain and infrastructure when they explicitly translate across the boundary

---

## Common Code Smells (Beyond SOLID)

| Smell | Signs |
|-------|-------|
| **Long method** | Function > 30 lines, multiple levels of nesting, mixing multiple levels of abstraction in one function |
| **Feature envy** | Method uses more data from another class than its own |
| **Data clumps** | Same group of parameters passed together repeatedly |
| **Primitive obsession** | Domain concepts expressed as raw primitive types (`string email`, `int orderId`, `float64 money`) instead of purpose-built value types — forces callers to know which string is an email and which is a name |
| **Flag argument** | A boolean parameter that makes the function do two fundamentally different things depending on its value — a sign the function has two responsibilities and should be split |
| **Shallow module** | The interface or documentation of a component is more complex relative to the functionality it provides — the abstraction is not earning its keep |
| **Shotgun surgery** | One change requires edits across many files |
| **Divergent change** | One file changes for many unrelated reasons |
| **Dead code** | Unreachable or never-called code |
| **Speculative generality** | Abstractions for hypothetical future needs |
| **Magic numbers/strings** | Hardcoded values without named constants |

### Cognitive Overload Thresholds

Use these as hints, not hard rules — context and readability matter more than line counts:

- Function > 50 lines or nesting > 5 levels: likely Critical
- Function 20–50 lines or nesting 4–5 levels: likely Warning
- Parameter list with > 4 parameters: consider introducing a config/options struct
- Boolean expressions with 3+ combined conditions: extract into a named predicate
- Train-wreck chains (`a.GetB().GetC().DoD()`): Law of Demeter violation, consider restructuring

---

## Domain Model Health

For DDD-based projects, check whether the diff preserves or degrades domain model integrity.

| Smell | Signs |
|-------|-------|
| **Anemic Domain Model** | Domain objects have only getters/setters; all business logic lives in service/usecase layer instead of on the entity that owns the data |
| **Ubiquitous Language drift** | Code names (variables, types, methods) diverge from what business stakeholders call the same concept |
| **Bounded Context violation** | Direct cross-context imports without an anti-corruption layer or translation |
| **Value Object as Entity** | A concept defined entirely by its attributes (e.g., Money, Email, Address) is given a mutable ID and lifecycle instead of being immutable and replaceable |
| **Feature Envy in wrong layer** | Domain logic appears in adapter/handler/delivery layer rather than in the domain or usecase layer that owns the concept |

- **Ask**: "Does the code faithfully represent the problem domain, or does it model schemas and infrastructure instead?"
- **Ask**: "Is business logic living on the entity that owns the data, or scattered across service layers?"

### What Not to Flag

- CRUD-heavy workflows may legitimately use transaction scripts instead of rich domain objects
- DTOs, persistence records, and API payload models are allowed to be data-only structures
- Thin entities are acceptable when the business domain itself is simple
- Shared infrastructure language (e.g., `Repository`, `Handler`) should not be mistaken for domain drift

---

## Refactor Heuristics

1. **Split by responsibility, not by size** - A small file can still violate SRP
2. **Introduce abstraction only when needed** - Wait for the second use case
3. **Keep refactors incremental** - Isolate behavior before moving
4. **Preserve behavior first** - Add tests before restructuring
5. **Name things by intent** - If naming is hard, the abstraction might be wrong
6. **Prefer composition over inheritance** - Inheritance creates tight coupling
7. **Make illegal states unrepresentable** - Use types to enforce invariants

---

## Code-Judo Mindset

Be ambitious about structural simplification. Do not stop at "this could be a bit cleaner."

### Principles

- Look for opportunities to **reframe** the change so that whole branches, helpers, modes, conditionals, or layers disappear entirely.
- Prefer the solution that makes the code feel **inevitable in hindsight**.
- Assume there is often a re-organization that uses the existing architecture more effectively and makes the change dramatically simpler.
- If you see a path to **delete complexity** rather than rearrange it, push hard for that path.
- Do not be satisfied with a merely cleaner version of the same messy idea if there is a plausible path to a much simpler idea.

### Code-Judo Remedies

When a SOLID violation or code smell is found, prefer these high-impact remedies over surface-level fixes:

| Instead of... | Try... |
|---------------|--------|
| Polishing a layer of indirection | Deleting the layer entirely |
| Centralizing conditionals | Reframing the state model so conditionals disappear |
| Adding a helper to wrap messy code | Changing the ownership boundary so the feature becomes a natural extension |
| Adding special-case branches | Turning special-case logic into a simpler default flow with fewer exceptions |
| Moving complexity to a new location | Restructuring so the complexity no longer exists |
| Renaming things for clarity | Questioning whether the abstraction should exist at all |

### Questions to Ask Before Proposing a Fix

- Is there a code-judo move that would make this dramatically simpler?
- Can we reframe this so fewer concepts, branches, or helper layers are needed?
- Does the proposed fix delete complexity, or merely move it?
- Is the refactor making the number of concepts a reader must hold smaller, or just rearranging them?
