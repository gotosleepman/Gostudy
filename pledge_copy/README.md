# pledge_copy

`pledge_copy` 是一个借贷项目拆解示例，包含：

- 一个简化版借贷池合约：`contracts/SimpleLendingPool.sol`
- 一个 Go 后端服务：`lending-backend`
- 一个定时同步任务入口：`lending-backend/cmd/lending_task`


## 目录结构

```text
pledge_copy/
├─ contracts/
│  ├─ mocks/MockERC20.sol
│  ├─ SimpleLendingPool.sol
│  ├─ SimpleLendingPoolProxy.sol
│  └─ SimpleLendingPoolV2.sol
├─ ignition/modules/      # Ignition 模块化部署/升级
├─ scripts/               # 普通部署脚本
├─ test/                  # Hardhat 测试
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

## 合约开发（Hardhat）

项目根目录已接入 Hardhat，可直接编译、测试、部署 `contracts/` 下合约。
当前部署模式已改为 **UUPS 可升级代理**。

### 安装依赖

在项目根目录执行：

```bash
npm install
```

### 常用命令

```bash
# 编译合约
npm run compile

# 运行测试
npm run test

# 启动本地 Hardhat 节点
npm run node

# 本地网络部署（默认 hardhat 网络）
npm run deploy:local

# 部署到 BSC 测试网
npm run deploy:bscTestnet

# Ignition 模块化部署（UUPS Proxy）
npm run ignition:deploy:local
npm run ignition:deploy:bscTestnet

# Ignition 模块化升级到 V2
npm run ignition:upgrade:v2:local
npm run ignition:upgrade:v2:bscTestnet
```

### 环境变量

1. 复制 `.env.example` 为 `.env`
2. 按需填写 RPC 与私钥

> 注意：`.env` 已加入忽略列表，避免泄露私钥。

### Ignition 模块说明

- `ignition/modules/SimpleLendingPoolProxy.js`
  - 部署 `SimpleLendingPool` 实现合约
  - 编码 `initialize(uint256,uint256)` 初始化数据
  - 部署 `SimpleLendingPoolProxy`（基于 ERC1967）并绑定 `SimpleLendingPool` 接口
- `ignition/modules/SimpleLendingPoolUpgradeV2.js`
  - 部署新实现 `SimpleLendingPoolV2`
  - 通过代理执行 `upgradeToAndCall` 完成升级

示例（本地）：

```bash
# 先部署代理
npm run ignition:deploy:local

# 再升级到 V2
npm run ignition:upgrade:v2:local
```

### UUPS 升级实现说明

- `SimpleLendingPool` 已使用 `Initializable + OwnableUpgradeable + UUPSUpgradeable`
- 通过 `initialize` 完成初始化，构造函数中已 `_disableInitializers()`
- 通过 `_authorizeUpgrade` 限制只有 `owner` 可升级
- `SimpleLendingPoolV2` 提供升级演示新增状态与 `version()` 方法

### 测试说明

当前 `npm run test` 已覆盖 9 个核心场景，包括：

- 代理初始化参数与 owner 校验
- owner 权限（设置费率/创建池子/升级）校验
- UUPS 升级到 V2 成功并可调用新增方法
- `depositLend` 状态流转（`MATCH -> EXECUTION`）
- 抵押借款成功路径与超额借款失败路径
- 还款后提取抵押物成功路径

测试中使用 `contracts/mocks/MockERC20.sol` 作为借贷与抵押代币，便于本地快速验证业务流程。

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
