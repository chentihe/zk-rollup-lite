VERSION="$1"

time snarkjs powersoftau new bn128 "$VERSION" "pot${VERSION}_0000.ptau" -v
ENTROPY1=$(head -c 1024 /dev/urandom | LC_CTYPE=C tr -dc 'a-zA-Z0-9' | head -c 128)
time snarkjs powersoftau contribute "pot${VERSION}_0000.ptau" "pot${VERSION}_0001.ptau" --name="First contribute" -e="$ENTROPY1"
time snarkjs powersoftau prepare phase2 "pot${VERSION}_0001.ptau" "pot${VERSION}_final.ptau" -v

mkdir -p "trusted_setup"
mv "pot${VERSION}_final.ptau" "trusted_setup/pot${VERSION}_final.ptau"
rm "pot${VERSION}_0000.ptau"
rm "pot${VERSION}_0001.ptau"