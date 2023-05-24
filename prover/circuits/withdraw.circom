pragma circom  2.1.5;

include "../node_modules/circomlib/circuits/eddsamimc.circom";

template Withdraw() {
    signal input publicKey[2];
    signal input signature[3];
    signal input nullifier;

    var SIGNATURE_R8X_IDX = 0;
    var SIGNATURE_R8Y_IDX = 1;
    var SIGNATURE_S_IDX = 2;

    var ENABLED = 1;

    EdDSAMiMCVerifier()(
        ENABLED, 
        publicKey[0], 
        publicKey[1], 
        signature[SIGNATURE_S_IDX], 
        signature[SIGNATURE_R8X_IDX], 
        signature[SIGNATURE_R8Y_IDX], 
        nullifier
    );
}

component main {public [publicKey, nullifier]} = Withdraw();