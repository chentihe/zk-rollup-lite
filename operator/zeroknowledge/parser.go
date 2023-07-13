package zeroknowledge

import (
	"bytes"
	"fmt"
	"math/big"
	"strings"
)

func stringsToArrayBigInt(publicInputs []string) ([]*big.Int, error) {
	p := make([]*big.Int, 0, len(publicInputs))
	for _, s := range publicInputs {
		sb, err := stringToBigInt(s)
		if err != nil {
			return nil, err
		}
		p = append(p, sb)
	}
	return p, nil
}

func stringToBigInt(s string) (*big.Int, error) {
	base := 10
	if bytes.HasPrefix([]byte(s), []byte("0x")) {
		base = 16
		s = strings.TrimPrefix(s, "0x")
	}
	n, ok := new(big.Int).SetString(s, base)
	if !ok {
		return nil, fmt.Errorf("can not parse string to *big.Int: %s", s)
	}
	return n, nil
}
