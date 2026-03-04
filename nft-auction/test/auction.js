const { expect } = require("chai");
const { ethers, upgrades } = require("hardhat");

describe("Auction", function () {
  async function createAuctionFixture() {
    const [seller, bidder1, bidder2] = await ethers.getSigners();

    const Feed = await ethers.getContractFactory("MockV3Aggregator");
    const ethFeed = await Feed.deploy(8, 2000n * 10n ** 8n); // 1 ETH = 2000 USD
    const tokenFeed = await Feed.deploy(8, 1n * 10n ** 8n); // 1 Token = 1 USD

    const NFT = await ethers.getContractFactory("MyNFT");
    const nft = await NFT.deploy();
    await nft.waitForDeployment();

    const Token = await ethers.getContractFactory("MyToken");
    const token = await Token.deploy();
    await token.waitForDeployment();

    const Auction = await ethers.getContractFactory("AuctionUpgradeable");
    const auction = await upgrades.deployProxy(
      Auction,
      [await ethFeed.getAddress()],
      { kind: "uups" },
    );
    await auction.waitForDeployment();

    await auction.setTokenPriceFeed(
      await token.getAddress(),
      await tokenFeed.getAddress(),
    );

    const nowBlock = await ethers.provider.getBlock("latest");
    const endTime = BigInt(nowBlock.timestamp) + 3600n;

    await nft.mintNFT(await seller.getAddress());
    await nft.approve(await auction.getAddress(), 0);
    await auction.createAuction(await nft.getAddress(), 0, endTime);

    return { seller, bidder1, bidder2, nft, token, auction, endTime };
  }

  it("should create auction with endTime", async function () {
    const { seller, nft, auction, endTime } = await createAuctionFixture();
    const data = await auction.auctions(0);
    expect(data.exists).to.equal(true);
    expect(data.seller).to.equal(await seller.getAddress());
    expect(data.nftAddress).to.equal(await nft.getAddress());
    expect(data.tokenId).to.equal(0n);
    expect(data.endTime).to.equal(endTime);
  });

  it("should accept higher ETH bid and reject lower ETH bid", async function () {
    const { bidder1, bidder2, auction } = await createAuctionFixture();

    await auction.connect(bidder1).bidWithETH(0, { value: ethers.parseEther("1") });

    await expect(
      auction.connect(bidder2).bidWithETH(0, { value: ethers.parseEther("0.5") }),
    ).to.be.revertedWith("Bid too low");

    const data = await auction.auctions(0);
    expect(data.highestBidder).to.equal(await bidder1.getAddress());
    expect(data.highestBid).to.equal(ethers.parseEther("1"));
  });

  it("should accept ERC20 bid and reject lower ERC20 bid", async function () {
    const { bidder1, bidder2, token, auction } = await createAuctionFixture();

    await token.mint(await bidder1.getAddress(), ethers.parseEther("2000"));
    await token.mint(await bidder2.getAddress(), ethers.parseEther("2000"));

    await token.connect(bidder1).approve(await auction.getAddress(), ethers.parseEther("1200"));
    await auction
      .connect(bidder1)
      .bidWithERC20(0, await token.getAddress(), ethers.parseEther("1200"));

    await token.connect(bidder2).approve(await auction.getAddress(), ethers.parseEther("1000"));
    await expect(
      auction
        .connect(bidder2)
        .bidWithERC20(0, await token.getAddress(), ethers.parseEther("1000")),
    ).to.be.revertedWith("Bid too low");
  });

  it("should compare ETH and ERC20 by USD value", async function () {
    const { bidder1, bidder2, token, auction } = await createAuctionFixture();

    // 1 ETH = 2000 USD
    await auction.connect(bidder1).bidWithETH(0, { value: ethers.parseEther("1") });

    await token.mint(await bidder2.getAddress(), ethers.parseEther("3000"));
    await token.connect(bidder2).approve(await auction.getAddress(), ethers.parseEther("1500"));
    await expect(
      auction
        .connect(bidder2)
        .bidWithERC20(0, await token.getAddress(), ethers.parseEther("1500")),
    ).to.be.revertedWith("Bid too low");

    await token.connect(bidder2).approve(await auction.getAddress(), ethers.parseEther("2500"));
    await auction
      .connect(bidder2)
      .bidWithERC20(0, await token.getAddress(), ethers.parseEther("2500"));

    const data = await auction.auctions(0);
    expect(data.highestBidder).to.equal(await bidder2.getAddress());
    expect(data.bidToken).to.equal(await token.getAddress());
  });

  it("should queue and withdraw refunds for ETH and ERC20", async function () {
    const { seller, bidder1, bidder2, nft, token, auction } = await createAuctionFixture();

    // ETH refund
    await auction.connect(bidder1).bidWithETH(0, { value: ethers.parseEther("1") });
    await auction.connect(bidder2).bidWithETH(0, { value: ethers.parseEther("2") });
    expect(await auction.pendingEthReturns(await bidder1.getAddress())).to.equal(
      ethers.parseEther("1"),
    );
    await auction.connect(bidder1).withdrawEthRefund();
    expect(await auction.pendingEthReturns(await bidder1.getAddress())).to.equal(0n);

    // ERC20 refund
    const nowBlock = await ethers.provider.getBlock("latest");
    const secondEndTime = BigInt(nowBlock.timestamp) + 3600n;
    await nft.mintNFT(await seller.getAddress());
    await nft.approve(await auction.getAddress(), 1);
    await auction.createAuction(await nft.getAddress(), 1, secondEndTime);

    await token.mint(await bidder1.getAddress(), ethers.parseEther("3000"));
    await token.mint(await bidder2.getAddress(), ethers.parseEther("3000"));
    await token.connect(bidder1).approve(await auction.getAddress(), ethers.parseEther("2500"));
    await auction
      .connect(bidder1)
      .bidWithERC20(1, await token.getAddress(), ethers.parseEther("2500"));

    await token.connect(bidder2).approve(await auction.getAddress(), ethers.parseEther("2600"));
    await auction
      .connect(bidder2)
      .bidWithERC20(1, await token.getAddress(), ethers.parseEther("2600"));

    expect(
      await auction.pendingTokenReturns(await bidder1.getAddress(), await token.getAddress()),
    ).to.equal(ethers.parseEther("2500"));

    await auction.connect(bidder1).withdrawTokenRefund(await token.getAddress());
    expect(
      await auction.pendingTokenReturns(await bidder1.getAddress(), await token.getAddress()),
    ).to.equal(0n);
  });

  it("should end auction and transfer NFT with payout", async function () {
    const { seller, bidder1, token, nft, auction, endTime } = await createAuctionFixture();

    await token.mint(await bidder1.getAddress(), ethers.parseEther("2200"));
    await token.connect(bidder1).approve(await auction.getAddress(), ethers.parseEther("2200"));
    await auction
      .connect(bidder1)
      .bidWithERC20(0, await token.getAddress(), ethers.parseEther("2200"));

    await ethers.provider.send("evm_setNextBlockTimestamp", [Number(endTime + 1n)]);
    await ethers.provider.send("evm_mine", []);

    const sellerBefore = await token.balanceOf(await seller.getAddress());
    await auction.endAuction(0);
    const sellerAfter = await token.balanceOf(await seller.getAddress());

    expect(await nft.ownerOf(0)).to.equal(await bidder1.getAddress());
    expect(sellerAfter - sellerBefore).to.equal(ethers.parseEther("2200"));

    await expect(auction.endAuction(0)).to.be.revertedWith("Already ended");
  });

  it("should reject invalid auction id", async function () {
    const { bidder1, auction } = await createAuctionFixture();
    await expect(
      auction.connect(bidder1).bidWithETH(99, { value: ethers.parseEther("1") }),
    ).to.be.revertedWith("Auction not found");
  });

  it("should reject bids after auction end time", async function () {
    const { bidder1, token, auction, endTime } = await createAuctionFixture();

    await ethers.provider.send("evm_setNextBlockTimestamp", [Number(endTime + 1n)]);
    await ethers.provider.send("evm_mine", []);

    await expect(
      auction.connect(bidder1).bidWithETH(0, { value: ethers.parseEther("1") }),
    ).to.be.revertedWith("Auction expired");

    await token.mint(await bidder1.getAddress(), ethers.parseEther("2000"));
    await token.connect(bidder1).approve(await auction.getAddress(), ethers.parseEther("1500"));
    await expect(
      auction
        .connect(bidder1)
        .bidWithERC20(0, await token.getAddress(), ethers.parseEther("1500")),
    ).to.be.revertedWith("Auction expired");
  });

  it("should reject ending auction without any bids", async function () {
    const { auction, endTime } = await createAuctionFixture();
    await ethers.provider.send("evm_setNextBlockTimestamp", [Number(endTime + 1n)]);
    await ethers.provider.send("evm_mine", []);
    await expect(auction.endAuction(0)).to.be.revertedWith("No bids");
  });

  it("should reject ERC20 bid when token feed is not configured", async function () {
    const { bidder1, auction } = await createAuctionFixture();

    const Token = await ethers.getContractFactory("MyToken");
    const anotherToken = await Token.deploy();
    await anotherToken.waitForDeployment();
    await anotherToken.mint(await bidder1.getAddress(), ethers.parseEther("1000"));
    await anotherToken
      .connect(bidder1)
      .approve(await auction.getAddress(), ethers.parseEther("800"));

    await expect(
      auction
        .connect(bidder1)
        .bidWithERC20(0, await anotherToken.getAddress(), ethers.parseEther("800")),
    ).to.be.revertedWith("Feed not set");
  });
});