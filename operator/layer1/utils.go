package layer1

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/core/types"
)

func DecodeTxHash(txHash string) (*types.Transaction, error) {
	rawTxBytes, err := hex.DecodeString(txHash)
	if err != nil {
		return nil, err
	}

	tx := new(types.Transaction)
	if err = tx.UnmarshalBinary(rawTxBytes); err != nil {
		return nil, err
	}

	return tx, nil
}
