import { expect } from "chai";
import { network } from "hardhat";

describe("SimpleLendingPool", function () {
  async function setupPoolFixture() {
    const { ethers } = await network.connect();
    const [owner, lender, borrower, other] = await ethers.getSigners();
    const { pool, proxy } = await deployProxy(ethers);

    const lendToken = await ethers.deployContract("MockERC20", ["Lend", "LND"]);
    const borrowToken = await ethers.deployContract("MockERC20", ["Borrow", "BRW"]);
    await lendToken.waitForDeployment();
    await borrowToken.waitForDeployment();

    const mintAmount = 1_000_000n;
    await lendToken.mint(lender.address, mintAmount);
    await lendToken.mint(owner.address, mintAmount);
    await borrowToken.mint(borrower.address, mintAmount);

    await pool.createPoolInfo(
      1n,
      1000n,
      500n,
      100_000n,
      50_000_000n,
      await lendToken.getAddress(),
      await borrowToken.getAddress(),
      owner.address,
      owner.address,
      0n
    );

    return { ethers, owner, lender, borrower, other, pool, proxy, lendToken, borrowToken };
  }

  async function deployProxy(ethers) {
    const implementation = await ethers.deployContract("SimpleLendingPool");
    await implementation.waitForDeployment();
    const initData = implementation.interface.encodeFunctionData("initialize", [
      11n,
      22n
    ]);
    const proxy = await ethers.deployContract("SimpleLendingPoolProxy", [
      await implementation.getAddress(),
      initData
    ]);
    await proxy.waitForDeployment();
    const pool = await ethers.getContractAt(
      "SimpleLendingPool",
      await proxy.getAddress()
    );
    return { implementation, proxy, pool };
  }

  it("通过代理初始化后 owner 与费率正确", async function () {
    const { ethers } = await network.connect();
    const [deployer] = await ethers.getSigners();
    const { pool } = await deployProxy(ethers);

    expect(await pool.owner()).to.equal(deployer.address);
    expect(await pool.lendFee()).to.equal(11n);
    expect(await pool.borrowFee()).to.equal(22n);
  });

  it("只有 owner 可以设置手续费", async function () {
    const { ethers } = await network.connect();
    const [, other] = await ethers.getSigners();
    const { pool } = await deployProxy(ethers);

    await expect(pool.connect(other).setFee(100, 200))
      .to.be.revertedWithCustomError(pool, "OwnableUnauthorizedAccount")
      .withArgs(other.address);
  });

  it("UUPS 可升级到 V2 并调用新方法", async function () {
    const { ethers } = await network.connect();
    const { proxy, pool } = await deployProxy(ethers);

    const implementationV2 = await ethers.deployContract("SimpleLendingPoolV2");
    await implementationV2.waitForDeployment();

    await pool.upgradeToAndCall(await implementationV2.getAddress(), "0x");
    const poolV2 = await ethers.getContractAt(
      "SimpleLendingPoolV2",
      await proxy.getAddress()
    );

    await poolV2.setProtocolNote("upgraded");
    expect(await poolV2.version()).to.equal("v2");
    expect(await poolV2.protocolNote()).to.equal("upgraded");
  });

  it("只有 owner 可以创建池子", async function () {
    const { owner, other, pool, lendToken, borrowToken } = await setupPoolFixture();

    await expect(
      pool.connect(other).createPoolInfo(
        1n,
        1000n,
        500n,
        100_000n,
        50_000_000n,
        await lendToken.getAddress(),
        await borrowToken.getAddress(),
        owner.address,
        owner.address,
        0n
      )
    )
      .to.be.revertedWithCustomError(pool, "OwnableUnauthorizedAccount")
      .withArgs(other.address);
  });

  it("depositLend 后池状态从 MATCH 进入 EXECUTION", async function () {
    const { lender, pool, lendToken } = await setupPoolFixture();

    await lendToken.connect(lender).approve(await pool.getAddress(), 20_000n);
    await pool.connect(lender).depositLend(0n, 20_000n);

    expect(await pool.poolLength()).to.equal(1n);
    expect(await pool.getPoolState(0n)).to.equal(1n);
    expect(await pool.supplied(0n, lender.address)).to.equal(20_000n);
  });

  it("borrower 可按抵押率借款并记录头寸", async function () {
    const { lender, borrower, pool, lendToken, borrowToken } = await setupPoolFixture();

    await lendToken.connect(lender).approve(await pool.getAddress(), 40_000n);
    await pool.connect(lender).depositLend(0n, 40_000n);

    await borrowToken.connect(borrower).approve(await pool.getAddress(), 20_000n);
    await pool.connect(borrower).depositBorrow(0n, 20_000n, 10_000n);

    expect(await pool.collateral(0n, borrower.address)).to.equal(20_000n);
    expect(await pool.borrowed(0n, borrower.address)).to.equal(10_000n);
  });

  it("超抵押率借款会被拒绝", async function () {
    const { lender, borrower, pool, lendToken, borrowToken } = await setupPoolFixture();

    await lendToken.connect(lender).approve(await pool.getAddress(), 50_000n);
    await pool.connect(lender).depositLend(0n, 50_000n);

    await borrowToken.connect(borrower).approve(await pool.getAddress(), 20_000n);
    await expect(pool.connect(borrower).depositBorrow(0n, 20_000n, 12_000n)).to.be.revertedWith(
      "collateral"
    );
  });

  it("还款后可在健康因子允许下提取抵押物", async function () {
    const { lender, borrower, pool, lendToken, borrowToken } = await setupPoolFixture();

    await lendToken.connect(lender).approve(await pool.getAddress(), 30_000n);
    await pool.connect(lender).depositLend(0n, 30_000n);

    await borrowToken.connect(borrower).approve(await pool.getAddress(), 20_000n);
    await pool.connect(borrower).depositBorrow(0n, 20_000n, 10_000n);

    await lendToken.connect(borrower).approve(await pool.getAddress(), 4_000n);
    await pool.connect(borrower).repay(0n, 4_000n);

    await pool.connect(borrower).withdrawCollateral(0n, 2_000n);
    expect(await pool.borrowed(0n, borrower.address)).to.equal(6_000n);
    expect(await pool.collateral(0n, borrower.address)).to.equal(18_000n);
  });

  it("非 owner 不能升级实现合约", async function () {
    const { ethers } = await network.connect();
    const [, other] = await ethers.getSigners();
    const { pool } = await deployProxy(ethers);
    const implementationV2 = await ethers.deployContract("SimpleLendingPoolV2");
    await implementationV2.waitForDeployment();

    await expect(
      pool.connect(other).upgradeToAndCall(await implementationV2.getAddress(), "0x")
    )
      .to.be.revertedWithCustomError(pool, "OwnableUnauthorizedAccount")
      .withArgs(other.address);
  });
});
