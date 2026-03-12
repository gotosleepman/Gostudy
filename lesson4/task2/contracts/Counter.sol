// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract Counter {
    uint256 private count;

    event Incremented(uint256 newValue);

    function getCount() external view returns (uint256) {
        return count;
    }

    function increment() external {
        count += 1;
        emit Incremented(count);
    }
}
