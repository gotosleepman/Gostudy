// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";

contract AuctionUpgradeable is Initializable, UUPSUpgradeable, OwnableUpgradeable {

    struct Auction {
        address seller;
        address nftAddress;
        uint256 tokenId;
        uint256 highestBid;
        address highestBidder;
        bool ended;
    }

    uint256 public auctionId;
    mapping(uint256 => Auction) public auctions;

    AggregatorV3Interface public priceFeed;

    function initialize(address _priceFeed) public initializer {
        __Ownable_init(msg.sender);
        priceFeed = AggregatorV3Interface(_priceFeed);
    }

    function _authorizeUpgrade(address newImplementation)
        internal
        override
        onlyOwner
    {}

    function createAuction(address nftAddress, uint256 tokenId) public {
        IERC721(nftAddress).transferFrom(
            msg.sender,
            address(this),
            tokenId
        );

        auctions[auctionId] = Auction(
            msg.sender,
            nftAddress,
            tokenId,
            0,
            address(0),
            false
        );

        auctionId++;
    }

    function bidWithETH(uint256 _auctionId) public payable {
        Auction storage auction = auctions[_auctionId];
        require(msg.value > auction.highestBid, "Bid too low");

        auction.highestBid = msg.value;
        auction.highestBidder = msg.sender;
    }

    function endAuction(uint256 _auctionId) public {
        Auction storage auction = auctions[_auctionId];
        require(!auction.ended, "Already ended");

        auction.ended = true;

        IERC721(auction.nftAddress).transferFrom(
            address(this),
            auction.highestBidder,
            auction.tokenId
        );

        payable(auction.seller).transfer(auction.highestBid);
    }

    function getLatestPrice() public view returns (int) {
        (, int price,,,) = priceFeed.latestRoundData();
        return price;
    }
}