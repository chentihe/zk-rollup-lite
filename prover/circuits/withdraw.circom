pragma circom  2.1.5;

include "../node_modules/circomlib/circuits/eddsamimc.circom";
include "../node_modules/circomlib/circuits/poseidon.circom";
include "../node_modules/circomlib/circuits/comparators.circom";
include "../node_modules/circomlib/circuits/smt/smtprocessor.circom";

template Withdraw(depth) {
    signal input publicKey[2];
    signal input signature[3];
    signal input nullifier;

    signal input balanceTreeRoot;
    signal output newBalanceTreeRoot;

    // [public_key_x, public_key_y, balance, nonce]
    var BALANCE_TREE_LEAF_DATA_LENGTH = 4;

    // The balance is the amount after withdrawal
    signal input balance;
    signal input nonce;
    signal input pathElements[depth];
    signal input oldKey;
    signal input oldValue;
    signal input isOld0;
    signal input newKey;

    signal input func[2];
    var SIGNATURE_R8X_IDX = 0;
    var SIGNATURE_R8Y_IDX = 1;
    var SIGNATURE_S_IDX = 2;

    var ENABLED = 1;

    // 1. Check the signature is valid
    EdDSAMiMCVerifier()(
        ENABLED, 
        publicKey[0], 
        publicKey[1], 
        signature[SIGNATURE_S_IDX], 
        signature[SIGNATURE_R8X_IDX], 
        signature[SIGNATURE_R8Y_IDX], 
        nullifier
    );


    // 1 Make sure the balance is valid
    var TRUE = 1;
    signal isBalanceValid <== GreaterThan(252)([balance, -1]);
    TRUE === isBalanceValid;

    // 2. Process the SMT
    signal accountLeaf <== Poseidon(BALANCE_TREE_LEAF_DATA_LENGTH)
        ([publicKey[0], publicKey[1], balance, nonce]);

    newBalanceTreeRoot <== SMTProcessor(depth)(
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

component main {public [balanceTreeRoot, publicKey, nullifier]} = Withdraw(6);