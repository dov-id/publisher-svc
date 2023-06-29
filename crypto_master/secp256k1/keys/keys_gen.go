package keys

import (
	"fmt"
	"math/big"

	"github.com/dov-id/publisher-svc/crypto_master/secp256k1/ecc_math"
	"github.com/dov-id/publisher-svc/crypto_master/secp256k1/random"
)

type KeyPair struct {
	privateKey big.Int
	PublicKey  ecc_math.ECPoint
}

func GetPrivateKey(pair KeyPair) big.Int {
	return pair.privateKey
}

func GenKeyPair() (pair KeyPair) {
	private := random.GenerateRandomBigInt()
	public := ecc_math.ScalarMult(ecc_math.BasePointGGet(), *private)
	pair.privateKey = *private
	pair.PublicKey = public
	return
}

func KeyPairToString(pair KeyPair) string {
	return fmt.Sprintf("%X", &pair.privateKey) + " " + ecc_math.ECPointToString(pair.PublicKey)
}

func PrintKeyPair(pair KeyPair) {
	fmt.Println("Private key:\t", fmt.Sprintf("%X", &pair.privateKey), "\nPublic key:")
	ecc_math.PrintECPoint(pair.PublicKey)
}
