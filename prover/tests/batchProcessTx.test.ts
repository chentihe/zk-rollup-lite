import path from "path"
const { smt, poseidon } = require("circomlibjs")
const { Scalar } = require("ffjavascript")

const wasm_tester = require("circom_tester").wasm

const { 
    generateUser,
    insertTree,
    generateTransaction,
    fillSiblings,
    formatTx,
    stringifyBigInts 
} = require("./utils/helpers")

describe("batchProcessTx.circom", () => {
    it("BatchProcessTx(2, 6)", async () => {
        // Compile the circuit
        const circuit = await wasm_tester(
            path.join(__dirname, "circuits", "batchProcessTxTest.circom"),
            {
                output: path.join(__dirname, "../../operator/build", "batchProcessTxTest"),
                recompile: true,
                reduceConstraints: false,
            }
        )

        const numberOfUsers = 10
        const batchSize = 2
        const depth = 6

        // Create users
        const users = []
        for (let i = 0; i < numberOfUsers; i++) {
            users.push(generateUser(i))
        }

        // Create leaves
        const balanceTreeLeaves: any[] = []
        users.forEach(user => {
            const leaf = {
                publicKey: user.publicKey,
                balance: BigInt(1000e18),
                nonce: 0,
            }
            balanceTreeLeaves.push(leaf)
        })

        // Create the balance tree
        const balanceTree = await smt.newMemEmptyTrie()
        for (let i = 0; i < numberOfUsers; i++) {
            await insertTree(balanceTree, i, balanceTreeLeaves[i])
        }

        const circuitInputs: any = {
            balanceTreeRoots: [],
            txData: [],
            txSendersPublicKey: [],
            txSendersBalance: [],
            txSendersNonce: [],
            txSendersPathElements: [],
            txRecipientsPublicKey: [],
            txRecipientsBalance: [],
            txRecipientsNonce: [],
            txRecipientsPathElements: [],
            intermediateBalanceTreeRoots: [],
            intermediateBalanceTreesPathElements: []
        }
        for (let i = 0; i < batchSize; i++) {
            // Create a new transaction
            const tx = generateTransaction(balanceTreeLeaves, users)

            // Get current balance tree root
            const balanceTreeRoot = balanceTree.root

            // Get sender and recipient siblings
            const senderSiblings = fillSiblings(depth, (await balanceTree.find(tx.from)).siblings)
            const recipientSiblings = fillSiblings(depth, (await balanceTree.find(tx.to)).siblings)

            // Update txSender
            const senderLeaf = balanceTreeLeaves[tx.from]
            const oldSenderLeaf = Object.assign({}, senderLeaf)
            senderLeaf.nonce = tx.nonce
            senderLeaf.balance = Scalar.sub(Scalar.sub(senderLeaf.balance, tx.amount), tx.fee)

            const intermediateLeaf = [...users[tx.from].publicKey, senderLeaf.balance, tx.nonce]
            const hashedIntermediateLeaf = poseidon(intermediateLeaf)
            
            await balanceTree.update(tx.from, hashedIntermediateLeaf)
            balanceTreeLeaves[tx.from] = senderLeaf

            // Get intermediate balance tree root & siblings
            const intermediateBalanceTreeRoot = balanceTree.root
            const intermediateBalanceTreeSiblings = fillSiblings(depth, (await balanceTree.find(tx.to)).siblings)

            // Update txRecipient
            const recipientLeaf = balanceTreeLeaves[tx.to]
            const oldRecipientLeaf = Object.assign({}, recipientLeaf)
            recipientLeaf.balance = Scalar.add(recipientLeaf.balance, tx.amount)
            const finalRecipientLeaf = [...users[tx.to].publicKey, recipientLeaf.balance, recipientLeaf.nonce]
            const hashedFinalRecipientLeaf = poseidon(finalRecipientLeaf)

            await balanceTree.update(tx.to, hashedFinalRecipientLeaf)
            balanceTreeLeaves[tx.to] = recipientLeaf

            // Update circuit inputs
            circuitInputs.balanceTreeRoots.push(balanceTreeRoot)
            circuitInputs.txData.push(formatTx(tx))
            circuitInputs.txSendersPublicKey.push(users[tx.from].publicKey)
            circuitInputs.txSendersBalance.push(oldSenderLeaf.balance)
            circuitInputs.txSendersNonce.push(oldSenderLeaf.nonce)
            circuitInputs.txSendersPathElements.push(senderSiblings)
            circuitInputs.txRecipientsPublicKey.push(users[tx.to].publicKey)
            circuitInputs.txRecipientsBalance.push(oldRecipientLeaf.balance)
            circuitInputs.txRecipientsNonce.push(oldRecipientLeaf.nonce)
            circuitInputs.txRecipientsPathElements.push(recipientSiblings)
            circuitInputs.intermediateBalanceTreeRoots.push(intermediateBalanceTreeRoot)
            circuitInputs.intermediateBalanceTreesPathElements.push(intermediateBalanceTreeSiblings)
        }

        const witness = await circuit.calculateWitness(stringifyBigInts(circuitInputs))
        await circuit.checkConstraints(witness)
    })
})