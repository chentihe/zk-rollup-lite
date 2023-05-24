// SPDX-License-Identifier: MIT
pragma solidity >0.8.0 <=0.9;

library Constants {
    uint256 internal constant SNARK_SCALAR_FIELD =
        21888242871839275222246405745257275088548364400416034343698204186575808495617;
    uint256 constant ZERO_VALUE = uint256(keccak256(abi.encodePacked("ROLLUP_LITE"))) % SNARK_SCALAR_FIELD;
    uint256 constant UINT256_MAX = type(uint256).max;
}