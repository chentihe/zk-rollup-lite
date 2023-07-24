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

	for i := 0; i < 2; i++ {
		p.A[i] = a[i]
	}

	for i, arr := range proof.B {
		b, err := stringsToArrayBigInt(arr)
		if err != nil {
			return err
		}

		if i < 2 {
			p.B[i] = ([2]*big.Int)(b)
		}
	}

	c, err := stringsToArrayBigInt(proof.C)
	if err != nil {
		return err
	}

	for i := 0; i < 2; i++ {
		p.C[i] = c[i]
	}

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
