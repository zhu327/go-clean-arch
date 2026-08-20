# Go Clean Architecture Template — Agent Index

Entry point for project rules, read by both pi and Cursor. Mandatory invariants live here; task workflows live in `.agents/skills/` (discovered by both harnesses); agent roles live in `.pi/agents/` (pi) and `.cursor/agents/` (Cursor). Concrete patterns follow existing code.

## Stack & Commands

Go API template: DDD + Clean Architecture, Wire DI, GORM, gomock, HTTP + Swagger.

```bash
make mock   # regen mocks — after changing usecase ports (interfaces.go)
make di     # regen Wire DI — after changing providers
make doc    # regen Swagger docs — after changing handler annotations
make lint   # golangci-lint; lint-dupl and nilaway are separate targets
make test   # tests + coverage profile (make cov to open)
make build  # build cmd/api/main.go
make all    # lint + test + build
make fmt    # gofumpt + golines
make e2e    # scripts/e2e.sh; starts its own docker compose stack
```

If AGENTS.md conflicts with the code or repository tooling, follow the code for behavior and report the mismatch instead of silently editing code or docs.

## Architecture

Dependencies point inward only.

```text
internal/
├── {domain}/domain/     # entities, domain errors, validation — pure Go, no IO/framework
├── {domain}/usecase/    # ports (interfaces.go), DTOs, usecase errors, mock/ (generated)
├── {domain}/adapter/    # delivery/http (handlers) + repository (GORM, domain↔model mapping)
├── shared/              # cross-cutting adapters (delivery server composition)
└── di/                  # Wire providers
```

- domain imports nothing from other layers; usecase depends only on domain and its own ports.
- Handlers stay thin: bind/validate → usecase → respond. Swagger annotations live on handlers.
- Never edit generated files (`wire_gen.go`, `mock/`, Swagger docs) — change the source and regenerate.

## Boundary invariants

- Handlers delegate errors via `c.Error()` to the shared ErrorHandler middleware — never hand-build error JSON responses.
- Repositories use `db.WithContext(ctx)` for every query, map `gorm.ErrRecordNotFound` to domain/usecase errors, and wrap multi-table writes in transactions.

## Verification ladder

Pick the narrowest strategy that gives credible evidence, then widen with the change's scope and risk; a failing check must be fixed and rerun before moving on:

| Change | Verification |
|--------|--------------|
| Reproducible bug | regression test that fails first, then fix |
| Business rule / edge-heavy logic | red → green, one behavior at a time |
| Refactor / mechanical move | existing suite as the safety net |
| Ports / DI / Swagger inputs | `make mock` / `make di` / `make doc` |
| Public HTTP behavior | `make e2e` when the harness covers it |

Tests assert observable behavior (returns, persisted state, HTTP responses) — no test-only production hooks, no shared mutable state, no order/timing coupling. An applicable failing check blocks completion; a missing tool is a reported gap, not a blocker.

## Workflow

| Mode | Signals | Execution |
|------|---------|-----------|
| Direct (default) | clear, local, reversible | implement → narrowest checks → self-review diff → report |
| Planned | dependent steps, cross-layer contracts | lightweight plan (Goal/Assumptions/Steps/Validation), then implement |
| Parallel | independent slices, disjoint write sets | skill `plan-execute` |
| High-risk | auth, security, privacy, credentials, money, migrations, destructive writes, concurrency, public contracts | design agreement + rollback plan + one independent `code-reviewer` review |

File count is a signal, not the rule — a one-file auth change is high-risk. Ask only questions whose answers materially change the result; state minor assumptions and proceed. Never invent project tooling.

## Output hygiene

- Execute directly; do not restate the plan before acting.
- No reasoning leakage in code, docs, or commits (decision numbering, review dialogue, change narratives).
- Read a file once per session; prefer grep over re-reading. Never re-dispatch a subagent with an unchanged task.
- Lean replies: no preamble, no restating the diff, no narrating the next step before doing it.

## Skills & Agents

Skills (`.agents/skills/`, shared by pi and Cursor): `grill-me` (relentless design interview before consequential work), `plan-execute` (parallel/high-risk planning + wave execution), `code-review-expert` (changeset review methodology), `improve-codebase-architecture` (deepening assessment, explicit invocation).

Agents (`.pi/agents/`): `code-reviewer` (read-only independent review), `code-simplifier` (targeted behavior-preserving simplification). At most one full changeset review; re-review fixes only when they reshape the design.
