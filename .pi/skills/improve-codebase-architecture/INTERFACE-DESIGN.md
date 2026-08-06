# Interface Design

When the user wants to explore alternative interfaces for a chosen deepening candidate, compare genuinely credible designs. Use parallel subagents only when the decision is high-impact or independent exploration is likely to improve it.

Uses the vocabulary in [LANGUAGE.md](LANGUAGE.md) — **module**, **interface**, **seam**, **adapter**, **leverage**.

## Process

### 1. Frame the problem space

Before spawning sub-agents, write a user-facing explanation of the problem space for the chosen candidate:

- The constraints any new interface would need to satisfy
- The dependencies it would rely on, and which category they fall into (see [DEEPENING.md](DEEPENING.md))
- A rough illustrative code sketch to ground the constraints — not a proposal, just a way to make the constraints concrete

Show this to the user, then immediately proceed to Step 2. The user reads and thinks while the sub-agents work in parallel.

### 2. Generate alternatives

Default to two credible alternatives produced in the current context. Do not manufacture a weak design merely to increase option count.

For a long-lived public interface, irreversible migration, disputed architecture decision, or highly uncertain design, spawn 2–3 subagents in parallel. Give each the same technical brief (file paths, coupling details, dependency category from [DEEPENING.md](DEEPENING.md), and what sits behind the seam) but a relevant design emphasis, such as:

- minimize the interface and maximize leverage;
- optimize the common caller and safe defaults;
- preserve a required extension or ports-and-adapters seam.

Avoid "maximize flexibility" as a default goal unless concrete future consumers justify it.

Include both [LANGUAGE.md](LANGUAGE.md) vocabulary and the project's architecture conventions in the brief so each sub-agent names things consistently with the architecture language and the project's domain language.

Each sub-agent outputs:

1. Interface (types, methods, params — plus invariants, ordering, error modes)
2. Usage example showing how callers use it
3. What the implementation hides behind the seam
4. Dependency strategy and adapters (see [DEEPENING.md](DEEPENING.md))
5. Trade-offs — where leverage is high, where it's thin

### 3. Present and compare

Present designs sequentially so the user can absorb each one, then compare them in prose. Contrast by **depth** (leverage at the interface), **locality** (where change concentrates), and **seam placement**.

After comparing, give your own recommendation: which design you think is strongest and why. If elements from different designs would combine well, propose a hybrid. Be opinionated — the user wants a strong read, not a menu.
