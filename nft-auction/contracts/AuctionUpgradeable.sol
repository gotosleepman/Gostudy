// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";

contract AuctionUpgradeable is
    Initializable,
    UUPSUpgradeable,
    OwnableUpgradeable
{
    using SafeERC20 for IERC20;

    struct Auction {
        bool exists;
        address seller;
        address nftAddress;
        uint256 tokenId;
        uint256 endTime;
        address bidToken; // address(0) means ETH
        uint256 highestBid;
        uint256 highestBidUsd;
        address highestBidder;
        bool ended;
    }

    uint256 public auctionId;
    mapping(uint256 => Auction) public auctions;

    AggregatorV3Interface public ethUsdFeed;
    mapping(address => AggregatorV3Interface) public tokenUsdFeeds;

    mapping(address => uint256) public pendingEthReturns;
    mapping(address => mapping(address => uint256)) public pendingTokenReturns;
    uint256 private _lockState;

    event AuctionCreated(
        uint256 indexed auctionId,
        address indexed seller,
        address indexed nftAddress,
        uint256 tokenId,
        uint256 endTime
    );
    event BidPlaced(
        uint256 indexed auctionId,
        address indexed bidder,
        address indexed bidToken,
        uint256 amount,
        uint256 usdValue
    );
    event AuctionEnded(
        uint256 indexed auctionId,
        address indexed seller,
        address indexed winner,
        address bidToken,
        uint256 amount
    );
    event RefundQueued(
        address indexed bidder,
        address indexed refundToken,
        uint256 amount
    );
    event RefundWithdrawn(
        address indexed bidder,
        address indexed refundToken,
        uint256 amount
    );

    function initialize(address _ethUsdFeed) public initializer {
        require(_ethUsdFeed != address(0), "Invalid ETH feed");
        __Ownable_init(msg.sender);
        _lockState = 1;
        ethUsdFeed = AggregatorV3Interface(_ethUsdFeed);
    }

    function _authorizeUpgrade(address newImplementation)
        internal
        override
        onlyOwner
    {}

    modifier nonReentrant() {
        require(_lockState == 1, "Reentrancy");
        _lockState = 2;
        _;
        _lockState = 1;
    }

    function setTokenPriceFeed(address token, address feed) external onlyOwner {
        require(token != address(0), "Invalid token");
        require(feed != address(0), "Invalid feed");
        tokenUsdFeeds[token] = AggregatorV3Interface(feed);
    }

    function createAuction(
        address nftAddress,
        uint256 tokenId,
        uint256 endTime
    ) public {
        require(nftAddress != address(0), "Invalid NFT");
        require(endTime > block.timestamp, "Invalid end time");

        IERC721(nftAddress).transferFrom(msg.sender, address(this), tokenId);

        auctions[auctionId] = Auction({
            exists: true,
            seller: msg.sender,
            nftAddress: nftAddress,
            tokenId: tokenId,
            endTime: endTime,
            bidToken: address(0),
            highestBid: 0,
            highestBidUsd: 0,
            highestBidder: address(0),
            ended: false
        });

        emit AuctionCreated(auctionId, msg.sender, nftAddress, tokenId, endTime);
        auctionId++;
    }

    function bidWithETH(uint256 _auctionId) public payable nonReentrant {
        Auction storage auction = auctions[_auctionId];
        require(auction.exists, "Auction not found");
        require(!auction.ended, "Already ended");
        require(block.timestamp < auction.endTime, "Auction expired");
        require(msg.value > 0, "Bid is zero");

        uint256 bidUsd = _ethToUsd(msg.value);
        require(bidUsd > auction.highestBidUsd, "Bid too low");

        _queueRefund(auction);
        auction.bidToken = address(0);
        auction.highestBid = msg.value;
        auction.highestBidUsd = bidUsd;
        auction.highestBidder = msg.sender;

        emit BidPlaced(_auctionId, msg.sender, address(0), msg.value, bidUsd);
    }

    function bidWithERC20(
        uint256 _auctionId,
        address token,
        uint256 amount
    ) public nonReentrant {
        Auction storage auction = auctions[_auctionId];
        require(auction.exists, "Auction not found");
        require(!auction.ended, "Already ended");
        require(block.timestamp < auction.endTime, "Auction expired");
        require(token != address(0), "Invalid token");
        require(amount > 0, "Bid is zero");
        require(address(tokenUsdFeeds[token]) != address(0), "Feed not set");

        uint256 bidUsd = _tokenToUsd(token, amount);
        require(bidUsd > auction.highestBidUsd, "Bid too low");

        IERC20(token).safeTransferFrom(msg.sender, address(this), amount);
        _queueRefund(auction);

        auction.bidToken = token;
        auction.highestBid = amount;
        auction.highestBidUsd = bidUsd;
        auction.highestBidder = msg.sender;

        emit BidPlaced(_auctionId, msg.sender, token, amount, bidUsd);
    }

    function endAuction(uint256 _auctionId) public nonReentrant {
        Auction storage auction = auctions[_auctionId];
        require(auction.exists, "Auction not found");
        require(!auction.ended, "Already ended");
        require(block.timestamp >= auction.endTime, "Auction not ended");
        require(auction.highestBidder != address(0), "No bids");

        auction.ended = true;
        IERC721(auction.nftAddress).safeTransferFrom(
            address(this),
            auction.highestBidder,
            auction.tokenId
        );

        if (auction.bidToken == address(0)) {
            (bool ok, ) = payable(auction.seller).call{value: auction.highestBid}("");
            require(ok, "ETH payout failed");
        } else {
            IERC20(auction.bidToken).safeTransfer(auction.seller, auction.highestBid);
        }

        emit AuctionEnded(
            _auctionId,
            auction.seller,
            auction.highestBidder,
            auction.bidToken,
            auction.highestBid
        );
    }

    function withdrawEthRefund() external nonReentrant {
        uint256 amount = pendingEthReturns[msg.sender];
        require(amount > 0, "No ETH refund");
        pendingEthReturns[msg.sender] = 0;

        (bool ok, ) = payable(msg.sender).call{value: amount}("");
        require(ok, "ETH refund failed");

        emit RefundWithdrawn(msg.sender, address(0), amount);
    }

    function withdrawTokenRefund(address token) external nonReentrant {
        uint256 amount = pendingTokenReturns[msg.sender][token];
        require(amount > 0, "No token refund");
        pendingTokenReturns[msg.sender][token] = 0;

        IERC20(token).safeTransfer(msg.sender, amount);
        emit RefundWithdrawn(msg.sender, token, amount);
    }

    function getLatestPrice() public view returns (int256) {
        (, int256 price, , , ) = ethUsdFeed.latestRoundData();
        return price;
    }

    function _queueRefund(Auction storage auction) internal {
        if (auction.highestBidder == address(0) || auction.highestBid == 0) {
            return;
        }

        if (auction.bidToken == address(0)) {
            pendingEthReturns[auction.highestBidder] += auction.highestBid;
            emit RefundQueued(auction.highestBidder, address(0), auction.highestBid);
            return;
        }

        pendingTokenReturns[auction.highestBidder][auction.bidToken] += auction
            .highestBid;
        emit RefundQueued(auction.highestBidder, auction.bidToken, auction.highestBid);
    }

    function _ethToUsd(uint256 amount) internal view returns (uint256) {
        return _toUsdValue(ethUsdFeed, amount);
    }

    function _tokenToUsd(
        address token,
        uint256 amount
    ) internal view returns (uint256) {
        return _toUsdValue(tokenUsdFeeds[token], amount);
    }

    function _toUsdValue(
        AggregatorV3Interface feed,
        uint256 amount
    ) internal view returns (uint256) {
        require(address(feed) != address(0), "Feed not set");
        (, int256 answer, , , ) = feed.latestRoundData();
        require(answer > 0, "Invalid price");

        uint256 price = uint256(answer);
        uint256 decimals = uint256(feed.decimals());
        return (amount * price) / (10 ** decimals);
    }
}