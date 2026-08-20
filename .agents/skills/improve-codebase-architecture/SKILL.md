---
name: improve-codebase-architecture
description: Find deepening opportunities — refactors that turn shallow modules into deep ones for testability and AI-navigability. Use when the user wants architecture assessment, refactoring candidates, or module consolidation.
disable-model-invocation: true
---

# Improve Codebase Architecture

Surface architectural friction and propose **deepening** refactors. Vocabulary (module, interface, seam, adapter, depth, leverage, locality) is defined in [LANGUAGE.md](LANGUAGE.md). Key tests: the **deletion test** (deleting a pass-through makes complexity vanish; deleting a deep module pushes it into N callers) and **one adapter = hypothetical seam, two adapters = real seam**.

## Process

1. **Context**: read the architecture rules in AGENTS.md and the target area's code. Use the project's own terminology (domain, usecase, adapter, handler, repository), not generic glossaries.
2. **Explore** — lightest mode that covers the request (direct inspection for a quick assessment; a read-only exploration subagent for a formal review). Note friction: bouncing between many small modules to follow one concept; interfaces nearly as complex as implementations; pure functions extracted for testability while the real bugs hide at call sites; layer divergence from the prescribed conventions; areas untestable through their current interface. Apply the deletion test to every shallow suspect.
3. **Candidates**: quick assessment → concise findings in conversation. Formal review → `docs/architecture-reviews/YYYY-MM-DD-<topic>.md`, one section per candidate: Files, Problem (friction evidence), Solution (plain English), Benefits (locality / leverage / testability), recommendation strength (`Strong` / `Worth exploring` / `Speculative`), ending with a **Top recommendation**. Diagrams only when clearer than prose.
4. **Grilling**: after the user picks a candidate, walk the design tree — constraints, dependencies, the deepened module's shape, what sits behind the seam, which tests survive. Use `question` for structured choices. Do not design the full replacement interface until a candidate is selected. If the user rejects a candidate with a load-bearing reason, record it in a comment or rule so future reviews do not re-suggest it.
5. **Handoff**: agreed design → skill `plan-execute` for an implementation plan. New tests target the deepened module's interface; delete old shallow-module tests the new interface subsumes.

## Dependency categories (deepening safety)

| Dependency | Strategy |
|------------|----------|
| In-process | Merge freely; test through the new interface |
| Local-substitutable (in-memory stand-in exists) | Deepen; the seam stays internal |
| Remote but owned | Port at the seam; production adapter + in-memory test adapter |
| True external (third-party) | Injected port; mock adapter in tests |

Replace, don't layer: new tests assert observable outcomes at the interface; old shallow-module tests become waste once interface tests exist. A test that must change when the implementation changes is testing past the interface.
