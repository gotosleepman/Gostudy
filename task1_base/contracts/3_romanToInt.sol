//SPDX-License-Identifier:MIT
pragma solidity ^0.8.20;

contract RomanToInt {
    mapping (bytes1 => uint256) private romanValues;

    constructor() {
        romanValues[bytes1('I')] = 1;
        romanValues[bytes1('V')] = 5;
        romanValues[bytes1('X')] = 10;
        romanValues[bytes1('L')] = 50;
        romanValues[bytes1('C')] = 100;
        romanValues[bytes1('D')] = 500;
        romanValues[bytes1('M')] = 1000;
    }

    function romanToInt(string memory s) public view returns(uint256) {
        require(bytes(s).length > 0 && bytes(s).length <= 15, "The string length does not meet the requirements ");


        bytes memory romanBytes = bytes(s);
        uint256 result = 0;
        uint256 prevValue = 0;


        for (uint256 i = romanBytes.length ; i > 0 ; i--) {
            bytes1  currentChar = romanBytes[i - 1];
            uint256 currentValue = getValue(currentChar);

            if (currentValue < prevValue) {
                if (currentValue * 10 < prevValue) {
                    revert("Invalid Roman numeral: invalid subtraction");
                }
                result -= currentValue;
                
            }else {
                result += currentValue;
            }

            prevValue = currentValue;

        }
        require(result >= 1 && result <= 3999, "Result out of range [1, 3999]");

        return result;

    }


    function getValue(bytes1 c) public view returns(uint256) {
        uint256 value = romanValues[c];
        require(value > 0 , "Invalid Roman character");
        return value;
    }




    function isValidRoman(string memory s) public view returns (bool) {
    if (bytes(s).length == 0) return false;
    
    bytes memory romanBytes = bytes(s);

    for (uint256 i = 0; i < romanBytes.length; i++) {
        if (romanValues[romanBytes[i]] == 0) {
            return false;
        }
    }
    
    return true;
}
}