pragma circom  2.1.5;

include "../node_modules/circomlib/circuits/eddsamimc.circom";
include "../node_modules/circomlib/circuits/poseidon.circom";
include "../node_modules/circomlib/circuits/comparators.circom";
include "../node_modules/circomlib/circuits/smt/smtverifier.circom";
include "../node_modules/circomlib/circuits/smt/smtprocessor.circom";
include "../node_modules/circomlib/circuits/mux1.circom";

// use powersOfTau 2**15 constraints

template ProcessTx(depth) {
    signal input balanceTreeRoot;
    signal output newBalanceTreeRoot;

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
    // Path is the txRecipient path where txSender leaf is updated
    signal input intermediateBalanceTreeRoot;
    signal input intermediateBalanceTreePathElements[depth];

    // 1.1 Validate the signature
    // sign the tx with eddsa.signMiMC()
    var ENABLED = 1;
    signal hashedMsg <== Poseidon(TX_DATA_WITHOUT_SIG_LENGTH)(
        [
            txData[TX_DATA_FROM_IDX], 
            txData[TX_DATA_TO_IDX],
            txData[TX_DATA_AMOUNT_WEI_IDX],
            txData[TX_DATA_FEE_WEI_IDX],
            txData[TX_DATA_NONCE_IDX]
        ]
    );
    EdDSAMiMCVerifier()(
        ENABLED,
        txSenderPublicKey[0],
        txSenderPublicKey[1],
        txData[TX_DATA_SIGNATURE_S_IDX],
        txData[TX_DATA_SIGNATURE_R8X_IDX],
        txData[TX_DATA_SIGNATURE_R8Y_IDX],
        hashedMsg
    );

    // 1.2 Make sure the nonce, amount and fee are valid
    var TRUE = 1;
    var NOT_OLD = 0;
    var VERIFY_INCLUSION = 0;

    txData[TX_DATA_NONCE_IDX] === txSenderNonce + 1;
    signal isAmountValid <== GreaterThan(252)([txData[TX_DATA_AMOUNT_WEI_IDX], 0]);
    signal isFeeValid <== GreaterThan(252)([txData[TX_DATA_FEE_WEI_IDX], 0]);

    TRUE === isAmountValid;
    TRUE === isFeeValid;
    
    // 2. Make sure sender balance > amount + fee
    signal isBalanceValid <== GreaterThan(252)([txSenderBalance, txData[TX_DATA_AMOUNT_WEI_IDX] + txData[TX_DATA_FEE_WEI_IDX]]);

    TRUE === isBalanceValid;

    // 3. Make sure sender exists in the balance tree
    signal txSenderLeaf <== Poseidon(BALANCE_TREE_LEAF_DATA_LENGTH)
        ([txSenderPublicKey[0], txSenderPublicKey[1], txSenderBalance, txSenderNonce]);

    SMTVerifier(depth)(
        ENABLED,
        balanceTreeRoot,
        txSenderPathElements,
        NOT_OLD,
        0,
        0,
        txData[TX_DATA_FROM_IDX],
        txSenderLeaf,
        VERIFY_INCLUSION
    );

    // 4. Create new txSender and txRecipient leaves
    var newTxSenderBalance = txSenderBalance - txData[TX_DATA_AMOUNT_WEI_IDX] - txData[TX_DATA_FEE_WEI_IDX];
    signal newTxSenderLeaf <== Poseidon(BALANCE_TREE_LEAF_DATA_LENGTH)
        ([txSenderPublicKey[0], txSenderPublicKey[1], newTxSenderBalance, txData[TX_DATA_NONCE_IDX]]);

    // In case the sender is the recipient, the balance and the nonce should be updated simultaneously
    signal isSenderRecipentEqual <== IsEqual()([txData[TX_DATA_FROM_IDX], txData[TX_DATA_TO_IDX]]);

    signal offsetTxRecipientBalance <== Mux1()([txRecipientBalance, newTxSenderBalance], isSenderRecipentEqual);
    signal offsetTxRecipientNonce <== Mux1()([txRecipientNonce, txData[TX_DATA_NONCE_IDX]], isSenderRecipentEqual);

    var newTxRecipientBalance = offsetTxRecipientBalance + txData[TX_DATA_AMOUNT_WEI_IDX];
    signal newTxRecipientLeaf <== Poseidon(BALANCE_TREE_LEAF_DATA_LENGTH)
        ([txRecipientPublicKey[0], txRecipientPublicKey[1], newTxRecipientBalance, offsetTxRecipientNonce]);
    
    // 5.1 Update txSender
    var UPDATE_FUNCTION[2] = [0, 1];

    signal computedIntermediateBalanceTreeRoot <== SMTProcessor(depth)(
        balanceTreeRoot,
        txSenderPathElements,
        txData[TX_DATA_FROM_IDX],
        txSenderLeaf, NOT_OLD,
        txData[TX_DATA_FROM_IDX],
        newTxSenderLeaf,
        UPDATE_FUNCTION
    );
    intermediateBalanceTreeRoot === computedIntermediateBalanceTreeRoot;

    // 5.2 Update txRecipient
    signal txRecipientLeaf <== Poseidon(BALANCE_TREE_LEAF_DATA_LENGTH)
        ([txRecipientPublicKey[0], txRecipientPublicKey[1], txRecipientBalance, txRecipientNonce]);

    newBalanceTreeRoot <== SMTProcessor(depth)(
        intermediateBalanceTreeRoot,
        intermediateBalanceTreePathElements,
        txData[TX_DATA_TO_IDX],
        txRecipientLeaf,
        NOT_OLD,
        txData[TX_DATA_TO_IDX],
        newTxRecipientLeaf,
        UPDATE_FUNCTION
    );
}