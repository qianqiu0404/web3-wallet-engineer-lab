---
type: project-doc
topic: web3-wallet-engineer-lab
status: active
updated: 2026-07-03
---

# 第四阶段：测试与验收

## 已有测试

```bash
go test ./...
```

覆盖点：

- 黑名单提现被风控拒绝
- 合法提现审核通过后分配 nonce 并生成链上交易
- 充值必须打入用户已分配地址
- HTTP happy path：创建用户、分配地址、模拟充值、提现、审核

## 接口测试说明

启动服务：

```bash
go run ./cmd/api
```

按 [examples/requests.http](../examples/requests.http) 顺序执行。

建议验收：

- `POST /api/users` 返回用户 ID
- `POST /api/users/{id}/addresses` 返回链地址
- `POST /api/deposits/simulate` 返回 `CREDITED`
- `POST /api/withdrawals` 返回 `PENDING_REVIEW`
- `POST /api/admin/withdrawals/{id}/approve` 返回 `BROADCASTED`，包含 `nonce` 和 `chain_tx_id`
- `POST /api/chain/tx/{id}/confirm` 将链上交易置为 `CONFIRMED`
- `GET /api/admin/audit-logs` 能看到完整审计轨迹
- `GET /metrics` 输出 Prometheus 指标

## 异常场景

- 未分配地址充值：返回 400
- 黑名单地址提现：状态为 `RISK_REJECTED`
- 非正金额提现：状态为 `RISK_REJECTED`
- 超大额提现：状态为 `RISK_REJECTED`
- 重复审核已广播提现：返回 400
- 不存在的链上交易确认：返回 404

## TODO

- 引入 SQLite/PostgreSQL repository，实现真实持久化
- 引入 Redis，实现 nonce 锁、幂等 key、风控缓存
- 增加账户余额和冻结余额
- 增加提现幂等号 `client_request_id`
- 增加链适配器接口：ETH/BTC/TRON
- 增加异步 worker：充值扫描、提现广播、确认同步、归集执行
- 增加 OpenAPI 文档
- 增加 Docker Compose
- 增加 structured logging 和 trace id
- 增加更完整的 Prometheus histogram/gauge/counter

## 生产级升级路线

### 1. 数据库层

将 `internal/store.MemoryStore` 替换为 repository interface + SQL 实现。关键表包括 users、addresses、deposits、withdrawals、chain_transactions、blacklist_entries、wallets、collection_tasks、audit_logs、nonce_allocations。

### 2. 资产账本

引入双入口账本：充值入账增加可用余额，提现申请冻结余额，审核拒绝解冻，提现确认扣减冻结余额。所有资产变更必须有 ledger entry 和幂等键。

### 3. 链服务

抽象 `ChainAdapter`：

```go
type ChainAdapter interface {
    BuildTransfer(...)
    Broadcast(...)
    GetTxStatus(...)
    CurrentNonce(...)
}
```

### 4. 签名安全

将签名服务从业务服务隔离，生产中使用 KMS、HSM、MPC 或 TSS。业务系统只提交签名请求，不直接接触私钥。

### 5. 风控体系

规则引擎从硬编码升级为配置化策略，支持名单、额度、频率、用户等级、设备、IP、KYT、人工复核队列和告警。

### 6. 可观测性

指标、日志、trace 三件套：

- 指标看趋势和告警
- 日志追单笔请求
- trace 串起 API、DB、队列、签名、链节点
