# GORM Preload 模式

**提取日期:** 2026-01-28
**适用场景:** 需要查询关联数据时避免 N+1 问题

## 问题

在循环中逐个查询关联数据会导致 N+1 查询问题，严重影响性能。

## 解决方案

使用 GORM 的 Preload 方法一次性加载关联数据。

## 示例

```go
// ❌ 错误 - N+1 查询
certs, _ := repo.List(ctx)
for _, cert := range certs {
    db.Where("cert_id = ?", cert.ID).Find(&cert.Deployments)
}

// ✅ 正确 - Preload 一次查询
db.Preload("Deployments").Find(&certs)

// ✅ 正确 - 条件 Preload
db.Preload("Deployments", "status = ?", "active").Find(&certs)

// ✅ 正确 - 嵌套 Preload
db.Preload("Deployments.LoadBalancer").Find(&certs)
```

## 何时使用

- 列表接口需要返回关联数据
- SQL 日志中发现大量重复查询
- 接口响应时间过长
