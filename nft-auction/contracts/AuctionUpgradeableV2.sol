// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "./AuctionUpgradeable.sol";

contract AuctionUpgradeableV2 is AuctionUpgradeable {
    uint256 public platformFeeBps;

    event PlatformFeeUpdated(uint256 feeBps);

    function setPlatformFeeBps(uint256 feeBps) external onlyOwner {
        require(feeBps <= 10000, "Fee too high");
        platformFeeBps = feeBps;
        emit PlatformFeeUpdated(feeBps);
    }

    function version() external pure returns (string memory) {
        return "v2";
    }
}
