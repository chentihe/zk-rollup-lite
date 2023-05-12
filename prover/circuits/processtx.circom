include "../node_modules/circomlib/circuits/eddsamimic.circom";

template ProcessTx(depth) {
    signal output newBalanceTreeRoot;

    signal input balanceTreeRoot;


    // [public_key_x, public_key_y, balance, nonce]
    var BALANCE_TREE_LEAF_DATA_LENGTH = 4;

    /*
        Transaction
        [0]: from
        [1]: to
        [2]: amount to send (in wei)
        [3]: fee (in wei)
        [4]: nonce
        [5]: signature R8X
        [6]: signature R8Y
        [7]: signature S
    */
    var TX_DATA_FROM_IDX = 0;
    var TX_DATA_TO_IDX = 1;
    var TX_DATA_AMOUNT_WEI_IDX = 2;
    var TX_DATA_FEE_WEI_IDX = 3;
    var TX_DATA_NONCE_IDX = 4;
    var TX_DATA_SIGNATURE_R8X_IDX = 5;
    var TX_DATA_SIGNATURE_R8Y_IDX = 6;
    var TX_DATA_SIGNATURE_S_IDX = 7;

    var TX_DATA_WITHOUT_SIG_LENGTH = 5;
    var TX_DATA_WITH_SIG_LENGTH = 8;

    sinal input txData[TX_DATA_WITH_SIG_LENGTH];

    // Transaction sender
    signal input txSenderPublicKey[2];
    signal input txSenderBalance;
    signal input txSenderNonce;
    signal input txSenderPathElements[depth];
    
    // Transaction recipient
    signal input txRecipientPublicKey[2];
    signal input txRecipientBalance;
    signal input txRecipientNonce;
    signal input txRecipientPathElements[depth];

    // 1. Validate the signature
    // TODO: hash txData using MiMCSponge
    // Message is hashed tx data
    var preimages;
    EdDSAMiMCVerifier()(txSenderPublicKey[0], txSenderPublicKey[1], txData[TX_DATA_SIGNATURE_R8X_IDX], txData[TX_DATA_SIGNATURE_R8Y_IDX], txData[TX_DATA_SIGNATURE_S_IDX], preimages);
}