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
			p.B[i][0], p.B[i][1] = p.B[i][1], p.B[i][0]
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

type Data struct {
	Input  InputsMarshaller
	Output OutputsUnmarshaller
}

type InputsMarshaller interface {
	InputsMarshal() ([]byte, error)
}

type OutputsUnmarshaller interface {
	OutputsUnmarshal(zkp *types.ZKProof) error
}
