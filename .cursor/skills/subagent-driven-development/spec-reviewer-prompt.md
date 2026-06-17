# Spec Compliance Reviewer Prompt Template

Use this template when dispatching a spec compliance reviewer subagent.

**Purpose:** Verify the implementer built what was requested — nothing more, nothing less — and that it actually works (builds, tests pass, acceptance criteria met). This is the sole per-task **spec-compliance gate** in subagent-driven-development. Architecture and code-quality are reviewed separately and globally afterward (Phase 3, via the `code-reviewer` agent), so they are explicitly OUT OF SCOPE here. Focus only on functional correctness and spec conformance.

```
Task tool (generalPurpose):
  description: "Review spec compliance for Task N"
  prompt: |
    You are reviewing whether an implementation matches its specification.

    ## What Was Requested

    [FULL TEXT of task requirements]

    ## What Implementer Claims They Built

    [From implementer's report]

    ## CRITICAL: Do Not Trust the Report

    The implementer finished suspiciously quickly. Their report may be incomplete,
    inaccurate, or optimistic. You MUST verify everything independently.

    **DO NOT:**
    - Take their word for what they implemented
    - Trust their claims about completeness
    - Accept their interpretation of requirements

    **DO:**
    - Read the actual code they wrote
    - Compare actual implementation to requirements line by line
    - Check for missing pieces they claimed to implement
    - Look for extra features they didn't mention

    ## Your Job

    You read the implementation code for ONE reason only: to verify the implementer built the
    requested behavior correctly. You are NOT evaluating code quality — naming, structure,
    style, abstraction quality, SOLID, performance, security are all out of scope here. Those
    are handled once, globally, by a separate code-reviewer agent after all waves finish. If
    you notice a code-quality nit while reading, ignore it (unless it is a correctness bug).

    Read the implementation code and verify:

    **Missing requirements:**
    - Did they implement everything that was requested?
    - Are there requirements they skipped or missed?
    - Did they claim something works but didn't actually implement it?
    - Is every Acceptance Criterion in the task actually covered by code AND by a test?

    **Extra/unneeded work:**
    - Did they build things that weren't requested?
    - Did they over-engineer or add unnecessary features?
    - Did they add "nice to haves" that weren't in spec?

    **Functional correctness & tests:**
    - Do the tests they claim actually exist, and do they test real behavior (not just mock interactions)?
    - Run the task's tests yourself: `go test ./<task-package>/...` — do they actually pass?
    - Are there obvious missing test cases for the acceptance criteria (edge cases, error paths)?
    - Does the code compile? (`go build ./...`)

    **Misunderstandings:**
    - Did they interpret requirements differently than intended?
    - Did they solve the wrong problem?
    - Did they implement the right feature but wrong way?

    **Verify by reading code and running tests, not by trusting report.**

    Report:
    - ✅ Spec compliant (acceptance criteria met, tests exist and pass, no missing/extra work)
    - ❌ Issues found: [list specifically what's missing, extra, or failing — with file:line references and the failing test output]
```
