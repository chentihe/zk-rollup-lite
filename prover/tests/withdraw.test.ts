import { genPrivKey } from "maci-crypto"
import { genPublicKey, stringifyBigInts, arrayifySignature } from "./utils/helpers"
import path from "path"
const { eddsa } = require("circomlibjs");
const wasm_tester = require("circom_tester").wasm

describe("withdraw.circom", () => {
    it("Withdraw()", async () => {
        // Compile the circuit
        const circuit = await wasm_tester(
            path.join(__dirname, "../circuits", "withdraw.circom"),
            {
                output: path.join(__dirname, "../build", "withdraw"),
                recompile: true,
                reduceConstraints: false
            }
        )

        // Create the mock user & nullifier
        const privA = genPrivKey().toString()
        const pubA = genPublicKey(privA)
        const nullifier = genPrivKey()

        // Sign the nullifier
        const signature = arrayifySignature(eddsa.signMiMC(privA, nullifier))

        // Create circuit inputs
        // All input format must be string
        const circuitInputs = stringifyBigInts({
            publicKey: pubA,
            signature,
            nullifier
        })

        const witness = await circuit.calculateWitness(circuitInputs)
        await circuit.checkConstraints(witness)
    })
})