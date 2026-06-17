# Implementer Subagent Prompt Template

Use this template when dispatching an implementer subagent.

## Standard (Sequential or Single Task)

```
Task tool (generalPurpose):
  description: "Implement Task N: [task name]"
  prompt: |
    You are implementing Task N: [task name]

    ## Task Description

    [FULL TEXT of task from plan - paste it here, don't make subagent read file]

    ## Context

    [Scene-setting: where this fits, dependencies, architectural context]

    ## Before You Begin

    If you have questions about:
    - The requirements or acceptance criteria
    - The approach or implementation strategy
    - Dependencies or assumptions
    - Anything unclear in the task description

    **Ask them now.** Raise any concerns before starting work.

    ## Your Job

    Once you're clear on requirements, follow TDD (red-green-refactor):
    1. Write a failing test (RED)
    2. Write minimal code to pass (GREEN)
    3. Refactor if needed
    4. Repeat for each behavior in the task
    5. Self-review (see below)
    6. Report back

    Work from: [directory]

    **While you work:** If you encounter something unexpected or unclear, **ask questions**.
    It's always OK to pause and clarify. Don't guess or make assumptions.

    ## Before Reporting Back: Self-Review

    Review your work with fresh eyes. Ask yourself:

    **Completeness:**
    - Did I fully implement everything in the spec?
    - Did I miss any requirements?
    - Are there edge cases I didn't handle?

    **Quality:**
    - Is this my best work?
    - Are names clear and accurate (match what things do, not how they work)?
    - Is the code clean and maintainable?

    **Discipline:**
    - Did I avoid overbuilding (YAGNI)?
    - Did I only build what was requested?
    - Did I follow existing patterns in the codebase?

    **Testing:**
    - Did I write each test before its implementation?
    - Did I watch each test fail before writing code?
    - Do tests verify behavior (not just mock behavior)?
    - Are tests comprehensive?

    If you find issues during self-review, fix them now before reporting.

    ## Report Format

    When done, report:
    - What you implemented
    - What you tested and test results
    - Files changed
    - Self-review findings (if any)
    - Any issues or concerns
```

## Parallel Wave Variant

When dispatching multiple implementers in the same wave, **add these sections** to each prompt:

```
    ## Parallel Execution Notice

    You are running in parallel with other implementer agents in this wave.
    Other agents are working on: [list other task names in this wave]

    **DO NOT** modify any files outside your task's file list.
    Your files:
    - [Create/Modify file list from plan]

    If you discover you need to modify a file not in your list, STOP and report it
    instead of making the change. The controller will coordinate the conflict.
```

## Prior Wave Context (for Wave N > 0)

When dispatching implementers that depend on prior wave outputs, the Controller MUST read the actual generated source files and extract exact signatures. **Add this section** with real code:

```
    ## Prior Wave Dependencies

    You depend on the following components created in previous waves.
    Here are their **actual signatures** (read from generated source files):

    ### From Task X ([task name]):
    `[file path]`:
    ```go
    [paste actual exported interface/struct/function signatures from the file]
    ```

    ### From Task Y ([task name]):
    `[file path]`:
    ```go
    [paste actual exported interface/struct/function signatures from the file]
    ```

    You MUST use these exact type names and method signatures.
    If you find a mismatch between the plan and the actual generated code,
    STOP and report it instead of guessing.
```

**Controller workflow for artifact passing:**
1. After a wave completes, look at next wave's tasks and their "Blocked by" fields
2. For each dependency, read the generated source files using Read tool
3. Extract exported interfaces, structs, and function signatures
4. Paste the actual code into the downstream implementer's prompt
