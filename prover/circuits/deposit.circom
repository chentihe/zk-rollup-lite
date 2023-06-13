pragma circom  2.1.5;

include "../node_modules/circomlib/circuits/poseidon.circom";
include "../node_modules/circomlib/circuits/comparators.circom";
include "../node_modules/circomlib/circuits/smt/smtprocessor.circom";

// use powersOfTau 2**15 constraints

template Deposit(depth) {
    signal input balanceTreeRoot;
    signal output newBalanceTreeRoot;

    // [public_key_x, public_key_y, balance, nonce]
    var BALANCE_TREE_LEAF_DATA_LENGTH = 4;

    // Account info
    signal input publicKey[2];
    signal input balance;
    signal input nonce;
    signal input pathElements[depth];
    signal input oldKey;
    signal input oldValue;
    signal input isOld0;
    signal input newKey;
    signal input func[2];

    // 1 Make sure the balance is valid
    var TRUE = 1;
    signal isBalanceValid <== GreaterThan(252)([balance, 0]);
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

component main {public [balanceTreeRoot, newKey, publicKey]} = Deposit(6);