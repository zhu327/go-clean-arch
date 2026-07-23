---
name: go
description: Orchestrates an end-to-end workflow from approved requirements through planning, implementation, validation, review, and simplification. Use when the user requests a full design-to-delivery workflow.
disable-model-invocation: true
---

# Go

Run an end-to-end delivery workflow without assuming a programming language, framework, architecture, or build tool.

## Flow

1. **Design:** use `brainstorming` to understand intent, constraints, alternatives, and success criteria. Wait for user approval.
2. **Choose scope:** use a direct TDD path for a small, localized change; use planning and subagents for a multi-area, cross-cutting, or architectural change.
3. **Plan:** for multi-step work, use `writing-plans` and complete its coverage check.
4. **Implement:** use `test-driven-development`; for a written multi-task plan, use `subagent-driven-development`.
5. **Validate:** run the repository's build, focused/full test, static-analysis, and integration commands discovered from its documentation and CI.
6. **Review:** dispatch `code-reviewer` after functional validation and the acceptance audit pass. Fix material findings and revalidate.
7. **Simplify:** dispatch `code-simplifier` after review. Preserve observable behavior and revalidate any resulting change.

## Constraints

- Lock approved requirements after design. Stop and ask before changing scope.
- Do not prescribe commands, test frameworks, package layouts, or E2E tooling. Discover them from the current repository.
- Do not skip validation, review, or simplification merely because a change is small; if a stage cannot be performed, report the blocker.
- Use the project's own architecture and vocabulary after inspecting its instructions and existing code.

## Completion

Report the approved goal, changed files, validation actually run and its results, review outcome, simplification outcome, and any remaining risks or blockers.
