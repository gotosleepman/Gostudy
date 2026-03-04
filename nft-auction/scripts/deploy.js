const { ethers, upgrades } = require("hardhat");
require("dotenv").config();

async function main() {
  const required = ["SEPOLIA_RPC_URL", "SEPOLIA_PRIVATE_KEY"];
  const missing = required.filter((name) => !process.env[name] || process.env[name].trim() === "");
  if (network.name === "sepolia" && missing.length > 0) {
    throw new Error(
      `Missing .env variables for sepolia: ${missing.join(", ")}. Please update your .env before deploy.`,
    );
  }

  const ethUsdFeed =
    process.env.ETH_USD_FEED || "0x694AA1769357215DE4FAC081bf1f309aDC325306";

  const Auction = await ethers.getContractFactory("AuctionUpgradeable");

  const auction = await upgrades.deployProxy(
      Auction,
      [ethUsdFeed],
      { initializer: "initialize", kind: "uups" }
  );

  await auction.waitForDeployment();

  console.log("Auction proxy deployed to:", await auction.getAddress());
  console.log("ETH/USD feed:", ethUsdFeed);

  const tokenAddress = process.env.BID_TOKEN_ADDRESS;
  const tokenUsdFeed = process.env.TOKEN_USD_FEED;
  if (tokenAddress && tokenUsdFeed) {
    const tx = await auction.setTokenPriceFeed(tokenAddress, tokenUsdFeed);
    await tx.wait();
    console.log("Configured token feed:", tokenAddress, "=>", tokenUsdFeed);
  } else {
    console.log(
      "Skip token feed config: set BID_TOKEN_ADDRESS and TOKEN_USD_FEED to enable."
    );
  }
}

main().catch((error) => {
  console.error(error);
  process.exit(1);
});