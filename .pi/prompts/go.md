---
description: Risk-adaptive end-to-end delivery
argument-hint: "<task description>"
---

# go

Deliver the requested change end to end, following the workflow in AGENTS.md:

1. **Assess**: read the relevant code and rules; classify Direct / Planned / Parallel / High-risk (table in AGENTS.md). Tell the user the mode when it affects interaction or artifacts.
2. **Clarify** only material unknowns (one grouped `question` call); load skill `grill-me` for consequential design decisions; otherwise state assumptions and proceed.
3. **Plan** at the depth the mode needs: Direct → `todo` if useful; Planned → lightweight plan; Parallel / High-risk → skill `plan-execute`.
4. **Implement** using the verification ladder in AGENTS.md; the coordinator may edit code directly for Direct / Planned work.
5. **Format & validate**: run `make fmt`, then validate narrow → broad: focused tests → affected package suite → `make all` / `make e2e` as risk warrants. Fix failures and rerun the failing check.
6. **Review** proportionally: always self-review the diff; dispatch `code-reviewer` at most once when high risk, large or cross-layer impact, multiple implementers, a public contract changes, or the user explicitly requests review.
7. **Report**: what changed, files touched, validation actually run with results, checks not run and why, residual risk.

Core rules: never invent project tooling; keep approved scope stable — stop and ask if discovery changes public behavior, data safety, or the agreed design; do not create plans, E2E infrastructure, abstractions, or cleanup solely to satisfy this workflow.

Task from user:

$@
