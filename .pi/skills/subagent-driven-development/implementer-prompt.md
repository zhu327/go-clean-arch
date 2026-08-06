# Implementer Subagent Prompt Template

Use this template when dispatching an implementer.

```markdown
You are implementing Task N: [task name].

## Task Description

[FULL text of the approved task. Do not ask the implementer to rediscover the plan.]

## Context

[Where this fits, completed dependencies, and relevant project conventions]

## Before You Begin

If requirements, acceptance criteria, dependencies, assumptions, or the implementation approach are unclear, ask before changing files. Do not guess.

## Your Job

Follow the verification strategy specified by the task. Load `test-driven-development` when choosing among regression-first, behavior test-first, characterization, existing-suite, or validator-based approaches. For a reproducible bug or suitable new behavior, prefer a focused failing test before the production change. Report any justified alternative explicitly.

Work only in: [allowed file list].

## Before Reporting

Check:

- every acceptance criterion is implemented;
- no unrequested behavior, dependencies, public contracts, or files changed;
- names and implementation follow existing project conventions;
- tests cover behavior, errors, and relevant edges;
- all validation actually run is reported truthfully.

## Report Format

- implemented behavior;
- tests and validation commands run, with results;
- files changed;
- self-review findings;
- unresolved concerns or blockers.
```

## Parallel Variant

Add this section to every parallel implementer prompt:

```markdown
## Parallel Execution Notice

Other agents are working on: [task names].
Do not create or modify files outside your allowed list. If additional files are necessary, stop and report the conflict to the controller.
```

## Prior-Wave Dependencies

The controller must read actual completed artifacts and provide the exact relevant declarations—not only a prose summary—in downstream prompts:

```text
## Prior-Wave Dependencies

From Task X: `path/to/artifact`
[exact relevant interface, schema, configuration, event, or exported declaration in the project’s syntax]

Use these declarations exactly. If they conflict with the plan, stop and report the mismatch.
```
