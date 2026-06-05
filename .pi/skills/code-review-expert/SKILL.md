---
name: code-review-expert
description: "Use for rigorous review of current git changes, PRs, or completed implementation work. Detects SOLID violations, structural regressions, spaghetti growth, security/reliability risks, removal candidates, and code-judo simplifications."
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

- Use `bash("git status -sb")`, `bash("git diff --stat")`, and `bash("git diff")` to scope changes.
- If needed, use `bash("rg ...")` or `bash("grep ...")` to find related modules, usages, and contracts.
- Identify entry points, ownership boundaries, and critical paths (auth, payments, data writes, network).
- Check file sizes before and after the diff — flag any file crossing 1000 lines.

**Edge cases:**
- **No changes**: If `git diff` is empty, inform user and ask if they want to review staged changes or a specific commit range.
- **Large diff (>500 lines)**: Summarize by file first, then review in batches by module/feature area.
- **Mixed concerns**: Group findings by logical feature, not just file order.

### 2) Structural quality + code-judo analysis

Resource paths below are relative to `.pi/skills/code-review-expert/`.

Before diving into checklist-driven review, step back and ask the big-picture question:

> Is there a "code judo" move that would make this change dramatically simpler?

- Look for opportunities to **reframe** the change so that whole branches, helpers, modes, conditionals, or layers disappear entirely.
- Do not stop at "this could be a bit cleaner" — push for solutions that **delete complexity** rather than rearrange it.
- Load `./references/code-quality-checklist.md` for structural review prompts.

### 3) SOLID + architecture smells

- Load `./references/solid-checklist.md` for specific prompts.
- Look for SRP, OCP, LSP, ISP, DIP violations.
- When you propose a refactor, explain *why* it improves cohesion/coupling and outline a minimal, safe split.

### 4) Removal candidates + iteration plan

- Load `./references/removal-plan.md` for template.
- Identify code that is unused, redundant, or feature-flagged off.
- Distinguish **safe delete now** vs **defer with plan**.

### 5) Security and reliability scan

- Load `./references/security-checklist.md` for coverage.
- Check for injection, authZ gaps, secret leakage, race conditions, resource exhaustion.

### 6) Code quality scan

- Load `./references/code-quality-checklist.md` for coverage.
- Check for error handling, performance (N+1 queries), boundary conditions, structural quality.

### 7) Output format

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

## Removal/Iteration Plan

## Additional Suggestions
```

### 8) Approval bar

Do not approve merely because behavior seems correct. The bar for approval is:

- No clear structural regression
- No obvious missed opportunity to make the implementation dramatically simpler
- No unjustified file-size explosion (especially crossing 1000 lines)
- No obvious spaghetti-growth from special-case branching
- No clear architecture-boundary leak

### 9) Next steps confirmation

After presenting findings, use `question` to ask user how to proceed:

```
question({
  questions: [{
    id: "next_action",
    prompt: "I found X issues (P0: _, P1: _, P2: _, P3: _). How would you like to proceed?",
    options: [
      { label: "Fix all - Implement all suggested fixes", value: "fix_all" },
      { label: "Fix P0/P1 only - Address critical and high priority issues", value: "fix_critical" },
      { label: "Fix specific items - I'll tell you which", value: "fix_specific" },
      { label: "No changes - Review complete", value: "no_changes" }
    ]
  }]
})
```

**Important**: Do NOT implement any changes until user explicitly confirms. This is a review-first workflow.

## Resources

### references/

| File | Purpose |
|------|---------|
| `solid-checklist.md` | SOLID smell prompts, code-judo remedies, and refactor heuristics |
| `security-checklist.md` | Web/app security and runtime risk checklist |
| `code-quality-checklist.md` | Error handling, performance, boundary conditions, structural quality |
| `removal-plan.md` | Template for deletion candidates and follow-up plan |
