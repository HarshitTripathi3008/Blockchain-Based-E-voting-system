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

// BindingsMetaData contains all meta data concerning the Bindings contract.
var BindingsMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"admin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_electionAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_title\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_winnerName\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_winningVotes\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_totalVoters\",\"type\":\"uint256\"}],\"name\":\"archiveResult\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"archivedResults\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"electionAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"title\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"winnerName\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"winningVotes\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalVoters\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561000f575f80fd5b50335f806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550610b168061005c5f395ff3fe608060405234801561000f575f80fd5b506004361061003f575f3560e01c806331bd3af7146100435780635c22ed3a14610078578063f851a44014610094575b5f80fd5b61005d60048036038101906100589190610448565b6100b2565b60405161006f96959493929190610524565b60405180910390f35b610092600480360381019061008d91906106e7565b610215565b005b61009c6103ba565b6040516100a99190610796565b60405180910390f35b6001602052805f5260405f205f91509050805f015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060010180546100f6906107dc565b80601f0160208091040260200160405190810160405280929190818152602001828054610122906107dc565b801561016d5780601f106101445761010080835404028352916020019161016d565b820191905f5260205f20905b81548152906001019060200180831161015057829003601f168201915b505050505090806002018054610182906107dc565b80601f01602080910402602001604051908101604052809291908181526020018280546101ae906107dc565b80156101f95780601f106101d0576101008083540402835291602001916101f9565b820191905f5260205f20905b8154815290600101906020018083116101dc57829003601f168201915b5050505050908060030154908060040154908060050154905086565b5f8054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146102a2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161029990610856565b60405180910390fd5b6040518060c001604052808673ffffffffffffffffffffffffffffffffffffffff1681526020018581526020018481526020018381526020018281526020014281525060015f8773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f820151815f015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550602082015181600101908161037b9190610a11565b5060408201518160020190816103919190610a11565b50606082015181600301556080820151816004015560a082015181600501559050505050505050565b5f8054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b5f604051905090565b5f80fd5b5f80fd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f610417826103ee565b9050919050565b6104278161040d565b8114610431575f80fd5b50565b5f813590506104428161041e565b92915050565b5f6020828403121561045d5761045c6103e6565b5b5f61046a84828501610434565b91505092915050565b61047c8161040d565b82525050565b5f81519050919050565b5f82825260208201905092915050565b5f5b838110156104b957808201518184015260208101905061049e565b5f8484015250505050565b5f601f19601f8301169050919050565b5f6104de82610482565b6104e8818561048c565b93506104f881856020860161049c565b610501816104c4565b840191505092915050565b5f819050919050565b61051e8161050c565b82525050565b5f60c0820190506105375f830189610473565b818103602083015261054981886104d4565b9050818103604083015261055d81876104d4565b905061056c6060830186610515565b6105796080830185610515565b61058660a0830184610515565b979650505050505050565b5f80fd5b5f80fd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b6105cf826104c4565b810181811067ffffffffffffffff821117156105ee576105ed610599565b5b80604052505050565b5f6106006103dd565b905061060c82826105c6565b919050565b5f67ffffffffffffffff82111561062b5761062a610599565b5b610634826104c4565b9050602081019050919050565b828183375f83830152505050565b5f61066161065c84610611565b6105f7565b90508281526020810184848401111561067d5761067c610595565b5b610688848285610641565b509392505050565b5f82601f8301126106a4576106a3610591565b5b81356106b484826020860161064f565b91505092915050565b6106c68161050c565b81146106d0575f80fd5b50565b5f813590506106e1816106bd565b92915050565b5f805f805f60a08688031215610700576106ff6103e6565b5b5f61070d88828901610434565b955050602086013567ffffffffffffffff81111561072e5761072d6103ea565b5b61073a88828901610690565b945050604086013567ffffffffffffffff81111561075b5761075a6103ea565b5b61076788828901610690565b9350506060610778888289016106d3565b9250506080610789888289016106d3565b9150509295509295909350565b5f6020820190506107a95f830184610473565b92915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f60028204905060018216806107f357607f821691505b602082108103610806576108056107af565b5b50919050565b7f4f6e6c792061646d696e2063616e206172636869766520726573756c747300005f82015250565b5f610840601e8361048c565b915061084b8261080c565b602082019050919050565b5f6020820190508181035f83015261086d81610834565b9050919050565b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f600883026108d07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82610895565b6108da8683610895565b95508019841693508086168417925050509392505050565b5f819050919050565b5f61091561091061090b8461050c565b6108f2565b61050c565b9050919050565b5f819050919050565b61092e836108fb565b61094261093a8261091c565b8484546108a1565b825550505050565b5f90565b61095661094a565b610961818484610925565b505050565b5b81811015610984576109795f8261094e565b600181019050610967565b5050565b601f8211156109c95761099a81610874565b6109a384610886565b810160208510156109b2578190505b6109c66109be85610886565b830182610966565b50505b505050565b5f82821c905092915050565b5f6109e95f19846008026109ce565b1980831691505092915050565b5f610a0183836109da565b9150826002028217905092915050565b610a1a82610482565b67ffffffffffffffff811115610a3357610a32610599565b5b610a3d82546107dc565b610a48828285610988565b5f60209050601f831160018114610a79575f8415610a67578287015190505b610a7185826109f6565b865550610ad8565b601f198416610a8786610874565b5f5b82811015610aae57848901518255600182019150602085019450602081019050610a89565b86831015610acb5784890151610ac7601f8916826109da565b8355505b6001600288020188555050505b50505050505056fea264697066735822122073096dab2fb05a2f4b567d68c17d697b9605be2d51c98cdc616df63292fc3df364736f6c63430008140033",
}

// BindingsABI is the input ABI used to generate the binding from.
// Deprecated: Use BindingsMetaData.ABI instead.
var BindingsABI = BindingsMetaData.ABI

// BindingsBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use BindingsMetaData.Bin instead.
var BindingsBin = BindingsMetaData.Bin

// DeployBindings deploys a new Ethereum contract, binding an instance of Bindings to it.
func DeployBindings(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Bindings, error) {
	parsed, err := BindingsMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BindingsBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Bindings{BindingsCaller: BindingsCaller{contract: contract}, BindingsTransactor: BindingsTransactor{contract: contract}, BindingsFilterer: BindingsFilterer{contract: contract}}, nil
}

// Bindings is an auto generated Go binding around an Ethereum contract.
type Bindings struct {
	BindingsCaller     // Read-only binding to the contract
	BindingsTransactor // Write-only binding to the contract
	BindingsFilterer   // Log filterer for contract events
}

// BindingsCaller is an auto generated read-only Go binding around an Ethereum contract.
type BindingsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BindingsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BindingsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BindingsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BindingsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BindingsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BindingsSession struct {
	Contract     *Bindings         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BindingsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BindingsCallerSession struct {
	Contract *BindingsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// BindingsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BindingsTransactorSession struct {
	Contract     *BindingsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// BindingsRaw is an auto generated low-level Go binding around an Ethereum contract.
type BindingsRaw struct {
	Contract *Bindings // Generic contract binding to access the raw methods on
}

// BindingsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BindingsCallerRaw struct {
	Contract *BindingsCaller // Generic read-only contract binding to access the raw methods on
}

// BindingsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BindingsTransactorRaw struct {
	Contract *BindingsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBindings creates a new instance of Bindings, bound to a specific deployed contract.
func NewBindings(address common.Address, backend bind.ContractBackend) (*Bindings, error) {
	contract, err := bindBindings(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Bindings{BindingsCaller: BindingsCaller{contract: contract}, BindingsTransactor: BindingsTransactor{contract: contract}, BindingsFilterer: BindingsFilterer{contract: contract}}, nil
}

// NewBindingsCaller creates a new read-only instance of Bindings, bound to a specific deployed contract.
func NewBindingsCaller(address common.Address, caller bind.ContractCaller) (*BindingsCaller, error) {
	contract, err := bindBindings(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BindingsCaller{contract: contract}, nil
}

// NewBindingsTransactor creates a new write-only instance of Bindings, bound to a specific deployed contract.
func NewBindingsTransactor(address common.Address, transactor bind.ContractTransactor) (*BindingsTransactor, error) {
	contract, err := bindBindings(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BindingsTransactor{contract: contract}, nil
}

// NewBindingsFilterer creates a new log filterer instance of Bindings, bound to a specific deployed contract.
func NewBindingsFilterer(address common.Address, filterer bind.ContractFilterer) (*BindingsFilterer, error) {
	contract, err := bindBindings(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BindingsFilterer{contract: contract}, nil
}

// bindBindings binds a generic wrapper to an already deployed contract.
func bindBindings(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BindingsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Bindings *BindingsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Bindings.Contract.BindingsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Bindings *BindingsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bindings.Contract.BindingsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Bindings *BindingsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Bindings.Contract.BindingsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Bindings *BindingsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Bindings.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Bindings *BindingsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bindings.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Bindings *BindingsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Bindings.Contract.contract.Transact(opts, method, params...)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_Bindings *BindingsCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Bindings.contract.Call(opts, &out, "admin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_Bindings *BindingsSession) Admin() (common.Address, error) {
	return _Bindings.Contract.Admin(&_Bindings.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_Bindings *BindingsCallerSession) Admin() (common.Address, error) {
	return _Bindings.Contract.Admin(&_Bindings.CallOpts)
}

// ArchivedResults is a free data retrieval call binding the contract method 0x31bd3af7.
//
// Solidity: function archivedResults(address ) view returns(address electionAddress, string title, string winnerName, uint256 winningVotes, uint256 totalVoters, uint256 timestamp)
func (_Bindings *BindingsCaller) ArchivedResults(opts *bind.CallOpts, arg0 common.Address) (struct {
	ElectionAddress common.Address
	Title           string
	WinnerName      string
	WinningVotes    *big.Int
	TotalVoters     *big.Int
	Timestamp       *big.Int
}, error) {
	var out []interface{}
	err := _Bindings.contract.Call(opts, &out, "archivedResults", arg0)

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
func (_Bindings *BindingsSession) ArchivedResults(arg0 common.Address) (struct {
	ElectionAddress common.Address
	Title           string
	WinnerName      string
	WinningVotes    *big.Int
	TotalVoters     *big.Int
	Timestamp       *big.Int
}, error) {
	return _Bindings.Contract.ArchivedResults(&_Bindings.CallOpts, arg0)
}

// ArchivedResults is a free data retrieval call binding the contract method 0x31bd3af7.
//
// Solidity: function archivedResults(address ) view returns(address electionAddress, string title, string winnerName, uint256 winningVotes, uint256 totalVoters, uint256 timestamp)
func (_Bindings *BindingsCallerSession) ArchivedResults(arg0 common.Address) (struct {
	ElectionAddress common.Address
	Title           string
	WinnerName      string
	WinningVotes    *big.Int
	TotalVoters     *big.Int
	Timestamp       *big.Int
}, error) {
	return _Bindings.Contract.ArchivedResults(&_Bindings.CallOpts, arg0)
}

// ArchiveResult is a paid mutator transaction binding the contract method 0x5c22ed3a.
//
// Solidity: function archiveResult(address _electionAddress, string _title, string _winnerName, uint256 _winningVotes, uint256 _totalVoters) returns()
func (_Bindings *BindingsTransactor) ArchiveResult(opts *bind.TransactOpts, _electionAddress common.Address, _title string, _winnerName string, _winningVotes *big.Int, _totalVoters *big.Int) (*types.Transaction, error) {
	return _Bindings.contract.Transact(opts, "archiveResult", _electionAddress, _title, _winnerName, _winningVotes, _totalVoters)
}

// ArchiveResult is a paid mutator transaction binding the contract method 0x5c22ed3a.
//
// Solidity: function archiveResult(address _electionAddress, string _title, string _winnerName, uint256 _winningVotes, uint256 _totalVoters) returns()
func (_Bindings *BindingsSession) ArchiveResult(_electionAddress common.Address, _title string, _winnerName string, _winningVotes *big.Int, _totalVoters *big.Int) (*types.Transaction, error) {
	return _Bindings.Contract.ArchiveResult(&_Bindings.TransactOpts, _electionAddress, _title, _winnerName, _winningVotes, _totalVoters)
}

// ArchiveResult is a paid mutator transaction binding the contract method 0x5c22ed3a.
//
// Solidity: function archiveResult(address _electionAddress, string _title, string _winnerName, uint256 _winningVotes, uint256 _totalVoters) returns()
func (_Bindings *BindingsTransactorSession) ArchiveResult(_electionAddress common.Address, _title string, _winnerName string, _winningVotes *big.Int, _totalVoters *big.Int) (*types.Transaction, error) {
	return _Bindings.Contract.ArchiveResult(&_Bindings.TransactOpts, _electionAddress, _title, _winnerName, _winningVotes, _totalVoters)
}
