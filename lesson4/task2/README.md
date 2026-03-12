# 任务 2：合约代码生成、部署与调用（Sepolia）

本任务在 `task2` 中完成以下目标：

- 编写 Solidity 计数器合约
- 编译生成 ABI 与字节码
- 使用 `abigen` 生成 Go 绑定代码
- 使用 Go 代码部署合约并调用 `increment()`

## 1. 项目结构

- `contracts/Counter.sol`：计数器合约
- `scripts/generate-bindings.ps1`：编译合约 + 生成 Go 绑定
- `bindings/counter/counter.go`：`abigen` 生成的绑定代码
- `main.go`：连接 Sepolia，支持部署与调用

## 2. 一键生成 ABI/字节码/Go 绑定

在 `task2` 目录执行：

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\generate-bindings.ps1
```

脚本会自动：

1. 安装 `solc`（使用 `solcjs`）
2. 编译 `Counter.sol` 生成：
   - `build/contracts_Counter_sol_Counter.abi`
   - `build/contracts_Counter_sol_Counter.bin`
3. 安装 `abigen`（若未安装）
4. 生成 `bindings/counter/counter.go`

## 3. 部署 Counter 合约

先准备环境变量（PowerShell）：

```powershell
$env:RPC_URL="https://sepolia.infura.io/v3/<INFURA_KEY>"
$env:PRIVATE_KEY="<YOUR_TEST_PRIVATE_KEY>"
```

执行：

```powershell
go mod tidy
go run . deploy
```

成功后会输出：

- 部署交易哈希
- 新部署的 `合约地址`（即 `<DEPLOYED_COUNTER_ADDRESS>`）
- 区块高度与 Etherscan 链接

## 4. 调用 increment

把上一步拿到的合约地址写入环境变量：

```powershell
$env:CONTRACT_ADDRESS="<DEPLOYED_COUNTER_ADDRESS>"
```

执行：

```powershell
go run . call
```

程序会读取调用前计数值，发送 `increment()` 交易，等待上链后输出调用后计数值。

## 5. 安全说明

- 仅使用测试网私钥，不要用于主网
- 不要将私钥写入代码或提交到仓库
- 推荐始终通过环境变量传入敏感信息
