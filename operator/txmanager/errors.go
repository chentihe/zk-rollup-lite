package txmanager

import "fmt"

var (
	ErrAmountNil     = fmt.Errorf("Amount should not be nil")
	ErrAmountTooLow  = fmt.Errorf("Amount should be larger than %s", minAmount.String())
	ErrAmountTooHigh = fmt.Errorf("Amount should not be larger than %s", maxAmount.String())

	ErrInvalidNonce = fmt.Errorf("Invalid nonce")

	ErrNonceTooLow = fmt.Errorf("Nonce should not be less than %d", minNonce)

	ErrFeeAmountNil     = fmt.Errorf("FeeAmount should not be nil")
	ErrFeeAmountTooLow  = fmt.Errorf("FeeAmount should not be less than %s", minFeeAmount.String())
	ErrFeeAmountTooHigh = fmt.Errorf("FeeAmount should not be larger than %s", maxFeeAmount.String())

	ErrFromAccountIndexTooLow  = fmt.Errorf("FromAccountIndex should not be less than %d", minAccountIndex)
	ErrFromAccountIndexTooHigh = fmt.Errorf("FromAccountIndex should not be larger than %d", maxAccountIndex)
	ErrToAccountIndexTooLow    = fmt.Errorf("ToAccountIndex should not be less than %d", minAccountIndex)
	ErrToAccountIndexTooHigh   = fmt.Errorf("ToAccountIndex should not be larger than %d", maxAccountIndex)

	ErrInvalidSignature = fmt.Errorf("Invalid signature")
)
