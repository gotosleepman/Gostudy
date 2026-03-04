const { expect } = require("chai");
const { ethers, upgrades } = require("hardhat");

describe("Auction", function () {

  it("Should create auction", async function () {
    const [deployer] = await ethers.getSigners();

    const NFT = await ethers.getContractFactory("MyNFT");
    const nft = await NFT.deploy();
    await nft.waitForDeployment();

    const Auction = await ethers.getContractFactory("AuctionUpgradeable");
    const auction = await upgrades.deployProxy(
      Auction,
      ["0x694AA1769357215DE4FAC081bf1f309aDC325306"],
      { kind: "uups" },
    );
    await auction.waitForDeployment();

    const auctionAddress = await auction.getAddress();
    const nftAddress = await nft.getAddress();

    await nft.mintNFT(await deployer.getAddress());
    await nft.approve(auctionAddress, 0);

    await auction.createAuction(nftAddress, 0);

    const data = await auction.auctions(0);

    expect(data.tokenId).to.equal(0n);
  });
});