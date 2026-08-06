# Interface Design

Use this process when the user wants to compare alternative interfaces for a selected architecture candidate. Preserve kfinops domain/layer terminology and use depth, seam, adapter, leverage, and locality when helpful.

## 1. Frame the problem

Explain:

- constraints every interface must satisfy;
- dependencies and their categories from [DEEPENING.md](DEEPENING.md);
- relevant `.cursor/rules/` and existing project conventions;
- a small illustrative sketch in Go that grounds constraints without becoming the proposal.

## 2. Generate alternatives

Default to two credible alternatives in the current context. Do not manufacture a weak design merely to increase option count.

For a long-lived public interface, irreversible migration, disputed architecture decision, or high uncertainty, dispatch 2–3 read-only design subagents in parallel. Give each the same concrete brief—files, coupling, dependency categories, layer rules, and what sits behind the seam—with relevant emphases such as:

- minimize the interface and maximize leverage;
- optimize the common caller and safe defaults;
- preserve a required extension or ports-and-adapters seam.

Do not optimize for generic flexibility without concrete consumers.

Each design should provide:

1. interface plus invariants, ordering, errors, and configuration;
2. a Go usage example;
3. behavior hidden behind the seam;
4. dependency and adapter strategy;
5. trade-offs in leverage, locality, migration, and testability.

## 3. Compare

Compare depth, locality, seam placement, compatibility, extension cost, and testability. Recommend one design—or a deliberate hybrid—with evidence. Do not prefer a design merely because it introduces more abstractions.
