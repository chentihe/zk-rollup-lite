package layer1

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

func DecodeTxHash(txHash string) (*types.Transaction, error) {
	rawTxBytes, err := hex.DecodeString(txHash)
	if err != nil {
		return nil, err
	}

	tx := new(types.Transaction)
	rlp.DecodeBytes(rawTxBytes, &tx)

	return tx, nil
}
