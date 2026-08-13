---
name: brainstorming
description: Clarify consequential product or architecture decisions before implementation. Use when requirements have materially different interpretations, durable trade-offs, or high-cost decisions; do not use for clear local changes.
---

# Brainstorming Consequential Changes

Use this skill to resolve decisions that would materially change behavior, architecture, compatibility, cost, or delivery risk. It is not a required prelude to every code change.

## 1. Ground the discussion

Read `.cursor/rules/`, relevant code, tests, recent history, and architecture records. Separate:

- facts established by the repository;
- assumptions that are safe and reversible;
- unresolved decisions that could change the solution.

Do not use AskQuestion for facts available in the repository.

## 2. Clarify only material unknowns

Group related blocking questions and present them with AskQuestion. Prefer concrete choices when the real alternatives are known; allow free text when they are not.

Proceed without a question when a reasonable assumption is low-risk and easy to reverse. State that assumption before implementation.

## 3. Compare real alternatives

When multiple credible approaches exist, present the smallest useful set—usually two—with AskQuestion. Explain the trade-off that matters for go-clean-arch and recommend a default.

Do not manufacture alternatives to satisfy a quota. If one approach follows clearly from project constraints, explain it directly.

## 4. Confirm the decision

Summarize the agreed behavior, affected DDD/Clean Architecture layers, important failure handling, validation strategy, and deferred scope. When the decision is consequential, use AskQuestion once for approval.

Save a design record only when it must survive the session, requires team approval, or captures migration/compatibility decisions. Use the existing project location; otherwise use `docs/plans/YYYY-MM-DD-<topic>-design.md`.

## Handoff

- A clear local change may proceed directly.
- A multi-step or cross-domain change may use `writing-plans`.
- A destructive or hard-to-reverse change must include verification and rollback before execution.
