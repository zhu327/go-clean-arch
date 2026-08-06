---
name: brainstorming
description: Clarify consequential product or architecture decisions before implementation. Use when requirements have materially different interpretations, durable trade-offs, or high-cost decisions; do not use for clear local changes.
---

# Brainstorming Consequential Changes

Use this skill to resolve decisions that would materially change behavior, architecture, compatibility, cost, or delivery risk. It is not a required prelude to every code change.

## 1. Ground the discussion

Read the relevant project instructions, code, tests, recent history, and architecture records. Separate:

- facts established by the repository;
- assumptions that are safe and reversible;
- unresolved decisions that could change the solution.

Do not ask the user for facts available in the repository.

## 2. Clarify only material unknowns

Group related blocking questions in one `question` call. Prefer concrete choices when the real alternatives are known; allow free text when they are not.

Proceed without a question when a reasonable assumption is low-risk and easy to reverse. State that assumption before implementation.

## 3. Compare real alternatives

When multiple credible approaches exist, present the smallest useful set—usually two. For each, explain the trade-off that matters in this project. Recommend a default.

Do not manufacture alternatives to satisfy a quota. If one approach clearly follows from project constraints, explain it directly.

## 4. Confirm the decision

Summarize the agreed behavior, boundaries, important failure handling, validation strategy, and explicitly deferred scope. Ask for approval once when the decision is consequential.

Save a design record only when it must survive this session, requires team approval, or captures migration/compatibility decisions. Use the project's existing location; otherwise use `docs/plans/YYYY-MM-DD-<topic>-design.md`.

## Handoff

- A clear local change may proceed directly.
- A multi-step or cross-module change may use `/skill:writing-plans`.
- A destructive or hard-to-reverse change must include verification and rollback before execution.
