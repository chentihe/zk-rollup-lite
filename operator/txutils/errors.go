package txutils

import "fmt"

var (
	ErrAmountNil     = fmt.Errorf("amount should not be nil")
	ErrAmountTooLow  = fmt.Errorf("amount should be larger than %s", minAmount.String())
	ErrAmountTooHigh = fmt.Errorf("amount should not be larger than %s", maxAmount.String())

	ErrInvalidNonce = fmt.Errorf("invalid nonce")
	ErrNonceTooLow  = fmt.Errorf("nonce should not be less than %d", minNonce)

	ErrFeeAmount        = fmt.Errorf("fee amount should be %s", Fee.String())
	ErrFeeAmountTooLow  = fmt.Errorf("fee amount should not be less than %s", minFeeAmount.String())
	ErrFeeAmountTooHigh = fmt.Errorf("fee amount should not be larger than %s", maxFeeAmount.String())

	ErrFromAccountIndexTooLow  = fmt.Errorf("from account index should not be less than %d", minAccountIndex)
	ErrFromAccountIndexTooHigh = fmt.Errorf("from account index should not be larger than %d", maxAccountIndex)
	ErrToAccountIndexTooLow    = fmt.Errorf("to account index should not be less than %d", minAccountIndex)
	ErrToAccountIndexTooHigh   = fmt.Errorf("to account index should not be larger than %d", maxAccountIndex)

	ErrInvalidSignature = fmt.Errorf("invalid signature")

	ErrAccountNotExist = fmt.Errorf("account doesn't exist")
)
