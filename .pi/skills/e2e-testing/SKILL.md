---
name: e2e-testing
description: Use when adding or changing end-to-end or integration tests that exercise a running application through real external interfaces. Derives the test strategy, tooling, commands, and fixtures from the current project.
---

# End-to-End Testing

Write tests that exercise the application through its public boundary (for example HTTP, CLI, queue, or RPC) with production-like dependencies where practical. This skill is intentionally independent of language, framework, database, and test runner.

## Before Writing Tests

Inspect the project rather than assuming its test stack:

1. Read the target entry point, route/command/consumer registration, request and response contracts, and relevant persistence or gateway code.
2. Read existing end-to-end tests, test helpers, CI configuration, container configuration, and documented test commands.
3. Identify the actual public contract, required authentication/authorization, data dependencies, external integrations, and cleanup requirements.
4. Reuse the project’s established runner, fixtures, helper functions, and commands. If none exist, propose the smallest suitable harness before creating one.

Never assume a database, container runtime, mock server, repository layout, package name, test count, or command.

## Test Design

For each changed public behavior, cover the applicable scenarios:

- successful primary flow;
- invalid input and contract validation;
- authentication, authorization, tenancy, or ownership boundaries;
- persisted state and observable side effects;
- important error mapping and external-service failure behavior;
- filters, pagination, idempotency, or lifecycle transitions where relevant.

Keep tests independent:

- create only the data each test needs;
- clean up before and after a test using project helpers or test-scoped resources;
- avoid fixed ports, shared mutable state, timing sleeps, and dependence on test order;
- assert stable public behavior, not internal implementation details.

## Implementation Process

1. Add the smallest failing end-to-end test for one behavior.
2. Run the project’s focused E2E command and verify it fails for the expected missing behavior.
3. Implement the minimal production change or test harness adjustment needed to pass.
4. Re-run the focused test, then the complete relevant E2E suite.
5. Run the project’s documented build, unit-test, and static-analysis commands when available.

When containers or external systems are required, use project-supported ephemeral resources and deterministic fakes for systems the test cannot safely control. Do not invent mocks that bypass the boundary being tested.

## Common Patterns

These patterns appear across projects regardless of language or framework. Adapt to the project's actual test runner, assertion library, and fixtures.

### CRUD Lifecycle

```
// Create → verify in list → get detail → delete → verify removed
const id = await client.create('/api/v1/things', { name: 'test' })

const list = await client.list('/api/v1/things')
expect(list.total).toBe(1)

const detail = await client.get(`/api/v1/things/${id}`)
expect(detail.name).toBe('test')

await client.delete(`/api/v1/things/${id}`)

const after = await client.list('/api/v1/things')
expect(after.total).toBe(0)
```

### Auth / Multi-Tenancy Data Seeding

Seed data with the test user's identity so authorization filters pass:

```
// Seed resource owned by the test user
await seed('things', { owner_id: TEST_USER.id, name: 'owned-by-me' })
await seed('things', { owner_id: 'other-user',   name: 'not-mine' })

// List endpoint should only return the owned resource
const result = await client.list('/api/v1/things')
expect(result.total).toBe(1)
expect(result.items[0].name).toBe('owned-by-me')
```

### Filter / Pagination

```
// Seed varied data
await seed('things', { type: 'A', status: 'active' })
await seed('things', { type: 'B', status: 'active' })
await seed('things', { type: 'A', status: 'archived' })

// Filter + paginate
const page1 = await client.list('/api/v1/things', { type: 'A', limit: 1, page: 1 })
expect(page1.total).toBe(2)
expect(page1.items.length).toBe(1)
```

## Quality Gate

Before completion, verify:

- [ ] The test invokes the public interface rather than private implementation details.
- [ ] Test data, credentials, and external responses are isolated and cleaned up.
- [ ] Assertions cover the documented contract and important failure cases.
- [ ] No project-specific URLs, secrets, images, users, or infrastructure assumptions were hard-coded without an existing project convention.
- [ ] The focused and full relevant E2E suites pass.
- [ ] The project’s available validation commands pass.

## Planning Integration

When a plan changes a public interface, include E2E work as a dependent task only when the project has an E2E harness or the approved requirements call for one. The plan must name the actual project command and harness after inspection; it must not prescribe a language- or framework-specific command by default.
