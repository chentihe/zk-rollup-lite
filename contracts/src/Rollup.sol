// SPDX-License-Identifier: MIT
pragma solidity >0.8.0 <=0.9;

import {IncrementalBinaryTree, IncrementalTreeData} from "zk-kit/incremental-merkle-tree.sol/contracts/IncrementalBinaryTree.sol";
import {Constants} from "./Constants.sol";

contract Rollup {
    using IncrementalBinaryTree for IncrementalTreeData;

    IncrementalTreeData balanceTree;

    TxVerifier txVerifier;
    WithdrawVerifier withdrawVerifier;

    event Deposit(
        uint256 balanceTreeIndex,
        uint256 publicKeyX,
        uint256 publicKeyY,
        uint256 balance,
        uint256 nonce
    );

    event Withdraw(
        uint256 balanceTreeIndex,
        uint256 publicKeyX,
        uint256 publicKeyY,
        uint256 balance,
        uint256 nonce
    );

    event RollUp(uint256 newBalanceTreeRoot);

    struct User {
        uint256 balanceTreeLeafIndex;
        uint256 publicKeyX;
        uint256 publicKeyY;
        uint256 balance;
        uint256 nonce;
    }

    // hashedPublicKey => User
    mapping(uint256 => User) balanceTreeUsers;
    mapping(uint256 => bool) isPublicKeysRegistered;
    mapping(uint256 => bool) usedNullifiers;

    // index => hashedPublicKey
    mapping(uint256 => uint256) balanceTreeKeys;

    uint256 accruedFees;

    constructor(TxVerifier _txVerifier, WithdrawVerifier _withdrawVerifier, uint256 _depth) {
        txVerifier = _txVerifier;
        withdrawVerifier = _withdrawVerifier;
        balanceTree.initWithDefaultZeroes(_depth);
    }
}
