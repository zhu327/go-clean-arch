---
name: executing-plans
description: Use when executing a written implementation plan task-by-task with todo tracking, TDD, verification commands, and human review checkpoints between batches. Outside /go; /go uses subagent-driven-development instead.
disable-model-invocation: true
---

# Executing Plans

## Overview

Load plan, review critically, execute tasks in batches, report for review between batches.

**Core principle:** Batch execution with checkpoints for architect review.

**Announce at start:** "I'm using the executing-plans skill to implement this plan."

## The Process

### Step 1: Load and Review Plan
1. Read plan file (use `read` tool)
2. Review critically - identify any questions or concerns about the plan
3. If concerns: Raise them with your human partner before starting
4. If no concerns: Create `todo` items and proceed

### Step 2: Execute Batch
**Default: First 3 tasks**

For each task:
1. Mark as in_progress in `todo`
2. Follow each step exactly (plan has bite-sized steps), using /skill:test-driven-development
3. Run verifications as specified (use `bash` tool)
4. Mark as completed in `todo`

### Step 3: Report
When batch complete:
- Show what was implemented
- Show verification output
- Say: "Ready for feedback."

### Step 4: Continue
Based on feedback:
- Apply changes if needed (use `edit` tool)
- Execute next batch
- Repeat until complete

## When to Stop and Ask for Help

**STOP executing immediately when:**
- Hit a blocker mid-batch (missing dependency, test fails, instruction unclear)
- Plan has critical gaps preventing starting
- You don't understand an instruction
- Verification fails repeatedly

**Ask for clarification rather than guessing.**

## When to Revisit Earlier Steps

**Return to Review (Step 1) when:**
- Partner updates the plan based on your feedback
- Fundamental approach needs rethinking

**Don't force through blockers** - stop and ask.

## Remember
- Review plan critically first
- Follow TDD: no production code without a failing test first
- Follow plan steps exactly
- Don't skip verifications
- Reference skills when plan says to (use /skill: syntax)
- Between batches: just report and wait
- Stop when blocked, don't guess
- Never start implementation on main/master branch without explicit user consent

## Integration

`/go` does not use this skill; `/go` always routes execution through `/skill:subagent-driven-development`. Use this skill outside `/go` when you want checkpointed, human-reviewed batch execution.

**Required workflow skills:**
- **/skill:writing-plans** - Creates the plan this skill executes
- **/skill:test-driven-development** - TDD discipline for all implementation
- **/skill:e2e-testing** - E2E API tests for new/modified endpoints (run `make e2e` after implementation)
