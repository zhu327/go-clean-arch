---
name: improve-codebase-architecture
description: Find evidence-backed deepening opportunities in kfinops. Use for architecture assessment, consolidation of tightly coupled modules, testability improvements, or navigation friction; supports quick conversational assessment and formal review.
disable-model-invocation: true
---

# Improve Codebase Architecture

Surface architectural friction and propose refactors that place useful behavior behind smaller, clearer interfaces. Preserve kfinops domain and layer terminology; use depth, seam, leverage, and locality when they add precision.

## 1. Gather context

Read relevant `.cursor/rules/`, source, tests, recent history, and architecture records. Determine actual DDD/Clean Architecture ownership, dependency direction, conventions, and constraints. Do not invent a module boundary or generic abstraction that conflicts with the repository.

## 2. Choose review depth

- **Quick assessment:** inspect the relevant domain/layers directly and return a few evidence-backed candidates in conversation.
- **Formal review:** for broad scope or a durable artifact, use a read-only exploration subagent and save a report to the project's existing architecture-review location, or `docs/architecture-reviews/YYYY-MM-DD-<topic>.md`.

Do not require a report, subagent, or diagram for a focused question.

## 3. Explore friction

Look for:

- one domain concept requiring jumps across many shallow modules;
- callers knowing repository/gateway details that a usecase interface should hide;
- business rules leaking into handlers, repositories, gateways, tasks, or Wire setup;
- tests crossing internal details instead of a stable behavior surface;
- duplicated decisions across domains or adapters;
- pass-through interfaces with no meaningful variation;
- actual dependencies diverging from `.cursor/rules/`.

Apply the deletion test: if removing a module makes complexity vanish, it may be pass-through indirection; if complexity reappears across callers, it may provide leverage.

## 4. Present candidates

For each candidate include files, evidence, proposed change, expected locality/leverage or test-surface gain, risks, and recommendation strength: **Strong**, **Worth exploring**, or **Speculative**. Use Mermaid or HTML only when it clarifies relationships better than prose.

Do not design a detailed replacement interface until the user selects a candidate, unless an end-to-end proposal was explicitly requested. When several candidates remain credible, use AskQuestion to choose which to explore.

## 5. Design and plan

Use [DEEPENING.md](DEEPENING.md) for dependency strategy and [INTERFACE-DESIGN.md](INTERFACE-DESIGN.md) when alternatives need comparison. Preserve existing kfinops terms such as domain, usecase, repository, gateway, handler, task, and Wire.

Once a consequential design is agreed, use `writing-plans` when a durable multi-step plan is justified. Test through the selected behavior surface and remove superseded shallow tests only when replacement coverage is credible.
