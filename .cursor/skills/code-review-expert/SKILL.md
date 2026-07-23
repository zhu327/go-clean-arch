---
name: code-review-expert
description: Expert code review of current git changes with a senior engineer lens. Detects structural regressions, spaghetti growth, security risks, and code-judo simplification opportunities. Use when user asks for code review, PR review, or when a development pipeline reaches the review stage.
disable-model-invocation: true
---

# Code Review Expert

Review-only by default. Do not implement until the user confirms.

**Philosophy**: Be ambitious about structure. Hunt for **code judo** — restructurings that preserve behavior while *deleting* complexity (not rearranging it). Prefer the solution that feels inevitable in hindsight.

**Iron Law**: Every P0/P1/P2 finding is **Symptom → Consequence → Remedy**. No consequence = noise. No remedy = complaint.

**Severity** (use engineering judgment; do not restate textbook definitions):
- **P0** — security, data loss, correctness → block merge
- **P1** — logic error, significant structural/SOLID regression → fix before merge
- **P2** — maintainability / spaghetti growth → fix or follow-up
- **P3** — style / naming → optional

Assume standard SOLID, OWASP, concurrency, error-handling, and performance knowledge. Load the refs below only for structural bars, security focus, and false-positive filters.

## Workflow

### 1) Preflight

- Scope with `git status -sb`, `git diff --stat`, `git diff` (or staged / commit range if asked).
- Map entry points, ownership boundaries, critical paths (auth, writes, network).
- Flag any file crossing **1000 lines** after the diff.
- **Empty diff** → ask staged vs commit range. **>500 lines** → summarize by file, review by module. **Mixed concerns** → group findings by feature.

### 2) Structural / code-judo pass

Load [`references/team-rules.md`](references/team-rules.md). Before checklist thinking, ask: *is there a reframe that makes whole branches, helpers, modes, or layers disappear?* Prefer deleting complexity over polishing it.

### 3) Architecture smells

Apply SOLID / coupling / DDD judgment from model knowledge. Load [`references/anti-false-positives.md`](references/anti-false-positives.md) before filing architecture findings — skip listed cases. Non-trivial refactors → incremental plan, not big-bang rewrite.

### 4) Removal candidates

Safe-delete-now vs defer-with-plan (see team-rules). Evidence: no refs (incl. dynamic), no external consumers, tests/docs updated — or preconditions + migration + rollback if deferred.

### 5) Security

Load [`references/security-focus.md`](references/security-focus.md). State **exploitability** and **impact**. Skip speculative supply-chain/CVE audits unless the diff changes deps or trust boundaries.

### 6) Reliability (diff-visible only)

Surface issues in the change: swallowed errors, N+1 / hot-path cost, nil/empty/off-by-one, race/TOCTOU, wrong-layer logic. Prefer the project's boundary error type at API edges — do not leak internals to clients.

### 6.5) Tests (only if test files in the diff)

Flag: obscure names, assertion roulette, implementation-coupled mocks, happy-path-only coverage illusion.
Do **not** flag: coherent multi-assert, mocks for nondeterminism, shared setup used by nearly every test, concise-but-obvious names.

### 7) Output format

```markdown
## Code Review Summary

**Files reviewed**: X files, Y lines changed
**Overall assessment**: [APPROVE / REQUEST_CHANGES / COMMENT]

---

## Findings

### P0 - Critical
### P1 - High
- **[file:line]** Brief title
  - Symptom: [what was observed]
  - Consequence: [what breaks or degrades]
  - Remedy: [concrete action — prefer code-judo]

### P2 - Medium
### P3 - Low
(brief; Consequence optional)

---

## Structural Simplification Opportunities
## Removal/Iteration Plan
## Additional Suggestions
```

Inline comments:

```
::code-comment{file="path/to/file" line="42" severity="P1"}
Symptom: ...
Consequence: ...
Remedy: ...
::
```

**Clean review**: state what was checked, gaps (e.g. migrations not verified), residual risks.

**Tone**: direct and demanding. Do not soften structural regressions into nits. Call out missed code-judo when a simpler reframe is visible.

### 8) Approval bar

Do not approve merely because behavior seems correct. **Presumptive blockers** unless justified:

- Clear structural regression or spaghetti growth (ad-hoc branches in shared paths)
- Visible code-judo path ignored while preserving incidental complexity
- File pushed across 1000 lines without strong reason
- Unnecessary wrapper/cast/optionality churn; wrong-layer logic; duplicated canonical helper
- Architecture-boundary leak

### 9) Findings priority

1. Structural regressions / missed code-judo
2. Spaghetti / branching growth
3. Security
4. Boundary / abstraction / type-contract
5. File-size / decomposition
6. SOLID / performance / legibility

Fewer high-conviction comments > long nit lists.

### 10) Next steps

**Direct user invoke** — ask, then wait:

```
I found X issues (P0: _, P1: _, P2: _, P3: _). How would you like to proceed?

A) Fix all
B) Fix P0/P1 only
C) Fix specific items
D) No changes
```

**Subagent / pipeline invoke**: skip the prompt; return findings to caller. Never implement until explicitly confirmed.

## Resources

| File | Purpose |
|------|---------|
| [`team-rules.md`](references/team-rules.md) | Structural non-negotiables, remedies, removal triage |
| [`security-focus.md`](references/security-focus.md) | Diff-scoped security escalation focus |
| [`anti-false-positives.md`](references/anti-false-positives.md) | Common false positives to skip |
