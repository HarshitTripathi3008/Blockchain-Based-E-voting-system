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

// ElectionFactoryMetaData contains all meta data concerning the ElectionFactory contract.
var ElectionFactoryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"companyEmail\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"deployedAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"el_n\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"el_d\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"email\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"election_name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"election_description\",\"type\":\"string\"}],\"name\":\"createElection\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"email\",\"type\":\"string\"}],\"name\":\"getDeployedElection\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// ElectionFactoryABI is the input ABI used to generate the binding from.
// Deprecated: Use ElectionFactoryMetaData.ABI instead.
var ElectionFactoryABI = ElectionFactoryMetaData.ABI

// ElectionFactory is an auto generated Go binding around an Ethereum contract.
type ElectionFactory struct {
	ElectionFactoryCaller     // Read-only binding to the contract
	ElectionFactoryTransactor // Write-only binding to the contract
	ElectionFactoryFilterer   // Log filterer for contract events
}

// ElectionFactoryCaller is an auto generated read-only Go binding around an Ethereum contract.
type ElectionFactoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ElectionFactoryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ElectionFactoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ElectionFactoryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ElectionFactoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ElectionFactorySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ElectionFactorySession struct {
	Contract     *ElectionFactory  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ElectionFactoryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ElectionFactoryCallerSession struct {
	Contract *ElectionFactoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// ElectionFactoryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ElectionFactoryTransactorSession struct {
	Contract     *ElectionFactoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// ElectionFactoryRaw is an auto generated low-level Go binding around an Ethereum contract.
type ElectionFactoryRaw struct {
	Contract *ElectionFactory // Generic contract binding to access the raw methods on
}

// ElectionFactoryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ElectionFactoryCallerRaw struct {
	Contract *ElectionFactoryCaller // Generic read-only contract binding to access the raw methods on
}

// ElectionFactoryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ElectionFactoryTransactorRaw struct {
	Contract *ElectionFactoryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewElectionFactory creates a new instance of ElectionFactory, bound to a specific deployed contract.
func NewElectionFactory(address common.Address, backend bind.ContractBackend) (*ElectionFactory, error) {
	contract, err := bindElectionFactory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ElectionFactory{ElectionFactoryCaller: ElectionFactoryCaller{contract: contract}, ElectionFactoryTransactor: ElectionFactoryTransactor{contract: contract}, ElectionFactoryFilterer: ElectionFactoryFilterer{contract: contract}}, nil
}

// NewElectionFactoryCaller creates a new read-only instance of ElectionFactory, bound to a specific deployed contract.
func NewElectionFactoryCaller(address common.Address, caller bind.ContractCaller) (*ElectionFactoryCaller, error) {
	contract, err := bindElectionFactory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ElectionFactoryCaller{contract: contract}, nil
}

// NewElectionFactoryTransactor creates a new write-only instance of ElectionFactory, bound to a specific deployed contract.
func NewElectionFactoryTransactor(address common.Address, transactor bind.ContractTransactor) (*ElectionFactoryTransactor, error) {
	contract, err := bindElectionFactory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ElectionFactoryTransactor{contract: contract}, nil
}

// NewElectionFactoryFilterer creates a new log filterer instance of ElectionFactory, bound to a specific deployed contract.
func NewElectionFactoryFilterer(address common.Address, filterer bind.ContractFilterer) (*ElectionFactoryFilterer, error) {
	contract, err := bindElectionFactory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ElectionFactoryFilterer{contract: contract}, nil
}

// bindElectionFactory binds a generic wrapper to an already deployed contract.
func bindElectionFactory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ElectionFactoryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ElectionFactory *ElectionFactoryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ElectionFactory.Contract.ElectionFactoryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ElectionFactory *ElectionFactoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ElectionFactory.Contract.ElectionFactoryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ElectionFactory *ElectionFactoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ElectionFactory.Contract.ElectionFactoryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ElectionFactory *ElectionFactoryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ElectionFactory.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ElectionFactory *ElectionFactoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ElectionFactory.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ElectionFactory *ElectionFactoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ElectionFactory.Contract.contract.Transact(opts, method, params...)
}

// CompanyEmail is a free data retrieval call binding the contract method 0xd1642a8a.
//
// Solidity: function companyEmail(string ) view returns(address deployedAddress, string el_n, string el_d)
func (_ElectionFactory *ElectionFactoryCaller) CompanyEmail(opts *bind.CallOpts, arg0 string) (struct {
	DeployedAddress common.Address
	ElN             string
	ElD             string
}, error) {
	var out []interface{}
	err := _ElectionFactory.contract.Call(opts, &out, "companyEmail", arg0)

	outstruct := new(struct {
		DeployedAddress common.Address
		ElN             string
		ElD             string
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.DeployedAddress = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.ElN = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.ElD = *abi.ConvertType(out[2], new(string)).(*string)

	return *outstruct, err

}

// CompanyEmail is a free data retrieval call binding the contract method 0xd1642a8a.
//
// Solidity: function companyEmail(string ) view returns(address deployedAddress, string el_n, string el_d)
func (_ElectionFactory *ElectionFactorySession) CompanyEmail(arg0 string) (struct {
	DeployedAddress common.Address
	ElN             string
	ElD             string
}, error) {
	return _ElectionFactory.Contract.CompanyEmail(&_ElectionFactory.CallOpts, arg0)
}

// CompanyEmail is a free data retrieval call binding the contract method 0xd1642a8a.
//
// Solidity: function companyEmail(string ) view returns(address deployedAddress, string el_n, string el_d)
func (_ElectionFactory *ElectionFactoryCallerSession) CompanyEmail(arg0 string) (struct {
	DeployedAddress common.Address
	ElN             string
	ElD             string
}, error) {
	return _ElectionFactory.Contract.CompanyEmail(&_ElectionFactory.CallOpts, arg0)
}

// GetDeployedElection is a free data retrieval call binding the contract method 0x37f46549.
//
// Solidity: function getDeployedElection(string email) view returns(address, string, string)
func (_ElectionFactory *ElectionFactoryCaller) GetDeployedElection(opts *bind.CallOpts, email string) (common.Address, string, string, error) {
	var out []interface{}
	err := _ElectionFactory.contract.Call(opts, &out, "getDeployedElection", email)

	if err != nil {
		return *new(common.Address), *new(string), *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	out1 := *abi.ConvertType(out[1], new(string)).(*string)
	out2 := *abi.ConvertType(out[2], new(string)).(*string)

	return out0, out1, out2, err

}

// GetDeployedElection is a free data retrieval call binding the contract method 0x37f46549.
//
// Solidity: function getDeployedElection(string email) view returns(address, string, string)
func (_ElectionFactory *ElectionFactorySession) GetDeployedElection(email string) (common.Address, string, string, error) {
	return _ElectionFactory.Contract.GetDeployedElection(&_ElectionFactory.CallOpts, email)
}

// GetDeployedElection is a free data retrieval call binding the contract method 0x37f46549.
//
// Solidity: function getDeployedElection(string email) view returns(address, string, string)
func (_ElectionFactory *ElectionFactoryCallerSession) GetDeployedElection(email string) (common.Address, string, string, error) {
	return _ElectionFactory.Contract.GetDeployedElection(&_ElectionFactory.CallOpts, email)
}

// CreateElection is a paid mutator transaction binding the contract method 0x4a6cf3a9.
//
// Solidity: function createElection(string email, string election_name, string election_description) returns()
func (_ElectionFactory *ElectionFactoryTransactor) CreateElection(opts *bind.TransactOpts, email string, election_name string, election_description string) (*types.Transaction, error) {
	return _ElectionFactory.contract.Transact(opts, "createElection", email, election_name, election_description)
}

// CreateElection is a paid mutator transaction binding the contract method 0x4a6cf3a9.
//
// Solidity: function createElection(string email, string election_name, string election_description) returns()
func (_ElectionFactory *ElectionFactorySession) CreateElection(email string, election_name string, election_description string) (*types.Transaction, error) {
	return _ElectionFactory.Contract.CreateElection(&_ElectionFactory.TransactOpts, email, election_name, election_description)
}

// CreateElection is a paid mutator transaction binding the contract method 0x4a6cf3a9.
//
// Solidity: function createElection(string email, string election_name, string election_description) returns()
func (_ElectionFactory *ElectionFactoryTransactorSession) CreateElection(email string, election_name string, election_description string) (*types.Transaction, error) {
	return _ElectionFactory.Contract.CreateElection(&_ElectionFactory.TransactOpts, email, election_name, election_description)
}
