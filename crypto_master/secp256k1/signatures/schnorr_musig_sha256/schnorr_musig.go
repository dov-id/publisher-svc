package schnorr_musig_sha256

import (
	"fmt"
	"math/big"

	"github.com/dov-id/publisher-svc/crypto_master/hash"
	"github.com/dov-id/publisher-svc/crypto_master/secp256k1/ecc_math"
	"github.com/dov-id/publisher-svc/crypto_master/secp256k1/random"
	"github.com/dov-id/publisher-svc/crypto_master/secp256k1/signatures/schnorr_single_sha256"
)

/*
1. L = H(X1,X2,…)
2. X = sum(H(L,Xi)*Xi)
3. ri - rand, Ri = ri*G
4. R = sum(Ri)
5. Each signer computes si = ri + H(X,R,m)*H(L,Xi)*xi
6. The final signature is (R,s) where s is the sum of the si values
7. Verification sG = R + H(X,R,m)*X
*/

var (
	mod = ecc_math.Curve.Params().N
	G   = ecc_math.BasePointGGet()
)

const N = 5000

type RSinglePair struct {
	R ecc_math.ECPoint
	r big.Int
}

type CommonParameters struct {
	L                   big.Int
	AggregatedPublicKey ecc_math.ECPoint
	R                   ecc_math.ECPoint
}

func FormSingleRPair() (rPair RSinglePair) {
	r := random.GenerateRandomBigInt()
	R := ecc_math.ScalarMult(G, *r)
	rPair.R = R
	rPair.r = *r
	return
} // ri - rand, Ri = ri*G

func CommonParametersGen(publicKeySet [N]ecc_math.ECPoint, RSet [N]ecc_math.ECPoint) (params CommonParameters) {

	stringL := ""
	for i := 0; i < N; i++ {
		stringL += ecc_math.ECPointToString(publicKeySet[i])
	}

	L := new(big.Int)
	L.SetString(fmt.Sprintf("%X", hash.SHA256StringToString(stringL)), 16) // L = H(X1,X2,…)

	Hash := new(big.Int)
	Hash.SetString(fmt.Sprintf("%X", hash.SHA256StringToString(L.String()+ecc_math.ECPointToString(publicKeySet[0]))), 16)
	aggregatedPublicKey := ecc_math.ScalarMult(publicKeySet[0], *Hash)
	for i := 1; i < N; i++ {
		Hash.SetString(fmt.Sprintf("%X", hash.SHA256StringToString(L.String()+ecc_math.ECPointToString(publicKeySet[i]))), 16)
		aggregatedPublicKey = ecc_math.AddECPoints(aggregatedPublicKey, ecc_math.ScalarMult(publicKeySet[i], *Hash))
	}
	// X = sum(H(L,Xi)*Xi)

	aggregatedR := RSet[0]
	for i := 1; i < N; i++ {
		aggregatedR = ecc_math.AddECPoints(aggregatedR, RSet[i])
	}
	//R = sum(Ri)

	params.R = aggregatedR
	params.AggregatedPublicKey = aggregatedPublicKey
	params.L = *L

	return
}

func FormSignaturePart(message string, rPair RSinglePair, parameters CommonParameters, publicKey ecc_math.ECPoint, privateKey big.Int) (signaturePart big.Int) {
	messageHashL := new(big.Int)
	messageHashL.SetString(fmt.Sprintf("%x", hash.SHA256StringToString(ecc_math.ECPointToString(parameters.AggregatedPublicKey)+ecc_math.ECPointToString(parameters.R)+message)), 16) // H(X,R,m)

	messageHashR := new(big.Int)
	messageHashR.SetString(fmt.Sprintf("%x", hash.SHA256StringToString(parameters.L.String()+ecc_math.ECPointToString(publicKey))), 16) // H(L,Xi)

	signaturePart.Mul(messageHashL, messageHashR)
	signaturePart.Mod(&signaturePart, mod)
	signaturePart.Mul(&signaturePart, &privateKey)
	signaturePart.Mod(&signaturePart, mod)
	signaturePart.Add(&signaturePart, &rPair.r)
	signaturePart.Mod(&signaturePart, mod)

	return
}

//si = ri + H(X,R,m)*H(L,Xi)*xi

func AggregareSignature(parameters CommonParameters, sigList [N]big.Int) (signature schnorr_single_sha256.SchnorrSignature) {
	signature.S = sigList[0]
	for i := 1; i < len(sigList); i++ {
		signature.S.Add(&signature.S, &sigList[i])
		signature.S.Mod(&signature.S, mod)
	}
	signature.R = parameters.R
	return
}
