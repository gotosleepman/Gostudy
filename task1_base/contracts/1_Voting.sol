// SPDX-License-Identifier: MIT

pragma solidity ^0.8.7;

contract Voting {

    mapping(address => uint256) private _votes;

    address[] private _candidates;

    mapping(address => bool) private _hasVoted;

    address private _owner;


    event Voted(address indexed voter , address indexed candidate , uint256 newVoteCount);
    event VotesReset(address indexed resetBy , uint256 candidateCount);
    event CandidateAdded(address indexed candidate);


    modifier onlyOwner() {
        require(msg.sender == _owner,"Voting: caller is not the owner");
        _;
    }

    modifier nonZeroAddress(address addr) {
        require(addr != address(0),"Voting:zero address not allowed");
        _;
    }

    constructor() {
        _owner = msg.sender;
    }


    // function addCandidate(address candidate)external onlyOwner nonZeroAddress(candidate) {
    //     bool alreadyExists = false;
    //     for (uint256 i = 0 ; i < _candidates.length ; i++) {
    //         if (_candidates[i] == candidate) {
    //             alreadyExists = true;
    //             break;
    //         }
    //     }
    // }

    function addCandidate(address candidate) external onlyOwner nonZeroAddress(candidate) {
        require(!isCandidate(candidate), "Voting: candidate already exists");
        _candidates.push(candidate);  // ← 补全这行，实际添加候选人
        emit CandidateAdded(candidate);
}

    function vote(address candidate) external nonZeroAddress(candidate) {
        require(isCandidate(candidate),"Voting: candidate not registered");
        require(!_hasVoted[msg.sender],"Voting:already voted");


        _votes[candidate] += 1;
        _hasVoted[msg.sender] = true;

        emit Voted(msg.sender, candidate, _votes[candidate]);

    }

    function getVotes(address candidate )external view nonZeroAddress(candidate) returns(uint256) {
        return _votes[candidate];

    }

    function resetVotes() external onlyOwner {
        for (uint256 i = 0 ; i < _candidates.length ; i++) {
            _votes[_candidates[i]] = 0;
        }

        emit VotesReset(msg.sender,_candidates.length);
    }

    function resetAllVotingStatus() external onlyOwner {
        for (uint256 i = 0; i < _candidates.length; i++) {
            _votes[_candidates[i]] = 0;
        }
        
        emit VotesReset(msg.sender, _candidates.length);
    }
    
    function isCandidate(address candidate) public view returns (bool) {
        for (uint256 i = 0; i < _candidates.length; i++) {
            if (_candidates[i] == candidate) {
                return true;
            }
        }
        return false;
    }
    
    function getAllCandidates() external view returns (address[] memory) {
        return _candidates;
    }
    
    function getCandidateCount() external view returns (uint256) {
        return _candidates.length;
    }
    
    function hasVoted(address voter) external view returns (bool) {
        return _hasVoted[voter];
    }
    
    function owner() external view returns (address) {
        return _owner;
    }
    
    function transferOwnership(address newOwner) external onlyOwner nonZeroAddress(newOwner) {
        _owner = newOwner;
    }


}