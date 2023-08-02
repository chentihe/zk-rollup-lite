source .env.anvil
forge script script/Rollup.s.sol:RollupScript \
--rpc-url $RPC_URL \
--private-key $PRIVATE_KEY --broadcast