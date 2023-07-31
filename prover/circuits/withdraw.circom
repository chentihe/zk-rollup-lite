pragma circom  2.1.5;

include "../node_modules/circomlib/circuits/eddsamimc.circom";
include "../node_modules/circomlib/circuits/poseidon.circom";
include "../node_modules/circomlib/circuits/comparators.circom";
include "../node_modules/circomlib/circuits/smt/smtverifier.circom";

template Withdraw(depth) {
    signal input balanceTreeRoot;
    signal input signature[3];
    signal input nullifier;

    // [public_key_x, public_key_y, balance, nonce]
    var BALANCE_TREE_LEAF_DATA_LENGTH = 4;
    signal input publicKey[2];
    signal input balance;
    signal input nonce;

    signal input pathElements[depth];
    signal input oldKey;
    signal input oldValue;
    signal input newKey;

    // only verify the inclusion case
    signal isOld0 <== 0;
    signal func <== 0;

    var SIGNATURE_R8X_IDX = 0;
    var SIGNATURE_R8Y_IDX = 1;
    var SIGNATURE_S_IDX = 2;

    signal enabled <== 1;

    // 1. Check the signature is valid
    EdDSAMiMCVerifier()(
        enabled, 
        publicKey[0], 
        publicKey[1], 
        signature[SIGNATURE_S_IDX], 
        signature[SIGNATURE_R8X_IDX], 
        signature[SIGNATURE_R8Y_IDX], 
        nullifier
    );

    // 1 Make sure the balance is valid
    var TRUE = 1;
    signal isBalanceValid <== GreaterEqThan(252)([balance, 0]);
    TRUE === isBalanceValid;

    // 2. Verify the SMT
    signal accountLeaf <== Poseidon(BALANCE_TREE_LEAF_DATA_LENGTH)
        ([publicKey[0], publicKey[1], balance, nonce]);

    SMTVerifier(depth)(
        enabled,
        balanceTreeRoot,
        pathElements,
        oldKey,
        oldValue,
        isOld0,
        newKey,
        accountLeaf,
        func
    );
}

component main {public [publicKey, nullifier]} = Withdraw(6);