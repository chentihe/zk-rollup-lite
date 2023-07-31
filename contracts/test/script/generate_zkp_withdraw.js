const { eddsa, poseidon, smt } = require("circomlibjs")
const { genPrivKey } = require("maci-crypto")
const { ethers } = require("ethers")
const {arrayifySignature, loadFiles, generateGroth16Proof, fillSiblings } = require("./utils")

const generateZkp = async (action, amount) => {
    const fileNames = loadFiles("withdraw")

    const depth = 6

    const tree = await smt.newMemEmptyTrie()
    const encoder = ethers.AbiCoder.defaultAbiCoder();
    const privateKey = genPrivKey().toString();
    const publicKey = eddsa.prv2pub(privateKey)
    const nullifier = genPrivKey()
    const signature = arrayifySignature(eddsa.signMiMC(privateKey, nullifier))
    const balance = amount * 10**18 * 2
    
    const leaf = poseidon([publicKey[0], publicKey[1], balance, 0]);
    await tree.insert(1, leaf)
    
    const res = await tree.find(1)
    const siblings = fillSiblings(depth, res.siblings)

    const circuitInputs = {
        balanceTreeRoot: tree.root,
        signature,
        nullifier,
        publicKey,
        balance,
        nonce: 0,
        pathElements: siblings,
        oldKey: 0,
        oldValue: 0,
        newKey: 1,
    }

    const {proof, publicSignals} = await generateGroth16Proof(circuitInputs, fileNames)
    
    process.stdout.write(encoder.encode(['uint256[2]', 'uint256[2][2]', 'uint256[2]', 'uint256[3]'], [proof.pi_a, proof.pi_b, proof.pi_c, publicSignals]));
    process.exit(0)
}
        
const action = process.argv[2]

const amount = process.argv[3]

generateZkp(action, amount)