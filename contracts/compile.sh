solc --abi src/Rollup.sol -o out/Rollup.abi
solc --bin src/Rollup.sol -o out/Rollup.bin
abigen --abi out/Rollup.abi --pkg contracts --type Rollup --out out/Rollup.go --bin out/Rollup.bin/Rollup.bin