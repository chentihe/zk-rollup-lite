pragma circom  2.1.0;

include "../node_modules/circomlib/circuits/eddsamimic.circom";
include "../node_modules/circomlib/circuits/poseidon.circom";
include "../node_modules/circomlib/circuits/comparators.circom";
include "../node_modules/circomlib/circuits/smt/smtverifier.circom";
include "../node_modules/circomlib/circuits/mux1.circom";

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

    signal input txData[TX_DATA_WITH_SIG_LENGTH];

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

    // Intermediate balance tree root
    signal input intermediateBalanceTreeRoot;
    signal input intermediateBalanceTreePathElements[depth];

    // 1. Validate the signature
    // sign the tx with eddsa.signMiMc()
    signal hashedMsg <== Poseidon(TX_DATA_WITHOUT_SIG_LENGTH)(txData[TX_DATA_FROM_IDX], txData[TX_DATA_TO_IDX], txData[TX_DATA_AMOUNT_WEI_IDX], txData[TX_DATA_FEE_WEI_IDX], txData[TX_DATA_NONCE_IDX]);
    EdDSAMiMCVerifier()(txSenderPublicKey[0], txSenderPublicKey[1], txData[TX_DATA_SIGNATURE_R8X_IDX], txData[TX_DATA_SIGNATURE_R8Y_IDX], txData[TX_DATA_SIGNATURE_S_IDX], hashedMsg);

    // 1.2 Make sure the nonce, amount and fee are valid
    tx[TX_DATA_NONCE_IDX] === txSenderNonce;
    1 === GreaterThan(256)(txData[TX_DATA_AMOUNT_WEI_IDX], 0);
    1 === GreaterThan(256)(txData[TX_DATA_FEE_WEI_IDX], 0);
    
    // 2. Make sure sender balance > amount + fee
    1 === GreaterThan(256)(txSenderBalance, txData[TX_DATA_AMOUNT_WEI_IDX] + txData[TX_DATA_FEE_WEI_IDX]);

    // 3. Make sure sender exists in the balance tree
    signal txSenderLeaf <== Poseidon(BALANCE_TREE_LEAF_DATA_LENGTH)(txSenderPublicKey[0], txSenderPublicKey[1], txSenderBalance, txSenderNonce);

    SMTVerifier(depth)(1, balanceTreeRoot, txSenderPathElements, txData[TX_DATA_FROM_IDX], txSenderLeaf, 0, txData[TX_DATA_FROM_IDX], txSenderLeaf, 0);

    // 4. Create new txSender and txRecipient leaves
    // TODO: using SMTProcessor to compare with intermidiate root & create newBalanceRoot
    var newSenderBalance = txSenderBalance - txData[TX_DATA_AMOUNT_WEI_IDX] - txData[TX_DATA_FEE_WEI_IDX];
    signal newSenderLeaf <== Poseidon(BALANCE_TREE_LEAF_DATA_LENGTH)(txSenderPublicKey[0], txSenderPublicKey[1], newSenderBalance, txData[TX_DATA_NONCE_IDX]);

    signal isSenderRecipentEqual <== IsEqual()(txData[TX_DATA_FROM_IDX], txData[TX_DATA_TO_IDX]);

    signal selectedRecipientBalance <== Mux1()(txRecipientBalance, newSenderBalance, isSenderRecipentEqual);
    signal selectedRecipientNonce <== Mux1()(txRecipientNonce, txData[TX_DATA_NONCE_IDX], isSenderRecipentEqual);

    var newRecipientBalance = selectedRecipientBalance + txData[TX_DATA_AMOUNT_WEI_IDX];
    signal newRecipientLeaf <== Poseidon(txRecipientPublicKey[0], txRecipientPublicKey[1], newRecipientBalance, selectedRecipientNonce);

    
}