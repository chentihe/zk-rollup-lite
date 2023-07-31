const { ethers } = require("ethers");
const { genPrivKey } = require("maci-crypto");
const { eddsa } = require("circomlibjs");

const generateZkp = async () => {
    const encoder = ethers.AbiCoder.defaultAbiCoder();
    const privateKey = genPrivKey().toString();
    const publicKey = eddsa.prv2pub(privateKey)

    process.stdout.write(encoder.encode(['uint256[2]'], [publicKey]));
    process.exit(0)
}

generateZkp()