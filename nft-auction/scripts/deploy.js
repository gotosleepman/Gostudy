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

  const priceFeed = "0x694AA1769357215DE4FAC081bf1f309aDC325306";

  const Auction = await ethers.getContractFactory("AuctionUpgradeable");

  const auction = await upgrades.deployProxy(
      Auction,
      [priceFeed],
      { initializer: "initialize", kind: "uups" }
  );

  await auction.waitForDeployment();

  console.log("Auction deployed to:", await auction.getAddress());
}

main().catch((error) => {
  console.error(error);
  process.exit(1);
});