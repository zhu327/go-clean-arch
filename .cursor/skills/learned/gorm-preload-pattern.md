# GORM Preload Pattern

**Extracted:** 2026-01-28
**Use case:** Avoid N+1 query problems when loading associated data

## Problem

Querying associated data one-by-one inside a loop causes N+1 query issues, severely impacting performance.

## Solution

Use GORM's Preload method to load associated data in a single query.

## Example

```go
// Wrong - N+1 queries
certs, _ := repo.List(ctx)
for _, cert := range certs {
    db.Where("cert_id = ?", cert.ID).Find(&cert.Deployments)
}

// Correct - Preload in one query
db.Preload("Deployments").Find(&certs)

// Correct - Conditional Preload
db.Preload("Deployments", "status = ?", "active").Find(&certs)

// Correct - Nested Preload
db.Preload("Deployments.LoadBalancer").Find(&certs)
```

## When to Use

- List endpoints need to return associated data
- SQL logs show many repeated queries
- API response times are too long
