// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// ElectionArchiveMetaData contains all meta data concerning the ElectionArchive contract.
var ElectionArchiveMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"admin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_electionAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_title\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_winnerName\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_winningVotes\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_totalVoters\",\"type\":\"uint256\"}],\"name\":\"archiveResult\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"archivedResults\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"electionAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"title\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"winnerName\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"winningVotes\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalVoters\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550610b65806100606000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c806331bd3af7146100465780635c22ed3a1461007b578063f851a44014610097575b600080fd5b610060600480360381019061005b9190610462565b6100b5565b60405161007296959493929190610547565b60405180910390f35b61009560048036038101906100909190610717565b610221565b005b61009f6103cc565b6040516100ac91906107ca565b60405180910390f35b60016020528060005260406000206000915090508060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060010180546100fe90610814565b80601f016020809104026020016040519081016040528092919081815260200182805461012a90610814565b80156101775780601f1061014c57610100808354040283529160200191610177565b820191906000526020600020905b81548152906001019060200180831161015a57829003601f168201915b50505050509080600201805461018c90610814565b80601f01602080910402602001604051908101604052809291908181526020018280546101b890610814565b80156102055780601f106101da57610100808354040283529160200191610205565b820191906000526020600020905b8154815290600101906020018083116101e857829003601f168201915b5050505050908060030154908060040154908060050154905086565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146102af576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016102a690610891565b60405180910390fd5b6040518060c001604052808673ffffffffffffffffffffffffffffffffffffffff16815260200185815260200184815260200183815260200182815260200142815250600160008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008201518160000160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550602082015181600101908161038d9190610a5d565b5060408201518160020190816103a39190610a5d565b50606082015181600301556080820151816004015560a082015181600501559050505050505050565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000604051905090565b600080fd5b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061042f82610404565b9050919050565b61043f81610424565b811461044a57600080fd5b50565b60008135905061045c81610436565b92915050565b600060208284031215610478576104776103fa565b5b60006104868482850161044d565b91505092915050565b61049881610424565b82525050565b600081519050919050565b600082825260208201905092915050565b60005b838110156104d85780820151818401526020810190506104bd565b60008484015250505050565b6000601f19601f8301169050919050565b60006105008261049e565b61050a81856104a9565b935061051a8185602086016104ba565b610523816104e4565b840191505092915050565b6000819050919050565b6105418161052e565b82525050565b600060c08201905061055c600083018961048f565b818103602083015261056e81886104f5565b9050818103604083015261058281876104f5565b90506105916060830186610538565b61059e6080830185610538565b6105ab60a0830184610538565b979650505050505050565b600080fd5b600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6105f8826104e4565b810181811067ffffffffffffffff82111715610617576106166105c0565b5b80604052505050565b600061062a6103f0565b905061063682826105ef565b919050565b600067ffffffffffffffff821115610656576106556105c0565b5b61065f826104e4565b9050602081019050919050565b82818337600083830152505050565b600061068e6106898461063b565b610620565b9050828152602081018484840111156106aa576106a96105bb565b5b6106b584828561066c565b509392505050565b600082601f8301126106d2576106d16105b6565b5b81356106e284826020860161067b565b91505092915050565b6106f48161052e565b81146106ff57600080fd5b50565b600081359050610711816106eb565b92915050565b600080600080600060a08688031215610733576107326103fa565b5b60006107418882890161044d565b955050602086013567ffffffffffffffff811115610762576107616103ff565b5b61076e888289016106bd565b945050604086013567ffffffffffffffff81111561078f5761078e6103ff565b5b61079b888289016106bd565b93505060606107ac88828901610702565b92505060806107bd88828901610702565b9150509295509295909350565b60006020820190506107df600083018461048f565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b6000600282049050600182168061082c57607f821691505b60208210810361083f5761083e6107e5565b5b50919050565b7f4f6e6c792061646d696e2063616e206172636869766520726573756c74730000600082015250565b600061087b601e836104a9565b915061088682610845565b602082019050919050565b600060208201905081810360008301526108aa8161086e565b9050919050565b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b6000600883026109137fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff826108d6565b61091d86836108d6565b95508019841693508086168417925050509392505050565b6000819050919050565b600061095a6109556109508461052e565b610935565b61052e565b9050919050565b6000819050919050565b6109748361093f565b61098861098082610961565b8484546108e3565b825550505050565b600090565b61099d610990565b6109a881848461096b565b505050565b5b818110156109cc576109c1600082610995565b6001810190506109ae565b5050565b601f821115610a11576109e2816108b1565b6109eb846108c6565b810160208510156109fa578190505b610a0e610a06856108c6565b8301826109ad565b50505b505050565b600082821c905092915050565b6000610a3460001984600802610a16565b1980831691505092915050565b6000610a4d8383610a23565b9150826002028217905092915050565b610a668261049e565b67ffffffffffffffff811115610a7f57610a7e6105c0565b5b610a898254610814565b610a948282856109d0565b600060209050601f831160018114610ac75760008415610ab5578287015190505b610abf8582610a41565b865550610b27565b601f198416610ad5866108b1565b60005b82811015610afd57848901518255600182019150602085019450602081019050610ad8565b86831015610b1a5784890151610b16601f891682610a23565b8355505b6001600288020188555050505b50505050505056fea26469706673582212204f2976018faf140f92c41955e74fa914e7fccdef6de1be6b156a26299523cbb764736f6c63430008130033",
}

// ElectionArchiveABI is the input ABI used to generate the binding from.
// Deprecated: Use ElectionArchiveMetaData.ABI instead.
var ElectionArchiveABI = ElectionArchiveMetaData.ABI

// ElectionArchiveBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ElectionArchiveMetaData.Bin instead.
var ElectionArchiveBin = ElectionArchiveMetaData.Bin

// DeployElectionArchive deploys a new Ethereum contract, binding an instance of ElectionArchive to it.
func DeployElectionArchive(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ElectionArchive, error) {
	parsed, err := ElectionArchiveMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ElectionArchiveBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ElectionArchive{ElectionArchiveCaller: ElectionArchiveCaller{contract: contract}, ElectionArchiveTransactor: ElectionArchiveTransactor{contract: contract}, ElectionArchiveFilterer: ElectionArchiveFilterer{contract: contract}}, nil
}

// ElectionArchive is an auto generated Go binding around an Ethereum contract.
type ElectionArchive struct {
	ElectionArchiveCaller     // Read-only binding to the contract
	ElectionArchiveTransactor // Write-only binding to the contract
	ElectionArchiveFilterer   // Log filterer for contract events
}

// ElectionArchiveCaller is an auto generated read-only Go binding around an Ethereum contract.
type ElectionArchiveCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ElectionArchiveTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ElectionArchiveTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ElectionArchiveFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ElectionArchiveFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ElectionArchiveSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ElectionArchiveSession struct {
	Contract     *ElectionArchive  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ElectionArchiveCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ElectionArchiveCallerSession struct {
	Contract *ElectionArchiveCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// ElectionArchiveTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ElectionArchiveTransactorSession struct {
	Contract     *ElectionArchiveTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// ElectionArchiveRaw is an auto generated low-level Go binding around an Ethereum contract.
type ElectionArchiveRaw struct {
	Contract *ElectionArchive // Generic contract binding to access the raw methods on
}

// ElectionArchiveCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ElectionArchiveCallerRaw struct {
	Contract *ElectionArchiveCaller // Generic read-only contract binding to access the raw methods on
}

// ElectionArchiveTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ElectionArchiveTransactorRaw struct {
	Contract *ElectionArchiveTransactor // Generic write-only contract binding to access the raw methods on
}

// NewElectionArchive creates a new instance of ElectionArchive, bound to a specific deployed contract.
func NewElectionArchive(address common.Address, backend bind.ContractBackend) (*ElectionArchive, error) {
	contract, err := bindElectionArchive(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ElectionArchive{ElectionArchiveCaller: ElectionArchiveCaller{contract: contract}, ElectionArchiveTransactor: ElectionArchiveTransactor{contract: contract}, ElectionArchiveFilterer: ElectionArchiveFilterer{contract: contract}}, nil
}

// NewElectionArchiveCaller creates a new read-only instance of ElectionArchive, bound to a specific deployed contract.
func NewElectionArchiveCaller(address common.Address, caller bind.ContractCaller) (*ElectionArchiveCaller, error) {
	contract, err := bindElectionArchive(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ElectionArchiveCaller{contract: contract}, nil
}

// NewElectionArchiveTransactor creates a new write-only instance of ElectionArchive, bound to a specific deployed contract.
func NewElectionArchiveTransactor(address common.Address, transactor bind.ContractTransactor) (*ElectionArchiveTransactor, error) {
	contract, err := bindElectionArchive(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ElectionArchiveTransactor{contract: contract}, nil
}

// NewElectionArchiveFilterer creates a new log filterer instance of ElectionArchive, bound to a specific deployed contract.
func NewElectionArchiveFilterer(address common.Address, filterer bind.ContractFilterer) (*ElectionArchiveFilterer, error) {
	contract, err := bindElectionArchive(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ElectionArchiveFilterer{contract: contract}, nil
}

// bindElectionArchive binds a generic wrapper to an already deployed contract.
func bindElectionArchive(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ElectionArchiveMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ElectionArchive *ElectionArchiveRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ElectionArchive.Contract.ElectionArchiveCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ElectionArchive *ElectionArchiveRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ElectionArchive.Contract.ElectionArchiveTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ElectionArchive *ElectionArchiveRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ElectionArchive.Contract.ElectionArchiveTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ElectionArchive *ElectionArchiveCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ElectionArchive.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ElectionArchive *ElectionArchiveTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ElectionArchive.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ElectionArchive *ElectionArchiveTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ElectionArchive.Contract.contract.Transact(opts, method, params...)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_ElectionArchive *ElectionArchiveCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ElectionArchive.contract.Call(opts, &out, "admin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_ElectionArchive *ElectionArchiveSession) Admin() (common.Address, error) {
	return _ElectionArchive.Contract.Admin(&_ElectionArchive.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_ElectionArchive *ElectionArchiveCallerSession) Admin() (common.Address, error) {
	return _ElectionArchive.Contract.Admin(&_ElectionArchive.CallOpts)
}

// ArchivedResults is a free data retrieval call binding the contract method 0x31bd3af7.
//
// Solidity: function archivedResults(address ) view returns(address electionAddress, string title, string winnerName, uint256 winningVotes, uint256 totalVoters, uint256 timestamp)
func (_ElectionArchive *ElectionArchiveCaller) ArchivedResults(opts *bind.CallOpts, arg0 common.Address) (struct {
	ElectionAddress common.Address
	Title           string
	WinnerName      string
	WinningVotes    *big.Int
	TotalVoters     *big.Int
	Timestamp       *big.Int
}, error) {
	var out []interface{}
	err := _ElectionArchive.contract.Call(opts, &out, "archivedResults", arg0)

	outstruct := new(struct {
		ElectionAddress common.Address
		Title           string
		WinnerName      string
		WinningVotes    *big.Int
		TotalVoters     *big.Int
		Timestamp       *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ElectionAddress = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Title = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.WinnerName = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.WinningVotes = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.TotalVoters = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.Timestamp = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// ArchivedResults is a free data retrieval call binding the contract method 0x31bd3af7.
//
// Solidity: function archivedResults(address ) view returns(address electionAddress, string title, string winnerName, uint256 winningVotes, uint256 totalVoters, uint256 timestamp)
func (_ElectionArchive *ElectionArchiveSession) ArchivedResults(arg0 common.Address) (struct {
	ElectionAddress common.Address
	Title           string
	WinnerName      string
	WinningVotes    *big.Int
	TotalVoters     *big.Int
	Timestamp       *big.Int
}, error) {
	return _ElectionArchive.Contract.ArchivedResults(&_ElectionArchive.CallOpts, arg0)
}

// ArchivedResults is a free data retrieval call binding the contract method 0x31bd3af7.
//
// Solidity: function archivedResults(address ) view returns(address electionAddress, string title, string winnerName, uint256 winningVotes, uint256 totalVoters, uint256 timestamp)
func (_ElectionArchive *ElectionArchiveCallerSession) ArchivedResults(arg0 common.Address) (struct {
	ElectionAddress common.Address
	Title           string
	WinnerName      string
	WinningVotes    *big.Int
	TotalVoters     *big.Int
	Timestamp       *big.Int
}, error) {
	return _ElectionArchive.Contract.ArchivedResults(&_ElectionArchive.CallOpts, arg0)
}

// ArchiveResult is a paid mutator transaction binding the contract method 0x5c22ed3a.
//
// Solidity: function archiveResult(address _electionAddress, string _title, string _winnerName, uint256 _winningVotes, uint256 _totalVoters) returns()
func (_ElectionArchive *ElectionArchiveTransactor) ArchiveResult(opts *bind.TransactOpts, _electionAddress common.Address, _title string, _winnerName string, _winningVotes *big.Int, _totalVoters *big.Int) (*types.Transaction, error) {
	return _ElectionArchive.contract.Transact(opts, "archiveResult", _electionAddress, _title, _winnerName, _winningVotes, _totalVoters)
}

// ArchiveResult is a paid mutator transaction binding the contract method 0x5c22ed3a.
//
// Solidity: function archiveResult(address _electionAddress, string _title, string _winnerName, uint256 _winningVotes, uint256 _totalVoters) returns()
func (_ElectionArchive *ElectionArchiveSession) ArchiveResult(_electionAddress common.Address, _title string, _winnerName string, _winningVotes *big.Int, _totalVoters *big.Int) (*types.Transaction, error) {
	return _ElectionArchive.Contract.ArchiveResult(&_ElectionArchive.TransactOpts, _electionAddress, _title, _winnerName, _winningVotes, _totalVoters)
}

// ArchiveResult is a paid mutator transaction binding the contract method 0x5c22ed3a.
//
// Solidity: function archiveResult(address _electionAddress, string _title, string _winnerName, uint256 _winningVotes, uint256 _totalVoters) returns()
func (_ElectionArchive *ElectionArchiveTransactorSession) ArchiveResult(_electionAddress common.Address, _title string, _winnerName string, _winningVotes *big.Int, _totalVoters *big.Int) (*types.Transaction, error) {
	return _ElectionArchive.Contract.ArchiveResult(&_ElectionArchive.TransactOpts, _electionAddress, _title, _winnerName, _winningVotes, _totalVoters)
}
