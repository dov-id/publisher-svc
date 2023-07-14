package random

import (
	"crypto/rand"
	"math/big"

	"github.com/dov-id/publisher-svc/crypto_master/secp256r1/ecc_math"
)

func GenerateRandomBigInt() *big.Int {
	n, err := rand.Int(rand.Reader, ecc_math.Curve.Params().N)
	if err == nil {
		return n
	}
	return nil
}
