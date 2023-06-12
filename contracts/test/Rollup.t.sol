// SPDX-License-Identifier: UNLICENSED
pragma solidity >0.8.0 <=0.9;

import {Test} from "forge-std/Test.sol";
import {Rollup} from "../src/Rollup.sol";
import {TxVerifier} from "../src/verifiers/TxVerifier.sol";
import {WithdrawVerifier} from "../src/verifiers/WithdrawVerifier.sol";
import {DepositVerifier} from "../src/verifiers/DepositVerifier.sol";

contract RollupTest is Test {
    Rollup rollup;
    TxVerifier txVerifier;
    WithdrawVerifier withdrawVerifier;
    DepositVerifier depositVerifier;
    uint256 constant DEPTH = 6;

    function setUp() public {
        txVerifier = new TxVerifier();
        withdrawVerifier = new WithdrawVerifier();
        depositVerifier = new DepositVerifier();
        rollup = new Rollup(txVerifier, withdrawVerifier, depositVerifier);
    }

    function testDeposit() public {
        string[] memory runJsInputs = new string[](6);
        runJsInputs[0] = "npm";
        runJsInputs[1] = "--prefix";
        runJsInputs[2] = "contracts/test/script/";
        runJsInputs[3] = "--silent";
        runJsInputs[4] = "run";
        runJsInputs[5] = "generate-merkle-tree-proof";
        bytes memory jsResult = vm.ffi(runJsInputs);
        (uint256[2] memory a, uint256[2][2] memory b, uint256[2] memory c, uint256[4] memory input) = abi.decode(jsResult, (uint256[2], uint256[2][2], uint256[2], uint256[4]));
        
        uint256 newRoot = input[0];
        uint256 root = input[1];
        uint256 publicKeyX = input[2];
        uint256 publicKeyY = input[3];

        assertEq(rollup.balanceTreeRoot(), root);

        rollup.deposit{value: 1 ether}(a, b, c, input);

        uint256 keyHash = rollup.generateKeyHash(publicKeyX, publicKeyY);
        bool isRegistered = rollup.isPublicKeysRegistered(keyHash);
        assertEq(isRegistered, true);

        assertEq(rollup.balanceTreeRoot(), newRoot);
    }

    // function testWithdraw() public {
    //     string[] memory runJsInputs = new string[](8);
    //     runJsInputs[0] = "npm";
    //     runJsInputs[1] = "--prefix";
    //     runJsInputs[2] = "contracts/test/script/";
    //     runJsInputs[3] = "--silent";
    //     runJsInputs[4] = "run";
    //     runJsInputs[5] = "generate-zkp-withdraw";
    //     runJsInputs[6] = "withdraw";
    //     runJsInputs[7] = "0.5";
    //     bytes memory jsResult = vm.ffi(runJsInputs);
    //     (uint256[2] memory publicKey, uint8[6] memory pathIndices, uint256[6] memory siblings, uint256[2] memory a, uint256[2][2] memory b, uint256[2] memory c, uint256[3] memory input) = abi.decode(jsResult, (uint256[2], uint8[6], uint256[6], uint256[2], uint256[2][2], uint256[2], uint256[3]));
    //     uint256 publicKeyX = publicKey[0];
    //     uint256 publicKeyY = publicKey[1];
        
    //     rollup.deposit{value: 1 ether}(publicKeyX, publicKeyY);

    //     uint256 keyHash = rollup.generateKeyHash(publicKeyX, publicKeyY);
    //     bool isRegistered = rollup.isPublicKeysRegistered(keyHash);
    //     assertEq(isRegistered, true);

    //     uint256[] memory proofSiblings = new uint256[](6);
    //     uint8[] memory proofPathIndices = new uint8[](6);
    //     for (uint8 i = 0; i < 6; i++) {
    //         proofSiblings[i] = siblings[i];
    //         proofPathIndices[i] = pathIndices[i];
    //     }

    //     uint256 beforeWithdraw = address(this).balance;
    //     rollup.withdraw(0.5e18, a, b, c, input, proofSiblings, proofPathIndices);
    //     uint256 afterWithdraw = address(this).balance;

    //     assertEq(afterWithdraw - beforeWithdraw, 0.5e18);
    // }

    // function testWithdrawAll() public {
    //     string[] memory runJsInputs = new string[](8);
    //     runJsInputs[0] = "npm";
    //     runJsInputs[1] = "--prefix";
    //     runJsInputs[2] = "contracts/test/script/";
    //     runJsInputs[3] = "--silent";
    //     runJsInputs[4] = "run";
    //     runJsInputs[5] = "generate-zkp-withdraw";
    //     runJsInputs[6] = "withdraw-all";
    //     runJsInputs[7] = "0.5";
    //     bytes memory jsResult = vm.ffi(runJsInputs);
    //     (uint256[2] memory publicKey, uint8[6] memory pathIndices, uint256[6] memory siblings, uint256[2] memory a, uint256[2][2] memory b, uint256[2] memory c, uint256[3] memory input) = abi.decode(jsResult, (uint256[2], uint8[6], uint256[6], uint256[2], uint256[2][2], uint256[2], uint256[3]));
    //     uint256 publicKeyX = publicKey[0];
    //     uint256 publicKeyY = publicKey[1];
        
    //     rollup.deposit{value: 1 ether}(publicKeyX, publicKeyY);

    //     uint256 keyHash = rollup.generateKeyHash(publicKeyX, publicKeyY);
    //     bool isRegistered = rollup.isPublicKeysRegistered(keyHash);
    //     assertEq(isRegistered, true);

    //     uint256[] memory proofSiblings = new uint256[](6);
    //     uint8[] memory proofPathIndices = new uint8[](6);
    //     for (uint8 i = 0; i < 6; i++) {
    //         proofSiblings[i] = siblings[i];
    //         proofPathIndices[i] = pathIndices[i];
    //     }

    //     uint256 beforeWithdraw = address(this).balance;
    //     rollup.withdraw(a, b, c, input, proofSiblings, proofPathIndices);
    //     uint256 afterWithdraw = address(this).balance;

    //     assertEq(afterWithdraw - beforeWithdraw, 1e18);
    // }

    receive() external payable {}
    fallback() external payable {}
}