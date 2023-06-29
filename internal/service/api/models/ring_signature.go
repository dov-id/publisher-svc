package models

import (
	"fmt"

	"github.com/dov-id/publisher-svc/crypto_master/secp256k1/signatures/ring_sha256"
	"github.com/dov-id/publisher-svc/resources"
	"github.com/ethereum/go-ethereum/common"
)

func newRingSignature(signature ring_sha256.DynamicSizeRingSignature) resources.RingSig {
	iHash := common.BytesToHash(signature.I.Bytes())
	arraysLen := len(signature.C)

	var cHexArr = make([]string, arraysLen)
	var rHexArr = make([]string, arraysLen)

	for i := 0; i < arraysLen; i++ {
		cHexArr[i] = fmt.Sprintf("0x%s", common.Bytes2Hex(signature.C[i].Bytes()))
		rHexArr[i] = fmt.Sprintf("0x%s", common.Bytes2Hex(signature.R[i].Bytes()))
	}

	return resources.RingSig{
		Key: resources.Key{
			ID:   iHash.Hex(),
			Type: resources.RING_SIGNATURE,
		},
		Attributes: resources.RingSigAttr{
			I: iHash.Hex(),
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
