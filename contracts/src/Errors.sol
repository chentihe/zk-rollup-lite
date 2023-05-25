// SPDX-License-Identifier: MIT

pragma solidity >0.8.0 <=0.9;

library Errors {
    error INVALID_VALIE();
    error INVALID_MERKLE_TREE();
    error INVALID_ROLLUP_PROOFS();
    error INVALID_WITHDRAW_PROOFS();
    error INVALID_NULLIFIER();
    error INSUFFICIENT_BALANCE();
    error REENTRANT_CALL();
    error WITHDRAWAL_FAILED();
    error ONLY_OWNER();
}