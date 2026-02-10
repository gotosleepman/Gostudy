//SPDX-License-Identifier:MIT

pragma solidity ^0.8.20;

contract StringReverser {
    function reverseString(string memory str) public pure returns (string memory) {
        bytes memory strBytes = bytes(str);
        uint length = strBytes.length;

        for (uint i = 0 ; i < length / 2 ; i++) {
        //交换 strBytes[i] 和 strBytes[length - 1 -i]
            bytes1 temp = strBytes[i];
            strBytes[i] = strBytes[length - 1 -i]; 
            strBytes[length - 1 -i] = temp;

        }

        return string (strBytes) ;
    }
}