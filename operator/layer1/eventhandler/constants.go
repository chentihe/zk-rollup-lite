package eventhandler

import "github.com/ethereum/go-ethereum/crypto"

var (
	depositHash  = crypto.Keccak256Hash([]byte("Deposit((uint256,uint256,uint256,uint256,uint256))"))
	withdrawHash = crypto.Keccak256Hash([]byte("Withdraw((uint256,uint256,uint256,uint256,uint256))"))
)
