# Subagent-Driven Development Examples

These examples use neutral placeholders. Replace every path, command, artifact, and contract with values inspected from the current repository.

## Dependency-Safe Waves

```text
Plan tasks:
- Task 1: Add source import (no dependencies; modifies src/importer/*)
- Task 2: Add report export (no dependencies; modifies src/exporter/*)
- Task 3: Connect import to report (blocked by 1 and 2; modifies src/application/*)

Wave 0: Task 1 and Task 2 may run in parallel only after confirming their file lists do not overlap.
Wave 1: Task 3 starts after both Task 1 and Task 2 pass review and validation.
```

## Validation Failure

```text
Controller runs the validation commands recorded in the plan and project documentation.

Validation: <project build command> → FAIL
Output: <exact failure output>

Dispatch one integration fixer with:
- exact command output;
- changed files from Tasks 1 and 2;
- instruction to make the smallest integration correction.

Re-run: <project build command>, <project focused test command>, and any required static-analysis command.
Do not complete the wave until all pass.
```

## Artifact Passing

```text
Before Task 3, the controller reads actual artifacts from Tasks 1 and 2:

- path/to/import-contract: [exact exported declaration, schema, or configuration]
- path/to/export-contract: [exact exported declaration, schema, or configuration]

The controller includes those exact artifacts in Task 3’s prompt. It does not substitute an informal summary.

If an artifact differs from the approved plan, the controller stops and reports the mismatch.
```

## Parallel Prompt Addendum

```markdown
## Parallel Execution Notice

Other agents are implementing: [Task names].
Your permitted files are:
- [exact Create/Modify paths]

Do not modify files outside this list. If your task needs another file, stop and report the conflict.
```
