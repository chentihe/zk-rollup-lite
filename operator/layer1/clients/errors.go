package clients

import "fmt"

var (
	ErrPubKeyToECDSA = fmt.Errorf("cannot cast public key to ECDSA")
)
