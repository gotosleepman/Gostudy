import { buildModule } from "@nomicfoundation/hardhat-ignition/modules";

export default buildModule("SimpleLendingPoolProxyModule", (m) => {
  const owner = m.getAccount(0);
  const initialLendFee = m.getParameter("initialLendFee", 0n);
  const initialBorrowFee = m.getParameter("initialBorrowFee", 0n);

  const implementation = m.contract("SimpleLendingPool");
  const initData = m.encodeFunctionCall(implementation, "initialize", [
    initialLendFee,
    initialBorrowFee
  ]);

  const proxy = m.contract("SimpleLendingPoolProxy", [implementation, initData], {
    from: owner
  });

  const pool = m.contractAt("SimpleLendingPool", proxy);
  return { owner, implementation, proxy, pool };
});
