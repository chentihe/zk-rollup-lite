// SPDX-License-Identifier: UNLICENSED
pragma solidity >0.8.0 <=0.9;

import {Test} from "forge-std/Test.sol";
import {StdStorage, stdStorage} from "forge-std/StdStorage.sol";
import {Rollup} from "../src/Rollup.sol";
import {TxVerifier} from "../src/verifiers/TxVerifier.sol";
import {WithdrawVerifier} from "../src/verifiers/WithdrawVerifier.sol";

contract RollupTest is Test {
    using stdStorage for StdStorage;

    Rollup rollup;
    TxVerifier txVerifier;
    WithdrawVerifier withdrawVerifier;
    uint256 constant DEPTH = 6;

    function setUp() public {
        txVerifier = new TxVerifier();
        withdrawVerifier = new WithdrawVerifier();
        rollup = new Rollup(txVerifier, withdrawVerifier);
    }

    function testDeposit() public {
        string[] memory runJsInputs = new string[](6);
        runJsInputs[0] = "npm";
        runJsInputs[1] = "--prefix";
        runJsInputs[2] = "contracts/test/script/";
        runJsInputs[3] = "--silent";
        runJsInputs[4] = "run";
        runJsInputs[5] = "generate-public-key";
        bytes memory jsResult = vm.ffi(runJsInputs);
        (uint256[2] memory publicKey) = abi.decode(jsResult, (uint256[2]));
        
        uint256 publicKeyX = publicKey[0];
        uint256 publicKeyY = publicKey[1];

        rollup.deposit{value: 1 ether}(publicKeyX, publicKeyY);

        uint256 keyHash = rollup.generateKeyHash(publicKeyX, publicKeyY);
        bool isRegistered = rollup.isPublicKeysRegistered(keyHash);
        assertEq(isRegistered, true);
    }

    function testWithdraw() public {
        string[] memory runJsInputs = new string[](8);
        runJsInputs[0] = "npm";
        runJsInputs[1] = "--prefix";
        runJsInputs[2] = "contracts/test/script/";
        runJsInputs[3] = "--silent";
        runJsInputs[4] = "run";
        runJsInputs[5] = "generate-zkp-withdraw";
        runJsInputs[6] = "withdraw";
        runJsInputs[7] = "0.5";
        bytes memory jsResult = vm.ffi(runJsInputs);
        (uint256[2] memory a, uint256[2][2] memory b, uint256[2] memory c, uint256[3] memory input) = abi.decode(jsResult, (uint256[2], uint256[2][2], uint256[2], uint256[3]));

        uint256 nullifier = input[0];
        uint256 publicKeyX = input[1];
        uint256 publicKeyY = input[2];

        rollup.deposit{value: 1 ether}(publicKeyX, publicKeyY);

        uint256 beforeWithdraw = address(this).balance;
        rollup.withdraw(0.5e18, a, b, c, input);
        uint256 afterWithdraw = address(this).balance;

        assertEq(afterWithdraw - beforeWithdraw, 0.5e18);
        assertTrue(rollup.usedNullifiers(nullifier));
    }

    function testWithdrawAll() public {
        string[] memory runJsInputs = new string[](8);
        runJsInputs[0] = "npm";
        runJsInputs[1] = "--prefix";
        runJsInputs[2] = "contracts/test/script/";
        runJsInputs[3] = "--silent";
        runJsInputs[4] = "run";
        runJsInputs[5] = "generate-zkp-withdraw";
        runJsInputs[6] = "withdraw-all";
        runJsInputs[7] = "0.5";
        bytes memory jsResult = vm.ffi(runJsInputs);
        (uint256[2] memory a, uint256[2][2] memory b, uint256[2] memory c, uint256[3] memory input) = abi.decode(jsResult, (uint256[2], uint256[2][2], uint256[2], uint256[3]));

        uint256 nullifier = input[0];
        uint256 publicKeyX = input[1];
        uint256 publicKeyY = input[2];
        
        rollup.deposit{value: 1 ether}(publicKeyX, publicKeyY);

        uint256 beforeWithdraw = address(this).balance;
        rollup.withdraw(type(uint256).max, a, b, c, input);
        uint256 afterWithdraw = address(this).balance;

        assertEq(afterWithdraw - beforeWithdraw, 1e18);
        assertTrue(rollup.usedNullifiers(nullifier));
    }

    function testRollup() public {
        string[] memory runJsInputs = new string[](6);
        runJsInputs[0] = "npm";
        runJsInputs[1] = "--prefix";
        runJsInputs[2] = "contracts/test/script/";
        runJsInputs[3] = "--silent";
        runJsInputs[4] = "run";
        runJsInputs[5] = "generate-zkp-rollup";
        bytes memory jsResult = vm.ffi(runJsInputs);
        (uint256[2] memory senderKey, uint256[2] memory recipientKey, uint256[2] memory a, uint256[2][2] memory b, uint256[2] memory c, uint256[19] memory input) = abi.decode(jsResult, (uint256[2], uint256[2], uint256[2], uint256[2][2], uint256[2], uint256[19]));

        uint256 newRoot = input[0];
        uint256 root = input[1];
        
        // setup root
        stdstore
            .target(address(rollup))
            .sig(0xb4e7dddd)
            .checked_write(root);

        rollup.deposit{value: 10 ether}(senderKey[0], senderKey[1]);
        rollup.deposit{value: 10 ether}(recipientKey[0], recipientKey[1]);

        rollup.rollUp(a, b, c, input);
        Rollup.User memory sender = rollup.getUserByPublicKey(senderKey[0], senderKey[1]);
        Rollup.User memory recipient = rollup.getUserByPublicKey(recipientKey[0], recipientKey[1]);
        assertEq(newRoot, rollup.balanceTreeRoot());
        
        // sender transfer 1 ether to recipient twice
        // make a rollup to layer 1
        // fee is 0.5 ether
        assertEq(sender.balance, 7e18);
        assertEq(recipient.balance, 12e18);
    }

    receive() external payable {}
    fallback() external payable {}
}