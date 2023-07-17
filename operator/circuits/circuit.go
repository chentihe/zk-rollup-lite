package circuits

import (
	"math/big"

	"github.com/iden3/go-rapidsnark/types"
)

type ProofData struct {
	A [2]*big.Int
	B [2][2]*big.Int
	C [2]*big.Int
}

func (p *ProofData) ProofUnmarshal(proof *types.ProofData) error {
	a, err := stringsToArrayBigInt(proof.A)
	if err != nil {
		return err
	}

	p.A = ([2]*big.Int)(a)

	for i, arr := range proof.B {
		b, err := stringsToArrayBigInt(arr)
		if err != nil {
			return err
		}

		p.B[i] = ([2]*big.Int)(b)
	}

	c, err := stringsToArrayBigInt(proof.C)
	if err != nil {
		return err
	}

	p.C = ([2]*big.Int)(c)

	return nil
}

type Input interface {
	InputMarshaller
}

type InputMarshaller interface {
	InputMarshal() ([]byte, error)
}

type Output interface {
	OutputUnmarshaller
}

type OutputUnmarshaller interface {
	OutputUnmarshal(zkp *types.ZKProof) error
}
