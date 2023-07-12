package pubsub

import "fmt"

var (
	ErrPubKeyToECDSA = fmt.Errorf("Cannot cast public key to ECDSA")
)
