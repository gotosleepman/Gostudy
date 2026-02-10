//SPDX-License-Identifier:MIT
pragma solidity ^0.8.20;


contract intToRoman {
    function toRoman(uint256 num) public pure returns(string memory) {
        require(num >= 1 && num <= 3999, "Number out of range [1, 3999]");
        
        string[10] memory thousands = ["", "M", "MM", "MMM", "", "", "", "", "", ""];
        string[10] memory hundreds  = ["", "C", "CC", "CCC", "CD", "D", "DC", "DCC", "DCCC", "CM"];
        string[10] memory tens      = ["", "X", "XX", "XXX", "XL", "L", "LX", "LXX", "LXXX", "XC"];
        string[10] memory ones      = ["", "I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX"];
        

        uint256 thousand = num / 1000;
        uint256 hundred = (num % 1000) / 100;
        uint256 ten = (num % 100) / 10;
        uint256 one = num % 10;
        

        return string.concat(
            thousands[thousand],
            hundreds[hundred],
            tens[ten],
            ones[one]
        );
    }


}
