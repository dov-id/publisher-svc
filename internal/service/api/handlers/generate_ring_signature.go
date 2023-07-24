package handlers

import (
	"net/http"

	"github.com/dov-id/publisher-svc/crypto_master/secp256k1/signatures/ring_sha256"
	"github.com/dov-id/publisher-svc/internal/helpers"
	"github.com/dov-id/publisher-svc/internal/service/api/requests"
	"github.com/dov-id/publisher-svc/internal/service/api/responses"
	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GenerateRingSignature(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGenerateRingSignatureRequest(r)
	if err != nil {
		Log(r).WithError(err).Error("failed to parse generate proof request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	ECPoints, err := helpers.ConvertHexKeysToECPoints(request.Data.Attributes.PublicKeys)
	if err != nil {
		Log(r).WithError(err).Errorf("failed to convert hex public keys to EC points")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	privateKeyHash := common.HexToHash(request.Data.Attributes.PrivateKey)

	signature := ring_sha256.DynamicSizeRingSignatureGenBytes(request.Data.Attributes.Message, ECPoints, request.Data.Attributes.Index, *privateKeyHash.Big())

	ape.Render(w, responses.NewRingSignatureResponse(signature))
	return
}
