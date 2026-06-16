---
name: improve-codebase-architecture
description: Find deepening opportunities in a codebase, informed by the project's architecture rules. Use when the user wants to improve architecture, find refactoring opportunities, consolidate tightly-coupled modules, or make a codebase more testable and AI-navigable.
disable-model-invocation: true
---

# Improve Codebase Architecture

Surface architectural friction and propose **deepening opportunities** — refactors that turn shallow modules into deep ones. The aim is testability and AI-navigability.

## Glossary

Use these terms exactly in every suggestion. Consistent language is the point — don't drift into "component," "service," "API," or "boundary." Full definitions in [LANGUAGE.md](LANGUAGE.md).

- **Module** — anything with an interface and an implementation (function, class, package, slice).
- **Interface** — everything a caller must know to use the module: types, invariants, error modes, ordering, config. Not just the type signature.
- **Implementation** — the code inside.
- **Depth** — leverage at the interface: a lot of behaviour behind a small interface. **Deep** = high leverage. **Shallow** = interface nearly as complex as the implementation.
- **Seam** — where an interface lives; a place behaviour can be altered without editing in place. (Use this, not "boundary.")
- **Adapter** — a concrete thing satisfying an interface at a seam.
- **Leverage** — what callers get from depth.
- **Locality** — what maintainers get from depth: change, bugs, knowledge concentrated in one place.

Key principles (see [LANGUAGE.md](LANGUAGE.md) for the full list):

- **Deletion test**: imagine deleting the module. If complexity vanishes, it was a pass-through. If complexity reappears across N callers, it was earning its keep.
- **The interface is the test surface.**
- **One adapter = hypothetical seam. Two adapters = real seam.**

## Process

### 1. Gather Project Context

Read the project's architecture rules (e.g. `.cursor/rules/` if present, otherwise `.pi/` docs, ADRs, or `README.md`/`ARCHITECTURE.md`) to understand:

- Project structure and layering conventions (DDD + Clean Architecture)
- Naming conventions and code style
- Domain modules and their responsibilities
- Layer-specific patterns (domain, usecase, adapter, repository, gateway, handler)

Key documents to read (use whichever exist in this project):
- Project overview / architecture rules (e.g. `00-project-overview.mdc`, `README.md`, `ARCHITECTURE.md`) — architecture overview, domain modules, directory structure
- Code style / best practices (e.g. `01-code-style.mdc`, `02-best-practices.mdc`) — naming, error handling, context passing
- Layer rules (e.g. `10-domain-layer.mdc` through `15-task-layer.mdc`) — layer-specific constraints

Use these rules as the domain vocabulary when naming modules and describing seams.

### 2. Explore

Use the `subagent` tool with `subagent_type="Explore"` to walk the codebase. Don't follow rigid heuristics — explore organically and note where you experience friction:

- Where does understanding one concept require bouncing between many small modules?
- Where are modules **shallow** — interface nearly as complex as the implementation?
- Where have pure functions been extracted just for testability, but the real bugs hide in how they're called (no **locality**)?
- Where do tightly-coupled modules leak across their seams?
- Which parts of the codebase are untested, or hard to test through their current interface?
- Where do the actual layer boundaries diverge from the project's prescribed conventions?

Apply the **deletion test** to anything you suspect is shallow: would deleting it concentrate complexity, or just move it? A "yes, concentrates" is the signal you want.

### 3. Present candidates

Present findings as a markdown report using the `write` tool. Save it to `docs/architecture-reviews/YYYY-MM-DD-<topic>.md` (or another location the user prefers). If the project has a Canvas/visual-report workflow, you may render the same content there; otherwise keep the report in markdown with Mermaid diagrams.

The report should include for each candidate:

- **Files** — which files/modules are involved
- **Problem** — why the current architecture is causing friction
- **Solution** — plain English description of what would change
- **Benefits** — explained in terms of locality and leverage, and how tests would improve
- **Before / After diagram** — side-by-side, illustrating the shallowness and the deepening
- **Recommendation strength** — one of `Strong`, `Worth exploring`, `Speculative`

End the report with a **Top recommendation** section: which candidate you'd tackle first and why.

**Use the project's architecture vocabulary for the domain** (e.g. "the Certificate usecase manager" not "the cert handler"), and [LANGUAGE.md](LANGUAGE.md) vocabulary for the architecture.

Do NOT propose interfaces yet. After the report is presented, ask the user: "Which of these would you like to explore?"

### 4. Grilling loop

Once the user picks a candidate, drop into a grilling conversation. Walk the design tree with them — constraints, dependencies, the shape of the deepened module, what sits behind the seam, what tests survive.

Use the `question` tool to present design choices as structured options when appropriate.

Side effects happen inline as decisions crystallize:

- **Naming a deepened module?** Use vocabulary consistent with the project's conventions — PascalCase for structs/interfaces, snake_case for files.
- **User rejects the candidate with a load-bearing reason?** Suggest recording it in a comment or rule so future reviews don't re-suggest it.
- **Want to explore alternative interfaces for the deepened module?** See [INTERFACE-DESIGN.md](INTERFACE-DESIGN.md).

### 5. Plan the refactor

Once the design is agreed, hand off to the `writing-plans` skill to create an implementation plan. The plan should:

- Follow the vertical-slice structure from `writing-plans`
- Respect the project's layer conventions
- Include TDD steps with tests at the deepened module's interface
- Delete old shallow-module tests that the new interface subsumes

## Diagram Patterns

When creating the report, use these visual patterns (markdown-friendly, or a simple HTML file if the project already supports it):

- **Mermaid graph** — for dependency/call flow diagrams
- **Cross-section** — stacked horizontal bands showing layers a call passes through
- **Mass diagram** — two rectangles per module (interface vs implementation area)
- **Call-graph collapse** — before: tree of calls; after: one deep module

See [MARKDOWN-REPORT.md](MARKDOWN-REPORT.md) for detailed diagram pattern guidance. If your project has a Canvas workflow, you can render the same structures there; otherwise stay in markdown/Mermaid.

## Tone

Plain English, concise — but the architectural nouns and verbs come straight from [LANGUAGE.md](LANGUAGE.md). Concision is not an excuse to drift.

**Use exactly:** module, interface, implementation, depth, deep, shallow, seam, adapter, leverage, locality.

**Never substitute:** component, service, unit (for module) · API, signature (for interface) · boundary (for seam) · layer, wrapper (for module, when you mean module).
