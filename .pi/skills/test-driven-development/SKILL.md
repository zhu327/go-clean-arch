---
name: test-driven-development
description: Select and apply behavior-focused verification for features, bug fixes, and refactors. Prefer test-first where it produces a meaningful failure signal; adapt for legacy characterization, generated code, configuration, and mechanical changes.
---

# Behavior-Focused Development

The goal is credible, repeatable evidence that the requested behavior works—not compliance with a ritual.

## Choose the verification strategy

- **Reproducible bug:** add a regression test that fails for the reported defect before fixing it.
- **Business rule, parser/validator, state transition, or edge-heavy logic:** prefer red → green → refactor, one behavior at a time.
- **Ordinary new behavior:** add or update behavior-focused tests; test-first is the default when practical.
- **Legacy behavior:** characterization tests may first capture current behavior before refactoring.
- **Mechanical refactor:** use the existing suite as the safety net; add tests only for an uncovered behavior risk.
- **Generated code or declarative configuration:** use the project's generator, schema check, build, or validator.
- **Exploratory prototype:** exploration may precede tests; agree on verification before treating it as production work.
- **Not meaningfully automatable:** obtain agreement on an explicit alternative check.

## Test-first loop

When test-first applies:

1. Add the smallest test for an observable behavior.
2. Run the focused command discovered from the repository.
3. Confirm it fails because the behavior is missing, not because setup is broken.
4. Make the smallest production change that passes.
5. Run the focused test, then the relevant broader suite.
6. Refactor only while checks remain green.

## Test quality

Prefer tests that:

- assert returned values, persisted state, emitted events, rendered output, or documented external interactions;
- remain valid through internal refactors;
- use real lightweight collaborators or fakes before interaction-heavy mocks;
- cover important errors and boundaries implied by the requirement;
- create isolated data and avoid order, timing, or shared-state dependencies.

Do not add a public production hook solely for tests. Build fixtures from the real contract rather than guessed partial objects.

Load `testing-anti-patterns.md` when adding mocks, fakes, fixtures, or test-only seams. Load `tdd-rationale.md` only when deciding whether test-first adds value for a disputed case.

## Completion evidence

Report:

- the behavior covered and at what test level;
- commands actually run and their results;
- any alternative verification used;
- untested or unavailable checks and their residual risk.

An applicable failing check blocks completion. A test tool the project does not provide is a gap to report, not a reason to invent infrastructure automatically.
