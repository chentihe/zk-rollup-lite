// SPDX-License-Identifier: UNLICENSED
pragma solidity >0.8.0 <=0.9;

import {Test} from "forge-std/Test.sol";
import {StdStorage, stdStorage} from "forge-std/StdStorage.sol";
import {Rollup} from "../src/Rollup.sol";
import {TxVerifier} from "../src/verifiers/TxVerifier.sol";
import {WithdrawVerifier} from "../src/verifiers/WithdrawVerifier.sol";
import {DepositVerifier} from "../src/verifiers/DepositVerifier.sol";

contract RollupTest is Test {
    using stdStorage for StdStorage;

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
        runJsInputs[5] = "generate-zkp-deposit";
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
        (uint256[2] memory a, uint256[2][2] memory b, uint256[2] memory c, uint256[5] memory input) = abi.decode(jsResult, (uint256[2], uint256[2][2], uint256[2], uint256[5]));

        uint256 newRoot = input[0];
        uint256 root = input[1];
        uint256 publicKeyX = input[2];
        uint256 publicKeyY = input[3];
        uint256 nullifier = input[4];

        // setup root
        stdstore
            .target(address(rollup))
            .sig(0xb4e7dddd)
            .checked_write(root);

        // setup balance
        stdstore.target(address(rollup))
            .sig("getUserByPublicKey(uint256,uint256)")
            .with_key(publicKeyX)
            .with_key(publicKeyY)
            .depth(2)
            .checked_write(1e18);

        // transfer token
        (bool success, ) = address(rollup).call{value: 1e18}("");
        require(success, "transfer ether failed");

        assertEq(rollup.balanceTreeRoot(), root);

        uint256 beforeWithdraw = address(this).balance;
        rollup.withdraw(0.5e18, a, b, c, input);
        uint256 afterWithdraw = address(this).balance;

        assertEq(afterWithdraw - beforeWithdraw, 0.5e18);
        assertEq(newRoot, rollup.balanceTreeRoot());
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
        (uint256[2] memory a, uint256[2][2] memory b, uint256[2] memory c, uint256[5] memory input) = abi.decode(jsResult, (uint256[2], uint256[2][2], uint256[2], uint256[5]));

        uint256 newRoot = input[0];
        uint256 root = input[1];
        uint256 publicKeyX = input[2];
        uint256 publicKeyY = input[3];
        uint256 nullifier = input[4];
        
        // setup root
        stdstore
            .target(address(rollup))
            .sig(0xb4e7dddd)
            .checked_write(root);

        // setup balance
        stdstore.target(address(rollup))
            .sig("getUserByPublicKey(uint256,uint256)")
            .with_key(publicKeyX)
            .with_key(publicKeyY)
            .depth(2)
            .checked_write(1e18);

        // transfer token
        (bool success, ) = address(rollup).call{value: 1e18}("");
        require(success, "transfer ether failed");

        uint256 beforeWithdraw = address(this).balance;
        rollup.withdraw(type(uint256).max, a, b, c, input);
        uint256 afterWithdraw = address(this).balance;

        assertEq(afterWithdraw - beforeWithdraw, 1e18);
        assertEq(newRoot, rollup.balanceTreeRoot());
        assertTrue(rollup.usedNullifiers(nullifier));
    }

    receive() external payable {}
    fallback() external payable {}
}