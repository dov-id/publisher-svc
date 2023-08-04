package ring_sha256

import (
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/dov-id/publisher-svc/crypto_master/secp256k1/ecc_math"
	"github.com/dov-id/publisher-svc/crypto_master/secp256k1/random"
	"github.com/ethereum/go-ethereum/common"
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

type RingSignature struct {
	I big.Int
	C []big.Int
	R []big.Int
}

func RingSignatureGenerate(message string, ringPublicKeys []ecc_math.ECPoint, index int, privateKey big.Int) (signature RingSignature) {
	participantsLen := len(ringPublicKeys)

	hash := sha256.Sum256([]byte(message))
	messageHash := new(big.Int).SetBytes(hash[:])

	// PubKeyHash
	hash = sha256.Sum256(append(ringPublicKeys[index].X.Bytes(), ringPublicKeys[index].Y.Bytes()...))
	pubKeyHash := new(big.Int).SetBytes(hash[:])

	// I <- a[i]*H(A[i])
	I := new(big.Int)
	I.Mul(&privateKey, pubKeyHash)
	I.Mod(I, Mod)

	var c = make([]big.Int, participantsLen)
	var r = make([]big.Int, participantsLen)

	// c[j], r[j] [j=0..n-1, j!=i] <- random
	for i := 0; i < participantsLen; i++ {
		if i != index {
			c[i] = *random.GenerateRandomBigInt()
			r[i] = *random.GenerateRandomBigInt()
		}
	}

	// k <- random
	k := *random.GenerateRandomBigInt()

	var X = make([]ecc_math.ECPoint, participantsLen)
	var Y = make([]big.Int, participantsLen)

	// For j <- 0..n-1, j!=i
	for i := 0; i < participantsLen; i++ {
		if i != index {
			// X[j] <- c[j]*A[j]+r[j]*G
			X[i] = ecc_math.AddECPoints(ecc_math.ScalarMult(ringPublicKeys[i], c[i]), ecc_math.ScalarMult(G, r[i]))

			hash = sha256.Sum256(append(ringPublicKeys[i].X.Bytes(), ringPublicKeys[i].Y.Bytes()...))
			pubKeyHashA := new(big.Int).SetBytes(hash[:])

			// Y[j] <- c[j]*I+r[j]*H(A[j])
			Y[i].Mul(&r[i], pubKeyHashA)
			Y[i].Add(&Y[i], new(big.Int).Mul(&c[i], I))
			Y[i].Mod(&Y[i], Mod)
		}
	}

	// X[i] <- k*G
	X[index] = ecc_math.ScalarMult(G, k)

	// Y[i] <- k*H(A[i])
	Y[index].Mul(&k, pubKeyHash)
	Y[index].Mod(&Y[index], Mod)

	tempBytes := messageHash.Bytes() // H(M)
	for i := 0; i < participantsLen; i++ {
		// H(M) || X[0] || Y[0] || X[1] || Y[1] || ... || X[n-1] || Y[n-1]
		tempBytes = append(tempBytes, X[i].X.Bytes()...)
		tempBytes = append(tempBytes, X[i].Y.Bytes()...)
		tempBytes = append(tempBytes, Y[i].Bytes()...)
	}

	// H(H(M) || X[0] || Y[0] || X[1] || Y[1] || ... || X[n-1] || Y[n-1])
	hash = sha256.Sum256(tempBytes)
	c[index].SetBytes(hash[:])

	var sum big.Int
	// Sum[j=0..n-1, j!=i](c[j])
	for i := 0; i < participantsLen; i++ {
		if i != index {
			sum.Add(&sum, &c[i])
		}
		sum.Mod(&sum, Mod)
	}

	//c[i] <- H(H(M) || X[0] || Y[0] || X[1] || Y[1] || ... || X[n-1] || Y[n-1]) - Sum[j=0..n-1, j!=i](c[j])
	c[index].Sub(&c[index], &sum)
	c[index].Mod(&c[index], Mod)

	// r[i] <- k-a[i]*c[i]
	r[index].Mul(&privateKey, &c[index])
	r[index].Sub(&k, &r[index])
	r[index].Mod(&r[index], Mod)

	signature.I = *I
	signature.C = c
	signature.R = r

	return
}

func RingSignatureVerify(message string, ringPublicKeys []ecc_math.ECPoint, signature RingSignature) bool {
	participantsLen := len(ringPublicKeys)

	hash := sha256.Sum256([]byte(message))
	messageHash := new(big.Int).SetBytes(hash[:])
	fmt.Println("Message hash: ", common.BigToHash(messageHash).Hex())

	var X = make([]ecc_math.ECPoint, participantsLen)
	var Y = make([]big.Int, participantsLen)

	// For i <- 0..n-1
	fmt.Println("\n\nFor i <- 0..n-1")
	fmt.Println("I : ", common.BigToHash(&signature.I).Hex())
	for i := 0; i < participantsLen; i++ {
		// X[i] <- c[i]*A[i]+r[i]*G
		X[i] = ecc_math.AddECPoints(ecc_math.ScalarMult(ringPublicKeys[i], signature.C[i]), ecc_math.ScalarMult(G, signature.R[i]))

		fmt.Printf("i : %d\n", i)
		fmt.Println("c[i] : ", common.BigToHash(&signature.C[i]).Hex())
		fmt.Println("PubKey [i] : ")
		fmt.Printf("X: %s || Y: %s\n", common.BigToHash(ringPublicKeys[i].X).Hex(), common.BigToHash(ringPublicKeys[i].Y).Hex())
		fmt.Println("c[i]*A[i] : ")
		fmt.Printf("X: %s || Y: %s\n",
			common.BigToHash(ecc_math.ScalarMult(ringPublicKeys[i], signature.C[i]).X).Hex(),
			common.BigToHash(ecc_math.ScalarMult(ringPublicKeys[i], signature.C[i]).Y).Hex(),
		)
		fmt.Println("R[i] : ", common.BigToHash(&signature.R[i]).Hex())
		fmt.Println("r[i]*G : ")
		fmt.Printf("X: %s || Y: %s\n",
			common.BigToHash(ecc_math.ScalarMult(G, signature.R[i]).X).Hex(),
			common.BigToHash(ecc_math.ScalarMult(G, signature.R[i]).Y).Hex(),
		)

		fmt.Println("X[i] : ")
		fmt.Printf("X: %s || Y: %s\n",
			common.BigToHash(X[i].X).Hex(),
			common.BigToHash(X[i].Y).Hex(),
		)

		hash = sha256.Sum256(append(ringPublicKeys[i].X.Bytes(), ringPublicKeys[i].Y.Bytes()...))
		pubKeyHashA := new(big.Int).SetBytes(hash[:])
		fmt.Println("PubKeyHash : ", common.BigToHash(pubKeyHashA).Hex())

		// Y[i] <- c[i]*I+r[i]*H(A[i])
		Y[i].Mul(&signature.R[i], pubKeyHashA)
		Y[i].Add(&Y[i], new(big.Int).Mul(&signature.C[i], &signature.I))
		Y[i].Mod(&Y[i], Mod)
		fmt.Println("Y[i] : ", common.BigToHash(&Y[i]).Hex())

	}

	fmt.Print("\n\n")
	left := new(big.Int)
	tempBytes := messageHash.Bytes()

	for i := 0; i < participantsLen; i++ {
		// H(M) || X[0] || Y[0] || X[1] || Y[1] || ... || X[n-1] || Y[n-1]
		tempBytes = append(tempBytes, X[i].X.Bytes()...)
		tempBytes = append(tempBytes, X[i].Y.Bytes()...)
		tempBytes = append(tempBytes, Y[i].Bytes()...)
	}

	hash = sha256.Sum256(tempBytes)
	left.SetBytes(hash[:])

	fmt.Println("\nConcated bytes hash : ", common.BigToHash(left).Hex())
	fmt.Print("\n\n")

	var sum big.Int
	// Sum[j=0..n-1, j!=i](c[j])
	fmt.Println("Sum[j=0..n-1, j!=i](c[j])")
	for i := 0; i < participantsLen; i++ {
		fmt.Printf("i : %d\n", i)

		sum.Add(&sum, &signature.C[i])
		sum.Mod(&sum, Mod)
		fmt.Println("sum[i] : ", common.BigToHash(&sum).Hex())
	}

	return left.Cmp(&sum) == 0
}
