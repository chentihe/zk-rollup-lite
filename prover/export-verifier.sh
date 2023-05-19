# TODO: replace with script to generate the verifier

VERIFIER_NAME="$1"
PARENT_FOLDER="$(dirname "$(pwd)")"
echo "${PARENT_FOLDER}"
time snarkjs zkey export solidityverifier "build/${VERIFIER_NAME}/circuit_0000.zkey" "${VERIFIER_NAME}Verifier.sol"

mv "${VERIFIER_NAME}Verifier.sol" "${PARENT_FOLDER}/contracts/src/verifiers/${VERIFIER_NAME}Verifier.sol"
