// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Test} from "forge-std/Test.sol";
import {Rollup} from "../src/Rollup.sol";
import {TxVerifier} from "../src/verifiers/TxVerifier.sol";
import {WithdrawVerifier} from "../src/verifiers/WithdrawVerifier.sol";

contract RollupTest is Test {
    Rollup rollup;
    TxVerifier txVerifier;
    WithdrawVerifier withdrawVerifier;
    uint256 constant DEPTH = 6;

    function setUp() public {
        txVerifier = new TxVerifier();
        withdrawVerifier = new WithdrawVerifier();
        rollup = new Rollup(txVerifier, withdrawVerifier, DEPTH);
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
        uint256[2] memory jsGeneratedPublicKey = abi.decode(jsResult, (uint256[2]));
        uint256 publicKeyX = jsGeneratedPublicKey[0];
        uint256 publicKeyY = jsGeneratedPublicKey[1];
        rollup.deposit{value: 1 ether}(publicKeyX, publicKeyY);

        uint256 keyHash = rollup.generateKeyHash(publicKeyX, publicKeyY);
        bool isRegistered = rollup.isPublicKeysRegistered(keyHash);
        assertEq(isRegistered, true);
    }

    function testDepositWithRegisteredUser() public {
        uint256 depth = 6;
        string[] memory runJsInputs = new string[](6);
        runJsInputs[0] = "npm";
        runJsInputs[1] = "--prefix";
        runJsInputs[2] = "contracts/test/script/";
        runJsInputs[3] = "--silent";
        runJsInputs[4] = "run";
        runJsInputs[5] = "generate-merkle-tree-proof";
        bytes memory jsResult = vm.ffi(runJsInputs);
        (uint256[2] memory publicKey, uint8[6] memory pathIndices, uint256[6] memory siblings) = abi.decode(jsResult, (uint256[2], uint8[6], uint256[6]));
        uint256 publicKeyX = publicKey[0];
        uint256 publicKeyY = publicKey[1];
        
        rollup.deposit{value: 1 ether}(publicKeyX, publicKeyY);

        uint256 keyHash = rollup.generateKeyHash(publicKeyX, publicKeyY);
        bool isRegistered = rollup.isPublicKeysRegistered(keyHash);
        assertEq(isRegistered, true);

        uint256[] memory proofSiblings = new uint256[](depth);
        uint8[] memory proofPathIndices = new uint8[](depth);
        for (uint8 i = 0; i < depth; i++) {
            proofSiblings[i] = siblings[i];
            proofPathIndices[i] = pathIndices[i];
        }

        uint256 beforeDeposit = address(this).balance;
        rollup.deposit{value: 1 ether}(publicKeyX, publicKeyY, proofSiblings, proofPathIndices);
        uint256 afterDeposit = address(this).balance;
        assertEq(beforeDeposit - afterDeposit, 1 ether);
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
        (uint256[2] memory publicKey, uint8[6] memory pathIndices, uint256[6] memory siblings, uint256[2] memory a, uint256[2][2] memory b, uint256[2] memory c, uint256[3] memory input) = abi.decode(jsResult, (uint256[2], uint8[6], uint256[6], uint256[2], uint256[2][2], uint256[2], uint256[3]));
        uint256 publicKeyX = publicKey[0];
        uint256 publicKeyY = publicKey[1];
        
        rollup.deposit{value: 1 ether}(publicKeyX, publicKeyY);

        uint256 keyHash = rollup.generateKeyHash(publicKeyX, publicKeyY);
        bool isRegistered = rollup.isPublicKeysRegistered(keyHash);
        assertEq(isRegistered, true);

        uint256[] memory proofSiblings = new uint256[](6);
        uint8[] memory proofPathIndices = new uint8[](6);
        for (uint8 i = 0; i < 6; i++) {
            proofSiblings[i] = siblings[i];
            proofPathIndices[i] = pathIndices[i];
        }

        uint256 beforeWithdraw = address(this).balance;
        rollup.withdraw(0.5e18, a, b, c, input, proofSiblings, proofPathIndices);
        uint256 afterWithdraw = address(this).balance;

        assertEq(afterWithdraw - beforeWithdraw, 0.5e18);
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
        (uint256[2] memory publicKey, uint8[6] memory pathIndices, uint256[6] memory siblings, uint256[2] memory a, uint256[2][2] memory b, uint256[2] memory c, uint256[3] memory input) = abi.decode(jsResult, (uint256[2], uint8[6], uint256[6], uint256[2], uint256[2][2], uint256[2], uint256[3]));
        uint256 publicKeyX = publicKey[0];
        uint256 publicKeyY = publicKey[1];
        
        rollup.deposit{value: 1 ether}(publicKeyX, publicKeyY);

        uint256 keyHash = rollup.generateKeyHash(publicKeyX, publicKeyY);
        bool isRegistered = rollup.isPublicKeysRegistered(keyHash);
        assertEq(isRegistered, true);

        uint256[] memory proofSiblings = new uint256[](6);
        uint8[] memory proofPathIndices = new uint8[](6);
        for (uint8 i = 0; i < 6; i++) {
            proofSiblings[i] = siblings[i];
            proofPathIndices[i] = pathIndices[i];
        }

        uint256 beforeWithdraw = address(this).balance;
        rollup.withdraw(a, b, c, input, proofSiblings, proofPathIndices);
        uint256 afterWithdraw = address(this).balance;

        assertEq(afterWithdraw - beforeWithdraw, 1e18);
    }

    receive() external payable {}
    fallback() external payable {}
}