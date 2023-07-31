const { poseidon, smt } = require("circomlibjs")
const { ethers } = require("ethers")
const { Scalar } = require("ffjavascript")
const {generateUser, generateTransaction, fillSiblings, formatTx, loadFiles, generateGroth16Proof} = require("./utils")

const generateZkpRollup = async () => {
    const fileNames = loadFiles("tx")
    const depth = 6;

    // Create mock users
    const users = []
    for (let i = 0; i < 3; i++) {
        users.push(generateUser(i))
    }

    // Create mock leaves
    const balanceTreeLeaves = []
    users.forEach(user => {
        var leaf
        if (user.index == 0) {
            leaf = 0
        } else {
            leaf = {
                publicKey: user.publicKey,
                balance: BigInt(10e18),
                nonce: 0
            }
        }
        balanceTreeLeaves.push(leaf)
    })

    // Init merkle tree
    const tree = await smt.newMemEmptyTrie()
    for (let i = 0; i < 3; i++) {
        if (i == 0) {
            await tree.insert(i, 0)
            continue
        }
        const leaf = balanceTreeLeaves[i]
        const hashedLeaf =poseidon([leaf.publicKey[0], leaf.publicKey[1], leaf.balance, leaf.nonce])
        await tree.insert(i, hashedLeaf)
    }

    // Get abi encoder
    const encoder = ethers.AbiCoder.defaultAbiCoder();

    // Circuit inputs
    const circuitInputs = {
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

    for (let i = 0; i < 2; i++) {
        // Create a new transaction
        const tx = generateTransaction(balanceTreeLeaves, users)

        // Get current balance tree root
        const balanceTreeRoot = tree.root

        // Get sender and recipient siblings
        const senderSiblings = fillSiblings(depth, (await tree.find(tx.from)).siblings)
        const recipientSiblings = fillSiblings(depth, (await tree.find(tx.to)).siblings)

        // Update txSender
        const senderLeaf = balanceTreeLeaves[tx.from]
        const oldSenderLeaf = Object.assign({}, senderLeaf)
        senderLeaf.nonce = tx.nonce
        senderLeaf.balance = Scalar.sub(Scalar.sub(senderLeaf.balance, tx.amount), tx.fee)

        const intermediateLeaf = [...users[tx.from].publicKey, senderLeaf.balance, tx.nonce]
        const hashedIntermediateLeaf = poseidon(intermediateLeaf)
        
        await tree.update(tx.from, hashedIntermediateLeaf)
        balanceTreeLeaves[tx.from] = senderLeaf

        // Get intermediate balance tree root & siblings
        const intermediateBalanceTreeRoot = tree.root
        const intermediateBalanceTreeSiblings = fillSiblings(depth, (await tree.find(tx.to)).siblings)

        // Update txRecipient
        const recipientLeaf = balanceTreeLeaves[tx.to]
        const oldRecipientLeaf = Object.assign({}, recipientLeaf)
        recipientLeaf.balance = Scalar.add(recipientLeaf.balance, tx.amount)
        const finalRecipientLeaf = [...users[tx.to].publicKey, recipientLeaf.balance, recipientLeaf.nonce]
        const hashedFinalRecipientLeaf = poseidon(finalRecipientLeaf)

        await tree.update(tx.to, hashedFinalRecipientLeaf)
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

    const {proof, publicSignals} = await generateGroth16Proof(circuitInputs, fileNames)
    
    process.stdout.write(encoder.encode(['uint256[2]', 'uint256[2]', 'uint256[2]', 'uint256[2][2]', 'uint256[2]', 'uint256[19]'], [users[0].publicKey, users[1].publicKey, proof.pi_a, proof.pi_b, proof.pi_c, publicSignals]));
    process.exit(0)
}

generateZkpRollup()