// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract BeggingContract {
    // 合约所有者
    address public owner;
    
    // 记录每个地址的捐赠金额
    mapping(address => uint256) public donations;
    
    // 捐赠总额
    uint256 public totalDonations;
    
    // 捐赠时间限制（可选功能）
    uint256 public donationStartTime;
    uint256 public donationEndTime;
    
    // 捐赠事件
    event Donation(address indexed donor, uint256 amount);
    
    // 提款事件
    event Withdrawal(address indexed owner, uint256 amount);
    
    // 捐赠者结构体（用于排行榜）
    struct Donor {
        address donorAddress;
        uint256 amount;
    }
    
    // 捐赠者数组
    Donor[] private allDonors;
    
    // 修改器：仅合约所有者可调用
    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner can call this function");
        _;
    }
    
    // 构造函数：设置合约所有者
    constructor() payable {
        owner = msg.sender;
        donationStartTime = block.timestamp; // 立即开始
        donationEndTime = block.timestamp + 30 days; // 30天捐赠期
    }
    
    // 捐赠函数：接收以太币
    function donate() external payable {
        require(block.timestamp >= donationStartTime, "Donation period has not started");
        require(block.timestamp <= donationEndTime, "Donation period has ended");
        require(msg.value > 0, "Donation amount must be greater than 0");
        
        // 更新捐赠记录
        donations[msg.sender] += msg.value;
        totalDonations += msg.value;
        
        // 更新或添加捐赠者到数组
        bool donorExists = false;
        for (uint256 i = 0; i < allDonors.length; i++) {
            if (allDonors[i].donorAddress == msg.sender) {
                allDonors[i].amount += msg.value;
                donorExists = true;
                break;
            }
        }
        
        if (!donorExists) {
            allDonors.push(Donor(msg.sender, msg.value));
        }
        
        // 触发捐赠事件
        emit Donation(msg.sender, msg.value);
    }
    
    // 获取指定地址的捐赠金额
    function getDonation(address donor) external view returns (uint256) {
        return donations[donor];
    }
    
    // 获取合约余额
    function getContractBalance() external view returns (uint256) {
        return address(this).balance;
    }
    
    // 提款函数：仅所有者可调用
    function withdraw() external onlyOwner {
        uint256 balance = address(this).balance;
        require(balance > 0, "No funds to withdraw");
        
        // 使用transfer发送以太币
        payable(owner).transfer(balance);
        
        // 触发提款事件
        emit Withdrawal(owner, balance);
    }
    
    // 获取前3名捐赠者（简化版）
    function getTop3Donors() external view returns (address[3] memory topAddresses, uint256[3] memory topAmounts) {
        // 初始化数组
        for (uint256 i = 0; i < 3; i++) {
            topAddresses[i] = address(0);
            topAmounts[i] = 0;
        }
        
        // 简单的前3名查找算法
        for (uint256 i = 0; i < allDonors.length; i++) {
            address currentAddress = allDonors[i].donorAddress;
            uint256 currentAmount = allDonors[i].amount;
            
            // 检查是否进入前3名
            if (currentAmount > topAmounts[0]) {
                topAmounts[2] = topAmounts[1];
                topAddresses[2] = topAddresses[1];
                topAmounts[1] = topAmounts[0];
                topAddresses[1] = topAddresses[0];
                topAmounts[0] = currentAmount;
                topAddresses[0] = currentAddress;
            } else if (currentAmount > topAmounts[1]) {
                topAmounts[2] = topAmounts[1];
                topAddresses[2] = topAddresses[1];
                topAmounts[1] = currentAmount;
                topAddresses[1] = currentAddress;
            } else if (currentAmount > topAmounts[2]) {
                topAmounts[2] = currentAmount;
                topAddresses[2] = currentAddress;
            }
        }
        
        return (topAddresses, topAmounts);
    }
    
    // 获取捐赠时间信息
    function getDonationPeriod() external view returns (uint256 startTime, uint256 endTime, bool isActive) {
        isActive = (block.timestamp >= donationStartTime && block.timestamp <= donationEndTime);
        return (donationStartTime, donationEndTime, isActive);
    }
    
    // 接收以太币的回退函数
    receive() external payable {
        require(block.timestamp >= donationStartTime, "Donation period has not started");
        require(block.timestamp <= donationEndTime, "Donation period has ended");
        require(msg.value > 0, "Donation amount must be greater than 0");
        
        donations[msg.sender] += msg.value;
        totalDonations += msg.value;
        
        emit Donation(msg.sender, msg.value);
    }
}