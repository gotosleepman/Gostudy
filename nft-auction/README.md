# NFT Auction Market (Hardhat)

本项目实现了一个可升级的 NFT 拍卖市场，支持：

- ERC721 NFT 上架拍卖
- 使用 ETH 或 ERC20 出价
- 使用 Chainlink 喂价将出价换算为 USD 统一比较
- 使用 UUPS 代理模式升级拍卖合约

## 项目结构

- `contracts/AuctionUpgradeable.sol`：主拍卖合约（UUPS）
- `contracts/MyNFT.sol`：示例 NFT 合约（ERC721）
- `contracts/MyToken.sol`：示例 ERC20 合约
- `contracts/mocks/MockV3Aggregator.sol`：本地测试喂价 Mock
- `test/auction.js`：功能测试
- `scripts/deploy.js`：Sepolia 部署脚本
- `scripts/check-env.js`：环境变量检查脚本

## 核心流程

1. 卖家调用 `createAuction(nftAddress, tokenId, endTime)` 创建拍卖并托管 NFT。
2. 买家可用以下方式出价：
   - `bidWithETH(auctionId)`
   - `bidWithERC20(auctionId, token, amount)`
3. 合约将 ETH/ERC20 出价换算成 USD，比较 `highestBidUsd`。
4. 到期后调用 `endAuction(auctionId)`，NFT 转给最高出价者，资金转给卖家。

## Chainlink 价格换算

- ETH/USD：`ethUsdFeed`
- ERC20/USD：`tokenUsdFeeds[token]`
- 管理函数：`setTokenPriceFeed(token, feed)`（owner）

## 升级方案（UUPS）

`AuctionUpgradeable` 使用 UUPS 代理，升级权限由 owner 控制（`_authorizeUpgrade`）。

## 升级验证说明（V1 -> V2）

项目已新增：

- `contracts/AuctionUpgradeableV2.sol`
- `test/upgrade.js`

验证命令：

```bash
npx hardhat test test/upgrade.js
```

测试验证点：

1. 使用 `deployProxy` 部署 V1（`AuctionUpgradeable`）。
2. 升级前写入业务状态（创建拍卖 + 出价）。
3. 使用 `upgradeProxy` 升级到 V2（`AuctionUpgradeableV2`）。
4. 断言代理地址不变（Proxy 地址升级前后一致）。
5. 断言 V1 存储状态保留（拍卖数据与最高出价人不丢失）。
6. 断言 V2 新功能可用（`version()` 返回 `v2`，`setPlatformFeeBps()` 可设置）。
7. 断言权限正确（非 owner 调用 `setPlatformFeeBps` 失败）。

## 本地运行与测试

```bash
npm install
npx hardhat compile
npx hardhat test
```

## 环境变量

在项目根目录创建 `.env`：

```env
SEPOLIA_RPC_URL=https://sepolia.infura.io/v3/<YOUR_KEY>
SEPOLIA_PRIVATE_KEY=0x<YOUR_PRIVATE_KEY>

# 可选：覆盖默认 ETH/USD 喂价地址
ETH_USD_FEED=0x...

# 可选：部署后自动配置一个 ERC20 喂价
BID_TOKEN_ADDRESS=0x...
TOKEN_USD_FEED=0x...
```

校验：

```bash
npm run env:check:sepolia
```

## Sepolia 部署

```bash
npm run deploy:sepolia
```

部署后请将以下地址补充到本文档：

- Auction Proxy:
- Auction Implementation:
- MyNFT:
- MyToken:
- ETH/USD Feed:
- ERC20/USD Feed:

## 测试覆盖率与结果

```bash
npm run coverage
npm run report:test
```


报告输出：

- `reports/test-report.md`
- `reports/test-output.log`

最近一次覆盖率结果（本地执行）：

| 指标 | 覆盖率 |
| --- | --- |
| Statements | 94.67% |
| Branches | 61.36% |
| Functions | 91.67% |
| Lines | 95.05% |

按合约文件统计（节选）：

| 文件 | Stmts | Branch | Funcs | Lines | Uncovered Lines |
| --- | ---: | ---: | ---: | ---: | --- |
| `contracts/AuctionUpgradeable.sol` | 93.94% | 61.25% | 93.33% | 95.4% | `188,189,224,225` |
| `contracts/AuctionUpgradeableV2.sol` | 100% | 75% | 100% | 100% | - |
| `contracts/MyNFT.sol` | 100% | 50% | 100% | 100% | - |
| `contracts/MyToken.sol` | 100% | 50% | 100% | 100% | - |
| `contracts/mocks/MockV3Aggregator.sol` | 100% | 100% | 66.67% | 75% | `22` |

> 注：覆盖率会随测试用例和合约改动变化，提交前建议重新执行 `npm run coverage` 并更新该表。


