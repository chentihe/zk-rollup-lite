const { groth16 } = require("snarkjs")
const path = require("path")
const fs = require("fs")
const { eddsa, poseidon } = require("circomlibjs")
const { genPrivKey } = require("maci-crypto")
const { ethers } = require("ethers")
const { IncrementalMerkleTree } = require("@zk-kit/incremental-merkle-tree");
const { Scalar } = require("ffjavascript")

const generateUser = (index) => {
    const privateKey = genPrivKey().toString()
    // Uint8Array => BigInt
    const publicKey = eddsa.prv2pub(privateKey).flatMap(axis => ethers.toBigInt(axis))

    return {
        index,
        privateKey,
        publicKey
    }
}

const generateTransaction = (leaves, users) => {
    const sender = users[0]
    const recipient = users[1]
    const amount = BigInt(1e18)
    const fee = BigInt(0.5e18)
    const senderNonce = leaves[sender.index].nonce

    const rawTx = {
        from: sender.index,
        to: recipient.index,
        amount,
        fee,
        nonce: senderNonce + 1
    }
    const hashedTx = poseidon([rawTx.from, rawTx.to, rawTx.amount, rawTx.fee, rawTx.nonce])
    const signature = eddsa.signMiMC(sender.privateKey, hashedTx)
    const tx = Object.assign({}, rawTx, {signature})

    return tx
}

const formatTx = (tx) => {
    return [
        tx.from,
        tx.to,
        tx.amount,
        tx.fee,
        tx.nonce,
        tx.signature !== undefined ? tx.signature.R8[0] : null,
        tx.signature !== undefined ? tx.signature.R8[1] : null,
        tx.signature !== undefined ? tx.signature.S : null
    ].filter((x) => x !== null)
    .map((x) => BigInt(x))
}

const generateZkpRollup = async (amount) => {
    // Read wasm, zkey, json files
    const dir = path.join(__dirname, "../../../prover/build", "tx")
    const fileNames = fs.readdirSync(dir)
        .map(file => {
            const concatFile = dir + "/" + file
            return concatFile
        }).filter(file => {
            const extension = file.split(".")[1]
            return (extension == "wasm" || extension == "zkey")
        })
    
    // Set zero value
    const ZERO_VALUE = BigInt(0)
    
    // Create mock users
    const users = []
    for (let i = 0; i < 2; i++) {
        users.push(generateUser(i))
    }

    // Create mock leaves
    const balanceTreeLeaves = []
    users.forEach(user => {
        const leaf = {
            publicKey: user.publicKey,
            balance: BigInt(10e18),
            nonce: 0
        }
        balanceTreeLeaves.push(leaf)
    })

    // Init merkle tree
    const tree = new IncrementalMerkleTree(poseidon, 6, ZERO_VALUE, 2)
    for (let i = 0; i < 2; i++) {
        const leaf = balanceTreeLeaves[i]
        const hashedLeaf = ethers.toBigInt(poseidon([leaf.publicKey[0], leaf.publicKey[1], leaf.balance, leaf.nonce]))
        tree.insert(hashedLeaf)
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

    // Pathindices for the verifier
    // TODO: add sender updated index & recipient updated index
    const pathIndices = []

    for (let i = 0; i < 2; i++) {
        // Create a new transaction
        const tx = generateTransaction(balanceTreeLeaves, users)

        // Get current balance tree root
        const balanceTreeRoot = tree.root

        // Update txSender
        const senderLeaf = balanceTreeLeaves[tx.from]
        const oldSenderLeaf = Object.assign({}, senderLeaf)
        senderLeaf.nonce = tx.nonce
        senderLeaf.balance = Scalar.sub(Scalar.sub(senderLeaf.balance, tx.amount), tx.fee)

        const intermediateLeaf = [...users[tx.from].publicKey, senderLeaf.balance, tx.nonce]
        const hashedIntermediateLeaf = ethers.toBigInt(poseidon(intermediateLeaf))

        tree.update(tx.from, hashedIntermediateLeaf)
        balanceTreeLeaves[tx.from] = senderLeaf

        const senderProof = tree.createProof(0)
        const senderSiblings = senderProof.siblings
        pathIndices.push(senderProof.pathIndices)

        // Get intermediate balance tree root & siblings
        const intermediateBalanceTreeRoot = tree.root
        const intermediateBalanceTreeSiblings = tree.createProof(tx.to).siblings

        // Update txRecipient
        const recipientLeaf = balanceTreeLeaves[tx.to]
        const oldRecipientLeaf = Object.assign({}, recipientLeaf)
        recipientLeaf.balance = Scalar.add(recipientLeaf.balance, tx.amount)
        const finalRecipientLeaf = [...users[tx.to].publicKey, recipientLeaf.balance, recipientLeaf.nonce]
        const hashedFinalRecipientLeaf = ethers.toBigInt(poseidon(finalRecipientLeaf))

        tree.update(tx.to, hashedFinalRecipientLeaf)
        balanceTreeLeaves[tx.to] = recipientLeaf

        const recipientProof = tree.createProof(1)
        const recipientSiblings = recipientProof.siblings
        pathIndices.push(recipientProof.pathIndices)

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

    const {proof, publicSignals} = await groth16.fullProve(circuitInputs, fileNames[0], fileNames[1])

    const packedSolidityProof = {
        pi_a: [proof.pi_a[0], proof.pi_a[1]],
        pi_b: [
            [proof.pi_b[0][1], proof.pi_b[0][0]],
            [proof.pi_b[1][1], proof.pi_b[1][0]]
        ],
        pi_c: [proof.pi_c[0], proof.pi_c[1]]
    }
    
    // const bigIntSiblings = []
    // bigIntSiblings.push(siblings[0][0])
    // siblings.shift()
    // siblings.forEach(sibling => {
    //     const bigIntSibling = ethers.toBigInt(sibling[0])
    //     bigIntSiblings.push(bigIntSibling)
    // })

    // TODO: modify encode elements
    process.stdout.write(encoder.encode(['uint256[2]', 'uint8[6]', 'uint256[6]', 'uint256[2]', 'uint256[2][2]', 'uint256[2]', 'uint256[3]'], [publicKey, pathIndices, bigIntSiblings, packedSolidityProof.pi_a, packedSolidityProof.pi_b, packedSolidityProof.pi_c, publicSignals]));
    process.exit(0)
}

const amount = process.argv[2]

generateZkpRollup(amount)