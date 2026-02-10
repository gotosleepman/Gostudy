// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;


contract mergeSortArray {

    function mergeSorted(uint256[] memory arr1 , uint256[] memory arr2)
        public
        pure 
        returns(uint256[] memory) 
    {

        if(arr1.length == 0) return arr2;
        if(arr2.length == 0) return arr1;


        uint256[] memory result = new uint256[] (arr1.length + arr2.length);


        uint256 i = 0; // arr1 的指针
        uint256 j = 0; // arr2 的指针
        uint256 k = 0; // result 的指针 

        while (i < arr1.length && j < arr2.length) {
        if (arr1[i] <= arr2[j]) {
            result[k] = arr1[i];
            i++;
        } else {
            result[k] = arr2[j];
            j++;
        }
        k++;
        }
        
        // 处理剩余元素（arr1 或 arr2 中可能还有剩余）
        while (i < arr1.length) {
            result[k] = arr1[i];
            i++;
            k++;
        }
        
        while (j < arr2.length) {
            result[k] = arr2[j];
            j++;
            k++;
        }
        
        return result;     

    }
}