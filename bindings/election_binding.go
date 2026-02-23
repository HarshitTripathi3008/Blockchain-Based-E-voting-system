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
	Bin: "0x60806040523480156200001157600080fd5b506040516200157738038062001577833981016040819052620000349162000149565b600080546001600160a01b0319166001600160a01b03851617905560016200005d838262000262565b5060026200006c828262000262565b50506003805460ff19166001179055506200032e9050565b634e487b7160e01b600052604160045260246000fd5b600082601f830112620000ac57600080fd5b81516001600160401b0380821115620000c957620000c962000084565b604051601f8301601f19908116603f01168101908282118183101715620000f457620000f462000084565b816040528381526020925086838588010111156200011157600080fd5b600091505b8382101562000135578582018301518183018401529082019062000116565b600093810190920192909252949350505050565b6000806000606084860312156200015f57600080fd5b83516001600160a01b03811681146200017757600080fd5b60208501519093506001600160401b03808211156200019557600080fd5b620001a3878388016200009a565b93506040860151915080821115620001ba57600080fd5b50620001c9868287016200009a565b9150509250925092565b600181811c90821680620001e857607f821691505b6020821081036200020957634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200025d57600081815260208120601f850160051c81016020861015620002385750805b601f850160051c820191505b81811015620002595782815560010162000244565b5050505b505050565b81516001600160401b038111156200027e576200027e62000084565b62000296816200028f8454620001d3565b846200020f565b602080601f831160018114620002ce5760008415620002b55750858301515b600019600386901b1c1916600185901b17855562000259565b600085815260208120601f198616915b82811015620002ff57888601518255948401946001909101908401620002de565b50858210156200031e5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b611239806200033e6000396000f3fe608060405234801561001057600080fd5b506004361061010b5760003560e01c80635216509a116100a257806382e15fcd1161007157806382e15fcd14610242578063a15148d11461026d578063e8685ba114610275578063ed35a5da1461027d578063ed836bc31461028557600080fd5b80635216509a146101d457806353fa2e64146101dd57806365fc783c1461022757806376874a7d1461022f57600080fd5b806335b8e820116100de57806335b8e8201461018457806339bfeae31461019757806342b03cc9146101aa5780634cbe32b8146101bd57600080fd5b8063044d5a9714610110578063200d2ed21461012e578063241084751461014b5780633477ee2e14610160575b600080fd5b61011861029b565b6040516101259190610db4565b60405180910390f35b60035461013b9060ff1681565b6040519015158152602001610125565b61015e610159366004610e71565b610329565b005b61017361016e366004610eb8565b6104a9565b604051610125959493929190610ed1565b610173610192366004610eb8565b6106f7565b61013b6101a5366004610f30565b6109e5565b61015e6101b8366004610f6d565b610a13565b6101c660075481565b604051908152602001610125565b6101c660065481565b6102126101eb366004610f30565b80516020818301810180516005825292820191909301209152805460019091015460ff1682565b60408051928352901515602083015201610125565b6007546101c6565b61021261023d366004610f30565b610aed565b600054610255906001600160a01b031681565b6040516001600160a01b039091168152602001610125565b6101c6610b38565b6006546101c6565b610118610c30565b61028d610c3d565b60405161012592919061101a565b600280546102a890611048565b80601f01602080910402602001604051908101604052809291908181526020018280546102d490611048565b80156103215780601f106102f657610100808354040283529160200191610321565b820191906000526020600020905b81548152906001019060200180831161030457829003601f168201915b505050505081565b6000546001600160a01b0316331461035c5760405162461bcd60e51b815260040161035390611082565b60405180910390fd5b60058160405161036c91906110b1565b9081526040519081900360200190206001015460ff16156103cf5760405162461bcd60e51b815260206004820152601d60248201527f4572726f723a20596f752063616e6e6f7420646f75626c6520766f74650000006044820152606401610353565b60065482106104205760405162461bcd60e51b815260206004820152601b60248201527f4572726f723a20496e76616c69642063616e64696461746520494400000000006044820152606401610353565b6040805180820182528381526001602082015290516005906104439084906110b1565b90815260405160209181900382019020825181559101516001909101805460ff19169115159190911790556007805490600061047e836110cd565b909155505060008281526004602052604081206003018054916104a0836110cd565b91905055505050565b6004602052600090815260409020805481906104c490611048565b80601f01602080910402602001604051908101604052809291908181526020018280546104f090611048565b801561053d5780601f106105125761010080835404028352916020019161053d565b820191906000526020600020905b81548152906001019060200180831161052057829003601f168201915b50505050509080600101805461055290611048565b80601f016020809104026020016040519081016040528092919081815260200182805461057e90611048565b80156105cb5780601f106105a0576101008083540402835291602001916105cb565b820191906000526020600020905b8154815290600101906020018083116105ae57829003601f168201915b5050505050908060020180546105e090611048565b80601f016020809104026020016040519081016040528092919081815260200182805461060c90611048565b80156106595780601f1061062e57610100808354040283529160200191610659565b820191906000526020600020905b81548152906001019060200180831161063c57829003601f168201915b50505050509080600301549080600401805461067490611048565b80601f01602080910402602001604051908101604052809291908181526020018280546106a090611048565b80156106ed5780601f106106c2576101008083540402835291602001916106ed565b820191906000526020600020905b8154815290600101906020018083116106d057829003601f168201915b5050505050905085565b60608060606000606060065486106107515760405162461bcd60e51b815260206004820152601b60248201527f4572726f723a20496e76616c69642063616e64696461746520494400000000006044820152606401610353565b600086815260046020526040808220815160a0810190925280548290829061077890611048565b80601f01602080910402602001604051908101604052809291908181526020018280546107a490611048565b80156107f15780601f106107c6576101008083540402835291602001916107f1565b820191906000526020600020905b8154815290600101906020018083116107d457829003601f168201915b5050505050815260200160018201805461080a90611048565b80601f016020809104026020016040519081016040528092919081815260200182805461083690611048565b80156108835780601f1061085857610100808354040283529160200191610883565b820191906000526020600020905b81548152906001019060200180831161086657829003601f168201915b5050505050815260200160028201805461089c90611048565b80601f01602080910402602001604051908101604052809291908181526020018280546108c890611048565b80156109155780601f106108ea57610100808354040283529160200191610915565b820191906000526020600020905b8154815290600101906020018083116108f857829003601f168201915b505050505081526020016003820154815260200160048201805461093890611048565b80601f016020809104026020016040519081016040528092919081815260200182805461096490611048565b80156109b15780601f10610986576101008083540402835291602001916109b1565b820191906000526020600020905b81548152906001019060200180831161099457829003601f168201915b5050509190925250508151602083015160408401516060850151608090950151929c919b5099509297509550909350505050565b60006005826040516109f791906110b1565b9081526040519081900360200190206001015460ff1692915050565b6000546001600160a01b03163314610a3d5760405162461bcd60e51b815260040161035390611082565b6006546040805160a08101825286815260208082018790528183018690526000606083018190526080830186905284815260049091529190912081518190610a859082611143565b5060208201516001820190610a9a9082611143565b5060408201516002820190610aaf9082611143565b506060820151600382015560808201516004820190610ace9082611143565b50506006805491506000610ae1836110cd565b91905055505050505050565b6000806000600584604051610b0291906110b1565b9081526040805160209281900383018120818301909252815480825260019092015460ff16151592018290529590945092505050565b600080546001600160a01b03163314610b635760405162461bcd60e51b815260040161035390611082565b600060065411610bac5760405162461bcd60e51b81526020600482015260146024820152734572726f723a204e6f2063616e6469646174657360601b6044820152606401610353565b600080805260046020527f17ef568e3e12ab5b9c7254a8d58478811de00f9e6eb34345acd53bf8fd09d3ef549060015b600654811015610c2957600081815260046020526040902060030154831015610c175760008181526004602052604090206003015492509050805b80610c21816110cd565b915050610bdc565b5091505090565b600180546102a890611048565b60608060016002818054610c5090611048565b80601f0160208091040260200160405190810160405280929190818152602001828054610c7c90611048565b8015610cc95780601f10610c9e57610100808354040283529160200191610cc9565b820191906000526020600020905b815481529060010190602001808311610cac57829003601f168201915b50505050509150808054610cdc90611048565b80601f0160208091040260200160405190810160405280929190818152602001828054610d0890611048565b8015610d555780601f10610d2a57610100808354040283529160200191610d55565b820191906000526020600020905b815481529060010190602001808311610d3857829003601f168201915b50505050509050915091509091565b60005b83811015610d7f578181015183820152602001610d67565b50506000910152565b60008151808452610da0816020860160208601610d64565b601f01601f19169290920160200192915050565b602081526000610dc76020830184610d88565b9392505050565b634e487b7160e01b600052604160045260246000fd5b600082601f830112610df557600080fd5b813567ffffffffffffffff80821115610e1057610e10610dce565b604051601f8301601f19908116603f01168101908282118183101715610e3857610e38610dce565b81604052838152866020858801011115610e5157600080fd5b836020870160208301376000602085830101528094505050505092915050565b60008060408385031215610e8457600080fd5b82359150602083013567ffffffffffffffff811115610ea257600080fd5b610eae85828601610de4565b9150509250929050565b600060208284031215610eca57600080fd5b5035919050565b60a081526000610ee460a0830188610d88565b8281036020840152610ef68188610d88565b90508281036040840152610f0a8187610d88565b90508460608401528281036080840152610f248185610d88565b98975050505050505050565b600060208284031215610f4257600080fd5b813567ffffffffffffffff811115610f5957600080fd5b610f6584828501610de4565b949350505050565b60008060008060808587031215610f8357600080fd5b843567ffffffffffffffff80821115610f9b57600080fd5b610fa788838901610de4565b95506020870135915080821115610fbd57600080fd5b610fc988838901610de4565b94506040870135915080821115610fdf57600080fd5b610feb88838901610de4565b9350606087013591508082111561100157600080fd5b5061100e87828801610de4565b91505092959194509250565b60408152600061102d6040830185610d88565b828103602084015261103f8185610d88565b95945050505050565b600181811c9082168061105c57607f821691505b60208210810361107c57634e487b7160e01b600052602260045260246000fd5b50919050565b60208082526015908201527422b93937b91d1020b1b1b2b9b9902232b734b2b21760591b604082015260600190565b600082516110c3818460208701610d64565b9190910192915050565b6000600182016110ed57634e487b7160e01b600052601160045260246000fd5b5060010190565b601f82111561113e57600081815260208120601f850160051c8101602086101561111b5750805b601f850160051c820191505b8181101561113a57828155600101611127565b5050505b505050565b815167ffffffffffffffff81111561115d5761115d610dce565b6111718161116b8454611048565b846110f4565b602080601f8311600181146111a6576000841561118e5750858301515b600019600386901b1c1916600185901b17855561113a565b600085815260208120601f198616915b828110156111d5578886015182559484019460019091019084016111b6565b50858210156111f35787850151600019600388901b60f8161c191681555b5050505050600190811b0190555056fea2646970667358221220aa1e9f032f28909aeb383b4cc8c8c079b6b944aa40854fe4102d7c7963aa4cf364736f6c63430008130033",
}

// ElectionABI is the input ABI used to generate the binding from.
// Deprecated: Use ElectionMetaData.ABI instead.
var ElectionABI = ElectionMetaData.ABI

// ElectionBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ElectionMetaData.Bin instead.
var ElectionBin = ElectionMetaData.Bin

// DeployElection deploys a new Ethereum contract, binding an instance of Election to it.
func DeployElection(auth *bind.TransactOpts, backend bind.ContractBackend, authority common.Address, name string, description string) (common.Address, *types.Transaction, *Election, error) {
	parsed, err := ElectionMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ElectionBin), backend, authority, name, description)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Election{ElectionCaller: ElectionCaller{contract: contract}, ElectionTransactor: ElectionTransactor{contract: contract}, ElectionFilterer: ElectionFilterer{contract: contract}}, nil
}

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
