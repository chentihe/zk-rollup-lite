const { ethers } = require("ethers");
const { genPrivKey } = require("maci-crypto");
const { eddsa, poseidon } = require("circomlibjs");
const { IncrementalMerkleTree } = require("@zk-kit/incremental-merkle-tree");

const generatePublicKey = async () => {
    const ZERO_VALUE = BigInt(0)
    const tree = new IncrementalMerkleTree(poseidon, 6, ZERO_VALUE, 2)
    const encoder = ethers.AbiCoder.defaultAbiCoder();
    const privateKey = genPrivKey().toString();
    // Uint8Array => BigInt
    const publicKey = eddsa.prv2pub(privateKey).flatMap(axis => ethers.toBigInt(axis))
    const leaf = ethers.toBigInt(poseidon([publicKey[0], publicKey[1], 1e18, 0]));
    tree.insert(leaf);

    let newLeaf = ethers.toBigInt(poseidon([publicKey[0], publicKey[1], 2e18, 0]));
            
    tree.update(0, newLeaf);
    const {pathIndices, siblings} = tree.createProof(0)
    const bigIntSiblings = []
    bigIntSiblings.push(siblings[0][0])
    siblings.shift()
    siblings.forEach(sibling => {
        const bigIntSibling = ethers.toBigInt(sibling[0])
        bigIntSiblings.push(bigIntSibling)
    })
    process.stdout.write(encoder.encode(['uint256[2]', 'uint8[6]', 'uint256[6]'], [publicKey, pathIndices, bigIntSiblings]));
}

generatePublicKey()