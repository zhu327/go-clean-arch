---
name: test-driven-development
description: Select and apply behavior-focused verification for go-clean-arch features, bug fixes, and refactors. Prefer test-first where it produces a meaningful failure signal; adapt for legacy characterization, generated code, configuration, and mechanical changes.
---

# Behavior-Focused Development

The goal is credible, repeatable evidence that requested behavior works—not compliance with a ritual. Follow project-specific testing conventions in `.cursor/rules/30-testing.mdc` where they apply.

## Choose the verification strategy

- **Reproducible bug:** add a regression test that fails for the defect before fixing it.
- **Business rule, parser/validator, state transition, or edge-heavy logic:** prefer red → green → refactor.
- **Ordinary new behavior:** add or update behavior-focused tests; test-first is the default when practical.
- **Legacy behavior:** characterization tests may capture current behavior before refactoring.
- **Mechanical refactor:** use the existing suite; add tests only for an uncovered behavior risk.
- **Generated code or declarative configuration:** use `make mock`, `make di`, `make doc`, build, or the relevant validator when inputs change.
- **Exploratory prototype:** exploration may precede tests; agree on verification before production use.
- **Not meaningfully automatable:** agree on an explicit alternative check.

## Test-first loop

When test-first applies:

1. Add the smallest observable behavior test.
2. Run a focused `go test` command derived from the package.
3. Confirm it fails for the missing behavior, not broken setup.
4. Make the smallest production change that passes.
5. Run the focused test and affected package/domain suite.
6. Refactor only while checks remain green.

## go-clean-arch defaults

- Prefer table-driven tests when several cases share setup and behavior.
- Use `testify` assertions and `gomock` following existing package conventions.
- Regenerate mocks only when their source interfaces change.
- Handler tests should exercise Gin/HTTP behavior through `httptest`; broader E2E is risk- and harness-driven.
- Treat coverage as a diagnostic signal. Do not add low-value assertions solely to reach a numeric target.

Prefer observable outputs, persisted state, events, response contracts, and documented external interactions over mock-call-only assertions. Do not add production APIs solely for tests or build fixtures from guessed partial contracts.

Load `testing-anti-patterns.md` when adding mocks, fakes, fixtures, or test-only seams. Load `tdd-rationale.md` only when test order is disputed.

## Completion evidence

Report behavior covered, test level, commands and results, alternative verification, and residual risk. An applicable failing check blocks completion; missing infrastructure is a gap to report, not a reason to invent it automatically.
