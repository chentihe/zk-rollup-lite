const { groth16 } = require("snarkjs")
const path = require("path")
const fs = require("fs")
const { eddsa, poseidon } = require("circomlibjs")
const { genPrivKey } = require("maci-crypto")
const { ethers } = require("ethers")
const { IncrementalMerkleTree } = require("@zk-kit/incremental-merkle-tree");

const arrayifySignature = (signature) => {
    return [...signature.R8, signature.S]
}

const generateZkp = async (action, amount) => {
    const dir = path.join(__dirname, "../../../prover/build", "withdraw")
    const fileNames = fs.readdirSync(dir)
        .map(file => {
            const concatFile = dir + "/" + file
            return concatFile
        }).filter(file => {
            const extension = file.split(".")[1]
            return (extension == "wasm" || extension == "zkey" || extension == "json")
        })
    const ZERO_VALUE = BigInt(0)
    const tree = new IncrementalMerkleTree(poseidon, 6, ZERO_VALUE, 2)
    const encoder = ethers.AbiCoder.defaultAbiCoder();
    const privateKey = genPrivKey().toString();
    // Uint8Array => BigInt
    const publicKey = eddsa.prv2pub(privateKey).flatMap(axis => ethers.toBigInt(axis))

    const nullifier = genPrivKey()
    const signature = arrayifySignature(eddsa.signMiMC(privateKey, nullifier))
    const circuitInputs = {
        publicKey,
        signature,
        nullifier
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

    const leaf = ethers.toBigInt(poseidon([publicKey[0], publicKey[1], amount * 10**18 * 2, 0]));
    tree.insert(leaf);
    let newLeaf
    switch (action) {
        case "withdraw":
            newLeaf = ethers.toBigInt(poseidon([publicKey[0], publicKey[1], amount * 10**18, 0]));
            break;
        case "withdraw-all":
            newLeaf = ethers.toBigInt(poseidon([publicKey[0], publicKey[1], amount * 10**18 * 2, 0]));
            break;
    }
        
    tree.update(0, newLeaf);
    
    const {pathIndices, siblings} = tree.createProof(0)
    const bigIntSiblings = []
    bigIntSiblings.push(siblings[0][0])
    siblings.shift()
    siblings.forEach(sibling => {
        const bigIntSibling = ethers.toBigInt(sibling[0])
        bigIntSiblings.push(bigIntSibling)
    })

    process.stdout.write(encoder.encode(['uint256[2]', 'uint8[6]', 'uint256[6]', 'uint256[2]', 'uint256[2][2]', 'uint256[2]', 'uint256[3]'], [publicKey, pathIndices, bigIntSiblings, packedSolidityProof.pi_a, packedSolidityProof.pi_b, packedSolidityProof.pi_c, publicSignals]));
    process.exit(0)
}

const action = process.argv[2]

const amount = process.argv[3]

generateZkp(action, amount)