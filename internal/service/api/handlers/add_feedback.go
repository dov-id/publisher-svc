package handlers

import (
	"fmt"
	"math/big"
	"net/http"

	"github.com/dov-id/publisher-svc/crypto_master/secp256k1/signatures/ring_sha256"
	"github.com/dov-id/publisher-svc/internal/service/api/requests"
	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func AddFeedback(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewAddFeedbackRequest(r)
	if err != nil {
		Log(r).WithError(err).Error("failed to parse add feedback request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	ECPoints, err := ConvertHexAddressesToECPoints(request.Data.Attributes.PublicKeys)
	if err != nil {
		Log(r).WithError(err).Errorf("failed to convert hex public keys to EC points")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	signature := newDynamicSizeRingSignature(request.Data.Attributes.Signature.I, request.Data.Attributes.Signature.C, request.Data.Attributes.Signature.R)

	isVerified := ring_sha256.DynamicSizeRingSignatureVerify(request.Data.Attributes.Feedback, ECPoints, signature)

	fmt.Println(isVerified)

	ape.Render(w, 200)
	return
}

func newDynamicSizeRingSignature(iHex string, cHexArr []string, rHexArr []string) ring_sha256.DynamicSizeRingSignature {
	var signature ring_sha256.DynamicSizeRingSignature
	arraysLen := len(cHexArr)

	signature.I = *common.HexToHash(iHex).Big()

	c := make([]big.Int, arraysLen)
	r := make([]big.Int, arraysLen)

	for i := 0; i < arraysLen; i++ {
		c[i] = *common.HexToHash(cHexArr[i]).Big()
		r[i] = *common.HexToHash(rHexArr[i]).Big()
	}

	signature.R = r
	signature.C = c

	return signature
}
