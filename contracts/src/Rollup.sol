// SPDX-License-Identifier: MIT
pragma solidity >0.8.0 <=0.9;

import {IncrementalBinaryTree, IncrementalTreeData} from "zk-kit/incremental-merkle-tree.sol/contracts/IncrementalBinaryTree.sol";
import {PoseidonT5} from "poseidon-solidity/PoseidonT5.sol";
import {PoseidonT6} from "poseidon-solidity/PoseidonT6.sol";
import {TxVerifier} from "./verifiers/TxVerifier.sol";
import {Constants} from "./Constants.sol";
import {Errors} from "./Errors.sol";

contract Rollup {
    using IncrementalBinaryTree for IncrementalTreeData;

    IncrementalTreeData balanceTree;

    TxVerifier txVerifier;
    WithdrawVerifier withdrawVerifier;
    PoseidonT6 poseidonT6;
    PoseidonT5 poseidonT5;

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

    constructor(TxVerifier _txVerifier, WithdrawVerifier _withdrawVerifier, PoseidonT6 _poseidonT6, PoseidonT5 _poseidonT5, uint256 _depth) {
        txVerifier = _txVerifier;
        withdrawVerifier = _withdrawVerifier;
        poseidonT5 = _poseidonT5;
        poseidonT6 = _poseidonT6;
        balanceTree.initWithDefaultZeroes(_depth);
    }

    function rollUp (
        uint[2] memory a,
        uint[2][2] memory b,
        uint[2] memory c,
        uint[65] memory input,
        uint[4][6] memory pathIndices
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

        uint256[depth] senderPathElements;
        uint256[depth] recipientPathElements;

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
            senderPublicKeyHash = keccak256(abi.encodePacked(input[curOffset], input[curOffset + 1]));

            // sendersPathElements[i]
            curOffset += (2 * (batchSize - i)) + (depth * i);
            for (uint8 j = 0; j < depth; j++) {
                senderPathElements[j] = input[curOffset + j];
            }

            // recipientPublicKey[i]
            curOffset += (depth * (batchSize - i)) + (2 * i);
            recipientPublicKeyHash = keccak256(abi.encodePacked(input[curOffset], input[curOffset + 1]));

            // recipientPathElements[i]
            curOffset += ((2 + depth) * (batchSize - i)) + (depth * i);
            for (uint8 k = 0; k < depth; k++) {
                recipientPathElements[k] = input[curOffset + k];
            }

            // update txSender
            User storage sender = balanceTreeUsers[senderPublicKeyHash];

            senderLeaf = poseidonT6.hash([sender.publicKeyX, sender.publicKeyY, sender.balance, sender.nonce]);

            // overflow / underflow can't happen
            // zkp verified all inputs
            unchecked {
                sender.balance -= amount;
                sender.balance -= fee; 
            }
            sender.nonce = nonce;


            accruedFees += fee;

            newSenderLeaf = poseidonT6.hash([sender.publicKeyX, sender.publicKeyY, sender.balance, sender.nonce]);

            balanceTree.update(leaf, newLeaf, senderPathElements, pathIndices[2 * i]);

            // update txRecipient
            User storage recipient = balanceTreeUsers[recipientPublicKeyHash];

            recipientLeaf = poseidonT6.hash([recipient.publicKeyX, recipient.publicKeyY, recipient.balance, recipient.nonce]);

            unchecked {
                recipient.balance += amount;
            }

            newRecipientLeaf = poseidonT6.hash([recipient.publicKeyX, recipient.publicKeyY, recipient.balance, recipient.nonce]);

            balanceTree.update(recipientLeaf, newRecipientLeaf, recipientPathElements, pathIndices[2 * i + 1]);
        }
    }

    // if the user is the first time to deposit, 
    // leave empty array for proofSiblings & proofPathindices
    function deposit(uint256 publicKeyX, uint256 publicKeyY, uint256[] calldata proofSiblings, uint8[] calldata proofPathIndices) public payable {
        uint256 publicKeyHash = keccak256(abi.encodePacked(publicKeyX, publicKeyY));
        User storage user = balanceTreeUsers[publicKeyHash];
        if (msg.value == 0) {
            revert Errors.INVALID_VALIE();
        }
        
        uint256 leaf = poseidonT5.hash([publicKeyX, publicKeyY, user.balance, user.nonce]);
        
        user.balance += msg.value;

        uint256 newLeaf = poseidonT5.hash([publicKeyX, publicKeyY, user.balance, user.nonce]);

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

    // TODO: withdraw need to generate zkp, need to write a circuit for withdraw
}
