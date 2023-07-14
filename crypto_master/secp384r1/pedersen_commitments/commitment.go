package pedersen_commitments

import (
	"math/big"

	"github.com/dov-id/publisher-svc/crypto_master/secp384r1/ecc_math"
)

type PedersenCommitment struct {
	Commitment ecc_math.ECPoint
}

type Witness struct {
	vote  big.Int
	nonce big.Int
}

func WitnessGen(vote, nonce big.Int) (witness Witness) {
	witness.vote = vote
	witness.nonce = nonce
	return
} // value & nonce

func CommitmentGen(witness Witness) (commitment PedersenCommitment) {
	var point ecc_math.ECPoint = ecc_math.AddECPoints(ecc_math.ScalarMult(ecc_math.BasePointGGet(), witness.nonce), ecc_math.ScalarMult(ecc_math.BasePointHGet(), witness.vote))
	commitment.Commitment = point
	return
} // C = nonce*G + value*H

func CommitmentVerify(witness Witness, commitment PedersenCommitment) bool {
	var point ecc_math.ECPoint = ecc_math.AddECPoints(ecc_math.ScalarMult(ecc_math.BasePointGGet(), witness.nonce), ecc_math.ScalarMult(ecc_math.BasePointHGet(), witness.vote))
	if point.X.Cmp(commitment.Commitment.X) == 0 && point.Y.Cmp(commitment.Commitment.Y) == 0 {
		return true
	} else {
		return false
	}
} // if nonce*G + value*H == C -> true
