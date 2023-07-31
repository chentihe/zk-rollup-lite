// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import "forge-std/Script.sol";
import {Rollup} from "../src/Rollup.sol";
import {DepositVerifier} from "../src/verifiers/DepositVerifier.sol";
import {TxVerifier} from "../src/verifiers/TxVerifier.sol";
import {WithdrawVerifier} from "../src/verifiers/WithdrawVerifier.sol";

contract RollupScript is Script {
    function run() public {
        vm.startBroadcast();
        TxVerifier txVerifier = new TxVerifier();
        WithdrawVerifier withdrawVerifier = new WithdrawVerifier();

        new Rollup(txVerifier, withdrawVerifier);
        vm.stopBroadcast();
    }
}
