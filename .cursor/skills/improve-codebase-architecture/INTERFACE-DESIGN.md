# Interface Design

Use this process when the user wants to compare alternative interfaces for a selected deepening candidate. It follows “design it twice”: the first plausible interface is rarely the only useful one.

Use the vocabulary from [LANGUAGE.md](LANGUAGE.md): module, interface, seam, adapter, leverage, and locality.

## 1. Frame the Problem

Before dispatching subagents, explain:

- the constraints every interface must satisfy;
- the dependencies and their categories from [DEEPENING.md](DEEPENING.md);
- a small illustrative sketch in the project’s actual language or notation.

This is context, not a proposed solution.

## 2. Explore Alternatives

Dispatch at least three read-only design subagents in parallel. Give each a concrete brief with relevant files, coupling details, dependency categories, existing project conventions, and a distinct objective:

1. minimize the interface and maximize leverage per entry point;
2. optimize for extension and variation;
3. optimize for the most common caller;
4. when remote or external dependencies matter, design the port-and-adapter strategy.

Each subagent returns:

1. the interface, including invariants, ordering, error modes, and configuration requirements;
2. an example of caller usage in the project’s language or notation;
3. what the implementation hides behind the seam;
4. the dependency and adapter strategy;
5. trade-offs in leverage, locality, and complexity.

## 3. Compare

Present alternatives sequentially, then compare depth, locality, seam placement, extension cost, and testability. Recommend one design—or a deliberate hybrid—with reasons. Do not select a design merely because it introduces more abstractions.
