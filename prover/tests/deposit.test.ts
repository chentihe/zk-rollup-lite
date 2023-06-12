import { genPrivKey } from "maci-crypto"
import { genPublicKey, stringifyBigInts, fillSiblings } from "./utils/helpers"
import path from "path"
const { smt, poseidon } = require("circomlibjs");
const wasm_tester = require("circom_tester").wasm

describe("deposit.circom", () => {
    it("Deposit()", async () => {
        // Compile the circuit
        const circuit = await wasm_tester(
            path.join(__dirname, "../circuits", "deposit.circom"),
            {
                output: path.join(__dirname, "../build", "deposit"),
                recompile: true,
                reduceConstraints: false
            }
        )

        const depth = 6;

        // Create the mock user
        const privA = genPrivKey().toString()
        const pubA = genPublicKey(privA)

        const balanceTree = await smt.newMemEmptyTrie()

        const leaf = poseidon([pubA[0], pubA[1], BigInt(1e18), 0])

        const res = await balanceTree.insert(0, leaf)

        const siblings = fillSiblings(depth, res.siblings)

        // Create circuit inputs
        // All input format must be string
        const circuitInputs = stringifyBigInts({
            balanceTreeRoot: res.oldRoot,
            publicKey: pubA,
            balance: BigInt(1e18),
            nonce: 0,
            pathElements: siblings,
            oldKey: res.isOld0 ? 0 : res.oldKey,
            oldValue: res.isOld0 ? 0 : res.oldValue,
            isOld0: res.isOld0 ? 1 : 0,
            newKey: res.oldKey,
            func: [1, 0]
        })

        const witness = await circuit.calculateWitness(circuitInputs)
        await circuit.checkConstraints(witness)
    })
})