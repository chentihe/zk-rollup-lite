package cache

import "fmt"

var (
	ErrPubKeyToECDSA   = fmt.Errorf("Cannot cast public key to ECDSA")
	ErrInsertedTxToInt = fmt.Errorf("Cannot cast inserted tx to int")
)
