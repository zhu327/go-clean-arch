# Testing Anti-Patterns

**Load this reference when:** writing or changing tests, adding mocks, or tempted to add test-only methods to production code.

## Overview

Tests must verify real behavior, not mock behavior. Mocks are a means to isolate, not the thing being tested.

**Core principle:** Test what the code does, not what the mocks do.

**Following strict TDD prevents these anti-patterns.**

## The Iron Laws

```
1. NEVER test mock behavior
2. NEVER add test-only methods to production classes
3. NEVER mock without understanding dependencies
```

## Anti-Pattern 1: Testing Mock Behavior

**The violation:**
```go
// ❌ BAD: Testing that the mock was called, not real behavior
func TestService_CreateUser(t *testing.T) {
    ctrl := gomock.NewController(t)
    mockRepo := mock.NewMockUserRepository(ctrl)
    mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

    svc := NewUserService(mockRepo)
    _ = svc.CreateUser(context.Background(), &User{})

    // Only verifying mock was called - proves nothing about real behavior!
}
```

**Why this is wrong:**
- You're verifying the mock works, not that the service works
- Test passes when mock is called, fails when it's not
- Tells you nothing about real behavior

**your human partner's correction:** "Are we testing the behavior of a mock?"

**The fix:**
```go
// ✅ GOOD: Test real behavior through assertions
func TestService_CreateUser_SetsCreatedAt(t *testing.T) {
    ctrl := gomock.NewController(t)
    mockRepo := mock.NewMockUserRepository(ctrl)

    var savedUser *User
    mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
        DoAndReturn(func(ctx context.Context, u *User) error {
            savedUser = u
            return nil
        })

    svc := NewUserService(mockRepo)
    err := svc.CreateUser(context.Background(), &User{Name: "test"})

    assert.NoError(t, err)
    assert.NotZero(t, savedUser.CreatedAt) // Testing real behavior!
}
```

### Gate Function

```
BEFORE asserting on any mock call:
  Ask: "Am I testing real behavior or just mock interaction?"

  IF testing mock interaction only:
    STOP - Add assertions on actual output or state changes

  Test real behavior instead
```

## Anti-Pattern 2: Test-Only Methods in Production

**The violation:**
```go
// ❌ BAD: Reset() only used in tests
type CertificateManager struct {
    cache map[string]*Certificate
}

// Reset 重置缓存 - 只在测试中使用！
func (m *CertificateManager) Reset() {
    m.cache = make(map[string]*Certificate)
}

// In tests
func TestSomething(t *testing.T) {
    defer manager.Reset() // Test-only cleanup
}
```

**Why this is wrong:**
- Production struct polluted with test-only code
- Dangerous if accidentally called in production
- Violates YAGNI and separation of concerns
- Confuses object lifecycle with entity lifecycle

**The fix:**
```go
// ✅ GOOD: Test utilities handle test cleanup
// CertificateManager has no Reset() - it's stateless in production

// In testutil/helpers.go
func CleanupManager(t *testing.T, m *CertificateManager) {
    t.Helper()
    // Use proper cleanup through interfaces or recreate manager
}

// In tests - create fresh instance per test
func TestSomething(t *testing.T) {
    manager := NewCertificateManager(mockRepo)
    // Each test gets fresh state
}
```

### Gate Function

```
BEFORE adding any method to production struct:
  Ask: "Is this only used by tests?"

  IF yes:
    STOP - Don't add it
    Put it in testutil/ or create fresh instances per test

  Ask: "Does this struct own this resource's lifecycle?"

  IF no:
    STOP - Wrong struct for this method
```

## Anti-Pattern 3: Mocking Without Understanding

**The violation:**
```go
// ❌ BAD: Mock breaks test logic
func TestManager_DetectsDuplicateCert(t *testing.T) {
    ctrl := gomock.NewController(t)
    // Mock prevents cache update that test depends on!
    mockRepo := mock.NewMockCertRepository(ctrl)
    mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

    manager := NewManager(mockRepo)
    _ = manager.CreateCert(ctx, cert)
    err := manager.CreateCert(ctx, cert) // Should detect duplicate - but won't!

    assert.Error(t, err) // FAILS: mock doesn't track state
}
```

**Why this is wrong:**
- Mocked method had side effect test depended on (updating state)
- Over-mocking to "be safe" breaks actual behavior
- Test passes for wrong reason or fails mysteriously

**The fix:**
```go
// ✅ GOOD: Mock at correct level with state tracking
func TestManager_DetectsDuplicateCert(t *testing.T) {
    ctrl := gomock.NewController(t)
    mockRepo := mock.NewMockCertRepository(ctrl)

    // Track state in mock
    created := make(map[string]bool)
    mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
        DoAndReturn(func(ctx context.Context, c *Cert) error {
            if created[c.ID] {
                return utils.ConflictError("already exists")
            }
            created[c.ID] = true
            return nil
        }).AnyTimes()

    manager := NewManager(mockRepo)
    _ = manager.CreateCert(ctx, cert)
    err := manager.CreateCert(ctx, cert) // Duplicate detected ✓

    assert.Error(t, err)
}
```

### Gate Function

```
BEFORE mocking any method:
  STOP - Don't mock yet

  1. Ask: "What side effects does the real method have?"
  2. Ask: "Does this test depend on any of those side effects?"
  3. Ask: "Do I fully understand what this test needs?"

  IF depends on side effects:
    Mock at lower level (the actual slow/external operation)
    OR use DoAndReturn to preserve necessary behavior
    NOT the high-level method the test depends on

  IF unsure what test depends on:
    Run test with real implementation FIRST
    Observe what actually needs to happen
    THEN add minimal mocking at the right level

  Red flags:
    - "I'll mock this to be safe"
    - "This might be slow, better mock it"
    - Mocking without understanding the dependency chain
```

## Anti-Pattern 4: Incomplete Mocks

**The violation:**
```go
// ❌ BAD: Partial mock - only fields you think you need
mockGateway.EXPECT().GetUserInfo(gomock.Any(), userID).Return(&UserInfo{
    ID:   userID,
    Name: "Alice",
    // Missing: Department, Role that downstream code uses
}, nil)

// Later: nil pointer panic when code accesses user.Department.Name
```

**Why this is wrong:**
- **Partial mocks hide structural assumptions** - You only mocked fields you know about
- **Downstream code may depend on fields you didn't include** - Silent failures or panics
- **Tests pass but integration fails** - Mock incomplete, real API complete
- **False confidence** - Test proves nothing about real behavior

**The Iron Rule:** Mock the COMPLETE data structure as it exists in reality, not just fields your immediate test uses.

**The fix:**
```go
// ✅ GOOD: Mirror real API completeness
mockGateway.EXPECT().GetUserInfo(gomock.Any(), userID).Return(&UserInfo{
    ID:         userID,
    Name:       "Alice",
    Department: &Department{ID: "dept-1", Name: "Engineering"},
    Role:       "developer",
    CreatedAt:  time.Now(),
    // All fields real API returns
}, nil)
```

### Gate Function

```
BEFORE creating mock responses:
  Check: "What fields does the real Gateway/API response contain?"

  Actions:
    1. Examine actual Gateway interface and domain struct
    2. Include ALL fields system might consume downstream
    3. Verify mock matches real response struct completely

  Critical:
    If you're creating a mock, you must understand the ENTIRE structure
    Partial mocks cause nil pointer panics when code depends on omitted fields

  If uncertain: Include all struct fields from domain definition
```

## Anti-Pattern 5: Integration Tests as Afterthought

**The violation:**
```
✅ Implementation complete
❌ No tests written
"Ready for testing"
```

**Why this is wrong:**
- Testing is part of implementation, not optional follow-up
- TDD would have caught this
- Can't claim complete without tests

**The fix:**
```
TDD cycle:
1. Write failing test
2. Implement to pass
3. Refactor
4. THEN claim complete
```

## When Mocks Become Too Complex

**Warning signs:**
- Mock setup longer than test logic
- Mocking everything to make test pass
- Mocks missing methods real components have
- Test breaks when mock changes

**your human partner's question:** "Do we need to be using a mock here?"

**Consider:** Integration tests with real components often simpler than complex mocks

## TDD Prevents These Anti-Patterns

**Why TDD helps:**
1. **Write test first** → Forces you to think about what you're actually testing
2. **Watch it fail** → Confirms test tests real behavior, not mocks
3. **Minimal implementation** → No test-only methods creep in
4. **Real dependencies** → You see what the test actually needs before mocking

**If you're testing mock behavior, you violated TDD** - you added mocks without watching test fail against real code first.

## Quick Reference

| Anti-Pattern | Fix |
|--------------|-----|
| Only verify mock calls | Add assertions on actual output/state |
| Test-only methods in production | Move to testutil/ or recreate per test |
| Mock without understanding | Understand dependencies first, use DoAndReturn |
| Incomplete mocks | Mirror real struct/API completely |
| Tests as afterthought | TDD - tests first |
| Over-complex mocks | Consider integration tests with real DB |

## Red Flags

- Only verifying `mockXxx.EXPECT()` calls without real assertions
- Methods only called in `*_test.go` files
- Mock setup with `DoAndReturn` is >50% of test
- Test fails when you remove mock
- Can't explain why mock is needed
- Mocking "just to be safe"
- Using `.AnyTimes()` without understanding call patterns

## The Bottom Line

**Mocks are tools to isolate, not things to test.**

If TDD reveals you're testing mock behavior, you've gone wrong.

Fix: Test real behavior or question why you're mocking at all.
