# TODO: replace with script to generate the verifier

VERIFIER_NAME="$1"
PARENT_FOLDER="$(dirname "$(pwd)")"
echo "${PARENT_FOLDER}"
time snarkjs zkey export solidityverifier "${PARENT_FOLDER}/operator/build/${VERIFIER_NAME}/circuit_final.zkey" "${PARENT_FOLDER}/contracts/src/verifiers/${VERIFIER_NAME}Verifier.sol"
