pragma circom 2.1.5;

include "processTx.circom";

// use powersOfTau 2**17 constraints

template BatchProcessTx(batchSize, depth) {
    signal output newBalanceTreeRoot;

    // Merkle root of the current balance tree
    signal input balanceTreeRoots[batchSize];

    var TX_DATA_WITH_SIG_LENGTH = 8;
    signal input txData[batchSize][TX_DATA_WITH_SIG_LENGTH];

    // Transaction senders
    signal input txSendersPublicKey[batchSize][2];
    signal input txSendersBalance[batchSize];
    signal input txSendersNonce[batchSize];
    signal input txSendersPathElements[batchSize][depth];

    // Transaction recipients
    signal input txRecipientsPublicKey[batchSize][2];
    signal input txRecipientsBalance[batchSize];
    signal input txRecipientsNonce[batchSize];
    signal input txRecipientsPathElements[batchSize][depth];

    // Intermediate balance tree roots
    // Path is the txRecipient path where txSender leaf is updated
    signal input intermediateBalanceTreeRoots[batchSize];
    signal input intermediateBalanceTreesPathElements[batchSize][depth];

    // Process Txs
    signal processTxs[batchSize];
    for (var i = 0; i < batchSize; i++) {
        processTxs[i] <== ProcessTx(depth)(
            balanceTreeRoots[i], 
            txData[i], 
            txSendersPublicKey[i], 
            txSendersBalance[i], 
            txSendersNonce[i], 
            txSendersPathElements[i], 
            txRecipientsPublicKey[i], 
            txRecipientsBalance[i], 
            txRecipientsNonce[i], 
            txRecipientsPathElements[i], 
            intermediateBalanceTreeRoots[i], 
            intermediateBalanceTreesPathElements[i]
        );
    }

    // TODO: this assertion doesn't make sense
    // if there is deposit or withdraw action between transactions
    // the calculated root won't be equal to the init root of the next tx
    // Make sure calculated roots are valid
    // for (var i = 1; i < batchSize; i++) {
    //     balanceTreeRoots[i] === processTxs[i - 1];
    // }

    newBalanceTreeRoot <== processTxs[batchSize - 1];
}