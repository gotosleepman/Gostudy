// SPDX-License-Identifier: MIT
pragma solidity ^0.8.23;

import {SimpleLendingPool} from "./SimpleLendingPool.sol";

/// @notice UUPS 升级演示版本，新增一个可读版本号与附加状态
contract SimpleLendingPoolV2 is SimpleLendingPool {
    string public protocolNote;

    function setProtocolNote(string calldata note) external onlyOwner {
        protocolNote = note;
    }

    function version() external pure returns (string memory) {
        return "v2";
    }
}
