# pledge_copy

`pledge_copy` 是一个借贷项目拆解示例，包含：

- 一个简化版借贷池合约：`contracts/SimpleLendingPool.sol`
- 一个 Go 后端服务：`lending-backend`
- 一个定时同步任务入口：`lending-backend/cmd/lending_task`

整体目标是对齐原有 Pledge 项目的核心数据结构（如 `poolBaseInfo` / `poolDataInfo`），便于后端扫链、入库和接口查询。

## 目录结构

```text
pledge_copy/
├─ contracts/
│  └─ SimpleLendingPool.sol
└─ lending-backend/
   ├─ api/                # HTTP 接口、参数校验、响应结构
   ├─ cmd/lending_task/   # 定时任务入口
   ├─ config/             # 配置与 config.toml
   ├─ contract/bindings/  # 合约 ABI 等绑定文件
   ├─ db/                 # MySQL / Redis 初始化
   ├─ schedule/           # 扫链与同步任务
   ├─ main.go             # API 服务入口
   └─ go.mod
```

## 合约说明

`SimpleLendingPool.sol` 提供了简化借贷流程，主要能力包括：

- 创建池子 `createPoolInfo`
- 出借资产 `depositLend`
- 抵押借款 `depositBorrow`
- 还款 `repay`
- 提取出借资产 `withdrawLend`
- 提取抵押资产 `withdrawCollateral`

合约中的池子字段命名尽量和后端表结构保持一致，便于后端直接同步和查询。

## 后端能力

后端基于 Gin + Gorm + Redis，主要提供：

- 池子基础信息查询
- 池子数据统计查询
- Token 列表查询
- 池子搜索能力

默认 API 前缀：`/api/v1`（由 `config/config.toml` 中 `env.version` 控制）

可用路由：

- `GET /api/v1/poolBaseInfo`
- `GET /api/v1/poolDataInfo`
- `GET /api/v1/token`
- `POST /api/v1/pool/search`

## 环境要求

- Go `1.17`（与 `go.mod` 保持一致）
- MySQL（默认库名：`lending_copy`）
- Redis

## 配置

后端配置文件：`lending-backend/config/config.toml`

启动前至少确认以下配置：

1. `mysql`：地址、账号、密码、数据库名
2. `redis`：地址、端口、DB
3. `test_net` / `main_net`：链节点地址、`lending_pool_addr`
4. `env.port`：服务端口（默认 `8081`）

## 启动方式

### 1) 启动 API 服务

在 `lending-backend` 目录执行：

```bash
go mod tidy
go run main.go
```

服务默认监听：`http://127.0.0.1:8081`

### 2) 启动定时任务（可选）

在 `lending-backend` 目录执行：

```bash
go run ./cmd/lending_task/main.go
```

定时任务会：

- 启动时刷新 Redis（当前 DB）
- 立即执行一次池子信息同步与余额监控
- 后续按固定周期继续执行

## 开发提示

- 本项目当前为拆解与学习用途，适合本地联调和流程理解
- 生产使用前建议补充：权限控制、异常恢复、测试用例、部署脚本与监控告警
