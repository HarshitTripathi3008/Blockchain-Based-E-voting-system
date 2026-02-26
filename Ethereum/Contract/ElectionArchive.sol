pragma solidity ^0.8.0;

contract ElectionArchive {
    address public admin;

    struct FinalResult {
        address electionAddress;
        string title;
        string winnerName;
        uint256 winningVotes;
        uint256 totalVoters;
        uint256 timestamp;
    }

    mapping(address => FinalResult) public archivedResults;

    modifier onlyAdmin() {
        require(msg.sender == admin, "Only admin can archive results");
        _;
    }

    constructor() {
        admin = msg.sender;
    }

    function archiveResult(
        address _electionAddress,
        string memory _title,
        string memory _winnerName,
        uint256 _winningVotes,
        uint256 _totalVoters
    ) public onlyAdmin {
        archivedResults[_electionAddress] = FinalResult({
            electionAddress: _electionAddress,
            title: _title,
            winnerName: _winnerName,
            winningVotes: _winningVotes,
            totalVoters: _totalVoters,
            timestamp: block.timestamp
        });
    }
}
