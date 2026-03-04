const { expect } = require("chai");
const { ethers, upgrades } = require("hardhat");

describe("Auction UUPS Upgrade", function () {
  it("should upgrade from V1 to V2 and keep state", async function () {
    const [owner, bidder] = await ethers.getSigners();

    const Feed = await ethers.getContractFactory("MockV3Aggregator");
    const ethFeed = await Feed.deploy(8, 2000n * 10n ** 8n);
    await ethFeed.waitForDeployment();

    const NFT = await ethers.getContractFactory("MyNFT");
    const nft = await NFT.deploy();
    await nft.waitForDeployment();

    const AuctionV1 = await ethers.getContractFactory("AuctionUpgradeable");
    const auctionV1 = await upgrades.deployProxy(
      AuctionV1,
      [await ethFeed.getAddress()],
      { kind: "uups" },
    );
    await auctionV1.waitForDeployment();
    const proxyAddressBefore = await auctionV1.getAddress();

    const nowBlock = await ethers.provider.getBlock("latest");
    const endTime = BigInt(nowBlock.timestamp) + 3600n;

    await nft.mintNFT(await owner.getAddress());
    await nft.approve(proxyAddressBefore, 0);
    await auctionV1.createAuction(await nft.getAddress(), 0, endTime);
    await auctionV1.connect(bidder).bidWithETH(0, { value: ethers.parseEther("1") });

    const before = await auctionV1.auctions(0);
    expect(before.highestBidder).to.equal(await bidder.getAddress());
    expect(before.highestBid).to.equal(ethers.parseEther("1"));

    const AuctionV2 = await ethers.getContractFactory("AuctionUpgradeableV2");
    const auctionV2 = await upgrades.upgradeProxy(proxyAddressBefore, AuctionV2);
    await auctionV2.waitForDeployment();
    const proxyAddressAfter = await auctionV2.getAddress();

    expect(proxyAddressAfter).to.equal(proxyAddressBefore);

    const after = await auctionV2.auctions(0);
    expect(after.highestBidder).to.equal(await bidder.getAddress());
    expect(after.highestBid).to.equal(ethers.parseEther("1"));
    expect(await auctionV2.version()).to.equal("v2");

    await expect(auctionV2.connect(bidder).setPlatformFeeBps(250)).to.be.reverted;
    await auctionV2.setPlatformFeeBps(250);
    expect(await auctionV2.platformFeeBps()).to.equal(250n);
  });
});
