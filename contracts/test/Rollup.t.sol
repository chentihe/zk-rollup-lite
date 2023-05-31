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
}
