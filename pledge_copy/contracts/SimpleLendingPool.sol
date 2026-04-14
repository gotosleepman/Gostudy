// SPDX-License-Identifier: MIT
pragma solidity ^0.8.23;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";

/// @notice 最小 ERC20 接口（与 Pledge 池子字段对齐，便于后端用同一套表结构同步）
interface IERC20 {
    function transfer(address to, uint256 amount) external returns (bool);
    function transferFrom(address from, address to, uint256 amount) external returns (bool);
    function balanceOf(address account) external view returns (uint256);
}

/// @title SimpleLendingPool
/// @dev 仿 ProjectBreakdown-Pledge 的池结构：多池、借出币/抵押币、费率、poolBaseInfo/poolDataInfo 视图对齐后端扫链字段
contract SimpleLendingPool is Initializable, OwnableUpgradeable, UUPSUpgradeable {
    uint256 public lendFee;
    uint256 public borrowFee;

    enum PoolState {
        MATCH,
        EXECUTION,
        FINISH,
        LIQUIDATION,
        UNDONE
    }

    struct PoolBaseInfo {
        uint256 settleTime;
        uint256 endTime;
        uint256 interestRate;
        uint256 maxSupply;
        uint256 lendSupply;
        uint256 borrowSupply;
        uint256 martgageRate;
        address lendToken;
        address borrowToken;
        PoolState state;
        address spCoin;
        address jpCoin;
        uint256 autoLiquidateThreshold;
    }

    struct PoolDataInfo {
        uint256 settleAmountLend;
        uint256 settleAmountBorrow;
        uint256 finishAmountLend;
        uint256 finishAmountBorrow;
        uint256 liquidationAmounLend;
        uint256 liquidationAmounBorrow;
    }

    PoolBaseInfo[] public poolBaseInfo;
    PoolDataInfo[] public poolDataInfo;

    mapping(uint256 => mapping(address => uint256)) public supplied;
    mapping(uint256 => mapping(address => uint256)) public collateral;
    mapping(uint256 => mapping(address => uint256)) public borrowed;

    event DepositLend(address indexed user, uint256 indexed pid, uint256 amount);
    event DepositBorrow(address indexed user, uint256 indexed pid, uint256 collateralAmt, uint256 borrowAmt);
    event Repay(address indexed user, uint256 indexed pid, uint256 amount);
    event WithdrawLend(address indexed user, uint256 indexed pid, uint256 amount);
    event WithdrawCollateral(address indexed user, uint256 indexed pid, uint256 amount);
    event SetFee(uint256 lendFee, uint256 borrowFee);
    event StateChange(uint256 indexed pid, uint256 beforeState, uint256 afterState);

    constructor() {
        _disableInitializers();
    }

    function initialize(uint256 _lendFee, uint256 _borrowFee) public initializer {
        __Ownable_init(msg.sender);
        lendFee = _lendFee;
        borrowFee = _borrowFee;
        emit SetFee(_lendFee, _borrowFee);
    }

    function _authorizeUpgrade(address) internal override onlyOwner {}

    function setFee(uint256 _lendFee, uint256 _borrowFee) external onlyOwner {
        lendFee = _lendFee;
        borrowFee = _borrowFee;
        emit SetFee(_lendFee, _borrowFee);
    }

    /// @notice 创建池（字段命名与 Pledge `createPoolInfo` 一致，便于后端沿用 poolbases 表）
    function createPoolInfo(
        uint256 _settleTime,
        uint256 _endTime,
        uint256 _interestRate,
        uint256 _maxSupply,
        uint256 _martgageRate,
        address _lendToken,
        address _borrowToken,
        address _spToken,
        address _jpToken,
        uint256 _autoLiquidateThreshold
    ) external onlyOwner {
        poolBaseInfo.push(
            PoolBaseInfo({
                settleTime: _settleTime,
                endTime: _endTime,
                interestRate: _interestRate,
                maxSupply: _maxSupply,
                lendSupply: 0,
                borrowSupply: 0,
                martgageRate: _martgageRate,
                lendToken: _lendToken,
                borrowToken: _borrowToken,
                state: PoolState.MATCH,
                spCoin: _spToken,
                jpCoin: _jpToken,
                autoLiquidateThreshold: _autoLiquidateThreshold
            })
        );
        poolDataInfo.push(PoolDataInfo(0, 0, 0, 0, 0, 0));
    }

    function poolLength() external view returns (uint256) {
        return poolBaseInfo.length;
    }

    function getPoolState(uint256 pid) external view returns (uint256) {
        return uint256(poolBaseInfo[pid].state);
    }

    /// @notice 供应借出资产（贷款池存款）
    function depositLend(uint256 pid, uint256 amount) external {
        PoolBaseInfo storage p = poolBaseInfo[pid];
        require(uint256(p.state) <= uint256(PoolState.EXECUTION), "bad state");
        require(amount > 0, "amount");
        require(p.lendSupply + amount <= p.maxSupply, "max supply");
        IERC20(p.lendToken).transferFrom(msg.sender, address(this), amount);
        supplied[pid][msg.sender] += amount;
        p.lendSupply += amount;
        if (p.state == PoolState.MATCH) {
            emit StateChange(pid, uint256(p.state), uint256(PoolState.EXECUTION));
            p.state = PoolState.EXECUTION;
        }
        emit DepositLend(msg.sender, pid, amount);
    }

    /// @notice 存入抵押并借出 lendToken（演示用 1:1 计价，抵押率 martgageRate 以 1e8 为基数）
    function depositBorrow(uint256 pid, uint256 collateralAmt, uint256 borrowAmt) external {
        PoolBaseInfo storage p = poolBaseInfo[pid];
        require(p.state == PoolState.EXECUTION, "state");
        require(collateralAmt > 0 && borrowAmt > 0, "amount");
        IERC20(p.borrowToken).transferFrom(msg.sender, address(this), collateralAmt);
        collateral[pid][msg.sender] += collateralAmt;
        uint256 maxBorrow = (collateral[pid][msg.sender] * p.martgageRate) / 1e8;
        require(borrowed[pid][msg.sender] + borrowAmt <= maxBorrow, "collateral");
        require(p.borrowSupply + borrowAmt <= p.lendSupply, "liquidity");
        require(IERC20(p.lendToken).balanceOf(address(this)) >= borrowAmt, "bal");
        IERC20(p.lendToken).transfer(msg.sender, borrowAmt);
        borrowed[pid][msg.sender] += borrowAmt;
        p.borrowSupply += borrowAmt;
        emit DepositBorrow(msg.sender, pid, collateralAmt, borrowAmt);
    }

    function repay(uint256 pid, uint256 amount) external {
        PoolBaseInfo storage p = poolBaseInfo[pid];
        require(amount > 0 && amount <= borrowed[pid][msg.sender], "repay");
        IERC20(p.lendToken).transferFrom(msg.sender, address(this), amount);
        borrowed[pid][msg.sender] -= amount;
        p.borrowSupply -= amount;
        poolDataInfo[pid].finishAmountBorrow += amount;
        emit Repay(msg.sender, pid, amount);
    }

    function withdrawLend(uint256 pid, uint256 amount) external {
        PoolBaseInfo storage p = poolBaseInfo[pid];
        require(amount > 0 && amount <= supplied[pid][msg.sender], "supply");
        require(p.lendSupply - p.borrowSupply >= amount, "reserved");
        supplied[pid][msg.sender] -= amount;
        p.lendSupply -= amount;
        IERC20(p.lendToken).transfer(msg.sender, amount);
        poolDataInfo[pid].finishAmountLend += amount;
        emit WithdrawLend(msg.sender, pid, amount);
    }

    function withdrawCollateral(uint256 pid, uint256 amount) external {
        PoolBaseInfo storage p = poolBaseInfo[pid];
        require(amount > 0 && amount <= collateral[pid][msg.sender], "col");
        uint256 newCol = collateral[pid][msg.sender] - amount;
        require(borrowed[pid][msg.sender] * 1e8 <= newCol * p.martgageRate, "hf");
        collateral[pid][msg.sender] = newCol;
        IERC20(p.borrowToken).transfer(msg.sender, amount);
        emit WithdrawCollateral(msg.sender, pid, amount);
    }
}
