import path from "path";
import { genPrivKey } from "maci-crypto";
const { eddsa, smt, poseidon } = require("circomlibjs");
const { Scalar } = require("ffjavascript")
const wasm_tester = require("circom_tester").wasm;
const { 
    fillSiblings, 
    stringifyBigInts, 
    formatTx,
    genPublicKey,
    insertTree
} = require("./utils/helpers")

describe("processTx.circom", () => {
    it("ProcessTx(4)",async () => {
        // Compile the circuit
        const circuit = await wasm_tester(
            path.join(__dirname, "circuits", "processTxTest.circom"),
            {
                output: path.join(__dirname, "../build", "processTxTest"),
                recompile: true,
                reduceConstraints: false,
            },
        )

        const depth = 4

        // Create mock users
        const userAIndex = 0
        const userABalance = BigInt(50e18)
        const userANonce = 1
        const privA = genPrivKey().toString()
        const pubA = genPublicKey(privA)

        const userBIndex = 1
        const userBBalance = BigInt(0)
        const userBNonce = 0
        const privB = genPrivKey().toString()
        const pubB = genPublicKey(privB)
    
        const sendAmount = BigInt(20e18)
        const feeAmount = BigInt(0.5e18)
        
        const balanceTree = await smt.newMemEmptyTrie()
        for (let i = 0; i < depth; i++) {
            const leaf = {
                publicKey: genPublicKey(genPrivKey().toString()),
                balance: BigInt(0),
                nonce: 0
            }
            if (i === userAIndex) {
                leaf.publicKey = pubA
                leaf.balance = userABalance
                leaf.nonce = userANonce
            } else if (i === userBIndex) {
                leaf.publicKey = pubB
                leaf.balance = userBBalance
                leaf.nonce = userBNonce
            }
            await insertTree(balanceTree, i, leaf)
        }

        const balanceTreeRoot = balanceTree.root
        const userASiblings = fillSiblings(depth, (await balanceTree.find(userAIndex)).siblings)
        const userBSiblings = fillSiblings(depth, (await balanceTree.find(userBIndex)).siblings)

        // Create the transaction
        const rawTx = {
            from: userAIndex,
            to: userBIndex,
            amount: sendAmount,
            fee: feeAmount,
            nonce: userANonce + 1
        }
        const hashedTx = poseidon([userAIndex, userBIndex, sendAmount, feeAmount, rawTx.nonce])
        const signature = eddsa.signMiMC(privA, hashedTx)
        const tx = Object.assign({}, rawTx, {signature})

        // IntermediateBalanceRoot
        // Update txSender
        const intermediateUserALeaf = [...pubA, Scalar.sub(Scalar.sub(userABalance, sendAmount), feeAmount), rawTx.nonce]
        const hashedIntermediateUserALeaf = poseidon(intermediateUserALeaf)

        await balanceTree.update(rawTx.from, hashedIntermediateUserALeaf)

        // Update txRecipient
        const intermediateBalanceTreeRoot = balanceTree.root
        const intermediateBalanceTreeSiblings = fillSiblings(depth, (await balanceTree.find(rawTx.to)).siblings)
        const finalUserBLeaf = [...pubB, Scalar.add(userBBalance, sendAmount), userBNonce]
        const hashedFinalUserBLeaf = poseidon(finalUserBLeaf)
        await balanceTree.update(rawTx.to, hashedFinalUserBLeaf)

        // Create circuit inputs
        // All input format must be string
        const circuitInputs = stringifyBigInts({
            balanceTreeRoot: balanceTreeRoot,
            txData: formatTx(tx),
            txSenderPublicKey: pubA,
            txSenderBalance: userABalance,
            txSenderNonce: userANonce,
            txSenderPathElements: userASiblings,
            txRecipientPublicKey: pubB,
            txRecipientBalance: userBBalance,
            txRecipientNonce: userBNonce,
            txRecipientPathElements: userBSiblings,
            intermediateBalanceTreeRoot: intermediateBalanceTreeRoot,
            intermediateBalanceTreePathElements: intermediateBalanceTreeSiblings
        })

        const witness = await circuit.calculateWitness(circuitInputs)
        await circuit.checkConstraints(witness)
    });
})