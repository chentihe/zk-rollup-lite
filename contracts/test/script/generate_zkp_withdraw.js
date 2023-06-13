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
    
    const leaf = poseidon([publicKey[0], publicKey[1], amount * 10**18 * 2, 0]);
    await tree.insert(0, leaf)
    
    let withdrawal
    let balance
    switch (action) {
        case "withdraw":
            withdrawal = amount * 10**18
            balance = amount * 10**18
            break;
        case "withdraw-all":
            withdrawal = amount * 10**18 * 2
            balance = 0
            break;
    }
            
    const newLeaf = poseidon([publicKey[0], publicKey[1], withdrawal, 0]);
    const res = await tree.update(0, newLeaf);
    const siblings = fillSiblings(depth, res.siblings)

    const circuitInputs = {
        publicKey,
        signature,
        nullifier,
        balanceTreeRoot: res.oldRoot,
        balance,
        nonce: 0,
        pathElements: siblings,
        oldKey: res.isOld0 ? 0 : res.oldKey,
        oldValue: res.isOld0 ? 0 : res.oldValue,
        isOld0: res.isOld0 ? 1 : 0,
        newKey: res.isOld0 ? res.oldKey : res.newKey,
        func: [0, 1]
    }

    const {proof, publicSignals} = await generateGroth16Proof(circuitInputs, fileNames)
    
    process.stdout.write(encoder.encode(['uint256[2]', 'uint256[2][2]', 'uint256[2]', 'uint256[5]'], [proof.pi_a, proof.pi_b, proof.pi_c, publicSignals]));
    process.exit(0)
}
        
const action = process.argv[2]

const amount = process.argv[3]

generateZkp(action, amount)