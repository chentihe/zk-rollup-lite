const { ethers } = require("ethers");
const { genPrivKey } = require("maci-crypto");
const { eddsa, poseidon, smt } = require("circomlibjs");
const { groth16 } = require("snarkjs");
const path = require("path")
const fs = require("fs")

const generateZkp = async () => {
    // Read wasm, zkey, json files
    const dir = path.join(__dirname, "../../../prover/build", "deposit")
    const fileNames = fs.readdirSync(dir)
        .map(file => {
            const concatFile = dir + "/" + file
            return concatFile
        }).filter(file => {
            const extension = file.split(".")[1]
            return (extension == "wasm" || extension == "zkey")
        })
    
    const depth = 6;
    const tree = await smt.newMemEmptyTrie()
    const encoder = ethers.AbiCoder.defaultAbiCoder();
    const privateKey = genPrivKey().toString();
    const publicKey = eddsa.prv2pub(privateKey)
    const leaf = poseidon([publicKey[0], publicKey[1], 1e18, 0]);
    const res = await tree.insert(0, leaf);
    const siblings = res.siblings
    while (siblings.length < depth) siblings.push(0)

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

    const {proof, publicSignals} = await groth16.fullProve(circuitInputs, fileNames[0], fileNames[1])
    
    const packedSolidityProof = {
        pi_a: [proof.pi_a[0], proof.pi_a[1]],
        pi_b: [
            [proof.pi_b[0][1], proof.pi_b[0][0]],
            [proof.pi_b[1][1], proof.pi_b[1][0]]
        ],
        pi_c: [proof.pi_c[0], proof.pi_c[1]]
    }

    process.stdout.write(encoder.encode(['uint256[2]', 'uint256[2][2]', 'uint256[2]', 'uint256[4]'], [packedSolidityProof.pi_a, packedSolidityProof.pi_b, packedSolidityProof.pi_c, publicSignals]));
    process.exit(0)
}

generateZkp()