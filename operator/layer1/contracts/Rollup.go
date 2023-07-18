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

// RollupUser is an auto generated low-level Go binding around an user-defined struct.
type RollupUser struct {
	Index      *big.Int
	PublicKeyX *big.Int
	PublicKeyY *big.Int
	Balance    *big.Int
	Nonce      *big.Int
}

// RollupMetaData contains all meta data concerning the Rollup contract.
var RollupMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractTxVerifier\",\"name\":\"_txVerifier\",\"type\":\"address\"},{\"internalType\":\"contractWithdrawVerifier\",\"name\":\"_withdrawVerifier\",\"type\":\"address\"},{\"internalType\":\"contractDepositVerifier\",\"name\":\"_depositVerifier\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"INSUFFICIENT_BALANCE\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"INVALID_DEPOSIT_PROOFS\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"INVALID_MERKLE_TREE\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"INVALID_NULLIFIER\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"INVALID_ROLLUP_PROOFS\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"INVALID_USER\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"INVALID_VALUE\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"INVALID_WITHDRAW_PROOFS\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ONLY_OWNER\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"REENTRANT_CALL\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WITHDRAWAL_FAILED\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"publicKeyX\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"publicKeyY\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structRollup.User\",\"name\":\"user\",\"type\":\"tuple\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBalanceTreeRoot\",\"type\":\"uint256\"}],\"name\":\"RollUp\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"publicKeyX\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"publicKeyY\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structRollup.User\",\"name\":\"user\",\"type\":\"tuple\"}],\"name\":\"Withdraw\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"balanceTreeKeys\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"balanceTreeRoot\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"balanceTreeUsers\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"publicKeyX\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"publicKeyY\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"a\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2][2]\",\"name\":\"b\",\"type\":\"uint256[2][2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"c\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[5]\",\"name\":\"input\",\"type\":\"uint256[5]\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"publicKeyX\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"publicKeyY\",\"type\":\"uint256\"}],\"name\":\"generateKeyHash\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"getUserByIndex\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"publicKeyX\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"publicKeyY\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"internalType\":\"structRollup.User\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"publicKeyX\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"publicKeyY\",\"type\":\"uint256\"}],\"name\":\"getUserByPublicKey\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"publicKeyX\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"publicKeyY\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"internalType\":\"structRollup.User\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"isPublicKeysRegistered\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"a\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2][2]\",\"name\":\"b\",\"type\":\"uint256[2][2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"c\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[19]\",\"name\":\"input\",\"type\":\"uint256[19]\"}],\"name\":\"rollUp\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"usedNullifiers\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256[2]\",\"name\":\"a\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2][2]\",\"name\":\"b\",\"type\":\"uint256[2][2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"c\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[5]\",\"name\":\"input\",\"type\":\"uint256[5]\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawAccruedFees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
}

// RollupABI is the input ABI used to generate the binding from.
// Deprecated: Use RollupMetaData.ABI instead.
var RollupABI = RollupMetaData.ABI

// Rollup is an auto generated Go binding around an Ethereum contract.
type Rollup struct {
	RollupCaller     // Read-only binding to the contract
	RollupTransactor // Write-only binding to the contract
	RollupFilterer   // Log filterer for contract events
}

// RollupCaller is an auto generated read-only Go binding around an Ethereum contract.
type RollupCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RollupTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RollupTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RollupFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RollupFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RollupSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RollupSession struct {
	Contract     *Rollup           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RollupCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RollupCallerSession struct {
	Contract *RollupCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// RollupTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RollupTransactorSession struct {
	Contract     *RollupTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RollupRaw is an auto generated low-level Go binding around an Ethereum contract.
type RollupRaw struct {
	Contract *Rollup // Generic contract binding to access the raw methods on
}

// RollupCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RollupCallerRaw struct {
	Contract *RollupCaller // Generic read-only contract binding to access the raw methods on
}

// RollupTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RollupTransactorRaw struct {
	Contract *RollupTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRollup creates a new instance of Rollup, bound to a specific deployed contract.
func NewRollup(address common.Address, backend bind.ContractBackend) (*Rollup, error) {
	contract, err := bindRollup(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Rollup{RollupCaller: RollupCaller{contract: contract}, RollupTransactor: RollupTransactor{contract: contract}, RollupFilterer: RollupFilterer{contract: contract}}, nil
}

// NewRollupCaller creates a new read-only instance of Rollup, bound to a specific deployed contract.
func NewRollupCaller(address common.Address, caller bind.ContractCaller) (*RollupCaller, error) {
	contract, err := bindRollup(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RollupCaller{contract: contract}, nil
}

// NewRollupTransactor creates a new write-only instance of Rollup, bound to a specific deployed contract.
func NewRollupTransactor(address common.Address, transactor bind.ContractTransactor) (*RollupTransactor, error) {
	contract, err := bindRollup(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RollupTransactor{contract: contract}, nil
}

// NewRollupFilterer creates a new log filterer instance of Rollup, bound to a specific deployed contract.
func NewRollupFilterer(address common.Address, filterer bind.ContractFilterer) (*RollupFilterer, error) {
	contract, err := bindRollup(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RollupFilterer{contract: contract}, nil
}

// bindRollup binds a generic wrapper to an already deployed contract.
func bindRollup(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RollupMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Rollup *RollupRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Rollup.Contract.RollupCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Rollup *RollupRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rollup.Contract.RollupTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Rollup *RollupRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Rollup.Contract.RollupTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Rollup *RollupCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Rollup.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Rollup *RollupTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rollup.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Rollup *RollupTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Rollup.Contract.contract.Transact(opts, method, params...)
}

// BalanceTreeKeys is a free data retrieval call binding the contract method 0x34344812.
//
// Solidity: function balanceTreeKeys(uint256 ) view returns(uint256)
func (_Rollup *RollupCaller) BalanceTreeKeys(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Rollup.contract.Call(opts, &out, "balanceTreeKeys", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceTreeKeys is a free data retrieval call binding the contract method 0x34344812.
//
// Solidity: function balanceTreeKeys(uint256 ) view returns(uint256)
func (_Rollup *RollupSession) BalanceTreeKeys(arg0 *big.Int) (*big.Int, error) {
	return _Rollup.Contract.BalanceTreeKeys(&_Rollup.CallOpts, arg0)
}

// BalanceTreeKeys is a free data retrieval call binding the contract method 0x34344812.
//
// Solidity: function balanceTreeKeys(uint256 ) view returns(uint256)
func (_Rollup *RollupCallerSession) BalanceTreeKeys(arg0 *big.Int) (*big.Int, error) {
	return _Rollup.Contract.BalanceTreeKeys(&_Rollup.CallOpts, arg0)
}

// BalanceTreeRoot is a free data retrieval call binding the contract method 0xb4e7dddd.
//
// Solidity: function balanceTreeRoot() view returns(uint256)
func (_Rollup *RollupCaller) BalanceTreeRoot(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Rollup.contract.Call(opts, &out, "balanceTreeRoot")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceTreeRoot is a free data retrieval call binding the contract method 0xb4e7dddd.
//
// Solidity: function balanceTreeRoot() view returns(uint256)
func (_Rollup *RollupSession) BalanceTreeRoot() (*big.Int, error) {
	return _Rollup.Contract.BalanceTreeRoot(&_Rollup.CallOpts)
}

// BalanceTreeRoot is a free data retrieval call binding the contract method 0xb4e7dddd.
//
// Solidity: function balanceTreeRoot() view returns(uint256)
func (_Rollup *RollupCallerSession) BalanceTreeRoot() (*big.Int, error) {
	return _Rollup.Contract.BalanceTreeRoot(&_Rollup.CallOpts)
}

// BalanceTreeUsers is a free data retrieval call binding the contract method 0xf6492213.
//
// Solidity: function balanceTreeUsers(uint256 ) view returns(uint256 index, uint256 publicKeyX, uint256 publicKeyY, uint256 balance, uint256 nonce)
func (_Rollup *RollupCaller) BalanceTreeUsers(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Index      *big.Int
	PublicKeyX *big.Int
	PublicKeyY *big.Int
	Balance    *big.Int
	Nonce      *big.Int
}, error) {
	var out []interface{}
	err := _Rollup.contract.Call(opts, &out, "balanceTreeUsers", arg0)

	outstruct := new(struct {
		Index      *big.Int
		PublicKeyX *big.Int
		PublicKeyY *big.Int
		Balance    *big.Int
		Nonce      *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Index = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.PublicKeyX = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.PublicKeyY = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Balance = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.Nonce = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// BalanceTreeUsers is a free data retrieval call binding the contract method 0xf6492213.
//
// Solidity: function balanceTreeUsers(uint256 ) view returns(uint256 index, uint256 publicKeyX, uint256 publicKeyY, uint256 balance, uint256 nonce)
func (_Rollup *RollupSession) BalanceTreeUsers(arg0 *big.Int) (struct {
	Index      *big.Int
	PublicKeyX *big.Int
	PublicKeyY *big.Int
	Balance    *big.Int
	Nonce      *big.Int
}, error) {
	return _Rollup.Contract.BalanceTreeUsers(&_Rollup.CallOpts, arg0)
}

// BalanceTreeUsers is a free data retrieval call binding the contract method 0xf6492213.
//
// Solidity: function balanceTreeUsers(uint256 ) view returns(uint256 index, uint256 publicKeyX, uint256 publicKeyY, uint256 balance, uint256 nonce)
func (_Rollup *RollupCallerSession) BalanceTreeUsers(arg0 *big.Int) (struct {
	Index      *big.Int
	PublicKeyX *big.Int
	PublicKeyY *big.Int
	Balance    *big.Int
	Nonce      *big.Int
}, error) {
	return _Rollup.Contract.BalanceTreeUsers(&_Rollup.CallOpts, arg0)
}

// GenerateKeyHash is a free data retrieval call binding the contract method 0x877f0640.
//
// Solidity: function generateKeyHash(uint256 publicKeyX, uint256 publicKeyY) pure returns(uint256)
func (_Rollup *RollupCaller) GenerateKeyHash(opts *bind.CallOpts, publicKeyX *big.Int, publicKeyY *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Rollup.contract.Call(opts, &out, "generateKeyHash", publicKeyX, publicKeyY)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GenerateKeyHash is a free data retrieval call binding the contract method 0x877f0640.
//
// Solidity: function generateKeyHash(uint256 publicKeyX, uint256 publicKeyY) pure returns(uint256)
func (_Rollup *RollupSession) GenerateKeyHash(publicKeyX *big.Int, publicKeyY *big.Int) (*big.Int, error) {
	return _Rollup.Contract.GenerateKeyHash(&_Rollup.CallOpts, publicKeyX, publicKeyY)
}

// GenerateKeyHash is a free data retrieval call binding the contract method 0x877f0640.
//
// Solidity: function generateKeyHash(uint256 publicKeyX, uint256 publicKeyY) pure returns(uint256)
func (_Rollup *RollupCallerSession) GenerateKeyHash(publicKeyX *big.Int, publicKeyY *big.Int) (*big.Int, error) {
	return _Rollup.Contract.GenerateKeyHash(&_Rollup.CallOpts, publicKeyX, publicKeyY)
}

// GetUserByIndex is a free data retrieval call binding the contract method 0xff5d32fe.
//
// Solidity: function getUserByIndex(uint256 index) view returns((uint256,uint256,uint256,uint256,uint256))
func (_Rollup *RollupCaller) GetUserByIndex(opts *bind.CallOpts, index *big.Int) (RollupUser, error) {
	var out []interface{}
	err := _Rollup.contract.Call(opts, &out, "getUserByIndex", index)

	if err != nil {
		return *new(RollupUser), err
	}

	out0 := *abi.ConvertType(out[0], new(RollupUser)).(*RollupUser)

	return out0, err

}

// GetUserByIndex is a free data retrieval call binding the contract method 0xff5d32fe.
//
// Solidity: function getUserByIndex(uint256 index) view returns((uint256,uint256,uint256,uint256,uint256))
func (_Rollup *RollupSession) GetUserByIndex(index *big.Int) (RollupUser, error) {
	return _Rollup.Contract.GetUserByIndex(&_Rollup.CallOpts, index)
}

// GetUserByIndex is a free data retrieval call binding the contract method 0xff5d32fe.
//
// Solidity: function getUserByIndex(uint256 index) view returns((uint256,uint256,uint256,uint256,uint256))
func (_Rollup *RollupCallerSession) GetUserByIndex(index *big.Int) (RollupUser, error) {
	return _Rollup.Contract.GetUserByIndex(&_Rollup.CallOpts, index)
}

// GetUserByPublicKey is a free data retrieval call binding the contract method 0x8ce66420.
//
// Solidity: function getUserByPublicKey(uint256 publicKeyX, uint256 publicKeyY) view returns((uint256,uint256,uint256,uint256,uint256))
func (_Rollup *RollupCaller) GetUserByPublicKey(opts *bind.CallOpts, publicKeyX *big.Int, publicKeyY *big.Int) (RollupUser, error) {
	var out []interface{}
	err := _Rollup.contract.Call(opts, &out, "getUserByPublicKey", publicKeyX, publicKeyY)

	if err != nil {
		return *new(RollupUser), err
	}

	out0 := *abi.ConvertType(out[0], new(RollupUser)).(*RollupUser)

	return out0, err

}

// GetUserByPublicKey is a free data retrieval call binding the contract method 0x8ce66420.
//
// Solidity: function getUserByPublicKey(uint256 publicKeyX, uint256 publicKeyY) view returns((uint256,uint256,uint256,uint256,uint256))
func (_Rollup *RollupSession) GetUserByPublicKey(publicKeyX *big.Int, publicKeyY *big.Int) (RollupUser, error) {
	return _Rollup.Contract.GetUserByPublicKey(&_Rollup.CallOpts, publicKeyX, publicKeyY)
}

// GetUserByPublicKey is a free data retrieval call binding the contract method 0x8ce66420.
//
// Solidity: function getUserByPublicKey(uint256 publicKeyX, uint256 publicKeyY) view returns((uint256,uint256,uint256,uint256,uint256))
func (_Rollup *RollupCallerSession) GetUserByPublicKey(publicKeyX *big.Int, publicKeyY *big.Int) (RollupUser, error) {
	return _Rollup.Contract.GetUserByPublicKey(&_Rollup.CallOpts, publicKeyX, publicKeyY)
}

// IsPublicKeysRegistered is a free data retrieval call binding the contract method 0xebe682eb.
//
// Solidity: function isPublicKeysRegistered(uint256 ) view returns(bool)
func (_Rollup *RollupCaller) IsPublicKeysRegistered(opts *bind.CallOpts, arg0 *big.Int) (bool, error) {
	var out []interface{}
	err := _Rollup.contract.Call(opts, &out, "isPublicKeysRegistered", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsPublicKeysRegistered is a free data retrieval call binding the contract method 0xebe682eb.
//
// Solidity: function isPublicKeysRegistered(uint256 ) view returns(bool)
func (_Rollup *RollupSession) IsPublicKeysRegistered(arg0 *big.Int) (bool, error) {
	return _Rollup.Contract.IsPublicKeysRegistered(&_Rollup.CallOpts, arg0)
}

// IsPublicKeysRegistered is a free data retrieval call binding the contract method 0xebe682eb.
//
// Solidity: function isPublicKeysRegistered(uint256 ) view returns(bool)
func (_Rollup *RollupCallerSession) IsPublicKeysRegistered(arg0 *big.Int) (bool, error) {
	return _Rollup.Contract.IsPublicKeysRegistered(&_Rollup.CallOpts, arg0)
}

// UsedNullifiers is a free data retrieval call binding the contract method 0xaad24061.
//
// Solidity: function usedNullifiers(uint256 ) view returns(bool)
func (_Rollup *RollupCaller) UsedNullifiers(opts *bind.CallOpts, arg0 *big.Int) (bool, error) {
	var out []interface{}
	err := _Rollup.contract.Call(opts, &out, "usedNullifiers", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// UsedNullifiers is a free data retrieval call binding the contract method 0xaad24061.
//
// Solidity: function usedNullifiers(uint256 ) view returns(bool)
func (_Rollup *RollupSession) UsedNullifiers(arg0 *big.Int) (bool, error) {
	return _Rollup.Contract.UsedNullifiers(&_Rollup.CallOpts, arg0)
}

// UsedNullifiers is a free data retrieval call binding the contract method 0xaad24061.
//
// Solidity: function usedNullifiers(uint256 ) view returns(bool)
func (_Rollup *RollupCallerSession) UsedNullifiers(arg0 *big.Int) (bool, error) {
	return _Rollup.Contract.UsedNullifiers(&_Rollup.CallOpts, arg0)
}

// Deposit is a paid mutator transaction binding the contract method 0xf40e1440.
//
// Solidity: function deposit(uint256[2] a, uint256[2][2] b, uint256[2] c, uint256[5] input) payable returns()
func (_Rollup *RollupTransactor) Deposit(opts *bind.TransactOpts, a [2]*big.Int, b [2][2]*big.Int, c [2]*big.Int, input [5]*big.Int) (*types.Transaction, error) {
	return _Rollup.contract.Transact(opts, "deposit", a, b, c, input)
}

// Deposit is a paid mutator transaction binding the contract method 0xf40e1440.
//
// Solidity: function deposit(uint256[2] a, uint256[2][2] b, uint256[2] c, uint256[5] input) payable returns()
func (_Rollup *RollupSession) Deposit(a [2]*big.Int, b [2][2]*big.Int, c [2]*big.Int, input [5]*big.Int) (*types.Transaction, error) {
	return _Rollup.Contract.Deposit(&_Rollup.TransactOpts, a, b, c, input)
}

// Deposit is a paid mutator transaction binding the contract method 0xf40e1440.
//
// Solidity: function deposit(uint256[2] a, uint256[2][2] b, uint256[2] c, uint256[5] input) payable returns()
func (_Rollup *RollupTransactorSession) Deposit(a [2]*big.Int, b [2][2]*big.Int, c [2]*big.Int, input [5]*big.Int) (*types.Transaction, error) {
	return _Rollup.Contract.Deposit(&_Rollup.TransactOpts, a, b, c, input)
}

// RollUp is a paid mutator transaction binding the contract method 0x9b66f706.
//
// Solidity: function rollUp(uint256[2] a, uint256[2][2] b, uint256[2] c, uint256[19] input) returns()
func (_Rollup *RollupTransactor) RollUp(opts *bind.TransactOpts, a [2]*big.Int, b [2][2]*big.Int, c [2]*big.Int, input [19]*big.Int) (*types.Transaction, error) {
	return _Rollup.contract.Transact(opts, "rollUp", a, b, c, input)
}

// RollUp is a paid mutator transaction binding the contract method 0x9b66f706.
//
// Solidity: function rollUp(uint256[2] a, uint256[2][2] b, uint256[2] c, uint256[19] input) returns()
func (_Rollup *RollupSession) RollUp(a [2]*big.Int, b [2][2]*big.Int, c [2]*big.Int, input [19]*big.Int) (*types.Transaction, error) {
	return _Rollup.Contract.RollUp(&_Rollup.TransactOpts, a, b, c, input)
}

// RollUp is a paid mutator transaction binding the contract method 0x9b66f706.
//
// Solidity: function rollUp(uint256[2] a, uint256[2][2] b, uint256[2] c, uint256[19] input) returns()
func (_Rollup *RollupTransactorSession) RollUp(a [2]*big.Int, b [2][2]*big.Int, c [2]*big.Int, input [19]*big.Int) (*types.Transaction, error) {
	return _Rollup.Contract.RollUp(&_Rollup.TransactOpts, a, b, c, input)
}

// Withdraw is a paid mutator transaction binding the contract method 0xb3531741.
//
// Solidity: function withdraw(uint256 amount, uint256[2] a, uint256[2][2] b, uint256[2] c, uint256[5] input) returns()
func (_Rollup *RollupTransactor) Withdraw(opts *bind.TransactOpts, amount *big.Int, a [2]*big.Int, b [2][2]*big.Int, c [2]*big.Int, input [5]*big.Int) (*types.Transaction, error) {
	return _Rollup.contract.Transact(opts, "withdraw", amount, a, b, c, input)
}

// Withdraw is a paid mutator transaction binding the contract method 0xb3531741.
//
// Solidity: function withdraw(uint256 amount, uint256[2] a, uint256[2][2] b, uint256[2] c, uint256[5] input) returns()
func (_Rollup *RollupSession) Withdraw(amount *big.Int, a [2]*big.Int, b [2][2]*big.Int, c [2]*big.Int, input [5]*big.Int) (*types.Transaction, error) {
	return _Rollup.Contract.Withdraw(&_Rollup.TransactOpts, amount, a, b, c, input)
}

// Withdraw is a paid mutator transaction binding the contract method 0xb3531741.
//
// Solidity: function withdraw(uint256 amount, uint256[2] a, uint256[2][2] b, uint256[2] c, uint256[5] input) returns()
func (_Rollup *RollupTransactorSession) Withdraw(amount *big.Int, a [2]*big.Int, b [2][2]*big.Int, c [2]*big.Int, input [5]*big.Int) (*types.Transaction, error) {
	return _Rollup.Contract.Withdraw(&_Rollup.TransactOpts, amount, a, b, c, input)
}

// WithdrawAccruedFees is a paid mutator transaction binding the contract method 0xada82c7d.
//
// Solidity: function withdrawAccruedFees() returns()
func (_Rollup *RollupTransactor) WithdrawAccruedFees(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rollup.contract.Transact(opts, "withdrawAccruedFees")
}

// WithdrawAccruedFees is a paid mutator transaction binding the contract method 0xada82c7d.
//
// Solidity: function withdrawAccruedFees() returns()
func (_Rollup *RollupSession) WithdrawAccruedFees() (*types.Transaction, error) {
	return _Rollup.Contract.WithdrawAccruedFees(&_Rollup.TransactOpts)
}

// WithdrawAccruedFees is a paid mutator transaction binding the contract method 0xada82c7d.
//
// Solidity: function withdrawAccruedFees() returns()
func (_Rollup *RollupTransactorSession) WithdrawAccruedFees() (*types.Transaction, error) {
	return _Rollup.Contract.WithdrawAccruedFees(&_Rollup.TransactOpts)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Rollup *RollupTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _Rollup.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Rollup *RollupSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Rollup.Contract.Fallback(&_Rollup.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Rollup *RollupTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Rollup.Contract.Fallback(&_Rollup.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Rollup *RollupTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rollup.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Rollup *RollupSession) Receive() (*types.Transaction, error) {
	return _Rollup.Contract.Receive(&_Rollup.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Rollup *RollupTransactorSession) Receive() (*types.Transaction, error) {
	return _Rollup.Contract.Receive(&_Rollup.TransactOpts)
}

// RollupDepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the Rollup contract.
type RollupDepositIterator struct {
	Event *RollupDeposit // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RollupDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RollupDeposit)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RollupDeposit)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RollupDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RollupDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RollupDeposit represents a Deposit event raised by the Rollup contract.
type RollupDeposit struct {
	User RollupUser
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0x03be40b3bcb1ca7bb1f0320436c492fae3f0b576cd339dc5690519236d0560f9.
//
// Solidity: event Deposit((uint256,uint256,uint256,uint256,uint256) user)
func (_Rollup *RollupFilterer) FilterDeposit(opts *bind.FilterOpts) (*RollupDepositIterator, error) {

	logs, sub, err := _Rollup.contract.FilterLogs(opts, "Deposit")
	if err != nil {
		return nil, err
	}
	return &RollupDepositIterator{contract: _Rollup.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0x03be40b3bcb1ca7bb1f0320436c492fae3f0b576cd339dc5690519236d0560f9.
//
// Solidity: event Deposit((uint256,uint256,uint256,uint256,uint256) user)
func (_Rollup *RollupFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *RollupDeposit) (event.Subscription, error) {

	logs, sub, err := _Rollup.contract.WatchLogs(opts, "Deposit")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RollupDeposit)
				if err := _Rollup.contract.UnpackLog(event, "Deposit", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDeposit is a log parse operation binding the contract event 0x03be40b3bcb1ca7bb1f0320436c492fae3f0b576cd339dc5690519236d0560f9.
//
// Solidity: event Deposit((uint256,uint256,uint256,uint256,uint256) user)
func (_Rollup *RollupFilterer) ParseDeposit(log types.Log) (*RollupDeposit, error) {
	event := new(RollupDeposit)
	if err := _Rollup.contract.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RollupRollUpIterator is returned from FilterRollUp and is used to iterate over the raw logs and unpacked data for RollUp events raised by the Rollup contract.
type RollupRollUpIterator struct {
	Event *RollupRollUp // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RollupRollUpIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RollupRollUp)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RollupRollUp)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RollupRollUpIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RollupRollUpIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RollupRollUp represents a RollUp event raised by the Rollup contract.
type RollupRollUp struct {
	NewBalanceTreeRoot *big.Int
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterRollUp is a free log retrieval operation binding the contract event 0xc3c2d7581ffe4d3d2ef00cefb41bedfaed37bfa5be059acffd469b7b200f5539.
//
// Solidity: event RollUp(uint256 newBalanceTreeRoot)
func (_Rollup *RollupFilterer) FilterRollUp(opts *bind.FilterOpts) (*RollupRollUpIterator, error) {

	logs, sub, err := _Rollup.contract.FilterLogs(opts, "RollUp")
	if err != nil {
		return nil, err
	}
	return &RollupRollUpIterator{contract: _Rollup.contract, event: "RollUp", logs: logs, sub: sub}, nil
}

// WatchRollUp is a free log subscription operation binding the contract event 0xc3c2d7581ffe4d3d2ef00cefb41bedfaed37bfa5be059acffd469b7b200f5539.
//
// Solidity: event RollUp(uint256 newBalanceTreeRoot)
func (_Rollup *RollupFilterer) WatchRollUp(opts *bind.WatchOpts, sink chan<- *RollupRollUp) (event.Subscription, error) {

	logs, sub, err := _Rollup.contract.WatchLogs(opts, "RollUp")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RollupRollUp)
				if err := _Rollup.contract.UnpackLog(event, "RollUp", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRollUp is a log parse operation binding the contract event 0xc3c2d7581ffe4d3d2ef00cefb41bedfaed37bfa5be059acffd469b7b200f5539.
//
// Solidity: event RollUp(uint256 newBalanceTreeRoot)
func (_Rollup *RollupFilterer) ParseRollUp(log types.Log) (*RollupRollUp, error) {
	event := new(RollupRollUp)
	if err := _Rollup.contract.UnpackLog(event, "RollUp", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RollupWithdrawIterator is returned from FilterWithdraw and is used to iterate over the raw logs and unpacked data for Withdraw events raised by the Rollup contract.
type RollupWithdrawIterator struct {
	Event *RollupWithdraw // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RollupWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RollupWithdraw)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RollupWithdraw)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RollupWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RollupWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RollupWithdraw represents a Withdraw event raised by the Rollup contract.
type RollupWithdraw struct {
	User RollupUser
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterWithdraw is a free log retrieval operation binding the contract event 0xad6251f53545fa57aa67137dc9192575d5e38d57358ed43b1d41a3fa4ced9c48.
//
// Solidity: event Withdraw((uint256,uint256,uint256,uint256,uint256) user)
func (_Rollup *RollupFilterer) FilterWithdraw(opts *bind.FilterOpts) (*RollupWithdrawIterator, error) {

	logs, sub, err := _Rollup.contract.FilterLogs(opts, "Withdraw")
	if err != nil {
		return nil, err
	}
	return &RollupWithdrawIterator{contract: _Rollup.contract, event: "Withdraw", logs: logs, sub: sub}, nil
}

// WatchWithdraw is a free log subscription operation binding the contract event 0xad6251f53545fa57aa67137dc9192575d5e38d57358ed43b1d41a3fa4ced9c48.
//
// Solidity: event Withdraw((uint256,uint256,uint256,uint256,uint256) user)
func (_Rollup *RollupFilterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *RollupWithdraw) (event.Subscription, error) {

	logs, sub, err := _Rollup.contract.WatchLogs(opts, "Withdraw")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RollupWithdraw)
				if err := _Rollup.contract.UnpackLog(event, "Withdraw", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseWithdraw is a log parse operation binding the contract event 0xad6251f53545fa57aa67137dc9192575d5e38d57358ed43b1d41a3fa4ced9c48.
//
// Solidity: event Withdraw((uint256,uint256,uint256,uint256,uint256) user)
func (_Rollup *RollupFilterer) ParseWithdraw(log types.Log) (*RollupWithdraw, error) {
	event := new(RollupWithdraw)
	if err := _Rollup.contract.UnpackLog(event, "Withdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
