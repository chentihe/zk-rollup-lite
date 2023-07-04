package txmanager

import (
	"encoding/hex"

	"github.com/iden3/go-iden3-crypto/babyjub"
)

// TODO: need two event handlers to update merkle tree
// 1. Deposit event
// 2. Withdraw event
// leaf: Poseidon([publicKey[0], publicKey[1], balance, nonce])

// Add: mt.Add(ctx, key, value)
// Update: mt.Update(ctx, key, value)

func DecodePublicKeyFromString(comp string) (*babyjub.PublicKey, error) {
	bytesPublicKey, err := hex.DecodeString(comp)

	publicKeyComp := babyjub.PublicKeyComp(bytesPublicKey)
	if err != nil {
		return nil, err
	}

	publicKey, err := publicKeyComp.Decompress()
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}
