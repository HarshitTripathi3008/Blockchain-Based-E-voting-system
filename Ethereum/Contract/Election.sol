// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract ElectionFact {
    
    struct ElectionDet {
        address deployedAddress;
        string el_n;
        string el_d;
    }
    
    mapping(string => ElectionDet) public companyEmail;
    
    function createElection(string memory email, string memory election_name, string memory election_description) public {
        address newElection = address(new Election(msg.sender, election_name, election_description));
        
        companyEmail[email] = ElectionDet({
            deployedAddress: newElection,
            el_n: election_name,
            el_d: election_description
        });
    }
    
    function getDeployedElection(string memory email) public view returns (address, string memory, string memory) {
        ElectionDet memory election = companyEmail[email];
        if(election.deployedAddress == address(0)) {
            return (address(0), "", "Create an election.");
        } else {
            return (election.deployedAddress, election.el_n, election.el_d);
        }
    }
}

contract Election {
    address public election_authority;
    string public election_name;
    string public election_description;
    bool public status;
    
    struct Candidate {
        string candidate_name;
        string candidate_description;
        string imgHash;
        uint256 voteCount;
        string email;
    }
    
    struct Voter {
        uint256 candidate_id_voted;
        bool voted;
    }
    
    mapping(uint256 => Candidate) public candidates;
    mapping(string => Voter) public voters;
    
    uint256 public numCandidates;
    uint256 public numVoters;
    
    constructor(address authority, string memory name, string memory description) {
        election_authority = authority;
        election_name = name;
        election_description = description;
        status = true;
    }
    
    modifier owner() {
        require(msg.sender == election_authority, "Error: Access Denied.");
        _;
    }
    
    function addCandidate(
        string memory candidate_name, 
        string memory candidate_description, 
        string memory imgHash, 
        string memory email
    ) public owner {
        uint256 candidateID = numCandidates;
        candidates[candidateID] = Candidate({
            candidate_name: candidate_name,
            candidate_description: candidate_description,
            imgHash: imgHash,
            voteCount: 0,
            email: email
        });
        numCandidates++;
    }
    
    function vote(uint256 candidateID, string memory e) public {
        require(!voters[e].voted, "Error: You cannot double vote");
        require(candidateID < numCandidates, "Error: Invalid candidate ID");
        
        voters[e] = Voter({
            candidate_id_voted: candidateID,
            voted: true
        });
        
        numVoters++;
        candidates[candidateID].voteCount++;
    }
    
    function getNumOfCandidates() public view returns(uint256) {
        return numCandidates;
    }
    
    function getNumOfVoters() public view returns(uint256) {
        return numVoters;
    }
    
    function getCandidate(uint256 candidateID) public view returns (
        string memory, 
        string memory, 
        string memory, 
        uint256, 
        string memory
    ) {
        require(candidateID < numCandidates, "Error: Invalid candidate ID");
        Candidate memory candidate = candidates[candidateID];
        return (
            candidate.candidate_name, 
            candidate.candidate_description, 
            candidate.imgHash, 
            candidate.voteCount, 
            candidate.email
        );
    }
    
    function winnerCandidate() public view owner returns (uint256) {
        require(numCandidates > 0, "Error: No candidates");
        
        uint256 largestVotes = candidates[0].voteCount;
        uint256 winningCandidateID = 0;
        
        for(uint256 i = 1; i < numCandidates; i++) {
            if(candidates[i].voteCount > largestVotes) {
                largestVotes = candidates[i].voteCount;
                winningCandidateID = i;
            }
        }
        return winningCandidateID;
    }
    
    function getElectionDetails() public view returns(string memory, string memory) {
        return (election_name, election_description);    
    }
    
    // Helper function to check if voter has voted
    function hasVoted(string memory email) public view returns (bool) {
        return voters[email].voted;
    }
    
    // Helper function to get voter details
    function getVoterDetails(string memory email) public view returns (uint256, bool) {
        Voter memory voter = voters[email];
        return (voter.candidate_id_voted, voter.voted);
    }
    

}