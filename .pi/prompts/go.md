---
description: Risk-adaptive end-to-end development workflow
argument-hint: "<task description>"
---

# go

Deliver the requested change end to end. Use the lightest workflow that gives credible confidence; increase process only when ambiguity, risk, or parallelism justifies it.

## Core rules

- Inspect the repository before choosing commands, architecture, or test strategy. Never invent project tooling.
- Ask only about unresolved decisions that can materially change the result. State minor assumptions and proceed.
- The coordinator may implement directly. Use subagents when independent work or an independent judgment is worth the context-transfer cost.
- Match validation and review depth to the change's risk. An applicable failing check is a blocker; a tool the project does not provide is a reported gap, not automatically a blocker.
- Do not create design documents, plans, E2E infrastructure, abstractions, or cleanup changes solely to satisfy this workflow.
- Keep approved scope stable. Stop and ask if a newly discovered requirement changes public behavior, data safety, or the agreed design.

## 1. Assess

Read the relevant code, tests, project instructions, and available validation commands. Classify the work using engineering judgment:

| Mode | Typical signals | Default execution |
|------|-----------------|-------------------|
| **Direct** | Clear, local, reversible; one coherent change | Implement in the current agent |
| **Planned** | Several dependent steps, cross-module contracts, or meaningful design choices | Create a concise plan, then implement |
| **Parallel** | Multiple independent slices with stable contracts and non-overlapping files | Plan dependency waves, then use subagents |
| **High-risk** | Auth/tenancy, security, migrations, destructive writes, concurrency, money, privacy, or hard-to-reverse public contracts | Obtain design agreement, plan rollback/verification, and require independent review |

File count is a signal, not the decision rule. A one-file authorization change may be high-risk; a many-file mechanical rename may be direct.

Briefly tell the user which mode you selected and why when the choice affects interaction, artifacts, or delivery time.

## 2. Clarify and design only as needed

- If the task is clear, proceed with explicit assumptions.
- If a blocking ambiguity exists, group related questions in one `question` call.
- Load `/skill:brainstorming` only when there are consequential alternatives, unclear product behavior, or a durable architecture decision.
- Save a design document only when the decision must survive the session, requires approval, or records migration/compatibility trade-offs.

## 3. Plan at the appropriate depth

- **Direct:** use `todo` only when useful; no plan document is required.
- **Planned:** write a lightweight sequence of changes and verification steps. Load `/skill:writing-plans` when exact contracts, dependencies, or handoff artifacts matter.
- **Parallel / High-risk:** use `/skill:writing-plans` for dependency-safe slices, acceptance criteria, actual file paths, and project-derived validation. Include rollback or recovery for destructive work.

Do not manufacture tasks or interfaces to fill a template.

## 4. Implement and test

Choose the verification approach based on the change:

- Reproducible bug, business rule, parser/validator, or state transition: prefer a failing regression/behavior test first.
- Ordinary new behavior: add or update behavior-focused tests; test-first is preferred when practical.
- Legacy characterization, mechanical refactor, generated code, or declarative configuration: existing tests, characterization tests, build/schema validation, or another explicit check may be more appropriate.
- Public-interface changes: use the repository's existing integration/E2E harness when it provides meaningful coverage. Create new E2E infrastructure only when requirements or risk justify it.

For **Direct** and **Planned** work, the coordinator may edit implementation code.

For **Parallel** work, load `/skill:subagent-driven-development`. Parallel tasks must have stable inputs and non-overlapping write sets. Pass downstream agents the real contracts produced by prerequisite tasks.

## 5. Validate from narrow to broad

Run the applicable checks discovered in the repository:

1. focused tests or validators for the changed behavior;
2. affected module/package suite;
3. build, typecheck, or lint relevant to the change;
4. integration/E2E checks when supported and warranted;
5. broader validation when risk or repository convention calls for it.

Fix relevant failures and rerun the failing check. Do not claim a command ran if it did not. Report unavailable or impractical checks and the resulting residual risk.

## 6. Review proportionally

Always inspect the final diff for scope, correctness, accidental files, and unnecessary complexity.

Use an independent `code-reviewer` when one or more apply:

- security, authorization, tenancy, privacy, money, migration, concurrency, or destructive behavior changed;
- public contracts or cross-module architecture changed materially;
- multiple implementers contributed;
- the diff is large or difficult to reason about in one pass;
- the user requested independent review.

Otherwise, self-review plus passing validation is sufficient.

Run at most one full changeset review. After fixes, use targeted re-review of affected findings unless the fixes materially reshape the design. Invoke `code-simplifier` only for a concrete simplification target found in review or requested by the user; it is not a mandatory delivery stage.

## 7. Deliver

Report concisely:

- what changed and any important design decision;
- files or areas changed;
- tests and validation actually run, with results;
- checks not run and why;
- residual risks or required follow-up.

Task from user:

$@
