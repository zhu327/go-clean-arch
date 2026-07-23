# Testing Anti-Patterns

Load this reference when writing or changing tests, adding fakes or mocks, or considering test-only production hooks.

## Test Observable Behavior

A mock is a tool for isolation, not the outcome under test. A test that only proves a collaborator was called provides weak evidence unless that interaction is itself the observable contract.

Prefer assertions on returned values, persisted state, emitted events, rendered output, or a documented interaction at a genuine external seam.

```
// ❌ Testing mock behavior — proves nothing about real outcome
repoMock.expect('save').resolves()
const result = service.create(data)
expect(repoMock.save).toHaveBeenCalled()  // only verifies mock plumbing

// ✅ Testing observable behavior — captures real outcome
let saved = null
repoMock.expect('save').do((record) => { saved = record })
const result = service.create(data)
expect(saved.createdAt).toBeRecent()       // verifies real behavior
```

## Do Not Add Test-Only Production APIs

Do not add a public reset, setter, bypass, or inspection method solely to make tests convenient. Prefer fresh instances, dependency injection, test fixtures, a test-only helper, or a lower-level seam that already exists.

A production API is justified only when production callers have the same lifecycle or observability need.

```
// ❌ Reset() exists only for test cleanup — pollutes production API
class Manager {
  cache: Map<string, Cert>
  reset() { this.cache.clear() }  // never called in production
}

// ✅ Each test creates a fresh instance — no production footprint
// test helper:
function newTestManager() { return new Manager(new FakeRepo()) }
```

## Mock Deliberately

Before replacing a collaborator, establish:

1. which behavior is slow, nondeterministic, external, destructive, or otherwise unsuitable for the test;
2. which side effects the test depends on;
3. whether a real lightweight implementation or fake is simpler than a mock;
4. what contract the substitute must preserve.

A substitute that omits required state transitions or response fields can make a test fail for the wrong reason—or pass while production fails.

```
// ❌ Partial fixture — downstream code panics on missing .department.name
gatewayMock.expect('getUser').resolves({ id: '1', name: 'Alice' })

// ✅ Contract-complete fixture — mirrors real response structure
gatewayMock.expect('getUser').resolves({
  id: '1', name: 'Alice',
  department: { id: 'd1', name: 'Engineering' },
  role: 'developer', createdAt: '2024-01-01',
})
```

## Avoid Incomplete Fixtures

Build fixtures from the real contract. Include every field, state, or behavior consumed by the code path under test, and use builders or defaults to keep this maintainable. Do not guess that omitted data is irrelevant.

## Do Not Treat Integration Testing as a Final Phase

Testing belongs to implementation. Add the smallest appropriate test before the production behavior, and add integration/E2E coverage when the changed public behavior or project conventions require it.

## Warning Signs

Reconsider the design when:

- setup is larger than the behavior being tested;
- every collaborator must be mocked;
- tests break after harmless internal refactors;
- a test needs a production-only escape hatch;
- assertions mostly inspect calls rather than outcomes;
- a fake needs to recreate a large hidden implementation.

These signs often indicate the interface is too coupled, the seam is in the wrong place, or the test should use a more integrated level.

## Quick Reference

| Smell | Prefer |
|---|---|
| Mock-call-only assertion | Observable output, state, event, or contract assertion |
| Test-only production hook | Fresh fixture, test helper, or existing seam |
| Mock added “just in case” | Identify the concrete isolation need first |
| Partial fake response | Contract-complete fixture with safe defaults |
| Tests written after implementation | Red-green-refactor from the first behavior |
