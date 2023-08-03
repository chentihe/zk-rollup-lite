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
	Bin: "0x608060405234801562000010575f80fd5b5060405162001ce238038062001ce28339818101604052810190620000369190620001bf565b8160035f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508060045f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060015f819055503360015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550505062000204565b5f80fd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f620001348262000109565b9050919050565b5f620001478262000128565b9050919050565b62000159816200013b565b811462000164575f80fd5b50565b5f8151905062000177816200014e565b92915050565b5f620001898262000128565b9050919050565b6200019b816200017d565b8114620001a6575f80fd5b50565b5f81519050620001b98162000190565b92915050565b5f8060408385031215620001d857620001d762000105565b5b5f620001e78582860162000167565b9250506020620001fa85828601620001a9565b9150509250929050565b611ad080620002125f395ff3fe6080604052600436106100a6575f3560e01c8063ada82c7d11610063578063ada82c7d146101ea578063b4e7dddd14610200578063e2bbb1581461022a578063ebe682eb14610246578063f649221314610282578063ff5d32fe146102c2576100a6565b806334344812146100aa5780635f2ee1bd146100e6578063877f06401461010e5780638ce664201461014a5780639b66f70614610186578063aad24061146101ae575b5f80fd5b3480156100b5575f80fd5b506100d060048036038101906100cb9190610e5a565b6102fe565b6040516100dd9190610e94565b60405180910390f35b3480156100f1575f80fd5b5061010c60048036038101906101079190611147565b610313565b005b348015610119575f80fd5b50610134600480360381019061012f91906111c0565b610682565b6040516101419190610e94565b60405180910390f35b348015610155575f80fd5b50610170600480360381019061016b91906111c0565b610695565b60405161017d9190611273565b60405180910390f35b348015610191575f80fd5b506101ac60048036038101906101a7919061133a565b6106eb565b005b3480156101b9575f80fd5b506101d460048036038101906101cf9190610e5a565b61098f565b6040516101e191906113ba565b60405180910390f35b3480156101f5575f80fd5b506101fe6109ac565b005b34801561020b575f80fd5b50610214610add565b6040516102219190610e94565b60405180910390f35b610244600480360381019061023f91906111c0565b610ae3565b005b348015610251575f80fd5b5061026c60048036038101906102679190610e5a565b610c2b565b60405161027991906113ba565b60405180910390f35b34801561028d575f80fd5b506102a860048036038101906102a39190610e5a565b610c48565b6040516102b99594939291906113d3565b60405180910390f35b3480156102cd575f80fd5b506102e860048036038101906102e39190610e5a565b610c7a565b6040516102f59190611273565b60405180910390f35b6009602052805f5260405f205f915090505481565b60025f540361034e576040517fdfc60d8500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60025f8190555060045f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166311479fea858585856040518563ffffffff1660e01b81526004016103b5949392919061162f565b602060405180830381865afa1580156103d0573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906103f4919061169e565b61042a576040517fae2f07de00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f815f6003811061043e5761043d6116c9565b5b602002015190505f8260016003811061045a576104596116c9565b5b602002015190505f83600260038110610476576104756116c9565b5b6020020151905060075f8481526020019081526020015f205f9054906101000a900460ff16156104d2576040517f6853f30700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f6104dd8383610d40565b90507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff890361050e57806003015498505b806003015489118061051f57505f89145b15610556576040517f50b1f35600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600160075f8681526020019081526020015f205f6101000a81548160ff02191690831515021790555088816003015f8282546105929190611723565b925050819055505f3373ffffffffffffffffffffffffffffffffffffffff168a6040516105be90611783565b5f6040518083038185875af1925050503d805f81146105f8576040519150601f19603f3d011682016040523d82523d5f602084013e6105fd565b606091505b5050905080610638576040517ffbefd20100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7fad6251f53545fa57aa67137dc9192575d5e38d57358ed43b1d41a3fa4ced9c48826040516106679190611863565b60405180910390a1505050505060015f819055505050505050565b5f61068d8383610dbc565b905092915050565b61069d610df0565b6106a78383610d40565b6040518060a00160405290815f8201548152602001600182015481526020016002820154815260200160038201548152602001600482015481525050905092915050565b60035f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663f3f22e72858585856040518563ffffffff1660e01b815260040161074b94939291906118fb565b602060405180830381865afa158015610766573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061078a919061169e565b6107c0576040517fc3c1f43a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f815f601381106107d4576107d36116c9565b5b602002015190505f805f805f805f600390505f600290505f5b818160ff16101561097857806008610805919061194c565b60ff16836108139190611988565b93508a8460138110610828576108276116c9565b5b602002015198508a60018561083d9190611988565b6013811061084e5761084d6116c9565b5b602002015197508a6002856108639190611988565b60138110610874576108736116c9565b5b602002015196508a6003856108899190611988565b6013811061089a576108996116c9565b5b602002015195508a6004856108af9190611988565b601381106108c0576108bf6116c9565b5b602002015194505f60055f60095f8d81526020019081526020015f205481526020019081526020015f20905087816003015f828254039250508190555086816003015f828254039250508190555085816004018190555086600a5f8282546109289190611988565b9250508190555060055f60095f8c81526020019081526020015f205481526020019081526020015f20905087816003015f8282540192505081905550508080610970906119bb565b9150506107ed565b508860028190555050505050505050505050505050565b6007602052805f5260405f205f915054906101000a900460ff1681565b60015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610a32576040517fd238ed5900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f3373ffffffffffffffffffffffffffffffffffffffff16600a54604051610a5990611783565b5f6040518083038185875af1925050503d805f8114610a93576040519150601f19603f3d011682016040523d82523d5f602084013e610a98565b606091505b5050905080610ad3576040517ffbefd20100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f600a8190555050565b60025481565b5f3403610b1c576040517ff289bba300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f610b278383610dbc565b90505f60055f8381526020019081526020015f20905034816003015f828254610b509190611988565b9250508190555060065f8381526020019081526020015f205f9054906101000a900460ff16610bee57600160065f8481526020019081526020015f205f6101000a81548160ff02191690831515021790555060085f815480929190610bb4906119e3565b9190505550600854815f01819055508381600101819055508281600201819055508160095f60085481526020019081526020015f20819055505b7f03be40b3bcb1ca7bb1f0320436c492fae3f0b576cd339dc5690519236d0560f981604051610c1d9190611863565b60405180910390a150505050565b6006602052805f5260405f205f915054906101000a900460ff1681565b6005602052805f5260405f205f91509050805f0154908060010154908060020154908060030154908060040154905085565b610c82610df0565b5f60095f8481526020019081526020015f2054905060065f8281526020019081526020015f205f9054906101000a900460ff16610ceb576040517f8986d85500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60055f8281526020019081526020015f206040518060a00160405290815f8201548152602001600182015481526020016002820154815260200160038201548152602001600482015481525050915050919050565b5f80610d4c8484610dbc565b905060065f8281526020019081526020015f205f9054906101000a900460ff16610da2576040517f8986d85500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60055f8281526020019081526020015f2091505092915050565b5f8282604051602001610dd0929190611a4a565b604051602081830303815290604052805190602001205f1c905092915050565b6040518060a001604052805f81526020015f81526020015f81526020015f81526020015f81525090565b5f604051905090565b5f80fd5b5f819050919050565b610e3981610e27565b8114610e43575f80fd5b50565b5f81359050610e5481610e30565b92915050565b5f60208284031215610e6f57610e6e610e23565b5b5f610e7c84828501610e46565b91505092915050565b610e8e81610e27565b82525050565b5f602082019050610ea75f830184610e85565b92915050565b5f80fd5b5f601f19601f8301169050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b610ef782610eb1565b810181811067ffffffffffffffff82111715610f1657610f15610ec1565b5b80604052505050565b5f610f28610e1a565b9050610f348282610eee565b919050565b5f67ffffffffffffffff821115610f5357610f52610ec1565b5b602082029050919050565b5f80fd5b5f610f74610f6f84610f39565b610f1f565b90508060208402830185811115610f8e57610f8d610f5e565b5b835b81811015610fb75780610fa38882610e46565b845260208401935050602081019050610f90565b5050509392505050565b5f82601f830112610fd557610fd4610ead565b5b6002610fe2848285610f62565b91505092915050565b5f67ffffffffffffffff82111561100557611004610ec1565b5b602082029050919050565b5f61102261101d84610feb565b610f1f565b9050806040840283018581111561103c5761103b610f5e565b5b835b8181101561106557806110518882610fc1565b84526020840193505060408101905061103e565b5050509392505050565b5f82601f83011261108357611082610ead565b5b6002611090848285611010565b91505092915050565b5f67ffffffffffffffff8211156110b3576110b2610ec1565b5b602082029050919050565b5f6110d06110cb84611099565b610f1f565b905080602084028301858111156110ea576110e9610f5e565b5b835b8181101561111357806110ff8882610e46565b8452602084019350506020810190506110ec565b5050509392505050565b5f82601f83011261113157611130610ead565b5b600361113e8482856110be565b91505092915050565b5f805f805f610180868803121561116157611160610e23565b5b5f61116e88828901610e46565b955050602061117f88828901610fc1565b94505060606111908882890161106f565b93505060e06111a188828901610fc1565b9250506101206111b38882890161111d565b9150509295509295909350565b5f80604083850312156111d6576111d5610e23565b5b5f6111e385828601610e46565b92505060206111f485828601610e46565b9150509250929050565b61120781610e27565b82525050565b60a082015f8201516112215f8501826111fe565b50602082015161123460208501826111fe565b50604082015161124760408501826111fe565b50606082015161125a60608501826111fe565b50608082015161126d60808501826111fe565b50505050565b5f60a0820190506112865f83018461120d565b92915050565b5f67ffffffffffffffff8211156112a6576112a5610ec1565b5b602082029050919050565b5f6112c36112be8461128c565b610f1f565b905080602084028301858111156112dd576112dc610f5e565b5b835b8181101561130657806112f28882610e46565b8452602084019350506020810190506112df565b5050509392505050565b5f82601f83011261132457611323610ead565b5b60136113318482856112b1565b91505092915050565b5f805f80610360858703121561135357611352610e23565b5b5f61136087828801610fc1565b94505060406113718782880161106f565b93505060c061138287828801610fc1565b92505061010061139487828801611310565b91505092959194509250565b5f8115159050919050565b6113b4816113a0565b82525050565b5f6020820190506113cd5f8301846113ab565b92915050565b5f60a0820190506113e65f830188610e85565b6113f36020830187610e85565b6114006040830186610e85565b61140d6060830185610e85565b61141a6080830184610e85565b9695505050505050565b5f60029050919050565b5f81905092915050565b5f819050919050565b5f61144c83836111fe565b60208301905092915050565b5f602082019050919050565b61146d81611424565b611477818461142e565b925061148282611438565b805f5b838110156114b25781516114998782611441565b96506114a483611458565b925050600181019050611485565b505050505050565b5f60029050919050565b5f81905092915050565b5f819050919050565b5f81905092915050565b6114ea81611424565b6114f481846114d7565b92506114ff82611438565b805f5b8381101561152f5781516115168782611441565b965061152183611458565b925050600181019050611502565b505050505050565b5f61154283836114e1565b60408301905092915050565b5f602082019050919050565b611563816114ba565b61156d81846114c4565b9250611578826114ce565b805f5b838110156115a857815161158f8782611537565b965061159a8361154e565b92505060018101905061157b565b505050505050565b5f60039050919050565b5f81905092915050565b5f819050919050565b5f602082019050919050565b6115e2816115b0565b6115ec81846115ba565b92506115f7826115c4565b805f5b8381101561162757815161160e8782611441565b9650611619836115cd565b9250506001810190506115fa565b505050505050565b5f610160820190506116435f830187611464565b611650604083018661155a565b61165d60c0830185611464565b61166b6101008301846115d9565b95945050505050565b61167d816113a0565b8114611687575f80fd5b50565b5f8151905061169881611674565b92915050565b5f602082840312156116b3576116b2610e23565b5b5f6116c08482850161168a565b91505092915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f61172d82610e27565b915061173883610e27565b92508282039050818111156117505761174f6116f6565b5b92915050565b5f81905092915050565b50565b5f61176e5f83611756565b915061177982611760565b5f82019050919050565b5f61178d82611763565b9150819050919050565b5f815f1c9050919050565b5f819050919050565b5f6117bd6117b883611797565b6117a2565b9050919050565b60a082015f8083015490506117d8816117ab565b6117e45f8601826111fe565b50600183015490506117f5816117ab565b61180260208601826111fe565b5060028301549050611813816117ab565b61182060408601826111fe565b5060038301549050611831816117ab565b61183e60608601826111fe565b506004830154905061184f816117ab565b61185c60808601826111fe565b5050505050565b5f60a0820190506118765f8301846117c4565b92915050565b5f60139050919050565b5f81905092915050565b5f819050919050565b5f602082019050919050565b6118ae8161187c565b6118b88184611886565b92506118c382611890565b805f5b838110156118f35781516118da8782611441565b96506118e583611899565b9250506001810190506118c6565b505050505050565b5f6103608201905061190f5f830187611464565b61191c604083018661155a565b61192960c0830185611464565b6119376101008301846118a5565b95945050505050565b5f60ff82169050919050565b5f61195682611940565b915061196183611940565b925082820261196f81611940565b9150808214611981576119806116f6565b5b5092915050565b5f61199282610e27565b915061199d83610e27565b92508282019050808211156119b5576119b46116f6565b5b92915050565b5f6119c582611940565b915060ff82036119d8576119d76116f6565b5b600182019050919050565b5f6119ed82610e27565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203611a1f57611a1e6116f6565b5b600182019050919050565b5f819050919050565b611a44611a3f82610e27565b611a2a565b82525050565b5f611a558285611a33565b602082019150611a658284611a33565b602082019150819050939250505056fea264697066735822122082abfaa6e41302f18a58953cb6d99c92ee258b84f5c61ec74d9e5440adff055f64736f6c637827302e382e32312d646576656c6f702e323032332e372e332b636f6d6d69742e32663435316131380058",
}

// RollupABI is the input ABI used to generate the binding from.
// Deprecated: Use RollupMetaData.ABI instead.
var RollupABI = RollupMetaData.ABI

// RollupBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use RollupMetaData.Bin instead.
var RollupBin = RollupMetaData.Bin

// DeployRollup deploys a new Ethereum contract, binding an instance of Rollup to it.
func DeployRollup(auth *bind.TransactOpts, backend bind.ContractBackend, _txVerifier common.Address, _withdrawVerifier common.Address, _depositVerifier common.Address) (common.Address, *types.Transaction, *Rollup, error) {
	parsed, err := RollupMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(RollupBin), backend, _txVerifier, _withdrawVerifier, _depositVerifier)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Rollup{RollupCaller: RollupCaller{contract: contract}, RollupTransactor: RollupTransactor{contract: contract}, RollupFilterer: RollupFilterer{contract: contract}}, nil
}

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
