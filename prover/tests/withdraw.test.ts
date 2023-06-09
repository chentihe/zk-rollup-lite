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

        await balanceTree.insert(0, leaf)

        const newLeaf = poseidon([pubA[0], pubA[1], BigInt(0.5e18), 0])
        const res = await balanceTree.update(0, newLeaf)

        const siblings = fillSiblings(depth, res.siblings)

        const nullifier = genPrivKey()

        // Sign the nullifier
        const signature = arrayifySignature(eddsa.signMiMC(privA, nullifier))

        // Create circuit inputs
        // All input format must be string
        const circuitInputs = stringifyBigInts({
            publicKey: pubA,
            signature,
            nullifier,
            balanceTreeRoot: res.oldRoot,
            balance: BigInt(0.5e18),
            nonce: 0,
            pathElements: siblings,
            oldKey: res.isOld0 ? 0 : res.oldKey,
            oldValue: res.isOld0 ? 0 : res.oldValue,
            isOld0: res.isOld0 ? 1 : 0,
            newKey: res.newKey,
            func: [0, 1]
        })

        const witness = await circuit.calculateWitness(circuitInputs)
        await circuit.checkConstraints(witness)
    })
})