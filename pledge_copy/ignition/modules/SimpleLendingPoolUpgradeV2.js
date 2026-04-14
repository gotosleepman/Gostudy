import { buildModule } from "@nomicfoundation/hardhat-ignition/modules";
import SimpleLendingPoolProxyModule from "./SimpleLendingPoolProxy.js";

export default buildModule("SimpleLendingPoolUpgradeV2Module", (m) => {
  const { owner, proxy } = m.useModule(SimpleLendingPoolProxyModule);

  const implementationV2 = m.contract("SimpleLendingPoolV2");
  const poolV2 = m.contractAt("SimpleLendingPoolV2", proxy);

  m.call(poolV2, "upgradeToAndCall", [implementationV2, "0x"], {
    from: owner
  });

  return { proxy, implementationV2, poolV2 };
});
