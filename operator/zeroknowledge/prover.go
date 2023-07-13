package zeroknowledge

import (
	"encoding/json"
	"math/big"
	"os"

	"github.com/iden3/go-rapidsnark/prover"
	"github.com/iden3/go-rapidsnark/types"
	"github.com/iden3/go-rapidsnark/witness/v2"
	"github.com/iden3/go-rapidsnark/witness/wasmer"
)

type BigIntProofData struct {
	A [2]*big.Int
	B [2][2]*big.Int
	C [2]*big.Int
}

type BigIntProof struct {
	Proof         *BigIntProofData
	PublicSignals [19]*big.Int
}

func GenerateGroth16Proof(circuitInput *CircuitInput) (*types.ZKProof, error) {
	bytesInput, err := json.Marshal(circuitInput)
	if err != nil {
		return nil, err
	}

	inputJSON, err := witness.ParseInputs(bytesInput)
	if err != nil {
		return nil, err
	}

	opts := witness.WithWasmEngine(wasmer.NewCircom2WitnessCalculator)

	wasmBytes, err := os.ReadFile(wasmFilePath)
	if err != nil {
		return nil, err
	}

	calc, err := witness.NewCalculator(wasmBytes, opts)
	if err != nil {
		return nil, err
	}

	wtns, err := calc.CalculateWTNSBin(inputJSON, true)
	if err != nil {
		return nil, err
	}

	zkeyBytes, err := os.ReadFile(zkeyFilePath)
	if err != nil {
		return nil, err
	}

	proof, err := prover.Groth16Prover(zkeyBytes, wtns)
	if err != nil {
		return nil, err
	}

	return proof, nil
}

func ParseProofToBigInt(proof *types.ZKProof) (*BigIntProof, error) {
	inputs, err := stringsToArrayBigInt(proof.PubSignals)
	if err != nil {
		return nil, err
	}

	publicSignals := ([19]*big.Int)(inputs)

	a, err := stringsToArrayBigInt(proof.Proof.A)
	if err != nil {
		return nil, err
	}

	piA := ([2]*big.Int)(a)

	var piB [2][2]*big.Int

	for i, arr := range proof.Proof.B {
		b, err := stringsToArrayBigInt(arr)
		if err != nil {
			return nil, err
		}

		piB[i] = ([2]*big.Int)(b)
	}

	c, err := stringsToArrayBigInt(proof.Proof.C)
	if err != nil {
		return nil, err
	}

	piC := ([2]*big.Int)(c)

	bigIntProof := &BigIntProofData{
		A: piA,
		B: piB,
		C: piC,
	}

	return &BigIntProof{
		Proof:         bigIntProof,
		PublicSignals: publicSignals,
	}, nil
}
