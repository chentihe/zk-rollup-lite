pragma circom  2.1.5;

include "batchProcessTx.circom";

component main {public 
[
    balanceTreeRoots, 
    txData, 
    txSendersPublicKey, 
    txSendersPathElements, 
    txRecipientsPublicKey, 
    txRecipientsPathElements, 
    intermediateBalanceTreeRoots, 
    intermediateBalanceTreesPathElements
]} = BatchProcessTx(2, 6);