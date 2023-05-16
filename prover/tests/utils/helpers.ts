import { genPrivKey } from "maci-crypto"
const { eddsa, poseidon } = require("circomlibjs")

export const genPublicKey = (privateKey: string) => {
    return eddsa.prv2pub(privateKey)
}

export const generateUser = (index: number) => {
    const privateKey = genPrivKey().toString()
    const publicKey = eddsa.prv2pub(privateKey)

    return {
        index,
        privateKey,
        publicKey
    }
}

const randomUser = (users: Array<any>) => {
    return users[Math.floor(Math.random() * users.length)]
}

export const generateTransaction = (leaves: any, users: Array<any>) => {
    const sender = randomUser(users)
    const recipient = randomUser(users)
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

export const fillSiblings = (depth: number, siblings: Array<BigInt>): Array<BigInt> => {
    const diff = depth - siblings.length
    if (diff > 0) {
        for (let i = 0; i < diff; i++) {
            siblings.push(BigInt(0))
        }
    }
    return siblings
}

export const insertTree = async (tree: any, index: number, leaf: any) => {
    const hashedLeaf = poseidon([leaf.publicKey[0], leaf.publicKey[1], leaf.balance, leaf.nonce])
    await tree.insert(index, hashedLeaf)
}

export const formatTx = (tx: any) => {
    return [
        tx.from,
        tx.to,
        tx.amount,
        tx.fee,
        tx.nonce,
        tx.signature !== undefined ? tx.signature.R8[0] : null,
        tx.signature !== undefined ? tx.signature.R8[1] : null,
        tx.signature !== undefined ? tx.signature.S : null
    ].filter((x: any) => x !== null)
    .map((x: any) => BigInt(x))
}

export const stringifyBigInts = (inputs: any): any => {
    if (typeof inputs == "bigint") {
        return inputs.toString(10)
    } else if (Array.isArray(inputs)) {
        return inputs.map(stringifyBigInts);
    } else if (typeof inputs == "object") {
        const keys = Object.keys(inputs)
        const res = Object.assign({}, inputs)
        keys.forEach((key: any) => {
            res[key] = stringifyBigInts(inputs[key])
        })
        return res
    } else {
        return inputs
    }
}

