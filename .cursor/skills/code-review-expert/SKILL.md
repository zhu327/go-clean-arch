---
name: code-review-expert
description: Expert code review of current git changes with a senior engineer lens. Detects SOLID violations, structural regressions, spaghetti growth, security risks, and proposes actionable improvements. Use when user asks for code review, PR review, or when /go pipeline reaches the review stage.
disable-model-invocation: true
---

# Code Review Expert

## Overview

Perform a structured review of the current git changes with focus on SOLID, architecture, structural simplification, removal candidates, and security risks. Default to review-only output unless the user asks to implement changes.

**Review philosophy**: Be ambitious about code structure. Do not merely identify local cleanup opportunities. Actively search for "code judo" moves — restructurings that preserve behavior while making the implementation dramatically simpler, smaller, more direct, and more elegant. Prefer the solution that makes the code feel inevitable in hindsight.

## Severity Levels

| Level | Name | Description | Action |
|-------|------|-------------|--------|
| **P0** | Critical | Security vulnerability, data loss risk, correctness bug | Must block merge |
| **P1** | High | Logic error, significant SOLID violation, structural regression, performance regression | Should fix before merge |
| **P2** | Medium | Code smell, maintainability concern, spaghetti growth, minor SOLID violation | Fix in this PR or create follow-up |
| **P3** | Low | Style, naming, minor suggestion | Optional improvement |

## Workflow

### 1) Preflight context

- Use `git status -sb`, `git diff --stat`, and `git diff` to scope changes.
- If needed, use `rg` or `grep` to find related modules, usages, and contracts.
- Identify entry points, ownership boundaries, and critical paths (auth, payments, data writes, network).
- Check file sizes before and after the diff — flag any file crossing 1000 lines.

**Edge cases:**
- **No changes**: If `git diff` is empty, inform user and ask if they want to review staged changes or a specific commit range.
- **Large diff (>500 lines)**: Summarize by file first, then review in batches by module/feature area.
- **Mixed concerns**: Group findings by logical feature, not just file order.

### 2) Structural quality + code-judo analysis

Before diving into checklist-driven review, step back and ask the big-picture question:

> Is there a "code judo" move that would make this change dramatically simpler?

- Look for opportunities to **reframe** the change so that whole branches, helpers, modes, conditionals, or layers disappear entirely.
- Do not stop at "this could be a bit cleaner" — push for solutions that **delete complexity** rather than rearrange it.
- Assume there is often a re-organization that uses the existing architecture more effectively.
- Load `./references/code-quality-checklist.md` for structural review prompts.

**For every meaningful change, ask:**
- Can this change be reframed so fewer concepts, branches, or helper layers are needed?
- Does this improve or worsen the local architecture?
- Did the diff add branching complexity where a better abstraction should exist?
- Did a previously cohesive module become more coupled, more stateful, or harder to scan?
- Is this logic living in the right file and layer?
- Did this change enlarge a file past a healthy size boundary?
- Are there repeated conditionals that signal a missing model or missing helper?
- Is the implementation direct and legible, or does it rely on special cases and incidental control flow?
- Is this abstraction actually earning its keep, or is it just a wrapper?
- Is this orchestration more sequential or less atomic than it needs to be?

### 3) SOLID + architecture smells

- Load `./references/solid-checklist.md` for specific prompts.
- Look for:
  - **SRP**: Overloaded modules with unrelated responsibilities.
  - **OCP**: Frequent edits to add behavior instead of extension points.
  - **LSP**: Subclasses that break expectations or require type checks.
  - **ISP**: Wide interfaces with unused methods.
  - **DIP**: High-level logic tied to low-level implementations.
- When you propose a refactor, explain *why* it improves cohesion/coupling and outline a minimal, safe split.
- If refactor is non-trivial, propose an incremental plan instead of a large rewrite.

### 4) Removal candidates + iteration plan

- Load `./references/removal-plan.md` for template.
- Identify code that is unused, redundant, or feature-flagged off.
- Distinguish **safe delete now** vs **defer with plan**.
- Provide a follow-up plan with concrete steps and checkpoints (tests/metrics).

### 5) Security and reliability scan

- Load `./references/security-checklist.md` for coverage.
- Check for:
  - XSS, injection (SQL/NoSQL/command), SSRF, path traversal
  - AuthZ/AuthN gaps, missing tenancy checks
  - Secret leakage or API keys in logs/env/files
  - Rate limits, unbounded loops, CPU/memory hotspots
  - Unsafe deserialization, weak crypto, insecure defaults
  - **Race conditions**: concurrent access, check-then-act, TOCTOU, missing locks
- Call out both **exploitability** and **impact**.

### 6) Code quality scan

- Load `./references/code-quality-checklist.md` for coverage.
- Check for:
  - **Error handling**: swallowed exceptions, overly broad catch, missing error handling, async errors
  - **Performance**: N+1 queries, CPU-intensive ops in hot paths, missing cache, unbounded memory
  - **Boundary conditions**: null/undefined handling, empty collections, numeric boundaries, off-by-one
  - **Structural quality**: file-size explosion, spaghetti growth, abstraction/boundary violations (see checklist)
- Flag issues that may cause silent failures or production incidents.

### 7) Output format

Structure your review as follows:

```markdown
## Code Review Summary

**Files reviewed**: X files, Y lines changed
**Overall assessment**: [APPROVE / REQUEST_CHANGES / COMMENT]

---

## Findings

### P0 - Critical
(none or list)

### P1 - High
- **[file:line]** Brief title
  - Description of issue
  - Suggested fix

### P2 - Medium
...

### P3 - Low
...

---

## Structural Simplification Opportunities
(code-judo moves, missed decompositions, complexity that can be deleted)

## Removal/Iteration Plan
(if applicable)

## Additional Suggestions
(optional improvements, not blocking)
```

**Inline comments**: Use this format for file-specific findings:
```
::code-comment{file="path/to/file.go" line="42" severity="P1"}
Description of the issue and suggested fix.
::
```

**Clean review**: If no issues found, explicitly state:
- What was checked
- Any areas not covered (e.g., "Did not verify database migrations")
- Residual risks or recommended follow-up tests

**Review tone**: Be direct, serious, and demanding about quality. Do not soften major maintainability issues into mild suggestions. If the code is making the codebase messier, say so clearly. If the implementation missed an opportunity for a dramatic simplification, say that clearly too.

Good phrasing examples:
- `this pushes the file past 1k lines. can we decompose this first?`
- `this adds another special-case branch into an already busy flow. can we move this behind its own abstraction?`
- `this works, but it makes the surrounding code more spaghetti. let's keep the behavior and restructure the implementation.`
- `this feels like feature logic leaking into a shared path. can we isolate it?`
- `this abstraction seems unnecessary. can we just keep the direct flow?`
- `i think there's a code-judo move here that makes this much simpler. can we reframe this so these branches disappear?`
- `this refactor moves complexity around, but doesn't really delete it. is there a way to make the model itself simpler?`

### 8) Approval bar

Do not approve merely because behavior seems correct. The bar for approval is:

- No clear structural regression
- No obvious missed opportunity to make the implementation dramatically simpler when such a path is visible
- No unjustified file-size explosion (especially crossing 1000 lines)
- No obvious spaghetti-growth from special-case branching
- No obviously hacky or magical abstraction that makes the code harder to reason about
- No unnecessary wrapper/cast/optionality churn obscuring the real design
- No clear architecture-boundary leak or avoidable canonical-helper duplication
- No missed opportunity for an obvious decomposition that would materially improve maintainability

**Presumptive blockers** (treat as blocking unless the author can justify):
- The PR preserves a lot of incidental complexity when there is a plausible code-judo move that would delete it
- The PR pushes a file from below 1000 lines to above 1000 lines
- The PR adds ad-hoc branching that makes an existing flow more tangled
- The PR solves a local problem by scattering feature checks across shared code
- The PR adds an unnecessary abstraction, wrapper, or cast-heavy contract that makes the design more indirect
- The PR duplicates an existing helper or puts logic in the wrong layer when there is a clear canonical home

### 9) Findings priority order

Prioritize findings in this order:

1. Structural code-quality regressions
2. Missed opportunities for dramatic simplification / code-judo restructuring
3. Spaghetti / branching complexity increases
4. Security vulnerabilities
5. Boundary / abstraction / type-contract problems
6. File-size and decomposition concerns
7. SOLID violations
8. Performance issues
9. Legibility and maintainability concerns

Do not flood the review with low-value nits if there are larger structural issues. Prefer a smaller number of high-conviction comments over a long list of cosmetic notes.

### 10) Next steps confirmation

**When invoked directly by user** (not as part of a subagent pipeline), use `AskQuestion`:

```
AskQuestion({
  title: "Code Review Next Steps",
  questions: [{
    id: "next_action",
    prompt: "I found X issues (P0: _, P1: _, P2: _, P3: _). How would you like to proceed?",
    options: [
      { id: "fix_all", label: "Fix all - Implement all suggested fixes" },
      { id: "fix_critical", label: "Fix P0/P1 only - Address critical and high priority issues" },
      { id: "fix_specific", label: "Fix specific items - I'll tell you which" },
      { id: "no_changes", label: "No changes - Review complete" }
    ]
  }]
})
```

**When loaded by `code-reviewer` agent (subagent context):** Skip AskQuestion — return findings to caller.

**Important**: Do NOT implement any changes until explicitly confirmed. This is a review-first workflow.

## Resources

### references/

| File | Purpose |
|------|---------|
| `solid-checklist.md` | SOLID smell prompts, code-judo remedies, and refactor heuristics |
| `security-checklist.md` | Web/app security and runtime risk checklist |
| `code-quality-checklist.md` | Error handling, performance, boundary conditions, structural quality |
| `removal-plan.md` | Template for deletion candidates and follow-up plan |
