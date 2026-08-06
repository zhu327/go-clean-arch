# Go Clean Architecture Template — AI Agent Configuration Index

> Project-specific architecture, coding, testing, and Swagger conventions live under `.cursor/rules/`. Scale the development process by risk instead of running the full pipeline for every change.

## Default Development Flow

```text
Inspect repository rules and code → Assess risk → Direct / Planned / Parallel / High-risk
→ Behavior verification → Applicable project checks → Proportionate review → Delivery report
```

| Mode | Typical use | Default execution |
|------|-------------|-------------------|
| Direct | Clear, local, reversible change | Current agent implements, tests, and self-reviews |
| Planned | Dependent steps, cross-layer contracts, or meaningful design choices | Lightweight plan; save a full plan only when useful |
| Parallel | Independent slices with stable contracts and non-overlapping files | Dependency-wave subagent execution |
| High-risk | Authentication/authorization, credentials, migrations, destructive writes, concurrency, privacy, or hard-to-reverse contracts | Design agreement, verification/rollback plan, independent review |

File count is only a signal. A one-file authentication change may be high-risk.

## Skills

| Skill | Purpose | Use when |
|-------|---------|----------|
| `go` | Risk-adaptive end-to-end delivery | User explicitly requests full delivery |
| `brainstorming` | Resolve consequential requirement or architecture choices | Real ambiguity can materially change the result; not for clear local work |
| `writing-plans` | Create a durable multi-step plan | Cross-layer/domain, high-risk, parallel, cross-session, or approval work |
| `test-driven-development` | Select a behavior verification strategy | Features, bugs, and refactors; allows characterization and validator-based alternatives |
| `subagent-driven-development` | Execute an approved plan in parallel | Slices are independent, contracts stable, write sets disjoint, and parallelism pays off |
| `e2e-testing` | Exercise the public interface | Existing harness adds value or requirements/risk justify it |
| `code-review-expert` | Independent changeset review method | High-risk, cross-domain, large diff, multiple implementers, or explicit request |
| `improve-codebase-architecture` | Quick or formal architecture assessment | Evidence-backed refactor candidates or interface design discussion |

## Agents

| Agent | Purpose | Trigger |
|-------|---------|---------|
| `code-reviewer` | Read-only independent changeset review | High-risk work, public contract/architecture changes, multiple implementers, difficult diff, or explicit request |
| `code-simplifier` | Behavior-preserving refactor for a concrete target | Review identified specific complexity or the user explicitly requested it; not a routine final stage |

## Validation Strategy

Run applicable checks from narrow to broad:

1. focused `go test` for changed behavior;
2. affected package/domain tests;
3. `make mock`, `make di`, or `make doc` when their inputs changed;
4. relevant `make build`, `make lint`, `make test`, or `make all`;
5. `make e2e` when the existing harness and risk warrant it.

An applicable failing check blocks completion. A missing or inapplicable tool is reported as a gap and residual risk; do not create infrastructure solely to satisfy a generic workflow.

## Project Rules

- `.cursor/rules/00-project-overview.mdc`: stack, architecture, layout, and commands
- `.cursor/rules/10-domain-layer.mdc` through `15-task-layer.mdc`: layer rules
- `.cursor/rules/20-wire-di.mdc`: Wire dependency injection
- `.cursor/rules/30-testing.mdc`: Go, testify, gomock, Handler tests, and coverage conventions
- `.cursor/rules/40-api-swagger.mdc`: HTTP and Swagger contracts

Project rules override generic skills. If a rule disagrees with actual code or tooling, verify repository reality and report the mismatch instead of applying it blindly.

## Common Commands

```bash
make build   # Build cmd/api/main.go
make serve   # Build and run with development config
make di      # Generate Wire dependency injection
make doc     # Generate Swagger docs
make lint    # Run golangci-lint
make test    # Run project tests with coverage profile
make cov     # Open coverage report
make mock    # Generate mocks
make e2e     # Run the existing E2E script
make fmt     # Run golines and gofumpt
make all     # Lint, test, and build
```

## Architecture Overview

The repository uses DDD + Clean Architecture with inward dependencies:

```text
internal/
├── {domain}/
│   ├── domain/
│   ├── usecase/
│   └── adapter/
├── shared/
└── di/
```

Use established terms such as domain, usecase, handler, router, repository, gateway, middleware, and Wire rather than replacing them to match a generic glossary.
