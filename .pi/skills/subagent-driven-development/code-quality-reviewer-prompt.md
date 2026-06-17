# Code Quality Reviewer Prompt Template

Use this template when dispatching a code quality reviewer subagent.

**Purpose:** Verify implementation is well-built (clean, tested, maintainable)

**Only dispatch after spec compliance review passes.**

```json
{
  "agent": "code-reviewer",
  "task": "Use /skill:code-review-expert\n\nWHAT_WAS_IMPLEMENTED: [from implementer's report]\nPLAN_OR_REQUIREMENTS: Task N from [plan-file]\nDESCRIPTION: [task summary]"
}
```

**Code reviewer returns:** Strengths, Issues (Critical/Important/Minor), Assessment
