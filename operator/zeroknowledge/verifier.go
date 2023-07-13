package zeroknowledge

import (
	"os"

	"github.com/iden3/go-rapidsnark/types"
	"github.com/iden3/go-rapidsnark/verifier"
)

func VerifierGroth16(proof *types.ZKProof) error {
	// verify zkp
	vkeyBytes, err := os.ReadFile(verficationKeyFilePath)
	if err != nil {
		return err
	}

	if err := verifier.VerifyGroth16(*proof, vkeyBytes); err != nil {
		return err
	}

	return nil
}
