---
name: using-superpowers
description: Use when starting any conversation - establishes how to find and use skills, requiring skill invocation BEFORE ANY response including clarifying questions
---

<EXTREMELY-IMPORTANT>
If you think there is even a 1% chance a skill might apply to what you are doing, you ABSOLUTELY MUST load and follow the skill.

IF A SKILL APPLIES TO YOUR TASK, YOU DO NOT HAVE A CHOICE. YOU MUST USE IT.

This is not negotiable. This is not optional. You cannot rationalize your way out of this.
</EXTREMELY-IMPORTANT>

# Using Skills

## The Rule

**Load and follow relevant skills BEFORE any response or action.** Even a 1% chance a skill might apply means that you should load the skill to check. If a loaded skill turns out to be wrong for the situation, you don't need to use it.

```
User message received
  ↓
Might any skill apply? (even 1%)
  ├─ yes → Load skill via /skill:name or read the SKILL.md
  │         ↓
  │       Announce: "Using [skill] to [purpose]"
  │         ↓
  │       Has checklist? → Create todo items per item
  │         ↓
  │       Follow skill exactly
  └─ no → Respond (including clarifications)
```

## How to Load Skills in Pi

Skills are loaded on-demand. To use a skill:
- Type `/skill:skill-name` in the prompt to force-load
- Or reference by name — the agent will `read` the SKILL.md automatically
- Skills are located in `.pi/skills/` (project) or `~/.pi/agent/skills/` (global)

## Red Flags

These thoughts mean STOP—you're rationalizing:

| Thought | Reality |
|---------|---------|
| "This is just a simple question" | Questions are tasks. Check for skills. |
| "I need more context first" | Skill check comes BEFORE clarifying questions. |
| "Let me explore the codebase first" | Skills tell you HOW to explore. Check first. |
| "I can check git/files quickly" | Files lack conversation context. Check for skills. |
| "Let me gather information first" | Skills tell you HOW to gather information. |
| "This doesn't need a formal skill" | If a skill exists, use it. |
| "I remember this skill" | Skills evolve. Read current version. |
| "This doesn't count as a task" | Action = task. Check for skills. |
| "The skill is overkill" | Simple things become complex. Use it. |
| "I'll just do this one thing first" | Check BEFORE doing anything. |

## Skill Priority

When multiple skills could apply, use this order:

1. **Process skills first** (brainstorming, writing-plans) - these determine HOW to approach the task
2. **Execution skills second** (executing-plans, subagent-driven-development, test-driven-development) - these guide implementation
3. **Review skills last** (code-review-expert) - these validate the result

"Let's build X" → brainstorming first, then execution skills.

## Skill Types

**Rigid** (TDD): Follow exactly. Don't adapt away discipline.

**Flexible** (patterns): Adapt principles to context.

The skill itself tells you which.

## User Instructions

Instructions say WHAT, not HOW. "Add X" or "Fix Y" doesn't mean skip workflows.
