---
name: code-review-expert
description: Expert changeset review with a senior engineer lens. Two axes — intent (does it do what was asked?) and code quality (structure, security, reliability, tests). Finds structural regressions and code-judo simplifications. Use for code review requests and pipeline review stages.
---

# Code Review Expert

Review-only until the user confirms. Hunt for **code judo** — restructurings that preserve behavior while *deleting* complexity. Prefer the solution that feels inevitable in hindsight.

**Iron Law**: every P0/P1/P2 finding is **Symptom → Consequence → Remedy**. No consequence = noise. No remedy = complaint.

- **P0** security / data loss / correctness — blocks merge
- **P1** logic error, structural or SOLID regression — fix before merge
- **P2** maintainability / spaghetti growth — fix or follow-up
- **P3** style / naming — optional

Assume standard SOLID, OWASP, concurrency, and error-handling knowledge; do not restate textbook definitions. **Don't duplicate green gates**: when `make lint` / `nilaway` / `make test` are green, don't re-report diagnostics those tools would catch — but still judge whether the gates actually cover the changed behavior. Hunt the semantic gaps tooling cannot see. Project standards live in AGENTS.md (layer rules, regen targets); cite them, do not restate them.

## Workflow

1. **Preflight**: pin the fixed point and confirm it resolves (`git rev-parse <base>`), confirm the diff is non-empty, then `git diff` (staged / commit range if asked). Map entry points, ownership boundaries, critical paths (auth, writes, network). Large diff → summarize by file, review by module. Mixed concerns → group findings by feature.
2. **Intent vs spec**: find what the change claims to implement — the task text, plan doc, or commit messages. Report (a) missing or partial requirements, (b) behavior beyond the ask (scope creep), (c) requirements that look implemented but are wrong. Quote the requirement for each finding. No spec available → state that and continue on code quality alone; keep the two axes separate in the report.
3. **Code-judo pass** (load [`references/team-rules.md`](references/team-rules.md)): before checklist thinking, ask *is there a reframe that deletes whole branches, helpers, modes, or layers?* Prefer deleting complexity over polishing it.
4. **Architecture** (load [`references/anti-false-positives.md`](references/anti-false-positives.md) before filing): boundary leaks, wrong-layer logic, duplicated canonical helpers. Non-trivial refactors → incremental plan, not big-bang rewrite.
5. **Removal candidates** (team-rules): safe-delete-now vs defer-with-plan. Evidence: no references (incl. dynamic), no external consumers, tests/docs updated — or preconditions + migration + rollback if deferred.
6. **Security** (load [`references/security-focus.md`](references/security-focus.md)): state **exploitability** and **impact**. Skip speculative supply-chain/CVE audits unless the diff changes deps or trust boundaries.
7. **Reliability (diff-visible only)**: swallowed errors, N+1 / hot-path cost, nil / off-by-one, race / TOCTOU, revertibility / blast radius (backward compatibility, irreversible migrations/config changes). Use the project's boundary error type at API edges — do not leak internals.
8. **Tests**: changed behavior needs new or updated verification — a missing test diff is a signal, not an automatic finding. Flag assertion roulette and implementation-coupled mocks; do not flag coherent multi-assert tests or shared setup used by nearly every test.

## Output

```markdown
## Code Review Summary
**Files reviewed**: X files, Y lines changed | **Assessment**: APPROVE / REQUEST_CHANGES / COMMENT

### Intent (vs spec)   — omit when no spec exists
### Findings
#### P1 — [file:line] title
- Symptom: what was observed
- Consequence: what breaks or degrades
- Remedy: concrete action — prefer code-judo

## Structural Simplification Opportunities
## Removal / Iteration Plan
## Residual Risk (checks not performed)
```

Inline findings: `::code-comment{file="path" line="42" severity="P1"} … ::`

Fewer high-conviction comments > long nit lists. Priority: intent gaps → structural regressions / missed code-judo → spaghetti growth → security → boundary / type contracts → decomposition → SOLID / performance.

## Approval bar

Do not approve merely because behavior seems correct. Presumptive blockers unless justified:

- requirement gaps or unrequested behavior found on the intent axis;
- structural regression or spaghetti growth (ad-hoc branches in shared paths);
- visible code-judo path ignored while preserving incidental complexity;
- file growth mixing responsibilities, hard to test and review;
- wrapper/cast/optionality churn; wrong-layer logic; duplicated helper; architecture-boundary leak.

## Next steps

Direct user invoke → use your harness's question tool (pi `question`, Cursor `AskUserQuestion`) with fix options (all / P0+P1 / specific / none) and wait for confirmation; without a question tool, ask in markdown. Subagent / pipeline invoke → return findings to the caller. Never implement until explicitly confirmed.

## Resources

| File | Purpose |
|------|---------|
| [`team-rules.md`](references/team-rules.md) | Structural non-negotiables, remedies, removal triage |
| [`security-focus.md`](references/security-focus.md) | Diff-scoped security escalation |
| [`anti-false-positives.md`](references/anti-false-positives.md) | Common false positives to skip |
