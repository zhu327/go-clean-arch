---
name: go
description: Risk-adaptive end-to-end development workflow. Use when the user explicitly requests full delivery from repository inspection through implementation, validation, and proportionate review.
disable-model-invocation: true
---

# Go

Deliver the requested change end to end using the lightest workflow that gives credible confidence. Increase process only when ambiguity, risk, or useful parallelism justifies it.

## Core rules

- Inspect repository rules, code, tests, and CI before choosing commands or architecture. Use go-clean-arch conventions under `.cursor/rules/` and never invent tooling.
- Use AskQuestion only for unresolved decisions that materially change the result. State minor reversible assumptions and proceed.
- The coordinator may implement directly. Use subagents only when independent work or judgment is worth the context-transfer cost.
- Match validation and review depth to risk. An applicable failing check blocks completion; a tool the repository does not provide is a reported gap, not automatically a blocker.
- Do not create design documents, plans, E2E infrastructure, abstractions, or cleanup changes solely to satisfy this workflow.
- Stop before changing approved public behavior, data safety, or a consequential design decision. Use AskQuestion to confirm before proceeding.

## 1. Assess

Classify the work using engineering judgment:

| Mode | Typical signals | Default execution |
|------|-----------------|-------------------|
| **Direct** | Clear, local, reversible, one coherent change | Implement in the current agent |
| **Planned** | Dependent steps, cross-layer contracts, or meaningful design choices | Make a concise plan, then implement |
| **Parallel** | Independent slices with stable contracts and non-overlapping files | Plan dependency waves, then use subagents |
| **High-risk** | Authentication/authorization, migrations, destructive writes, concurrency, credentials/privacy, or hard-to-reverse public contracts | Use AskQuestion to agree design and verification/rollback; require independent review |

File count is only a signal. Briefly state the selected mode when it changes interaction, artifacts, or delivery time.

## 2. Clarify and plan as needed

- Clear task: proceed with explicit assumptions.
- Blocking ambiguity: group related questions in AskQuestion.
- Consequential alternatives or durable architecture decisions: use `brainstorming`.
- Direct work needs no plan document.
- Planned work may use a lightweight in-session sequence.
- Parallel or high-risk work should use `writing-plans` for exact contracts, dependencies, acceptance criteria, files, and project-derived validation.
- Save a design/plan artifact only when it must survive the session, requires approval, or records migration/compatibility decisions.

Do not manufacture tasks or interfaces to fill a template.

## 3. Implement and test

Choose verification based on the change:

- Reproducible bug, business rule, parser/validator, or state transition: prefer a failing regression/behavior test first.
- Ordinary new behavior: add or update behavior-focused tests; test-first is preferred when practical.
- Legacy characterization, mechanical refactor, generated code, or configuration: existing tests, characterization, build/schema validation, or another explicit check may be more appropriate.
- HTTP/API behavior: update Swagger per `.cursor/rules/40-api-swagger.mdc`; use the existing integration/E2E harness when it provides meaningful coverage. Create new infrastructure only when requirements or risk justify it.
- Interface or dependency changes: regenerate mocks/Wire only when affected, using repository commands.

Direct and planned work may be implemented by the coordinator. For parallel work, use `subagent-driven-development` only with stable inputs and non-overlapping write sets; pass exact artifacts from prerequisite tasks.

## 4. Validate from narrow to broad

Run applicable repository checks:

1. focused `go test` for changed behavior;
2. affected package/domain suite;
3. generated artifact checks such as `make mock`, `make di`, or `make doc` when inputs changed;
4. relevant `make build`, `make lint`, or other documented checks;
5. integration/E2E and broader validation when supported and warranted.

Fix relevant failures and rerun the failing check. Report commands actually run, unavailable checks, and residual risk.

## 5. Review proportionally

Always inspect the final diff for scope, correctness, generated artifacts, accidental files, and unnecessary complexity.

Use `code-reviewer` when security, authentication/authorization, credentials/privacy, migration, concurrency, destructive behavior, public contracts, or cross-domain architecture changed; when multiple implementers contributed; when the diff is difficult to reason about; or when the user asks.

Otherwise self-review plus passing validation is sufficient. Run at most one full changeset review; after fixes, target the affected findings unless the design changed materially. Invoke `code-simplifier` only for a concrete simplification target or explicit user request.

## 6. Deliver

Report what changed, important decisions, affected areas, checks actually run and results, checks not run and why, and residual risks or follow-up.
