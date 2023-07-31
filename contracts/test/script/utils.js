const path = require("path")
const fs = require("fs")
const {eddsa, poseidon} = require("circomlibjs")
const {genPrivKey} = require("maci-crypto")
const { groth16 } = require("snarkjs")

const arrayifySignature = (signature) => {
    return [...signature.R8, signature.S]
}

const fillSiblings = (depth, siblings) => {
    const diff = depth - siblings.length
    if (diff > 0) {
        for (let i = 0; i < diff; i++) {
            siblings.push(BigInt(0))
        }
    }
    return siblings
}

const generateUser = (index) => {
    const privateKey = genPrivKey().toString()
    // Uint8Array => BigInt
    const publicKey = eddsa.prv2pub(privateKey)

    return {
        index,
        privateKey,
        publicKey
    }
}

const generateTransaction = (leaves, users) => {
    const sender = users[1]
    const recipient = users[2]
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

const loadFiles = (fileName) => {
    // Read wasm, zkey, json files
    const dir = path.join(__dirname, "../../../prover/build", fileName)
    const fileNames = fs.readdirSync(dir)
        .map(file => {
            const concatFile = dir + "/" + file
            return concatFile
        }).filter(file => {
            const extension = file.split(".")[1]
            return (extension == "wasm" || extension == "zkey")
        })
    return fileNames
}

const generateGroth16Proof = async (circuitInputs, fileNames) => {
    const {proof, publicSignals} = await groth16.fullProve(circuitInputs, fileNames[0], fileNames[1])

    const packedSolidityProof = {
        pi_a: [proof.pi_a[0], proof.pi_a[1]],
        pi_b: [
            [proof.pi_b[0][1], proof.pi_b[0][0]],
            [proof.pi_b[1][1], proof.pi_b[1][0]]
        ],
        pi_c: [proof.pi_c[0], proof.pi_c[1]]
    }

    return {
        proof: packedSolidityProof, 
        publicSignals
    }
}

module.exports = {
    arrayifySignature,
    fillSiblings,
    generateUser,
    generateTransaction,
    formatTx,
    loadFiles,
    generateGroth16Proof
};