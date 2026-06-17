# /learn - 提取可复用模式

分析当前会话，提取值得保存的经验作为 skill。

## 触发时机

在会话中解决了非平凡问题时运行 `/learn`。

## 提取类型

### 1. 错误解决模式

- 遇到了什么错误？
- 根本原因是什么？
- 如何修复的？
- 是否可用于类似错误？

### 2. 调试技巧

- 非显而易见的调试步骤
- 有效的工具组合
- 诊断模式

### 3. Go 特定模式

- GORM 查询模式（复杂查询、性能优化）
- Gin 中间件模式（认证、授权、错误处理）
- DDD 层间交互模式
- Wire 依赖注入模式
- 并发处理模式（goroutine、channel、sync）
- 错误处理模式（AppError 使用）

### 4. 项目特定模式

- kfinops 代码库约定
- 架构决策
- 集成模式（Gateway、Repository）
- 测试模式（表驱动测试、Mock）

### 5. 变通方案

- 库的特殊用法
- API 限制的解决方案
- 版本特定的修复

## 输出格式

在 `.cursor/skills/learned/` 创建 skill 文件：

```markdown
# [描述性模式名称]

**提取日期:** YYYY-MM-DD
**适用场景:** [简述何时应用此模式]

## 问题

[此模式解决什么问题 - 要具体]

## 解决方案

[模式/技术/变通方案]

## 示例

[代码示例]

```go
// 示例代码
```

## 何时使用

[触发条件 - 什么情况下应该激活此 skill]
```

## 示例 Skill

### GORM 查询优化

```markdown
# GORM 批量查询优化

**提取日期:** 2026-01-28
**适用场景:** 需要批量查询关联数据时

## 问题

循环中逐个查询关联数据导致 N+1 问题，性能极差。

## 解决方案

使用 Preload 或批量 IN 查询替代循环查询。

## 示例

```go
// ❌ 错误 - N+1 查询
for _, cert := range certs {
    db.Where("cert_id = ?", cert.ID).Find(&cert.Deployments)
}

// ✅ 正确 - Preload
db.Preload("Deployments").Find(&certs)

// ✅ 正确 - 批量 IN 查询
certIDs := make([]string, len(certs))
for i, c := range certs {
    certIDs[i] = c.ID
}
db.Where("cert_id IN ?", certIDs).Find(&deployments)
```

## 何时使用

- 查询列表数据时需要关联信息
- 发现 SQL 日志中有大量重复查询
- 接口响应时间过长
```

### Context 取消处理

```markdown
# Context 取消优雅处理

**提取日期:** 2026-01-28
**适用场景:** 长时间运行的操作需要支持取消

## 问题

长时间运行的操作（如批量同步）不响应 Context 取消，导致无法优雅停止。

## 解决方案

在循环中检查 ctx.Done()，及时响应取消请求。

## 示例

```go
func (m *Manager) SyncAll(ctx context.Context) error {
    items, err := m.repo.List(ctx)
    if err != nil {
        return err
    }

    for _, item := range items {
        // 检查 context 是否已取消
        select {
        case <-ctx.Done():
            log.Warn("Sync cancelled",
                zap.Int("processed", processed),
                zap.Int("total", len(items)),
            )
            return ctx.Err()
        default:
        }

        if err := m.syncOne(ctx, item); err != nil {
            log.Error("Sync failed", zap.Error(err))
            continue
        }
        processed++
    }
    return nil
}
```

## 何时使用

- 实现批量操作
- 长时间运行的任务
- 需要支持优雅关闭
```

## 执行流程

1. 审查会话中的可提取模式
2. 识别最有价值/可复用的洞察
3. 起草 skill 文件
4. 向用户确认后保存
5. 保存到 `.cursor/skills/learned/`

## Hook 自动注入

项目已配置 `.cursor/hooks.json`，在每次会话开始时：

1. `session-start.sh` 自动扫描 `.cursor/skills/learned/` 目录
2. 将所有已学习的 skill 列表注入到会话上下文
3. 无需手动更新索引文件

## 文件命名建议

使用 kebab-case 命名 skill 文件：

```
.cursor/skills/learned/
├── gorm-batch-query.md
├── context-cancellation.md
├── gin-middleware-chain.md
└── wire-provider-pattern.md
```

## 注意事项

- 不要提取平凡的修复（拼写错误、简单语法错误）
- 不要提取一次性问题（特定 API 故障等）
- 聚焦于能在未来会话中节省时间的模式
- 保持 skill 专注 - 每个 skill 一个模式
- 使用中文描述，代码注释可中英混用
