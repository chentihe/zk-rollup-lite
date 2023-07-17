package circuits

import (
	"os"

	"github.com/iden3/go-rapidsnark/prover"
	"github.com/iden3/go-rapidsnark/types"
	"github.com/iden3/go-rapidsnark/witness/v2"
	"github.com/iden3/go-rapidsnark/witness/wasmer"
)

func GenerateGroth16Proof(circuitInput []byte, circuitPath string) (*types.ZKProof, error) {
	inputJSON, err := witness.ParseInputs(circuitInput)
	if err != nil {
		return nil, err
	}

	// choose wasmer / wazoro as wasm engine
	opts := witness.WithWasmEngine(wasmer.NewCircom2WitnessCalculator)

	wasmBytes, err := os.ReadFile(circuitPath + wasmFilePath)
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

	zkeyBytes, err := os.ReadFile(circuitPath + zkeyFilePath)
	if err != nil {
		return nil, err
	}

	proof, err := prover.Groth16Prover(zkeyBytes, wtns)
	if err != nil {
		return nil, err
	}

	return proof, nil
}
