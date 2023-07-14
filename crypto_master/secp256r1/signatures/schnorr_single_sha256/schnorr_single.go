package schnorr_single_sha256

import (
	"fmt"
	"math/big"

	"github.com/dov-id/publisher-svc/crypto_master/hash"
	"github.com/dov-id/publisher-svc/crypto_master/secp256r1/ecc_math"
	"github.com/dov-id/publisher-svc/crypto_master/secp256r1/random"
)

/*Signatures are
			(R,s) = (r*G, r + H(X,R,m)*x)
Verification
			sG ?= R + H(X,R,m)X */

var (
	mod = ecc_math.Curve.Params().N
	G   = ecc_math.BasePointGGet()
)

type SchnorrSignature struct {
	R ecc_math.ECPoint
	S big.Int
}

func SchnorrSignatureGen(message string, publicKey ecc_math.ECPoint, privateKey big.Int) (signature SchnorrSignature) {
	r := random.GenerateRandomBigInt()
	R := ecc_math.ScalarMult(G, *r)

	messageHash := new(big.Int)
	messageHash.SetString(fmt.Sprintf("%x", hash.SHA256StringToString(ecc_math.ECPointToString(publicKey)+ecc_math.ECPointToString(R)+message)), 16) // messageHash

	tmp := new(big.Int)
	tmp.Mul(messageHash, &privateKey)
	tmp.Mod(tmp, mod)
	r.Add(r, tmp)
	r.Mod(r, mod) //	r + H(X,R,m)*x

	signature.R = R
	signature.S = *r

	return
}

func SchnorrSignatureVerify(message string, publicKey ecc_math.ECPoint, signature SchnorrSignature) bool {
	messageHash := new(big.Int)
	messageHash.SetString(fmt.Sprintf("%x", hash.SHA256StringToString(ecc_math.ECPointToString(publicKey)+ecc_math.ECPointToString(signature.R)+message)), 16) // messageHash
	// messageHash

	left := ecc_math.ScalarMult(G, signature.S)
	right := ecc_math.ScalarMult(publicKey, *messageHash)
	right = ecc_math.AddECPoints(right, signature.R)
	if left.X.Cmp(right.X) == 0 && left.Y.Cmp(right.Y) == 0 {
		return true
	} else {
		return false
	}
}
