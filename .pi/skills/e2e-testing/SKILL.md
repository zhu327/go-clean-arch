---
name: e2e-testing
description: Use when writing or updating kfinops end-to-end HTTP API tests backed by httpexpect and testcontainers. Trigger for E2E/end-to-end/integration coverage of new or changed API endpoints, not for unit-only handler tests.
disable-model-invocation: true
---

# E2E API Testing

Write E2E tests that hit real HTTP endpoints backed by Testcontainers (MongoDB + mockd), with external gateways mocked via real HTTP/gRPC calls to a mockd container.

## When to Use

- New API endpoint needs E2E coverage
- Existing endpoint lacks E2E tests
- Bug fix requires E2E regression test
- `/go` E2E gate finds missing or incomplete endpoint coverage

## Quick Reference

| Item | Value |
|------|-------|
| Test location | `test/e2e/<domain>_test.go` |
| Build tag | `//go:build e2e` |
| Package | `e2e` |
| Run command | `make e2e` (via `bash` tool) |
| HTTP client | `github.com/gavv/httpexpect/v2` |
| MongoDB image | `hub-mirror.wps.cn/sreopen/mongo:6.0` |
| mockd image | `hub-mirror.wps.cn/sreopen/mockd:latest` |
| mockd config | `test/e2e/mockd.yaml` |
| Proto files | `test/e2e/protos/` (copied from GOMODCACHE) |
| Auth | `SkipAuthMiddleware: true` → mock user `mock_dev` |
| Response codes | `utils.CodeSuccess(0)`, `CodeInvalidParams(1001)`, `CodeValidationError(1007)` |

## The Process

```
1. Identify target endpoints (read handler + router + DTO via `read` tool)
2. Determine data dependencies (MongoDB collection, BSON fields, gateway mocks)
3. Write test file following the template (use `write` or `edit` tool)
4. Compile: bash("go test -tags=e2e -c -o /dev/null ./test/e2e/")
5. Run: bash("make e2e")
6. All tests must pass (existing + new)
```

## Step 1: Identify Target Endpoints

Before writing any test, `read` these files for the target domain:

```
internal/{domain}/adapter/delivery/http/router/{domain}.go    → routes, middleware, admin checks
internal/{domain}/adapter/delivery/http/handler/{domain}.go   → request binding, response shape
internal/{domain}/adapter/delivery/http/dto/{domain}.go       → JSON/form field names, binding tags
internal/{domain}/adapter/repository/{domain}_mongo.go        → collection name, BSON field names
```

Key questions to answer:
- **Route path and HTTP method** (GET/POST/PUT/DELETE)
- **Request shape** — query params (`form:` tags) vs JSON body (`json:` tags), required fields (`binding:"required"`)
- **Response shape** — what keys are in `data`
- **MongoDB collection name** — constant at top of repository file
- **BSON field names** — may differ from JSON names
- **Middleware** — `CheckAdmin` (bypassed in E2E), `AssetServiceTreeAuth` / `BillingServiceTreeAuth` (affects data visibility)

## Step 2: Understand Service Tree Auth

Many routes use service-tree-based authorization middleware. With mock gateways:

| Middleware | Behavior with mocks | Workaround |
|-----------|-------------------|------------|
| `CheckAdmin` | Bypassed (`SkipAuthMiddleware: true`) | None needed |
| `AssetServiceTreeAuth` | Mock returns empty authorized projects | Seed data with `manager_id: "mock_dev"` or `tag.manager_id: "mock_dev"` |
| `BillingServiceTreeAuth` | Non-admin gets 403 (no departments) | Seed admin permission for `mock_dev` in `permission` collection |
| In-handler `IsAdmin` | Returns false for `mock_dev` | Seed data with `applicant_id: "mock_dev"` |

## Step 3: Write the Test

### File Template

```go
//go:build e2e

package e2e

import (
    "net/http"
    "testing"

    "ksogit.kingsoft.net/sre/kfinops/pkg/utils"
)

const myCollection = "my_collection"  // from repository constant

func TestMyEndpointListEmpty(t *testing.T) {
    client := getMongoClient(t)
    cleanCollection(t, client, myCollection)
    t.Cleanup(func() { cleanCollection(t, client, myCollection) })

    e := newExpect(t)
    obj := e.GET("/api/v1/my_endpoint").
        WithQuery("limit", 20).
        WithQuery("page_index", 1).
        Expect().
        Status(http.StatusOK).
        JSON().Object()

    obj.HasValue("code", utils.CodeSuccess)
    data := obj.Value("data").Object()
    data.HasValue("total", 0)
    data.Value("list").Array().Empty()
}
```

### Available Helpers

```go
newExpect(t)                                    // httpexpect client pointed at test server
getMongoClient(t)                               // MongoDB client (auto-disconnects via t.Cleanup)
cleanCollection(t, client, collectionName)      // Delete all docs in collection
seedDocuments(t, client, collectionName, docs)  // Insert []interface{} of bson.M documents
```

### Test Categories

Every domain should have at minimum:

| Category | What to test | Example |
|----------|-------------|---------|
| **List Empty** | Empty collection returns zero results | `total: 0`, `list: []` |
| **CRUD** | Create → List → Get → Update → Delete → Verify deleted | Full lifecycle |
| **Validation** | Missing required fields → 400 | `Status(http.StatusBadRequest)` |
| **With Data** | Seed MongoDB, query, verify response fields | `seedDocuments` + assertions |
| **Filter** | Seed multiple docs, filter, verify subset returned | Query params or body filters |

### Seeding MongoDB Data

Use BSON field names (not JSON names). Find them in the repository file's model struct tags.

```go
seedDocuments(t, client, myCollection, []interface{}{
    bson.M{
        "_id":            primitive.NewObjectID(),
        "bson_field_name": "value",        // Use bson:"..." tag, NOT json:"..." tag
        "created_at":     time.Now(),
        "updated_at":     time.Now(),
    },
})
```

### Assertion Patterns

```go
// Success response
obj.HasValue("code", utils.CodeSuccess)
obj.HasValue("message", "success")

// Nested data
data := obj.Value("data").Object()
data.HasValue("total", 1)
data.Value("list").Array().Length().Equal(1)

// Validation error
e.POST("/api/v1/endpoint").
    WithJSON(map[string]any{"incomplete": "body"}).
    Expect().
    Status(http.StatusBadRequest).
    JSON().Object().
    HasValue("code", utils.CodeInvalidParams)
```

### Mock Gateway Responses (mockd)

External gateways are mocked via a **mockd container**. If you need specific mock data, edit `test/e2e/mockd.yaml` (use `edit` tool).

## Step 4: Compile and Run

```bash
# Compile check (fast, no Docker needed)
go test -tags=e2e -c -o /dev/null ./test/e2e/

# Full run
make e2e
```

**All 44+ existing tests must still pass.** Never break existing tests.

## Quality Gate

Before marking E2E work as complete:

- [ ] Every new test file has `//go:build e2e` and `package e2e`
- [ ] Collection cleaned before test AND in `t.Cleanup`
- [ ] Seeded data uses BSON field names from repository model
- [ ] Happy path AND error/validation cases covered
- [ ] `bash("go test -tags=e2e -c")` compiles without errors
- [ ] `bash("make e2e")` passes — all tests green (existing + new)
- [ ] No hardcoded server URLs (use `newExpect(t)`)

## Integration with Other Skills

- **/skill:test-driven-development** — E2E tests complement unit tests, not replace them
- **/go** — validates E2E coverage in the conditional E2E gate
- **/skill:writing-plans** — identifies endpoint changes and includes E2E authoring in the plan
- **/skill:subagent-driven-development** — creates or updates E2E tests during implementation, then `/go` or the controller runs `make e2e`
