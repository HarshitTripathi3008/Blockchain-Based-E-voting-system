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

// ElectionMetaData contains all meta data concerning the Election contract.
var ElectionMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"authority\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"candidate_name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"candidate_description\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"imgHash\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"email\",\"type\":\"string\"}],\"name\":\"addCandidate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"candidates\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"candidate_name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"candidate_description\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"imgHash\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"voteCount\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"email\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"election_authority\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"election_description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"election_name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"candidateID\",\"type\":\"uint256\"}],\"name\":\"getCandidate\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getElectionDetails\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNumOfCandidates\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNumOfVoters\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"email\",\"type\":\"string\"}],\"name\":\"getVoterDetails\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"email\",\"type\":\"string\"}],\"name\":\"hasVoted\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"numCandidates\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"numVoters\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"status\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"candidateID\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"e\",\"type\":\"string\"}],\"name\":\"vote\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"voters\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"candidate_id_voted\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"voted\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"winnerCandidate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// ElectionABI is the input ABI used to generate the binding from.
// Deprecated: Use ElectionMetaData.ABI instead.
var ElectionABI = ElectionMetaData.ABI

// Election is an auto generated Go binding around an Ethereum contract.
type Election struct {
	ElectionCaller     // Read-only binding to the contract
	ElectionTransactor // Write-only binding to the contract
	ElectionFilterer   // Log filterer for contract events
}

// ElectionCaller is an auto generated read-only Go binding around an Ethereum contract.
type ElectionCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ElectionTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ElectionTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ElectionFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ElectionFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ElectionSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ElectionSession struct {
	Contract     *Election         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ElectionCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ElectionCallerSession struct {
	Contract *ElectionCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// ElectionTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ElectionTransactorSession struct {
	Contract     *ElectionTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ElectionRaw is an auto generated low-level Go binding around an Ethereum contract.
type ElectionRaw struct {
	Contract *Election // Generic contract binding to access the raw methods on
}

// ElectionCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ElectionCallerRaw struct {
	Contract *ElectionCaller // Generic read-only contract binding to access the raw methods on
}

// ElectionTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ElectionTransactorRaw struct {
	Contract *ElectionTransactor // Generic write-only contract binding to access the raw methods on
}

// NewElection creates a new instance of Election, bound to a specific deployed contract.
func NewElection(address common.Address, backend bind.ContractBackend) (*Election, error) {
	contract, err := bindElection(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Election{ElectionCaller: ElectionCaller{contract: contract}, ElectionTransactor: ElectionTransactor{contract: contract}, ElectionFilterer: ElectionFilterer{contract: contract}}, nil
}

// NewElectionCaller creates a new read-only instance of Election, bound to a specific deployed contract.
func NewElectionCaller(address common.Address, caller bind.ContractCaller) (*ElectionCaller, error) {
	contract, err := bindElection(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ElectionCaller{contract: contract}, nil
}

// NewElectionTransactor creates a new write-only instance of Election, bound to a specific deployed contract.
func NewElectionTransactor(address common.Address, transactor bind.ContractTransactor) (*ElectionTransactor, error) {
	contract, err := bindElection(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ElectionTransactor{contract: contract}, nil
}

// NewElectionFilterer creates a new log filterer instance of Election, bound to a specific deployed contract.
func NewElectionFilterer(address common.Address, filterer bind.ContractFilterer) (*ElectionFilterer, error) {
	contract, err := bindElection(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ElectionFilterer{contract: contract}, nil
}

// bindElection binds a generic wrapper to an already deployed contract.
func bindElection(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ElectionMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Election *ElectionRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Election.Contract.ElectionCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Election *ElectionRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Election.Contract.ElectionTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Election *ElectionRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Election.Contract.ElectionTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Election *ElectionCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Election.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Election *ElectionTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Election.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Election *ElectionTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Election.Contract.contract.Transact(opts, method, params...)
}

// Candidates is a free data retrieval call binding the contract method 0x3477ee2e.
//
// Solidity: function candidates(uint256 ) view returns(string candidate_name, string candidate_description, string imgHash, uint256 voteCount, string email)
func (_Election *ElectionCaller) Candidates(opts *bind.CallOpts, arg0 *big.Int) (struct {
	CandidateName        string
	CandidateDescription string
	ImgHash              string
	VoteCount            *big.Int
	Email                string
}, error) {
	var out []interface{}
	err := _Election.contract.Call(opts, &out, "candidates", arg0)

	outstruct := new(struct {
		CandidateName        string
		CandidateDescription string
		ImgHash              string
		VoteCount            *big.Int
		Email                string
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.CandidateName = *abi.ConvertType(out[0], new(string)).(*string)
	outstruct.CandidateDescription = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.ImgHash = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.VoteCount = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.Email = *abi.ConvertType(out[4], new(string)).(*string)

	return *outstruct, err

}

// Candidates is a free data retrieval call binding the contract method 0x3477ee2e.
//
// Solidity: function candidates(uint256 ) view returns(string candidate_name, string candidate_description, string imgHash, uint256 voteCount, string email)
func (_Election *ElectionSession) Candidates(arg0 *big.Int) (struct {
	CandidateName        string
	CandidateDescription string
	ImgHash              string
	VoteCount            *big.Int
	Email                string
}, error) {
	return _Election.Contract.Candidates(&_Election.CallOpts, arg0)
}

// Candidates is a free data retrieval call binding the contract method 0x3477ee2e.
//
// Solidity: function candidates(uint256 ) view returns(string candidate_name, string candidate_description, string imgHash, uint256 voteCount, string email)
func (_Election *ElectionCallerSession) Candidates(arg0 *big.Int) (struct {
	CandidateName        string
	CandidateDescription string
	ImgHash              string
	VoteCount            *big.Int
	Email                string
}, error) {
	return _Election.Contract.Candidates(&_Election.CallOpts, arg0)
}

// ElectionAuthority is a free data retrieval call binding the contract method 0x82e15fcd.
//
// Solidity: function election_authority() view returns(address)
func (_Election *ElectionCaller) ElectionAuthority(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Election.contract.Call(opts, &out, "election_authority")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ElectionAuthority is a free data retrieval call binding the contract method 0x82e15fcd.
//
// Solidity: function election_authority() view returns(address)
func (_Election *ElectionSession) ElectionAuthority() (common.Address, error) {
	return _Election.Contract.ElectionAuthority(&_Election.CallOpts)
}

// ElectionAuthority is a free data retrieval call binding the contract method 0x82e15fcd.
//
// Solidity: function election_authority() view returns(address)
func (_Election *ElectionCallerSession) ElectionAuthority() (common.Address, error) {
	return _Election.Contract.ElectionAuthority(&_Election.CallOpts)
}

// ElectionDescription is a free data retrieval call binding the contract method 0x044d5a97.
//
// Solidity: function election_description() view returns(string)
func (_Election *ElectionCaller) ElectionDescription(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Election.contract.Call(opts, &out, "election_description")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// ElectionDescription is a free data retrieval call binding the contract method 0x044d5a97.
//
// Solidity: function election_description() view returns(string)
func (_Election *ElectionSession) ElectionDescription() (string, error) {
	return _Election.Contract.ElectionDescription(&_Election.CallOpts)
}

// ElectionDescription is a free data retrieval call binding the contract method 0x044d5a97.
//
// Solidity: function election_description() view returns(string)
func (_Election *ElectionCallerSession) ElectionDescription() (string, error) {
	return _Election.Contract.ElectionDescription(&_Election.CallOpts)
}

// ElectionName is a free data retrieval call binding the contract method 0xed35a5da.
//
// Solidity: function election_name() view returns(string)
func (_Election *ElectionCaller) ElectionName(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Election.contract.Call(opts, &out, "election_name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// ElectionName is a free data retrieval call binding the contract method 0xed35a5da.
//
// Solidity: function election_name() view returns(string)
func (_Election *ElectionSession) ElectionName() (string, error) {
	return _Election.Contract.ElectionName(&_Election.CallOpts)
}

// ElectionName is a free data retrieval call binding the contract method 0xed35a5da.
//
// Solidity: function election_name() view returns(string)
func (_Election *ElectionCallerSession) ElectionName() (string, error) {
	return _Election.Contract.ElectionName(&_Election.CallOpts)
}

// GetCandidate is a free data retrieval call binding the contract method 0x35b8e820.
//
// Solidity: function getCandidate(uint256 candidateID) view returns(string, string, string, uint256, string)
func (_Election *ElectionCaller) GetCandidate(opts *bind.CallOpts, candidateID *big.Int) (string, string, string, *big.Int, string, error) {
	var out []interface{}
	err := _Election.contract.Call(opts, &out, "getCandidate", candidateID)

	if err != nil {
		return *new(string), *new(string), *new(string), *new(*big.Int), *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	out1 := *abi.ConvertType(out[1], new(string)).(*string)
	out2 := *abi.ConvertType(out[2], new(string)).(*string)
	out3 := *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	out4 := *abi.ConvertType(out[4], new(string)).(*string)

	return out0, out1, out2, out3, out4, err

}

// GetCandidate is a free data retrieval call binding the contract method 0x35b8e820.
//
// Solidity: function getCandidate(uint256 candidateID) view returns(string, string, string, uint256, string)
func (_Election *ElectionSession) GetCandidate(candidateID *big.Int) (string, string, string, *big.Int, string, error) {
	return _Election.Contract.GetCandidate(&_Election.CallOpts, candidateID)
}

// GetCandidate is a free data retrieval call binding the contract method 0x35b8e820.
//
// Solidity: function getCandidate(uint256 candidateID) view returns(string, string, string, uint256, string)
func (_Election *ElectionCallerSession) GetCandidate(candidateID *big.Int) (string, string, string, *big.Int, string, error) {
	return _Election.Contract.GetCandidate(&_Election.CallOpts, candidateID)
}

// GetElectionDetails is a free data retrieval call binding the contract method 0xed836bc3.
//
// Solidity: function getElectionDetails() view returns(string, string)
func (_Election *ElectionCaller) GetElectionDetails(opts *bind.CallOpts) (string, string, error) {
	var out []interface{}
	err := _Election.contract.Call(opts, &out, "getElectionDetails")

	if err != nil {
		return *new(string), *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	out1 := *abi.ConvertType(out[1], new(string)).(*string)

	return out0, out1, err

}

// GetElectionDetails is a free data retrieval call binding the contract method 0xed836bc3.
//
// Solidity: function getElectionDetails() view returns(string, string)
func (_Election *ElectionSession) GetElectionDetails() (string, string, error) {
	return _Election.Contract.GetElectionDetails(&_Election.CallOpts)
}

// GetElectionDetails is a free data retrieval call binding the contract method 0xed836bc3.
//
// Solidity: function getElectionDetails() view returns(string, string)
func (_Election *ElectionCallerSession) GetElectionDetails() (string, string, error) {
	return _Election.Contract.GetElectionDetails(&_Election.CallOpts)
}

// GetNumOfCandidates is a free data retrieval call binding the contract method 0xe8685ba1.
//
// Solidity: function getNumOfCandidates() view returns(uint256)
func (_Election *ElectionCaller) GetNumOfCandidates(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Election.contract.Call(opts, &out, "getNumOfCandidates")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNumOfCandidates is a free data retrieval call binding the contract method 0xe8685ba1.
//
// Solidity: function getNumOfCandidates() view returns(uint256)
func (_Election *ElectionSession) GetNumOfCandidates() (*big.Int, error) {
	return _Election.Contract.GetNumOfCandidates(&_Election.CallOpts)
}

// GetNumOfCandidates is a free data retrieval call binding the contract method 0xe8685ba1.
//
// Solidity: function getNumOfCandidates() view returns(uint256)
func (_Election *ElectionCallerSession) GetNumOfCandidates() (*big.Int, error) {
	return _Election.Contract.GetNumOfCandidates(&_Election.CallOpts)
}

// GetNumOfVoters is a free data retrieval call binding the contract method 0x65fc783c.
//
// Solidity: function getNumOfVoters() view returns(uint256)
func (_Election *ElectionCaller) GetNumOfVoters(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Election.contract.Call(opts, &out, "getNumOfVoters")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNumOfVoters is a free data retrieval call binding the contract method 0x65fc783c.
//
// Solidity: function getNumOfVoters() view returns(uint256)
func (_Election *ElectionSession) GetNumOfVoters() (*big.Int, error) {
	return _Election.Contract.GetNumOfVoters(&_Election.CallOpts)
}

// GetNumOfVoters is a free data retrieval call binding the contract method 0x65fc783c.
//
// Solidity: function getNumOfVoters() view returns(uint256)
func (_Election *ElectionCallerSession) GetNumOfVoters() (*big.Int, error) {
	return _Election.Contract.GetNumOfVoters(&_Election.CallOpts)
}

// GetVoterDetails is a free data retrieval call binding the contract method 0x76874a7d.
//
// Solidity: function getVoterDetails(string email) view returns(uint256, bool)
func (_Election *ElectionCaller) GetVoterDetails(opts *bind.CallOpts, email string) (*big.Int, bool, error) {
	var out []interface{}
	err := _Election.contract.Call(opts, &out, "getVoterDetails", email)

	if err != nil {
		return *new(*big.Int), *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(bool)).(*bool)

	return out0, out1, err

}

// GetVoterDetails is a free data retrieval call binding the contract method 0x76874a7d.
//
// Solidity: function getVoterDetails(string email) view returns(uint256, bool)
func (_Election *ElectionSession) GetVoterDetails(email string) (*big.Int, bool, error) {
	return _Election.Contract.GetVoterDetails(&_Election.CallOpts, email)
}

// GetVoterDetails is a free data retrieval call binding the contract method 0x76874a7d.
//
// Solidity: function getVoterDetails(string email) view returns(uint256, bool)
func (_Election *ElectionCallerSession) GetVoterDetails(email string) (*big.Int, bool, error) {
	return _Election.Contract.GetVoterDetails(&_Election.CallOpts, email)
}

// HasVoted is a free data retrieval call binding the contract method 0x39bfeae3.
//
// Solidity: function hasVoted(string email) view returns(bool)
func (_Election *ElectionCaller) HasVoted(opts *bind.CallOpts, email string) (bool, error) {
	var out []interface{}
	err := _Election.contract.Call(opts, &out, "hasVoted", email)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasVoted is a free data retrieval call binding the contract method 0x39bfeae3.
//
// Solidity: function hasVoted(string email) view returns(bool)
func (_Election *ElectionSession) HasVoted(email string) (bool, error) {
	return _Election.Contract.HasVoted(&_Election.CallOpts, email)
}

// HasVoted is a free data retrieval call binding the contract method 0x39bfeae3.
//
// Solidity: function hasVoted(string email) view returns(bool)
func (_Election *ElectionCallerSession) HasVoted(email string) (bool, error) {
	return _Election.Contract.HasVoted(&_Election.CallOpts, email)
}

// NumCandidates is a free data retrieval call binding the contract method 0x5216509a.
//
// Solidity: function numCandidates() view returns(uint256)
func (_Election *ElectionCaller) NumCandidates(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Election.contract.Call(opts, &out, "numCandidates")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NumCandidates is a free data retrieval call binding the contract method 0x5216509a.
//
// Solidity: function numCandidates() view returns(uint256)
func (_Election *ElectionSession) NumCandidates() (*big.Int, error) {
	return _Election.Contract.NumCandidates(&_Election.CallOpts)
}

// NumCandidates is a free data retrieval call binding the contract method 0x5216509a.
//
// Solidity: function numCandidates() view returns(uint256)
func (_Election *ElectionCallerSession) NumCandidates() (*big.Int, error) {
	return _Election.Contract.NumCandidates(&_Election.CallOpts)
}

// NumVoters is a free data retrieval call binding the contract method 0x4cbe32b8.
//
// Solidity: function numVoters() view returns(uint256)
func (_Election *ElectionCaller) NumVoters(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Election.contract.Call(opts, &out, "numVoters")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NumVoters is a free data retrieval call binding the contract method 0x4cbe32b8.
//
// Solidity: function numVoters() view returns(uint256)
func (_Election *ElectionSession) NumVoters() (*big.Int, error) {
	return _Election.Contract.NumVoters(&_Election.CallOpts)
}

// NumVoters is a free data retrieval call binding the contract method 0x4cbe32b8.
//
// Solidity: function numVoters() view returns(uint256)
func (_Election *ElectionCallerSession) NumVoters() (*big.Int, error) {
	return _Election.Contract.NumVoters(&_Election.CallOpts)
}

// Status is a free data retrieval call binding the contract method 0x200d2ed2.
//
// Solidity: function status() view returns(bool)
func (_Election *ElectionCaller) Status(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Election.contract.Call(opts, &out, "status")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Status is a free data retrieval call binding the contract method 0x200d2ed2.
//
// Solidity: function status() view returns(bool)
func (_Election *ElectionSession) Status() (bool, error) {
	return _Election.Contract.Status(&_Election.CallOpts)
}

// Status is a free data retrieval call binding the contract method 0x200d2ed2.
//
// Solidity: function status() view returns(bool)
func (_Election *ElectionCallerSession) Status() (bool, error) {
	return _Election.Contract.Status(&_Election.CallOpts)
}

// Voters is a free data retrieval call binding the contract method 0x53fa2e64.
//
// Solidity: function voters(string ) view returns(uint256 candidate_id_voted, bool voted)
func (_Election *ElectionCaller) Voters(opts *bind.CallOpts, arg0 string) (struct {
	CandidateIdVoted *big.Int
	Voted            bool
}, error) {
	var out []interface{}
	err := _Election.contract.Call(opts, &out, "voters", arg0)

	outstruct := new(struct {
		CandidateIdVoted *big.Int
		Voted            bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.CandidateIdVoted = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Voted = *abi.ConvertType(out[1], new(bool)).(*bool)

	return *outstruct, err

}

// Voters is a free data retrieval call binding the contract method 0x53fa2e64.
//
// Solidity: function voters(string ) view returns(uint256 candidate_id_voted, bool voted)
func (_Election *ElectionSession) Voters(arg0 string) (struct {
	CandidateIdVoted *big.Int
	Voted            bool
}, error) {
	return _Election.Contract.Voters(&_Election.CallOpts, arg0)
}

// Voters is a free data retrieval call binding the contract method 0x53fa2e64.
//
// Solidity: function voters(string ) view returns(uint256 candidate_id_voted, bool voted)
func (_Election *ElectionCallerSession) Voters(arg0 string) (struct {
	CandidateIdVoted *big.Int
	Voted            bool
}, error) {
	return _Election.Contract.Voters(&_Election.CallOpts, arg0)
}

// WinnerCandidate is a free data retrieval call binding the contract method 0xa15148d1.
//
// Solidity: function winnerCandidate() view returns(uint256)
func (_Election *ElectionCaller) WinnerCandidate(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Election.contract.Call(opts, &out, "winnerCandidate")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WinnerCandidate is a free data retrieval call binding the contract method 0xa15148d1.
//
// Solidity: function winnerCandidate() view returns(uint256)
func (_Election *ElectionSession) WinnerCandidate() (*big.Int, error) {
	return _Election.Contract.WinnerCandidate(&_Election.CallOpts)
}

// WinnerCandidate is a free data retrieval call binding the contract method 0xa15148d1.
//
// Solidity: function winnerCandidate() view returns(uint256)
func (_Election *ElectionCallerSession) WinnerCandidate() (*big.Int, error) {
	return _Election.Contract.WinnerCandidate(&_Election.CallOpts)
}

// AddCandidate is a paid mutator transaction binding the contract method 0x42b03cc9.
//
// Solidity: function addCandidate(string candidate_name, string candidate_description, string imgHash, string email) returns()
func (_Election *ElectionTransactor) AddCandidate(opts *bind.TransactOpts, candidate_name string, candidate_description string, imgHash string, email string) (*types.Transaction, error) {
	return _Election.contract.Transact(opts, "addCandidate", candidate_name, candidate_description, imgHash, email)
}

// AddCandidate is a paid mutator transaction binding the contract method 0x42b03cc9.
//
// Solidity: function addCandidate(string candidate_name, string candidate_description, string imgHash, string email) returns()
func (_Election *ElectionSession) AddCandidate(candidate_name string, candidate_description string, imgHash string, email string) (*types.Transaction, error) {
	return _Election.Contract.AddCandidate(&_Election.TransactOpts, candidate_name, candidate_description, imgHash, email)
}

// AddCandidate is a paid mutator transaction binding the contract method 0x42b03cc9.
//
// Solidity: function addCandidate(string candidate_name, string candidate_description, string imgHash, string email) returns()
func (_Election *ElectionTransactorSession) AddCandidate(candidate_name string, candidate_description string, imgHash string, email string) (*types.Transaction, error) {
	return _Election.Contract.AddCandidate(&_Election.TransactOpts, candidate_name, candidate_description, imgHash, email)
}

// Vote is a paid mutator transaction binding the contract method 0x24108475.
//
// Solidity: function vote(uint256 candidateID, string e) returns()
func (_Election *ElectionTransactor) Vote(opts *bind.TransactOpts, candidateID *big.Int, e string) (*types.Transaction, error) {
	return _Election.contract.Transact(opts, "vote", candidateID, e)
}

// Vote is a paid mutator transaction binding the contract method 0x24108475.
//
// Solidity: function vote(uint256 candidateID, string e) returns()
func (_Election *ElectionSession) Vote(candidateID *big.Int, e string) (*types.Transaction, error) {
	return _Election.Contract.Vote(&_Election.TransactOpts, candidateID, e)
}

// Vote is a paid mutator transaction binding the contract method 0x24108475.
//
// Solidity: function vote(uint256 candidateID, string e) returns()
func (_Election *ElectionTransactorSession) Vote(candidateID *big.Int, e string) (*types.Transaction, error) {
	return _Election.Contract.Vote(&_Election.TransactOpts, candidateID, e)
}
