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

---

## Common Code Smells (Beyond SOLID)

| Smell | Signs |
|-------|-------|
| **Long method** | Function > 30 lines, multiple levels of nesting |
| **Feature envy** | Method uses more data from another class than its own |
| **Data clumps** | Same group of parameters passed together repeatedly |
| **Primitive obsession** | Using strings/numbers instead of domain types |
| **Shotgun surgery** | One change requires edits across many files |
| **Divergent change** | One file changes for many unrelated reasons |
| **Dead code** | Unreachable or never-called code |
| **Speculative generality** | Abstractions for hypothetical future needs |
| **Magic numbers/strings** | Hardcoded values without named constants |

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
