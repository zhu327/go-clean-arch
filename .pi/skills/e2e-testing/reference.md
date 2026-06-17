# E2E Testing Reference

## Test Infrastructure

### TestMain (DO NOT MODIFY)

`test/e2e/e2e_test.go` bootstraps the entire test environment:

1. Starts MongoDB container via testcontainers
2. Builds config via `buildE2EConfig(mongoAddr)`
3. Initializes app via `di.InitializeAPI(cfg)`
4. Starts `httptest.NewServer` (in-process, no real port needed)
5. Runs all tests via `m.Run()`
6. Tears down via `defer` (container always cleaned up)

### E2E Config (`test/e2e/config.go`)

All external services are mocked via `UseMock: true`:

```
SRECore.UserCenterUseMock, SRECore.OrgUseMock
Kitsm.UseMock
CloudService.UseMock, CloudKit.UseMock
KLB.UseMock, KDNS.UseMock, KCMDB.UseMock
ServiceTree.UseMock, IngressGateway.UseMock
WOA.UseMock, Prom.UseMock, TfEngine.UseMock
```

Background tasks are disabled:
```
EnableResourceSnapshotTask: false
EnableKdnsInternalRecordsSync: false
EnableCertificateTopoSync: false
EnableCloudMetadataSync: false
EnableCloudOrderCheckTask: false
EnableSubDomainQpsSloSync: false
```

### Mock User Context

`SkipAuthMiddleware: true` injects:
```go
mockUser := &domain.UserInfo{
    ID:           "mock_dev",
    Name:         "Mock Developer",
    Account:      "mock_dev",
    Email:        "mock_dev@test.com",
    DepartmentId: "dept-001",
}
```

## Existing Test Files (44 tests)

| File | Tests | Domain | Data Source |
|------|-------|--------|------------|
| `health_test.go` | 2 | Health | None |
| `error_test.go` | 1 | Error | None |
| `tag_test.go` | 6 | Tag | Mock KCMDB gateway |
| `cloud_account_test.go` | 3 | Cloud Account | MongoDB `cld_account` |
| `view_test.go` | 2 | Asset View | MongoDB `view` |
| `identity_test.go` | 4 | Identity | Mock SRECore + ServiceTree |
| `workorder_test.go` | 3 | Work Order | MongoDB `work_order_records` |
| `certificate_test.go` | 5 | Certificate | MongoDB `certificate` |
| `subdomain_test.go` | 6 | Subdomain | MongoDB `resource` + Mock KLB/KDNS |
| `billing_test.go` | 2 | Billing | MongoDB `billing_instance` + `permission` |
| `cloud_metadata_test.go` | 4 | Cloud Metadata | MongoDB `meta_vpc/subnet` + Mock CloudKit |
| `cloud_order_test.go` | 3 | Cloud Order | MongoDB `cloud_order` |
| `restemplate_test.go` | 3 | Resource Template | MongoDB `res_template` |

## Response Envelope

All endpoints return:
```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
}
```

Error codes (`pkg/utils/response.go`):
```
CodeSuccess         = 0
CodeInvalidParams   = 1001  → HTTP 400 (ParamError, BadRequestError)
CodeUnauthorized    = 1002  → HTTP 401
CodeForbidden       = 1003  → HTTP 403
CodeNotFound        = 1004  → HTTP 404
CodeConflict        = 1005  → HTTP 409
CodeServiceError    = 1006  → HTTP 500
CodeValidationError = 1007  → HTTP 400 (ValidationError)
CodeDatabaseError   = 1008  → HTTP 500
```

## Mock Gateway Patterns

### gRPC Mock (example: KCMDB Tag)

```go
// internal/shared/gateway/grpc/kcmdb_tag_mock.go
type MockKcmdbTagGateway struct{}

func NewMockKcmdbTagGateway(_ *config.Config) *MockKcmdbTagGateway {
    log.Info("Using MockKcmdbTagGateway for KCMDB Tag")
    return &MockKcmdbTagGateway{}
}

func (m *MockKcmdbTagGateway) QueryDefinitions(ctx context.Context, req *dto.QueryTagDefinitionsReq) (*dto.QueryTagDefinitionsResp, error) {
    return &dto.QueryTagDefinitionsResp{Tags: []*dto.TagDefinition{}, Total: 0}, nil
}
```

### HTTP Mock (example: KLB with Full interface)

When multiple use-case interfaces share a gateway, define a "Full" interface:

```go
// internal/shared/gateway/http/klb_mock.go
type KlbGatewayFull interface {
    subdomainUsecase.KlbGateway
    certificateUsecase.KlbGateway
}

type MockKlbGateway struct{}

func NewMockKlbGateway(_ *config.Config) KlbGatewayFull {
    log.Info("Using MockKlbGateway for KLB")
    return &MockKlbGateway{}
}
```

### Provider Function Pattern

```go
// internal/di/providers.go
func provideKlbGateway(cfg *config.Config) (httpGateway.KlbGatewayFull, error) {
    if cfg.Integrations.KLB.UseMock {
        return httpGateway.NewMockKlbGateway(cfg), nil
    }
    return httpGateway.NewKlbGateway(cfg)
}
```

## Common Patterns

### CRUD Lifecycle Test

```go
func TestMyDomainCRUD(t *testing.T) {
    client := getMongoClient(t)
    cleanCollection(t, client, myCollection)
    t.Cleanup(func() { cleanCollection(t, client, myCollection) })

    e := newExpect(t)

    // Create
    createObj := e.POST("/api/v1/my_domain").
        WithJSON(map[string]any{...}).
        Expect().Status(http.StatusOK).JSON().Object()
    createObj.HasValue("code", utils.CodeSuccess)
    id := createObj.Value("data").String().Raw()

    // List — verify created
    listObj := e.POST("/api/v1/my_domain/list").
        WithJSON(map[string]any{"limit": 20, "page_index": 1}).
        Expect().Status(http.StatusOK).JSON().Object()
    listObj.Value("data").Object().HasValue("total", 1)

    // Get detail
    e.GET("/api/v1/my_domain/detail").
        WithQuery("id", id).
        Expect().Status(http.StatusOK).JSON().Object().
        HasValue("code", utils.CodeSuccess)

    // Delete
    e.DELETE("/api/v1/my_domain").
        WithJSON(map[string]any{"id": id}).
        Expect().Status(http.StatusOK).JSON().Object().
        HasValue("code", utils.CodeSuccess)

    // Verify deleted
    afterObj := e.POST("/api/v1/my_domain/list").
        WithJSON(map[string]any{"limit": 20, "page_index": 1}).
        Expect().Status(http.StatusOK).JSON().Object()
    afterObj.Value("data").Object().HasValue("total", 0)
}
```

### Seeding with Service Tree Visibility

```go
// For AssetServiceTreeAuth — seed resource with mock_dev as manager
seedDocuments(t, client, "resource", []interface{}{
    bson.M{
        "tag": bson.M{
            "res_type":   "domain_v2",
            "archived":   false,
            "manager_id": "mock_dev",  // matches mock user
        },
        "data": bson.M{...},
    },
})

// For BillingServiceTreeAuth — seed admin permission
seedDocuments(t, client, "permission", []interface{}{
    bson.M{
        "user_id": "mock_dev",
        "role":    "admin",
        ...
    },
})
```

### Filter Test

```go
func TestMyEndpointFilter(t *testing.T) {
    client := getMongoClient(t)
    cleanCollection(t, client, myCollection)
    t.Cleanup(func() { cleanCollection(t, client, myCollection) })

    // Seed 2 docs with different types
    seedDocuments(t, client, myCollection, []interface{}{
        bson.M{"type": "a", ...},
        bson.M{"type": "b", ...},
    })

    e := newExpect(t)

    // Filter for type "a" only
    obj := e.GET("/api/v1/my_endpoint").
        WithQuery("type", "a").
        WithQuery("limit", 20).
        WithQuery("page_index", 1).
        Expect().Status(http.StatusOK).JSON().Object()

    obj.Value("data").Object().HasValue("total_count", 1)
    obj.Value("data").Object().Value("records").Array().
        Value(0).Object().HasValue("type", "a")
}
```
