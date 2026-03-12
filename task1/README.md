# 任务 1：Sepolia 区块链读写（Go）

本项目实现了两个功能：

- 查询指定区块信息（哈希、时间戳、交易数量等）
- 发送一笔 Sepolia 测试网 ETH 转账交易并输出交易哈希

## 1. 环境搭建

### 1.1 安装 Go

建议 Go 1.22 或更高版本。安装后验证：

```bash
go version
```

### 1.2 获取 Sepolia RPC（Infura）

1. 注册 [Infura](https://www.infura.io/)
2. 创建项目并启用 Ethereum
3. 选择 Sepolia 网络，拿到 `https://sepolia.infura.io/v3/<INFURA_KEY>`

### 1.3 安装依赖

在项目根目录执行：

```bash
go mod tidy
```

## 2. 查询区块

推荐先设置环境变量（PowerShell）：

```bash
$env:RPC_URL="https://sepolia.infura.io/v3/<INFURA_KEY>"
```

命令：

```bash
go run . query-block -block 6000000
```

也支持直接传参：

```bash
go run . query-block -rpc "https://sepolia.infura.io/v3/<INFURA_KEY>" -block 6000000
```

输出示例（节选）：

- 区块号
- 区块哈希
- 父区块哈希
- 时间戳（Unix + UTC 时间）
- 交易数量
- GasUsed / GasLimit

## 3. 发送交易

### 3.1 准备账户

- 准备一个 Sepolia 账户私钥（仅测试用途）
- 通过 Sepolia Faucet 给该地址领取测试 ETH

### 3.2 发送交易命令

推荐先设置环境变量（PowerShell）：

```bash
$env:RPC_URL="https://sepolia.infura.io/v3/<INFURA_KEY>"
$env:PRIVATE_KEY="<PRIVATE_KEY_HEX>"
```

然后执行（更安全，不会在命令里暴露私钥）：

```bash
go run . send-tx -to "<TO_ADDRESS>" -amount 0.0001
```

也支持直接传参（不推荐）：

```bash
go run . send-tx -rpc "https://sepolia.infura.io/v3/<INFURA_KEY>" -pk "<PRIVATE_KEY_HEX>" -to "<TO_ADDRESS>" -amount 0.0001
```

参数说明：

- `-rpc`: Sepolia RPC 地址（可用环境变量 `RPC_URL` 代替）
- `-pk`: 发送方私钥（可用环境变量 `PRIVATE_KEY` 代替）
- `-to`: 接收方地址
- `-amount`: ETH 数量（最多 18 位小数）
- `-chainid`: 可选，默认 `11155111`（Sepolia）

成功后会输出：

- 交易哈希
- Etherscan 链接（Sepolia）

## 4. 安全提醒

- 私钥仅用于测试，不要在主网复用
- 不要把私钥提交到代码仓库
- 建议使用环境变量 `PRIVATE_KEY`，避免命令历史泄露私钥
