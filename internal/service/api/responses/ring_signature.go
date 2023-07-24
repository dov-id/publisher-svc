package responses

import (
	"github.com/dov-id/publisher-svc/crypto_master/secp256k1/signatures/ring_sha256"
	"github.com/dov-id/publisher-svc/resources"
)

func newRingSignature(signature ring_sha256.DynamicSizeRingSignature) resources.RingSig {
	arraysLen := len(signature.C)

	var cHexArr = make([]string, arraysLen)
	var rHexArr = make([]string, arraysLen)

	for i := 0; i < arraysLen; i++ {
		cHexArr[i] = signature.C[i].String()
		rHexArr[i] = signature.R[i].String()
	}

	return resources.RingSig{
		Key: resources.Key{
			ID:   signature.I.String(),
			Type: resources.RING_SIGNATURE,
		},
		Attributes: resources.RingSigAttr{
			I: signature.I.String(),
			C: cHexArr,
			R: rHexArr,
		},
	}
}

func NewRingSignatureResponse(signature ring_sha256.DynamicSizeRingSignature) resources.RingSigResponse {
	return resources.RingSigResponse{
		Data: newRingSignature(signature),
	}
}
