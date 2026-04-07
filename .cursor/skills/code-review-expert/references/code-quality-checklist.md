# Code Quality Checklist

## Error Handling

### Anti-patterns to Flag

- **Ignored errors**: Discarding error return values
  ```go
  result, _ := someFunction()  // Silent failure
  _ = db.Save(&entity)         // Ignoring DB error
  ```
- **Log and forget**: Logging error but not returning/handling it
- **Error information leakage**: Stack traces or internal details exposed to users
- **Missing error wrapping**: Use `utils.AppError` instead of `errors.New`
- **Goroutine errors**: Errors in goroutines not propagated to caller

### Best Practices to Check

- [ ] Errors are caught at appropriate boundaries
- [ ] Error messages are user-friendly (no internal details exposed)
- [ ] Errors are logged with sufficient context for debugging
- [ ] Async errors are properly propagated or handled
- [ ] Fallback behavior is defined for recoverable errors
- [ ] Critical errors trigger alerts/monitoring

### Questions to Ask
- "What happens when this operation fails?"
- "Will the caller know something went wrong?"
- "Is there enough context to debug this error?"

---

## Performance & Caching

### CPU-Intensive Operations

- **Expensive operations in hot paths**: Regex compilation, JSON parsing, crypto in loops
- **Blocking main thread**: Sync I/O, heavy computation without worker/async
- **Unnecessary recomputation**: Same calculation done multiple times
- **Missing memoization**: Pure functions called repeatedly with same inputs

### Database & I/O

- **N+1 queries**: Loop that makes a query per item instead of batch
  ```go
  // Bad: N+1
  for _, id := range ids {
      var user User
      db.First(&user, id)
  }
  // Good: Use Preload or batch query
  db.Preload("Users").Find(&orders)
  db.Where("id IN ?", ids).Find(&users)
  ```
- **Missing indexes**: Queries on unindexed columns
- **Over-fetching**: SELECT * when only few columns needed
- **No pagination**: Loading entire dataset into memory

### Caching Issues

- **Missing cache for expensive operations**: Repeated API calls, DB queries, computations
- **Cache without TTL**: Stale data served indefinitely
- **Cache without invalidation strategy**: Data updated but cache not cleared
- **Cache key collisions**: Insufficient key uniqueness
- **Caching user-specific data globally**: Security/privacy issue

### Memory

- **Unbounded collections**: Arrays/maps that grow without limit
- **Large object retention**: Holding references preventing GC
- **String concatenation in loops**: Use `strings.Builder` instead
- **Loading large files entirely**: Use streaming instead

### Questions to Ask
- "What's the time complexity of this operation?"
- "How does this behave with 10x/100x data?"
- "Is this result cacheable? Should it be?"
- "Can this be batched instead of one-by-one?"

---

## Boundary Conditions

### Nil/Zero Value Handling

- **Missing nil checks**: Accessing fields on potentially nil pointers
- **Zero value confusion**: `if value != ""` when empty string is valid
- **Pointer vs value**: Returning pointer when nil is not a valid state

### Empty Collections

- **Empty array not handled**: Code assumes array has items
- **Empty object edge case**: `for...in` or `Object.keys` on empty object
- **First/last element access**: `arr[0]` or `arr[arr.length-1]` without length check

### Numeric Boundaries

- **Division by zero**: Missing check before division
- **Integer overflow**: Large numbers exceeding safe integer range
- **Floating point comparison**: Using `===` instead of epsilon comparison
- **Negative values**: Index or count that shouldn't be negative
- **Off-by-one errors**: Loop bounds, array slicing, pagination

### String Boundaries

- **Empty string**: Not handled as edge case
- **Whitespace-only string**: Passes truthy check but is effectively empty
- **Very long strings**: No length limits causing memory/display issues
- **Unicode edge cases**: Emoji, RTL text, combining characters

### Common Patterns to Flag

```go
// Dangerous: no nil check
name := user.Profile.Name

// Dangerous: slice access without check
first := items[0]

// Dangerous: division without check
avg := total / count

// Dangerous: zero value check excludes valid values
if value != 0 { ... }  // fails when 0 is valid
```

### Questions to Ask
- "What if this pointer is nil?"
- "What if this slice is empty?"
- "What's the valid range for this number?"
- "What happens at the boundaries (0, -1, MaxInt)?"
