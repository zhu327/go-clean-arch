---
name: improve-codebase-architecture
description: Find deepening opportunities in any codebase. Use when the user wants to improve architecture, consolidate tightly coupled modules, improve testability, or make a codebase easier to navigate.
disable-model-invocation: true
---

# Improve Codebase Architecture

Surface architectural friction and propose refactors that make modules deeper: more useful behavior behind a smaller, clearer interface. The aim is testability, locality, and navigability.

Use the vocabulary defined in [LANGUAGE.md](LANGUAGE.md): module, interface, implementation, depth, seam, adapter, leverage, and locality.

## 1. Gather Context

Inspect project instructions, architecture documents, ADRs, source layout, build configuration, and tests. Determine the project’s real module boundaries, vocabulary, conventions, and constraints. If no explicit architecture documentation exists, state that limitation rather than inventing one.

Do not assume a programming language, framework, directory layout, delivery mechanism, report renderer, or architecture style.

## 2. Explore

Use a read-only exploration subagent to inspect the codebase. Look for friction:

- understanding one concept requires bouncing across shallow modules;
- callers must know implementation details that should be hidden;
- tests must cross an internal seam instead of the module’s external interface;
- behavior and knowledge are duplicated across callers;
- a seam leaks implementation concerns or has no meaningful variation;
- existing module boundaries diverge from documented project architecture.

Apply the deletion test: if removing a module makes complexity vanish, it was probably pass-through indirection; if complexity would reappear across callers, it may be earning its keep.

## 3. Report Candidates

Save a Markdown report to the project’s existing architecture-review location, or `docs/architecture-reviews/YYYY-MM-DD-<topic>.md` if none exists. For every candidate include:

- files/modules involved;
- problem and evidence;
- proposed change in plain language;
- expected leverage and locality gains;
- before/after diagram using Mermaid when useful;
- recommendation strength: **Strong**, **Worth exploring**, or **Speculative**.

End with a top recommendation and its rationale. Do not propose detailed interfaces until the user selects a candidate.

## 4. Design Discussion

After selection, use structured questions to clarify constraints, dependencies, seam placement, hidden implementation details, and the tests that survive the refactor. Use [DEEPENING.md](DEEPENING.md) and [INTERFACE-DESIGN.md](INTERFACE-DESIGN.md) as needed.

## 5. Plan

Once the design is agreed, hand off to `writing-plans`. The plan should preserve the project’s conventions, test through the selected module’s interface, and remove superseded shallow-module tests when safe.
