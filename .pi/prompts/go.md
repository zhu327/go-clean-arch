---
description: End-to-end development workflow — design, plan, implement, validate, review, simplify
argument-hint: "<task description>"
---

# go

Run a complete delivery workflow without assuming a programming language, framework, architecture, or build tool.

## Gates

Track these with the task-list tool:

- Design approved
- Design documented (when the project convention calls for it)
- Plan documented (for multi-step work)
- Implementation complete
- Project-derived validation passed
- Global review complete
- Simplification complete
- Final report delivered

## Workflow

1. Use `brainstorming` to clarify requirements, constraints, alternatives, and success criteria. Obtain approval before implementation.
2. Inspect the repository's instructions, architecture, CI, and commands. Determine whether the change is a small direct TDD task or requires a written multi-task plan.
3. For multi-step work, use `writing-plans`, then execute with `subagent-driven-development`. For smaller work, implement through `test-driven-development`.
4. Run the build, focused/full tests, static analysis, and integration/E2E checks actually supported by the project.
5. After validation and acceptance audit, dispatch `code-reviewer`. Resolve material findings through the appropriate implementer, then revalidate.
6. Dispatch `code-simplifier` against the changed files. Preserve behavior and rerun relevant validation if it modifies code.

## Constraints

- Requirements are locked after approval; stop and ask before a scope change.
- Never invent language-specific commands, paths, test frameworks, or infrastructure. Discover them from the repository.
- Do not skip a workflow gate silently. Report unavailable tools, missing tests, or other blockers explicitly.

Task from user:

$@
