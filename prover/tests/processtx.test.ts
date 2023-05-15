import path from "path";
import { genPrivKey } from "maci-crypto";
import { expect } from "chai";
const { eddsa, smt, poseidon } = require("circomlibjs");
const { Scalar } = require("ffjavascript")

const wasm_tester = require("circom_tester").wasm;

describe("processtx.circom", () => {
    it("ProcessTx(4)",async () => {
        const circuit = await wasm_tester(
            path.join(__dirname, "circuits", "processtxTest.circom"),
            {
                output: path.join(__dirname, "../build", "processtxTest"),
                recompile: true,
                reduceConstraints: false,
            },
        )

        const userAIndex = 0;
        const userABalance = BigInt(50e18);
        const userANonce = 1;
        const privA = genPrivKey().toString();
        const pubA = genPublicKey(privA);

        const userBIndex = 1;
        const userBBalance = BigInt(0);
        const userBNonce = 0;
        const privB = genPrivKey().toString();
        const pubB = genPublicKey(privB);
    
        const sendAmount = BigInt(20e18);
        const feeAmount = BigInt(0.5e18);
        
        const balanceTree = await smt.newMemEmptyTrie()
        await insertTree(balanceTree, userAIndex, pubA, userABalance, userANonce)
        await insertTree(balanceTree, userBIndex, pubB, userBBalance, userBNonce)
        await insertTree(balanceTree, 2, genPublicKey(genPrivKey().toString()), BigInt(0), 0)
        await insertTree(balanceTree, 3, genPublicKey(genPrivKey().toString()), BigInt(0), 0)

        const balanceTreeRoot = balanceTree.root;
        const userASiblings = (await balanceTree.find(userAIndex)).siblings
        const userBSiblings = (await balanceTree.find(userBIndex)).siblings

        const rawTx = {
            from: userAIndex,
            to: userBIndex,
            amount: sendAmount,
            fee: feeAmount,
            nonce: userANonce + 1
        }
        const hashedTx = poseidon([userAIndex, userBIndex, sendAmount, feeAmount, userANonce + 1])
        const signature = eddsa.signMiMC(privA, hashedTx)
        const tx = Object.assign({}, rawTx, {signature})

        const intermediateUserALeaf = [...pubA, Scalar.sub(Scalar.sub(userABalance, sendAmount), feeAmount), userANonce + 1]
        const hashedIntermediateUserALeaf = poseidon(intermediateUserALeaf)

        await balanceTree.update(userAIndex, hashedIntermediateUserALeaf)

        const intermediateBalanceTreeRoot = balanceTree.root;
        const intermediateBalanceTreeSiblings = (await balanceTree.find(userBIndex)).siblings
        const finalUserBLeaf = [...pubB, Scalar.add(userBBalance, sendAmount), userBNonce]
        const hashedFinalUserBLeaf = poseidon(finalUserBLeaf)
        await balanceTree.update(userBIndex, hashedFinalUserBLeaf)

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

        // TODO: fill the empty leaf to meet the path of siblings
        const witness = await circuit.calculateWitness(circuitInputs);
        const outputIdx = circuit.getSignalIdx("main.newBalanceTreeRoot")
        const newBalanceTreeRootCircom = witness[outputIdx]

        expect(balanceTree.root.toString()).equal(newBalanceTreeRootCircom.toString())
    });

    function genPublicKey(privateKey: string) {
        return eddsa.prv2pub(privateKey).map((each: { toString: () => any; }) => each.toString());
    }

    async function insertTree(tree: any , index: number, publicKey: any, balance: BigInt, nonce: number) {
        const leaf = poseidon([publicKey[0], publicKey[1], balance, nonce])
        await tree.insert(index, leaf)
    }

    function formatTx(tx: any) {
        return [
            tx.from,
            tx.to,
            tx.amount,
            tx.fee,
            tx.nonce,
            tx.signature !== undefined ? tx.signature.R8[0] : null,
            tx.signature !== undefined ? tx.signature.R8[1] : null,
            tx.signature !== undefined ? tx.signature.S : null
        ].filter((x: any) => x !== null)
        .map((x: any) => BigInt(x));
    }

    function stringifyBigInts(inputs: any): any {
        if (typeof inputs == "bigint") {
            return inputs.toString(10);
        } else if (Array.isArray(inputs)) {
            return inputs.map(stringifyBigInts);
        } else if (typeof inputs == "object") {
            const keys = Object.keys(inputs);
            const res = Object.assign({}, inputs)
            keys.forEach((key: any) => {
                res[key] = stringifyBigInts(inputs[key]);
            });
            return res;
        } else {
            return inputs
        }
    }
})