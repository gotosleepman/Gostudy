import { network } from "hardhat";

async function main() {
  const { ethers } = await network.connect();
  const implementation = await ethers.deployContract("SimpleLendingPool");
  await implementation.waitForDeployment();

  const initData = implementation.interface.encodeFunctionData("initialize", [
    0n,
    0n
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

  console.log("SimpleLendingPool implementation:", await implementation.getAddress());
  console.log("SimpleLendingPool proxy:", await proxy.getAddress());
  console.log("SimpleLendingPool owner:", await pool.owner());
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
