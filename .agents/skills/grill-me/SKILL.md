---
name: grill-me
description: Relentless interview to sharpen a plan or design. Use when the user wants to stress-test their thinking, asks to be grilled, or a task hinges on consequential unresolved decisions.
---

# Grill Me

Interview the user relentlessly until you reach shared understanding. Do not act on the result until the user confirms.

## Rounds and the frontier

Map the discussion as a **design tree**: every decision branches into the decisions hanging off it. Work the tree in rounds.

The **frontier** is every decision whose prerequisites are already settled — questions you can ask *now* without guessing at answers not yet heard. Ask the whole frontier in one round, capped at the four highest-consequence questions (the `question` tool's limit); the rest wait for the next round. A question whose answer depends on a still-open question belongs to a later round.

Each round:

1. Recompute the frontier from the user's latest answers — settled decisions push it outward and unblock downstream questions.
2. Ask every frontier question, numbered, with your recommended answer:

```markdown
❓ **Q1** — **<title>**: <body; concrete choices when alternatives are known>
➡️ <recommended answer>
```

3. Wait for all answers before the next round.

In pi, deliver rounds via the `question` tool (options = the credible alternatives, recommendation in the prompt); fall back to the markdown format when questions need long context or free-form answers.

## Pruning

Not every branch deserves a question. Prune decisions that are reversible, low-impact, derivable from project conventions, or safely handled by a stated assumption — record those as assumptions instead of asking. The interview is done when no unresolved decision can materially change the design or risk profile, not when every branch has been visited.

## Facts vs decisions

Finding *facts* is your job, never the user's. When a frontier question needs a fact from the environment (code, config, tool behavior), find it yourself — read/grep the repo or dispatch a read-only exploration subagent where the harness provides one. Don't block on it: a running lookup is an unsettled prerequisite, so only questions downstream of it wait; ask the rest of the frontier now. The *decisions* are the user's: put each to them and wait.

## Done

The session ends when no unresolved decision can materially change the design or risk profile — remaining branches are recorded as explicit assumptions. Summarize the settled decisions and explicitly deferred scope, get the user's confirmation, then hand off: multi-step work → skill `plan-execute`; a clear local change → implement directly.
