# NFT 拍卖市场 - 提交模板与自检清单

## 1) Sepolia 部署地址表（必填）

请将下面地址补全后提交：

| 项目 | 地址 |
| --- | --- |
| Auction Proxy | `0x...` |
| Auction Implementation | `0x...` |
| MyNFT | `0x...` |
| MyToken | `0x...` |
| ETH/USD Feed | `0x...` |
| ERC20/USD Feed（MyToken） | `0x...` |
| 部署交易哈希（Proxy） | `0x...` |
| 部署交易哈希（Implementation） | `0x...` |

> 说明：Proxy/Implementation 可从 `.openzeppelin/sepolia.json` 和部署日志提取。

## 2) 测试与覆盖率提交项

- 测试命令：`npm test`
- 覆盖率命令：`npm run coverage`
- 报告命令：`npm run report:test`

提交时请附：

- [ ] `reports/test-report.md`
- [ ] `reports/test-output.log`
- [ ] 覆盖率结果截图（或文本粘贴）
- [ ] 通过用例总数与失败数（例如：`10 passing, 0 failing`）

## 3) 功能验收自检（对应作业要求）

- [ ] 已实现 ERC721 NFT 铸造与转移（`MyNFT`）
- [ ] 已实现创建拍卖（`createAuction`）
- [ ] 已实现 ETH 出价（`bidWithETH`）
- [ ] 已实现 ERC20 出价（`bidWithERC20`）
- [ ] 已实现拍卖结束结算（`endAuction`）
- [ ] 已接入 Chainlink 喂价并进行 USD 比较
- [ ] 已实现 UUPS 升级模式（`_authorizeUpgrade` + Proxy 部署）
- [ ] 已有单元/集成测试覆盖主流程与异常路径
- [ ] README 已包含项目结构、运行、部署与说明

## 4) 提交包最终清单

- [ ] `contracts/` 合约源码完整
- [ ] `test/` 测试文件完整
- [ ] `scripts/` 部署与环境检查脚本完整
- [ ] `README.md` 文档完整
- [ ] 本文件地址表已填写

## 5) 可选加分（有时间再做）

- [ ] 增加更多事件断言测试
- [ ] 增加升级测试（V1 -> V2）
- [ ] 增加安全说明（重入、退款、权限、时间边界）
