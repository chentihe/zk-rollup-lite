// SPDX-License-Identifier: MIT
pragma solidity >0.8.0 <=0.9;

import {IncrementalBinaryTree, IncrementalTreeData} from "zk-kit/incremental-merkle-tree.sol/contracts/IncrementalBinaryTree.sol";

import {PoseidonT5} from "poseidon-solidity/PoseidonT5.sol";
import {PoseidonT6} from "poseidon-solidity/PoseidonT6.sol";

import {TxVerifier} from "./verifiers/TxVerifier.sol";
import {WithdrawVerifier} from "./verifiers/WithdrawVerifier.sol";

import {Constants} from "./Constants.sol";
import {Errors} from "./Errors.sol";

contract Rollup {
    using IncrementalBinaryTree for IncrementalTreeData;

    uint256 private constant _NOT_ENTERED = 1;
    uint256 private constant _ENTERED = 2;

    uint256 private _status;

    address owner;

    IncrementalTreeData balanceTree;

    TxVerifier txVerifier;
    WithdrawVerifier withdrawVerifier;
    
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
    mapping(uint256 => User) balanceTreeUsers;
    mapping(uint256 => bool) isPublicKeysRegistered;
    mapping(uint256 => bool) usedNullifiers;

    uint256 accruedFees;

    constructor(TxVerifier _txVerifier, WithdrawVerifier _withdrawVerifier, uint256 _depth) {
        txVerifier = _txVerifier;
        withdrawVerifier = _withdrawVerifier;
        balanceTree.initWithDefaultZeroes(_depth);
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

    function rollUp (
        uint[2] memory a,
        uint[2][2] memory b,
        uint[2] memory c,
        uint[65] memory input,
        uint8[][6] memory pathIndices
    ) public {
        uint256 balanceTreeRoot = input[1];
        uint256 depth = balanceTree.depth;

        if (balanceTree.root != balanceTreeRoot) {
            revert Errors.INVALID_MERKLE_TREE();
        }

        if (!txVerifier.verifyProof(a, b, c, input)) {
            revert Errors.INVALID_ROLLUP_PROOFS();
        }

        // Transaction
        uint256 from;
        uint256 to;
        uint256 amount;
        uint256 fee;
        uint256 nonce;
        uint256 curOffset;

        uint256 senderLeaf;
        uint256 newSenderLeaf;
        uint256 recipientLeaf;
        uint256 newRecipientLeaf;

        uint256 senderPublicKeyHash;
        uint256 recipientPublicKeyHash;

        uint256[] memory senderPathElements = new uint256[](depth);
        uint256[] memory recipientPathElements = new uint256[](depth);

        uint256 txDataOffset = 3;
        uint256 txDataLength = 8;
        uint256 batchSize = 2;

        for (uint8 i = 0; i < batchSize; i++) {
            // txData[i]
            curOffset = txDataOffset + (txDataLength * i);

            from = input[curOffset];
            to = input[curOffset + 1];
            amount = input[curOffset + 2];
            fee = input[curOffset + 3];
            nonce = input[curOffset + 4];

            // sendersPublicKey[i]
            curOffset += (txDataLength * (batchSize - i)) + (2 * i);
            senderPublicKeyHash = _generateKeyHash(input[curOffset], input[curOffset + 1]);

            // sendersPathElements[i]
            curOffset += (2 * (batchSize - i)) + (depth * i);
            for (uint8 j = 0; j < depth; j++) {
                senderPathElements[j] = input[curOffset + j];
            }

            // recipientPublicKey[i]
            curOffset += (depth * (batchSize - i)) + (2 * i);
            recipientPublicKeyHash = _generateKeyHash(input[curOffset], input[curOffset + 1]);

            // recipientPathElements[i]
            curOffset += ((2 + depth) * (batchSize - i)) + (depth * i);
            for (uint8 k = 0; k < depth; k++) {
                recipientPathElements[k] = input[curOffset + k];
            }

            // update txSender
            User storage sender = balanceTreeUsers[senderPublicKeyHash];

            senderLeaf = PoseidonT5.hash([sender.publicKeyX, sender.publicKeyY, sender.balance, sender.nonce]);

            // underflow can't happen
            // zkp verified all inputs
            unchecked {
                sender.balance -= amount;
                sender.balance -= fee; 
            }
            sender.nonce = nonce;

            accruedFees += fee;

            newSenderLeaf = PoseidonT5.hash([sender.publicKeyX, sender.publicKeyY, sender.balance, sender.nonce]);

            balanceTree.update(senderLeaf, newSenderLeaf, senderPathElements, pathIndices[2 * i]);

            // update txRecipient
            User storage recipient = balanceTreeUsers[recipientPublicKeyHash];

            recipientLeaf = PoseidonT5.hash([recipient.publicKeyX, recipient.publicKeyY, recipient.balance, recipient.nonce]);

            unchecked {
                recipient.balance += amount;
            }

            newRecipientLeaf = PoseidonT5.hash([recipient.publicKeyX, recipient.publicKeyY, recipient.balance, recipient.nonce]);

            balanceTree.update(recipientLeaf, newRecipientLeaf, recipientPathElements, pathIndices[2 * i + 1]);
        }
    }

    // if the user is the first time to deposit, 
    // leave empty array for proofSiblings & proofPathindices
    function deposit(
        uint256 publicKeyX,
        uint256 publicKeyY,
        uint256[] calldata proofSiblings,
        uint8[] calldata proofPathIndices
    ) external payable {
        if (msg.value == 0) {
            revert Errors.INVALID_VALIE();
        }

        uint256 publicKeyHash = _generateKeyHash(publicKeyX, publicKeyY);
        User storage user = balanceTreeUsers[publicKeyHash];
        
        uint256 leaf = PoseidonT5.hash([publicKeyX, publicKeyY, user.balance, user.nonce]);
        
        user.balance += msg.value;

        uint256 newLeaf = PoseidonT5.hash([publicKeyX, publicKeyY, user.balance, user.nonce]);

        if (isPublicKeysRegistered[publicKeyHash]) {
            balanceTree.update(leaf, newLeaf, proofSiblings, proofPathIndices);
        } else {
            isPublicKeysRegistered[publicKeyHash] = true;

            user.publicKeyX = publicKeyX;
            user.publicKeyY = publicKeyY;

            balanceTree.insert(newLeaf);
        }

        emit Deposit(user);
    }

    // withdraw all deposit
    function withdraw(
        uint256[] calldata proofSiblings,
        uint8[] calldata proofPathIndices,
        uint256[2] memory a,
        uint256[2][2] memory b,
        uint256[2] memory c,
        uint256[3] memory input
    ) external {
        _withdraw(Constants.UINT256_MAX, proofSiblings, proofPathIndices, a, b, c, input);
    }

    function withdraw(
        uint256 amount,
        uint256[] calldata proofSiblings,
        uint8[] calldata proofPathIndices,
        uint256[2] memory a,
        uint256[2][2] memory b,
        uint256[2] memory c,
        uint256[3] memory input
    ) external {
        _withdraw(amount, proofSiblings, proofPathIndices, a, b, c, input);
    }

    function _withdraw(
        uint256 amount,
        uint256[] calldata proofSiblings,
        uint8[] calldata proofPathIndices,
        uint256[2] memory a,
        uint256[2][2] memory b,
        uint256[2] memory c,
        uint256[3] memory input
    ) internal nonReentrant {
        uint256 publicKeyX = input[0];
        uint256 publicKeyY = input[1];
        uint256 nullifier = input[2];

        if (usedNullifiers[nullifier]) {
            revert Errors.INVALID_NULLIFIER();
        }

        if (!withdrawVerifier.verifyProof(a, b, c, input)) {
            revert Errors.INVALID_WITHDRAW_PROOFS();
        }

        User storage user = _getUserByPublicKey(publicKeyX, publicKeyY);
        if (amount == Constants.UINT256_MAX) {
            amount = user.balance;
        }

        if (amount >= user.balance || amount == 0) {
            revert Errors.INSUFFICIENT_BALANCE();
        }

        usedNullifiers[nullifier] = true;
        
        uint256 leaf = PoseidonT5.hash([publicKeyX, publicKeyY, user.balance, user.nonce]);

        user.balance -= amount;

        (bool success,) = msg.sender.call{value: amount}("");
        if (!success) {
            revert Errors.WITHDRAWAL_FAILED();
        } 

        uint256 newLeaf = PoseidonT5.hash([publicKeyX, publicKeyY, user.balance, user.nonce]);

        balanceTree.update(leaf, newLeaf, proofSiblings, proofPathIndices);

        emit Withdraw(user);
    }

    function _generateKeyHash(uint256 publicKeyX, uint256 publicKeyY) internal pure returns (uint256) {
        return uint256(keccak256(abi.encodePacked(publicKeyX, publicKeyY)));
    }

    function _getUserByPublicKey(uint256 publicKeyX, uint256 publicKeyY) internal view returns (User storage) {
        uint256 publicKeyHash = _generateKeyHash(publicKeyX, publicKeyY);
        return balanceTreeUsers[publicKeyHash];
    }

    function withdrawAccredFees() external {
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
