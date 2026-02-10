// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;


contract binarySearch {
    function binarySearchSafe(uint256[] memory arr, uint256 target)
        public
        pure
        returns (int256)
    {
        if (arr.length == 0) return -1;
        
        int256 left = 0;
        int256 right = int256(arr.length) - 1;
        
        while (left <= right) {
            // 使用 int256 计算，避免溢出问题
            int256 mid = left + (right - left) / 2;
            
            if (arr[uint256(mid)] == target) {
                return mid;
            } else if (arr[uint256(mid)] < target) {
                left = mid + 1;
            } else {
                right = mid - 1;
            }
        }
        
        return -1;
    }

}
