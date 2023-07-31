import { genPrivKey } from "maci-crypto"
import { genPublicKey, stringifyBigInts, arrayifySignature, insertTree, fillSiblings } from "./utils/helpers"
import path from "path"
const { eddsa, smt, poseidon } = require("circomlibjs");
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

        const depth = 6;

        // Create the mock user & nullifier
        const privA = genPrivKey().toString()
        const pubA = genPublicKey(privA)

        const balanceTree = await smt.newMemEmptyTrie()
        const leaf = poseidon([pubA[0], pubA[1], BigInt(1e18), 0])

        await balanceTree.insert(1, leaf)

        const newLeaf = poseidon([pubA[0], pubA[1], BigInt(0.5e18), 0])
        await balanceTree.update(1, newLeaf)
        const res = await balanceTree.find(1)

        const siblings = fillSiblings(depth, res.siblings)

        const nullifier = genPrivKey()

        // Sign the nullifier
        const signature = arrayifySignature(eddsa.signMiMC(privA, nullifier))

        // Create circuit inputs
        // All input format must be string
        const circuitInputs = stringifyBigInts({
            balanceTreeRoot: balanceTree.root,
            signature,
            nullifier,
            publicKey: pubA,
            balance: BigInt(0.5e18),
            nonce: 0,
            pathElements: siblings,
            oldKey: 0,
            oldValue: 0,
            newKey: 1
        })

        const witness = await circuit.calculateWitness(circuitInputs)
        await circuit.checkConstraints(witness)
    })
})