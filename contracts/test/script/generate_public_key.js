const { ethers } = require("ethers");
const { genPrivKey } = require("maci-crypto");
const { buildEddsa } = require("circomlibjs");

const generatePublicKey = async () => {
    const eddsa = await buildEddsa();
    const encoder = ethers.AbiCoder.defaultAbiCoder();
    const privateKey = genPrivKey().toString();
    // Uint8Array => BigInt
    const publicKey = eddsa.prv2pub(privateKey).flatMap(axis => ethers.toBigInt(axis))
    process.stdout.write(encoder.encode(['uint256[2]'], [publicKey]));
}

generatePublicKey()