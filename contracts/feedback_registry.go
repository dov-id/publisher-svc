// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

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

// FeedbackRegistryMetaData contains all meta data concerning the FeedbackRegistry contract.
var FeedbackRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"certIntegrator_\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"course_\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"i_\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"c_\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"r_\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"publicKeysX_\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"publicKeysY_\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes32[][]\",\"name\":\"merkleTreeProofs_\",\"type\":\"bytes32[][]\"},{\"internalType\":\"bytes32[]\",\"name\":\"keys_\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"values_\",\"type\":\"bytes32[]\"},{\"internalType\":\"string\",\"name\":\"ipfsHash_\",\"type\":\"string\"}],\"name\":\"addFeedback\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"contractFeedbacks\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllFeedbacks\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"courses_\",\"type\":\"bytes[]\"},{\"internalType\":\"string[][]\",\"name\":\"feedbacks_\",\"type\":\"string[][]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"course_\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"offset_\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"limit_\",\"type\":\"uint256\"}],\"name\":\"getFeedbacks\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// FeedbackRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use FeedbackRegistryMetaData.ABI instead.
var FeedbackRegistryABI = FeedbackRegistryMetaData.ABI

// FeedbackRegistry is an auto generated Go binding around an Ethereum contract.
type FeedbackRegistry struct {
	FeedbackRegistryCaller     // Read-only binding to the contract
	FeedbackRegistryTransactor // Write-only binding to the contract
	FeedbackRegistryFilterer   // Log filterer for contract events
}

// FeedbackRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type FeedbackRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeedbackRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FeedbackRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeedbackRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FeedbackRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeedbackRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FeedbackRegistrySession struct {
	Contract     *FeedbackRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FeedbackRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FeedbackRegistryCallerSession struct {
	Contract *FeedbackRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// FeedbackRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FeedbackRegistryTransactorSession struct {
	Contract     *FeedbackRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// FeedbackRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type FeedbackRegistryRaw struct {
	Contract *FeedbackRegistry // Generic contract binding to access the raw methods on
}

// FeedbackRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FeedbackRegistryCallerRaw struct {
	Contract *FeedbackRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// FeedbackRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FeedbackRegistryTransactorRaw struct {
	Contract *FeedbackRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFeedbackRegistry creates a new instance of FeedbackRegistry, bound to a specific deployed contract.
func NewFeedbackRegistry(address common.Address, backend bind.ContractBackend) (*FeedbackRegistry, error) {
	contract, err := bindFeedbackRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FeedbackRegistry{FeedbackRegistryCaller: FeedbackRegistryCaller{contract: contract}, FeedbackRegistryTransactor: FeedbackRegistryTransactor{contract: contract}, FeedbackRegistryFilterer: FeedbackRegistryFilterer{contract: contract}}, nil
}

// NewFeedbackRegistryCaller creates a new read-only instance of FeedbackRegistry, bound to a specific deployed contract.
func NewFeedbackRegistryCaller(address common.Address, caller bind.ContractCaller) (*FeedbackRegistryCaller, error) {
	contract, err := bindFeedbackRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FeedbackRegistryCaller{contract: contract}, nil
}

// NewFeedbackRegistryTransactor creates a new write-only instance of FeedbackRegistry, bound to a specific deployed contract.
func NewFeedbackRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*FeedbackRegistryTransactor, error) {
	contract, err := bindFeedbackRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FeedbackRegistryTransactor{contract: contract}, nil
}

// NewFeedbackRegistryFilterer creates a new log filterer instance of FeedbackRegistry, bound to a specific deployed contract.
func NewFeedbackRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*FeedbackRegistryFilterer, error) {
	contract, err := bindFeedbackRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FeedbackRegistryFilterer{contract: contract}, nil
}

// bindFeedbackRegistry binds a generic wrapper to an already deployed contract.
func bindFeedbackRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FeedbackRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FeedbackRegistry *FeedbackRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FeedbackRegistry.Contract.FeedbackRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FeedbackRegistry *FeedbackRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FeedbackRegistry.Contract.FeedbackRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FeedbackRegistry *FeedbackRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FeedbackRegistry.Contract.FeedbackRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FeedbackRegistry *FeedbackRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FeedbackRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FeedbackRegistry *FeedbackRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FeedbackRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FeedbackRegistry *FeedbackRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FeedbackRegistry.Contract.contract.Transact(opts, method, params...)
}

// ContractFeedbacks is a free data retrieval call binding the contract method 0x476199db.
//
// Solidity: function contractFeedbacks(bytes , uint256 ) view returns(string)
func (_FeedbackRegistry *FeedbackRegistryCaller) ContractFeedbacks(opts *bind.CallOpts, arg0 []byte, arg1 *big.Int) (string, error) {
	var out []interface{}
	err := _FeedbackRegistry.contract.Call(opts, &out, "contractFeedbacks", arg0, arg1)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// ContractFeedbacks is a free data retrieval call binding the contract method 0x476199db.
//
// Solidity: function contractFeedbacks(bytes , uint256 ) view returns(string)
func (_FeedbackRegistry *FeedbackRegistrySession) ContractFeedbacks(arg0 []byte, arg1 *big.Int) (string, error) {
	return _FeedbackRegistry.Contract.ContractFeedbacks(&_FeedbackRegistry.CallOpts, arg0, arg1)
}

// ContractFeedbacks is a free data retrieval call binding the contract method 0x476199db.
//
// Solidity: function contractFeedbacks(bytes , uint256 ) view returns(string)
func (_FeedbackRegistry *FeedbackRegistryCallerSession) ContractFeedbacks(arg0 []byte, arg1 *big.Int) (string, error) {
	return _FeedbackRegistry.Contract.ContractFeedbacks(&_FeedbackRegistry.CallOpts, arg0, arg1)
}

// GetAllFeedbacks is a free data retrieval call binding the contract method 0x5fa01010.
//
// Solidity: function getAllFeedbacks() view returns(bytes[] courses_, string[][] feedbacks_)
func (_FeedbackRegistry *FeedbackRegistryCaller) GetAllFeedbacks(opts *bind.CallOpts) (struct {
	Courses   [][]byte
	Feedbacks [][]string
}, error) {
	var out []interface{}
	err := _FeedbackRegistry.contract.Call(opts, &out, "getAllFeedbacks")

	outstruct := new(struct {
		Courses   [][]byte
		Feedbacks [][]string
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Courses = *abi.ConvertType(out[0], new([][]byte)).(*[][]byte)
	outstruct.Feedbacks = *abi.ConvertType(out[1], new([][]string)).(*[][]string)

	return *outstruct, err

}

// GetAllFeedbacks is a free data retrieval call binding the contract method 0x5fa01010.
//
// Solidity: function getAllFeedbacks() view returns(bytes[] courses_, string[][] feedbacks_)
func (_FeedbackRegistry *FeedbackRegistrySession) GetAllFeedbacks() (struct {
	Courses   [][]byte
	Feedbacks [][]string
}, error) {
	return _FeedbackRegistry.Contract.GetAllFeedbacks(&_FeedbackRegistry.CallOpts)
}

// GetAllFeedbacks is a free data retrieval call binding the contract method 0x5fa01010.
//
// Solidity: function getAllFeedbacks() view returns(bytes[] courses_, string[][] feedbacks_)
func (_FeedbackRegistry *FeedbackRegistryCallerSession) GetAllFeedbacks() (struct {
	Courses   [][]byte
	Feedbacks [][]string
}, error) {
	return _FeedbackRegistry.Contract.GetAllFeedbacks(&_FeedbackRegistry.CallOpts)
}

// GetFeedbacks is a free data retrieval call binding the contract method 0x190a500d.
//
// Solidity: function getFeedbacks(bytes course_, uint256 offset_, uint256 limit_) view returns(string[])
func (_FeedbackRegistry *FeedbackRegistryCaller) GetFeedbacks(opts *bind.CallOpts, course_ []byte, offset_ *big.Int, limit_ *big.Int) ([]string, error) {
	var out []interface{}
	err := _FeedbackRegistry.contract.Call(opts, &out, "getFeedbacks", course_, offset_, limit_)

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// GetFeedbacks is a free data retrieval call binding the contract method 0x190a500d.
//
// Solidity: function getFeedbacks(bytes course_, uint256 offset_, uint256 limit_) view returns(string[])
func (_FeedbackRegistry *FeedbackRegistrySession) GetFeedbacks(course_ []byte, offset_ *big.Int, limit_ *big.Int) ([]string, error) {
	return _FeedbackRegistry.Contract.GetFeedbacks(&_FeedbackRegistry.CallOpts, course_, offset_, limit_)
}

// GetFeedbacks is a free data retrieval call binding the contract method 0x190a500d.
//
// Solidity: function getFeedbacks(bytes course_, uint256 offset_, uint256 limit_) view returns(string[])
func (_FeedbackRegistry *FeedbackRegistryCallerSession) GetFeedbacks(course_ []byte, offset_ *big.Int, limit_ *big.Int) ([]string, error) {
	return _FeedbackRegistry.Contract.GetFeedbacks(&_FeedbackRegistry.CallOpts, course_, offset_, limit_)
}

// AddFeedback is a paid mutator transaction binding the contract method 0xde646e9f.
//
// Solidity: function addFeedback(bytes course_, uint256 i_, uint256[] c_, uint256[] r_, uint256[] publicKeysX_, uint256[] publicKeysY_, bytes32[][] merkleTreeProofs_, bytes32[] keys_, bytes32[] values_, string ipfsHash_) returns()
func (_FeedbackRegistry *FeedbackRegistryTransactor) AddFeedback(opts *bind.TransactOpts, course_ []byte, i_ *big.Int, c_ []*big.Int, r_ []*big.Int, publicKeysX_ []*big.Int, publicKeysY_ []*big.Int, merkleTreeProofs_ [][][32]byte, keys_ [][32]byte, values_ [][32]byte, ipfsHash_ string) (*types.Transaction, error) {
	return _FeedbackRegistry.contract.Transact(opts, "addFeedback", course_, i_, c_, r_, publicKeysX_, publicKeysY_, merkleTreeProofs_, keys_, values_, ipfsHash_)
}

// AddFeedback is a paid mutator transaction binding the contract method 0xde646e9f.
//
// Solidity: function addFeedback(bytes course_, uint256 i_, uint256[] c_, uint256[] r_, uint256[] publicKeysX_, uint256[] publicKeysY_, bytes32[][] merkleTreeProofs_, bytes32[] keys_, bytes32[] values_, string ipfsHash_) returns()
func (_FeedbackRegistry *FeedbackRegistrySession) AddFeedback(course_ []byte, i_ *big.Int, c_ []*big.Int, r_ []*big.Int, publicKeysX_ []*big.Int, publicKeysY_ []*big.Int, merkleTreeProofs_ [][][32]byte, keys_ [][32]byte, values_ [][32]byte, ipfsHash_ string) (*types.Transaction, error) {
	return _FeedbackRegistry.Contract.AddFeedback(&_FeedbackRegistry.TransactOpts, course_, i_, c_, r_, publicKeysX_, publicKeysY_, merkleTreeProofs_, keys_, values_, ipfsHash_)
}

// AddFeedback is a paid mutator transaction binding the contract method 0xde646e9f.
//
// Solidity: function addFeedback(bytes course_, uint256 i_, uint256[] c_, uint256[] r_, uint256[] publicKeysX_, uint256[] publicKeysY_, bytes32[][] merkleTreeProofs_, bytes32[] keys_, bytes32[] values_, string ipfsHash_) returns()
func (_FeedbackRegistry *FeedbackRegistryTransactorSession) AddFeedback(course_ []byte, i_ *big.Int, c_ []*big.Int, r_ []*big.Int, publicKeysX_ []*big.Int, publicKeysY_ []*big.Int, merkleTreeProofs_ [][][32]byte, keys_ [][32]byte, values_ [][32]byte, ipfsHash_ string) (*types.Transaction, error) {
	return _FeedbackRegistry.Contract.AddFeedback(&_FeedbackRegistry.TransactOpts, course_, i_, c_, r_, publicKeysX_, publicKeysY_, merkleTreeProofs_, keys_, values_, ipfsHash_)
}
