package ecc_math

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

var (
	Curve  = secp256k1.S256()
	HPoint = ScalarMult(BasePointGGet(), *big.NewInt(2))
) // it's recommend to change hPoint value (generate it via trusted setup procedure or via Nothing-In-My-Sleeve)

type ECPoint struct {
	X *big.Int
	Y *big.Int
}

func BasePointGGet() (point ECPoint) {
	point.X = Curve.Params().Gx
	point.Y = Curve.Params().Gy
	return
} // G-generator receiving

func BasePointHGet() (point ECPoint) {
	point.X = HPoint.X
	point.Y = HPoint.Y
	return
} // H-generator receiving

func ECPointGen(x, y *big.Int) (point ECPoint) {
	point.X = x
	point.Y = y
	return
} // ECPoint creation with pre-defined parameters

// Operations

func IsOnCurveCheck(a ECPoint) (c bool) {
	c = Curve.IsOnCurve(a.X, a.Y)
	return
} // P ∈ CURVE? 	- works fine, take your hands off, tested with reference vectors

func AddECPoints(a, b ECPoint) (c ECPoint) {
	c.X, c.Y = Curve.Add(a.X, a.Y, b.X, b.Y)
	return
} // P + Q 		- works fine, take your hands off, tested with reference vectors

func DoubleECPoints(a ECPoint) (c ECPoint) {
	c.X, c.Y = Curve.Double(a.X, a.Y)
	return
} // 2P 			- works fine, take your hands off, tested with reference vectors

func ScalarMult(a ECPoint, k big.Int) (c ECPoint) {
	c.X, c.Y = Curve.ScalarMult(a.X, a.Y, k.Bytes())
	return
} // k * P 		- works fine, take your hands off, tested with reference vectors

// ToString & Print

func ECPointToString(point ECPoint) (s string) {
	s = fmt.Sprintf("%X", point.X) + " " + fmt.Sprintf("%X", point.Y)
	return
}

func PrintECPoint(point ECPoint) {
	fmt.Println("X:\t", fmt.Sprintf("%X", point.X), "\nY:\t", fmt.Sprintf("%X", point.Y))
}
