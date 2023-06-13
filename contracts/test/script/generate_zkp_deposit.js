const { ethers } = require("ethers");
const { genPrivKey } = require("maci-crypto");
const { eddsa, poseidon, smt } = require("circomlibjs");
const {loadFiles, fillSiblings, generateGroth16Proof} = require("./utils")

const generateZkp = async () => {
    const fileNames = loadFiles("deposit")
    const depth = 6;

    const tree = await smt.newMemEmptyTrie()
    const encoder = ethers.AbiCoder.defaultAbiCoder();
    const privateKey = genPrivKey().toString();
    const publicKey = eddsa.prv2pub(privateKey)
    const leaf = poseidon([publicKey[0], publicKey[1], 1e18, 0]);
    const res = await tree.insert(0, leaf);
    const siblings = fillSiblings(depth, res.siblings)

    const circuitInputs = {
        balanceTreeRoot: res.oldRoot,
        publicKey: publicKey,
        balance: BigInt(1e18),
        nonce: 0,
        pathElements: siblings,
        oldKey: res.isOld0 ? 0 : res.oldKey,
        oldValue: res.isOld0 ? 0 : res.oldValue,
        isOld0: res.isOld0 ? 1 : 0,
        newKey: res.oldKey,
        func: [1, 0]
    }

    const {proof, publicSignals} = await generateGroth16Proof(circuitInputs, fileNames)

    process.stdout.write(encoder.encode(['uint256[2]', 'uint256[2][2]', 'uint256[2]', 'uint256[5]'], [proof.pi_a, proof.pi_b, proof.pi_c, publicSignals]));
    process.exit(0)
}

generateZkp()