// SPDX-License-Identifier: MIT
pragma solidity >0.8.0 <=0.9;

import {TxVerifier} from "./verifiers/TxVerifier.sol";
import {WithdrawVerifier} from "./verifiers/WithdrawVerifier.sol";

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
    
    event Deposit(User user);
    event Withdraw(User user);
    event RollUp(uint256 newBalanceTreeRoot);

    struct User {
        uint256 index;
        uint256 publicKeyX;
        uint256 publicKeyY;
        uint256 balance;
        uint256 nonce;
    }

    // hashedPublicKey => User
    mapping(uint256 => User) public balanceTreeUsers;
    mapping(uint256 => bool) public isPublicKeysRegistered;
    mapping(uint256 => bool) public usedNullifiers;

    uint256 currentUserIndex;

    mapping(uint256 => uint256) public balanceTreeKeys;

    uint256 accruedFees;

    constructor(TxVerifier _txVerifier, WithdrawVerifier _withdrawVerifier) {
        txVerifier = _txVerifier;
        withdrawVerifier = _withdrawVerifier;
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
        uint256[2] memory a,
        uint256[2][2] memory b,
        uint256[2] memory c,
        uint[19] memory input
    ) external {
        if (!txVerifier.verifyProof(a, b, c, input)) {
            revert Errors.INVALID_ROLLUP_PROOFS();
        }

        uint256 newRoot = input[0];

        // Transaction
        uint256 sender;
        uint256 recipient;
        uint256 amount;
        uint256 fee;
        uint256 nonce;
        uint256 curOffset;

        uint256 txDataOffset = 3;
        uint256 batchSize = 2;

        // batchSize = 2
        for (uint8 i = 0; i < batchSize; i++) {
            // txData[i] txDataOffset = 3
            // txDataLength = 8
            curOffset = txDataOffset + (8 * i);

            sender = input[curOffset];
            recipient = input[curOffset + 1];
            amount = input[curOffset + 2];
            fee = input[curOffset + 3];
            nonce = input[curOffset + 4];

            // update txSender
            User storage user = balanceTreeUsers[balanceTreeKeys[sender]];

            // underflow can't happen
            // zkp verified all inputs
            unchecked {
                user.balance -= amount;
                user.balance -= fee;
            }
            user.nonce = nonce;

            accruedFees += fee;

            // update txRecipient
            user = balanceTreeUsers[balanceTreeKeys[recipient]];

            unchecked {
                user.balance += amount;
            }
        }

        balanceTreeRoot = newRoot;
    }

    function deposit(uint256 publicKeyX, uint256 publicKeyY) external payable {
        if (msg.value == 0) {
            revert Errors.INVALID_VALUE();
        }

        uint256 publicKeyHash = _generateKeyHash(publicKeyX, publicKeyY);
        User storage user = balanceTreeUsers[publicKeyHash];
        user.balance += msg.value;

        if (!isPublicKeysRegistered[publicKeyHash]) {
            isPublicKeysRegistered[publicKeyHash] = true;

            // index 0 of mt is reserved
            currentUserIndex++;
            user.index = currentUserIndex;
            user.publicKeyX = publicKeyX;
            user.publicKeyY = publicKeyY;

            balanceTreeKeys[currentUserIndex] = publicKeyHash;
        }

        emit Deposit(user);
    }

    // withdraw all deposit
    function withdraw(
        uint256 amount,
        uint256[2] memory a,
        uint256[2][2] memory b,
        uint256[2] memory c,
        uint256[3] memory input
    ) external nonReentrant {
        if (!withdrawVerifier.verifyProof(a, b, c, input)) {
            revert Errors.INVALID_WITHDRAW_PROOFS();
        }

        uint256 nullifier = input[0];
        uint256 publicKeyX = input[1];
        uint256 publicKeyY = input[2];

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
        if (!isPublicKeysRegistered[publicKeyHash]) {
            revert Errors.INVALID_USER();
        }
        return balanceTreeUsers[publicKeyHash];
    }

    function getUserByIndex(uint256 index) external view returns (User memory) {
        uint256 publicKeyHash = balanceTreeKeys[index];
        if (!isPublicKeysRegistered[publicKeyHash]) {
            revert Errors.INVALID_USER();
        }
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

    receive() external payable {}
    fallback() external payable {}
}
