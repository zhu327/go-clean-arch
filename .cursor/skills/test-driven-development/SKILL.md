---
name: test-driven-development
description: Use before implementation code for a feature, bug fix, refactor, or behavior change. Enforces red-green-refactor without assuming a language or test framework.
---

# Test-Driven Development

## Iron Law

**Do not add production behavior before a test demonstrates the missing behavior.**

Exceptions—generated code, declarative configuration, exploratory prototypes, and changes that cannot be meaningfully tested—require explicit user agreement and a documented alternative verification method.

## Red → Green → Refactor

For one behavior at a time (or a closely related batch of at most four):

1. **Red:** write the smallest test that describes the desired observable behavior.
2. **Verify red:** run the project’s focused test command. Confirm it fails because the behavior is missing, not because of setup, syntax, or an unrelated failure.
3. **Green:** write the smallest production change that makes the test pass.
4. **Verify green:** rerun the focused test and then the relevant broader suite.
5. **Refactor:** improve duplication or clarity only while the tests remain green.

Discover test commands from the repository’s package manifest, build scripts, CI, Makefile, or existing test documentation. Never prescribe a runner, path convention, assertion library, or language syntax.

## Good Tests

Tests should:

- describe one externally meaningful behavior with a clear name;
- assert outcomes, contracts, state changes, or interactions at a genuine seam;
- use real collaborators when practical and fakes/mocks only where isolation is necessary;
- cover success, failure, and relevant boundary cases;
- remain valid through internal refactors.

Avoid tests that only mirror implementation details, rely on order or time, share uncontrolled mutable state, or require production-only test hooks.

## Debugging and Refactoring

- A bug fix starts with a failing regression test.
- When a test is difficult to write, first reconsider the interface and dependency shape; do not immediately add test-only public APIs.
- If a refactor changes behavior, use the same red-green-refactor cycle. Purely mechanical formatting may use existing tests as the safety net.

## Completion Checklist

- [ ] Each new behavior has a test or an explicitly approved alternative verification method.
- [ ] The new tests were observed failing for the intended reason.
- [ ] Production changes were minimal and followed the failing test.
- [ ] Focused and relevant full suites pass.
- [ ] Edge cases and errors implied by requirements are covered.
- [ ] The project’s available validation commands pass.

For rationale and common test smells, read `tdd-rationale.md` and `testing-anti-patterns.md`.
