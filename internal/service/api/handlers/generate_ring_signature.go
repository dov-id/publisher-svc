package handlers

import (
	"net/http"

	"github.com/dov-id/publisher-svc/crypto_master/secp256k1/ecc_math"
	"github.com/dov-id/publisher-svc/crypto_master/secp256k1/signatures/ring_sha256"
	"github.com/dov-id/publisher-svc/internal/helpers"
	"github.com/dov-id/publisher-svc/internal/service/api/models"
	"github.com/dov-id/publisher-svc/internal/service/api/requests"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
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

	newEC := new([5]ecc_math.ECPoint)
	for i, key := range request.Data.Attributes.PublicKeys {
		publicECDSA, err := crypto.UnmarshalPubkey(common.Hex2Bytes(key))
		if err != nil {
			panic(errors.Wrap(err, "failed to convert hex to ecdsa"))
		}
		newEC[i] = ecc_math.ECPoint{X: publicECDSA.X, Y: publicECDSA.Y}
	}

	signature := ring_sha256.DynamicSizeRingSignatureGen(request.Data.Attributes.Message, ECPoints, request.Data.Attributes.Index, *privateKeyHash.Big())

	ape.Render(w, models.NewRingSignatureResponse(signature))
	return
}
