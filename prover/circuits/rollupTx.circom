pragma circom  2.1.5;

include "batchProcessTx.circom";

component main {public [balanceTreeRoots, txData]} = BatchProcessTx(2, 6);