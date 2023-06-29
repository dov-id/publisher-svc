package ring_sha256

import (
	"fmt"
	"math/big"

	"github.com/dov-id/publisher-svc/crypto_master/secp256k1/ecc_math"
	"github.com/dov-id/publisher-svc/crypto_master/secp256k1/random"

	"github.com/dov-id/publisher-svc/crypto_master/hash"
)

/*
Ring secp256r1 signature algorithm

Signature input:
	M 			- message
	A[0...n-1] 	- public keys of ring participants,
	а[i]		- private signer key,
	А[i] 		- public signer key.

Procedure generate_signature (M, A[1], A[2], ..., A[n], i, a[i]):

1. I <- a[i]*H(A[i])											// Private key image
2. c[j], r[j] [j=0..n-1, j!=i] <- random
3. k <- random
4. For j <- 0..n-1, j!=i
		4.1. X[j] <- c[j]*A[j]+r[j]*G
		4.2. Y[j] <- c[j]*I+r[j]*H(A[j])
5. X[i] <- k*G
6. Y[i] <- k*H(A[i])
7. c[i] <- H(H(M) || X[0] || Y[0] || X[1] || Y[1] || ... || X[n-1] || Y[n-1]) - Sum[j=0..n-1, j!=i](c[j])
8. r[i] <- k-a[i]*c[i]

Return (I, c[0] || r[0] || c[1] || r[1] || ... || c[n-1] || r[n-1])


Procedure verify_signature(M, A[0], A[1], ..., A[n-1], I, c[0], r[0], c[1], r[1], ..., c[n-1], r[n-1]):

1. For i <- 0..n-1
		1.1. X[i] <- c[i]*A[i]+r[i]*G
		1.2. Y[i] <- c[i]*I+r[i]*H(A[i])
2. If H(H(M) || X[0] || Y[0] || X[1] || Y[1] || ... || X[n-1] || Y[n-1])  == Sum[i=0..n-1](c[i])
		Return "Correct"
Else
		Return "Incorrect"
*/

var (
	Mod = ecc_math.Curve.Params().N
	G   = ecc_math.BasePointGGet()
)

const N = 5 // ring size

type RingSignature struct {
	I big.Int
	c [N]big.Int
	r [N]big.Int
}

func RingSignatureGen(message string, ringPublicKeys [N]ecc_math.ECPoint, index int, privateKey big.Int) (signature RingSignature) {
	messageHash := new(big.Int)
	messageHash.SetString(fmt.Sprintf("%x", hash.SHA256StringToString(message)), 16) // messageHash

	pubKeyHash := new(big.Int)
	pubKeyHash.SetString(ecc_math.ECPointToString(ringPublicKeys[index]), 16) // PubKeyHash

	I := new(big.Int)
	I.Mul(&privateKey, pubKeyHash)
	I.Mod(I, Mod) // I <- a[i]*H(A[i])

	var c [N]big.Int
	var r [N]big.Int

	for i := 0; i < N; i++ {
		if i != index {
			c[i] = *random.GenerateRandomBigInt()
			r[i] = *random.GenerateRandomBigInt()
		}
	} // c[j], r[j] [j=0..n-1, j!=i] <- random

	k := *random.GenerateRandomBigInt() // k <- random

	var X [N]ecc_math.ECPoint
	var Y [N]big.Int

	for i := 0; i < N; i++ {
		if i != index {
			X[i] = ecc_math.AddECPoints(ecc_math.ScalarMult(ringPublicKeys[i], c[i]), ecc_math.ScalarMult(G, r[i]))
			tmp := new(big.Int)
			tmp.Mul(&c[i], I)
			pubKeyHashA := new(big.Int)
			pubKeyHashA.SetString(ecc_math.ECPointToString(ringPublicKeys[i]), 16)
			Y[i].Mul(&r[i], pubKeyHashA)
			Y[i].Add(&Y[i], tmp)
			Y[i].Mod(&Y[i], Mod)
		} // For j <- 0..n-1, j!=i
		// 		X[j] <- c[j]*A[j]+r[j]*G
		// 		Y[j] <- c[j]*I+r[j]*H(A[j])
	}

	X[index] = ecc_math.ScalarMult(G, k) // X[i] <- k*G

	Y[index].Mul(&k, pubKeyHash) // Y[i] <- k*H(A[i])
	Y[index].Mod(&Y[index], Mod)

	tempStringS := messageHash.String() // H(M)
	for i := 0; i < N; i++ {
		tempStringS += "" + ecc_math.ECPointToString(X[i])
		tempStringS += "" + Y[i].String() // H(M) || X[0] || Y[0] || X[1] || Y[1] || ... || X[n-1] || Y[n-1]
	}

	c[index].SetString(fmt.Sprintf("%X", hash.SHA256StringToString(tempStringS)), 16) // H(H(M) || X[0] || Y[0] || X[1] || Y[1] || ... || X[n-1] || Y[n-1])

	var sum big.Int
	for i := 0; i < N; i++ {
		if i != index {
			sum.Add(&sum, &c[i])
		}
		sum.Mod(&sum, Mod)
	} // Sum[j=0..n-1, j!=i](c[j])

	c[index].Sub(&c[index], &sum) //c[i] <- H(H(M) || X[0] || Y[0] || X[1] || Y[1] || ... || X[n-1] || Y[n-1]) - Sum[j=0..n-1, j!=i](c[j])
	c[index].Mod(&c[index], Mod)

	r[index].Mul(&privateKey, &c[index])
	r[index].Sub(&k, &r[index])
	r[index].Mod(&r[index], Mod) // r[i] <- k-a[i]*c[i]

	signature.I = *I
	signature.c = c
	signature.r = r

	return
}

func RingSignatureVerify(message string, ringPublicKeys [N]ecc_math.ECPoint, signature RingSignature) bool {
	messageHash := new(big.Int)
	messageHash.SetString(fmt.Sprintf("%x", hash.SHA256StringToString(message)), 16)

	var X [N]ecc_math.ECPoint
	var Y [N]big.Int

	for i := 0; i < N; i++ {
		X[i] = ecc_math.AddECPoints(ecc_math.ScalarMult(ringPublicKeys[i], signature.c[i]), ecc_math.ScalarMult(G, signature.r[i]))
		tmp := new(big.Int)
		tmp.Mul(&signature.c[i], &signature.I)
		pubKeyHashA := new(big.Int)
		pubKeyHashA.SetString(ecc_math.ECPointToString(ringPublicKeys[i]), 16)
		Y[i].Mul(&signature.r[i], pubKeyHashA)
		Y[i].Add(&Y[i], tmp)
		Y[i].Mod(&Y[i], Mod) //	For i <- 0..n-1
		//		X[i] <- c[i]*A[i]+r[i]*G
		//		Y[i] <- c[i]*I+r[i]*H(A[i])
	}

	left := new(big.Int)
	tempStringS := messageHash.String() // H(M)
	for i := 0; i < N; i++ {
		tempStringS += "" + ecc_math.ECPointToString(X[i])
		tempStringS += "" + Y[i].String() // H(M) || X[0] || Y[0] || X[1] || Y[1] || ... || X[n-1] || Y[n-1]
	}

	left.SetString(fmt.Sprintf("%X", hash.SHA256StringToString(tempStringS)), 16) // H(H(M) || X[0] || Y[0] || X[1] || Y[1] || ... || X[n-1] || Y[n-1])

	var sum big.Int
	for i := 0; i < N; i++ {
		sum.Add(&sum, &signature.c[i])
		sum.Mod(&sum, Mod)
	} // Sum[j=0..n-1, j!=i](c[j])
	if left.Cmp(&sum) == 0 {
		return true
	} else {
		return false
	}
}

func PrintSignatureValue(signature RingSignature) {
	fmt.Println("Signature: ")
	fmt.Print("\tPrivate key image:\t")
	fmt.Printf("%X", &signature.I)
	fmt.Println()
	fmt.Println("\n\tC array:")
	for i := 0; i < N; i++ {
		fmt.Println("\t\t\t\t", fmt.Sprintf("%X", &signature.c[i]), ";")
	}
	fmt.Println("\n\tR array: ")
	for i := 0; i < N; i++ {
		fmt.Println("\t\t\t\t", fmt.Sprintf("%X", &signature.r[i]), ";")
	}
}
