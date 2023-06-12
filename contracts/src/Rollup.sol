// SPDX-License-Identifier: MIT
pragma solidity >0.8.0 <=0.9;

import {TxVerifier} from "./verifiers/TxVerifier.sol";
import {WithdrawVerifier} from "./verifiers/WithdrawVerifier.sol";
import {DepositVerifier} from "./verifiers/DepositVerifier.sol";

import {Constants} from "./Constants.sol";
import {Errors} from "./Errors.sol";

contract Rollup {
    uint256 private constant _NOT_ENTERED = 1;
    uint256 private constant _ENTERED = 2;

    uint256 private _status;

    address owner;

    uint256 public balanceTreeRoot;

    TxVerifier txVerifier;
    WithdrawVerifier withdrawVerifier;
    DepositVerifier depositVerifier;
    
    event Deposit(User user);

    event Withdraw(User user);

    event RollUp(uint256 newBalanceTreeRoot);

    struct User {
        uint256 publicKeyX;
        uint256 publicKeyY;
        uint256 balance;
        uint256 nonce;
    }

    // hashedPublicKey => User
    mapping(uint256 => User) public balanceTreeUsers;
    mapping(uint256 => bool) public isPublicKeysRegistered;
    mapping(uint256 => bool) usedNullifiers;

    uint256 accruedFees;

    constructor(TxVerifier _txVerifier, WithdrawVerifier _withdrawVerifier, DepositVerifier _depositVerifier) {
        txVerifier = _txVerifier;
        withdrawVerifier = _withdrawVerifier;
        depositVerifier = _depositVerifier;
        _status = _NOT_ENTERED;
        owner = msg.sender;
    }

    modifier nonReentrant() {
        if (_status == _ENTERED) {
            revert Errors.REENTRANT_CALL();
        }
        _status = _ENTERED;
        _;
        _status = _NOT_ENTERED;
    }

    // function rollUp (
    //     uint256[2] memory a,
    //     uint256[2][2] memory b,
    //     uint256[2] memory c,
    //     uint8[][] memory pathIndices,
    //     uint[65] memory input
    // ) external {
    //     // depth = 6
    //     // balanceRoot = input[1]
    //     if (balanceTree.root != input[1]) {
    //         revert Errors.INVALID_MERKLE_TREE();
    //     }

    //     if (!txVerifier.verifyProof(a, b, c, input)) {
    //         revert Errors.INVALID_ROLLUP_PROOFS();
    //     }

    //     // Transaction
    //     uint256 amount;
    //     uint256 fee;
    //     uint256 nonce;
    //     uint256 curOffset;

    //     uint256 leaf;
    //     uint256 newLeaf;

    //     uint256 publicKeyHash;

    //     uint256[] memory pathElements = new uint256[](6);
    //     uint8[] memory pathIndex;

    //     uint256 txDataOffset = 3;
    //     uint256 batchSize = 2;

    //     // batchSize = 2
    //     for (uint8 i = 0; i < batchSize; i++) {
    //         // txData[i] txDataOffset = 3
    //         // txDataLength = 8
    //         curOffset = txDataOffset + (8 * i);

    //         amount = input[curOffset + 2];
    //         fee = input[curOffset + 3];
    //         nonce = input[curOffset + 4];

    //         // sendersPublicKey[i]
    //         curOffset += (8 * (batchSize - i)) + (2 * i);
    //         publicKeyHash = _generateKeyHash(input[curOffset], input[curOffset + 1]);

    //         // sendersPathElements[i]
    //         // 6 = depth
    //         curOffset += (2 * (batchSize - i)) + (6 * i);
    //         for (uint8 j = 0; j < 6; j++) {
    //             pathElements[j] = input[curOffset + j];
    //         }

    //         // update txSender
    //         User storage user = balanceTreeUsers[publicKeyHash];

    //         leaf = PoseidonT5.hash([user.publicKeyX, user.publicKeyY, user.balance, user.nonce]);

    //         // underflow can't happen
    //         // zkp verified all inputs
    //         unchecked {
    //             user.balance -= amount;
    //             user.balance -= fee; 
    //         }
    //         user.nonce = nonce;

    //         accruedFees += fee;

    //         newLeaf = PoseidonT5.hash([user.publicKeyX, user.publicKeyY, user.balance, user.nonce]);

    //         pathIndex = pathIndices[2 * i];

    //         balanceTree.update(leaf, newLeaf, pathElements, pathIndex);

    //         // recipientPublicKey[i]
    //         curOffset += (6 * (batchSize - i)) + (2 * i);
    //         publicKeyHash = _generateKeyHash(input[curOffset], input[curOffset + 1]);

    //         // recipientPathElements[i]
    //         curOffset += ((2 + 6) * (batchSize - i)) + (6 * i);
    //         for (uint8 j = 0; j < 6; j++) {
    //             pathElements[j] = input[curOffset + j];
    //         }

    //         // update txRecipient
    //         user = balanceTreeUsers[publicKeyHash];

    //         leaf = PoseidonT5.hash([user.publicKeyX, user.publicKeyY, user.balance, user.nonce]);

    //         unchecked {
    //             user.balance += amount;
    //         }

    //         newLeaf = PoseidonT5.hash([user.publicKeyX, user.publicKeyY, user.balance, user.nonce]);

    //         pathIndex = pathIndices[2 * i + 1];

    //         balanceTree.update(leaf, newLeaf, pathElements, pathIndex);
    //     }
    // }

    function deposit(
        uint256[2] memory a,
        uint256[2][2] memory b,
        uint256[2] memory c,
        uint256[4] memory input
    ) external payable {
        if (!depositVerifier.verifyProof(a, b, c, input)) {
            revert Errors.INVALID_DEPOSIT_PROOFS();
        }

        if (msg.value == 0) {
            revert Errors.INVALID_VALUE();
        }

        uint256 newRoot = input[0];
        uint256 root = input[1];
        uint256 publicKeyX = input[2];
        uint256 publicKeyY = input[3];

        if (root != balanceTreeRoot) {
            revert Errors.INVALID_MERKLE_TREE();
        }

        uint256 publicKeyHash = _generateKeyHash(publicKeyX, publicKeyY);
        User storage user = balanceTreeUsers[publicKeyHash];
        
        // zkp is valid, just update balance
        user.balance += msg.value;


        if (!isPublicKeysRegistered[publicKeyHash]) {
            isPublicKeysRegistered[publicKeyHash] = true;

            user.publicKeyX = publicKeyX;
            user.publicKeyY = publicKeyY;
        }

        balanceTreeRoot = newRoot;

        emit Deposit(user);
    }

    // withdraw all deposit
    function withdraw(
        uint256 amount,
        uint256[2] memory a,
        uint256[2][2] memory b,
        uint256[2] memory c,
        uint256[5] memory input
    ) external nonReentrant {
        if (!withdrawVerifier.verifyProof(a, b, c, input)) {
            revert Errors.INVALID_WITHDRAW_PROOFS();
        }

        uint256 newRoot = input[0];
        uint256 root = input[1];
        uint256 publicKeyX = input[2];
        uint256 publicKeyY = input[3];
        uint256 nullifier = input[4];

        if (balanceTreeRoot != root) {
            revert Errors.INVALID_MERKLE_TREE();
        }

        if (usedNullifiers[nullifier]) {
            revert Errors.INVALID_NULLIFIER();
        }

        User storage user = _getUserByPublicKey(publicKeyX, publicKeyY);
        if (amount == Constants.UINT256_MAX) {
            amount = user.balance;
        }

        if (amount > user.balance || amount == 0) {
            revert Errors.INSUFFICIENT_BALANCE();
        }

        usedNullifiers[nullifier] = true;
        
        user.balance -= amount;

        (bool success,) = msg.sender.call{value: amount}("");
        if (!success) {
            revert Errors.WITHDRAWAL_FAILED();
        } 

        balanceTreeRoot = newRoot;

        emit Withdraw(user);
    }

    function generateKeyHash(uint256 publicKeyX, uint256 publicKeyY) external pure returns (uint256) {
        return _generateKeyHash(publicKeyX, publicKeyY);
    }

    function _generateKeyHash(uint256 publicKeyX, uint256 publicKeyY) internal pure returns (uint256) {
        return uint256(keccak256(abi.encodePacked(publicKeyX, publicKeyY)));
    }

    function getUserByPublicKey(uint256 publicKeyX, uint256 publicKeyY) external view returns (User memory) {
        return _getUserByPublicKey(publicKeyX, publicKeyY);
    }

    function _getUserByPublicKey(uint256 publicKeyX, uint256 publicKeyY) internal view returns (User storage) {
        uint256 publicKeyHash = _generateKeyHash(publicKeyX, publicKeyY);
        return balanceTreeUsers[publicKeyHash];
    }

    function withdrawAccruedFees() external {
        if (msg.sender != owner) {
            revert Errors.ONLY_OWNER();
        }

        (bool success,) = msg.sender.call{value: accruedFees}("");
        if (!success) {
            revert Errors.WITHDRAWAL_FAILED();
        } 

        accruedFees = 0;
    }
}
