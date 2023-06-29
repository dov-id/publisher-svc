package handlers

import (
	"fmt"
	"net/http"

	"github.com/dov-id/publisher-svc/crypto_master/secp256k1/ecc_math"
	"github.com/dov-id/publisher-svc/crypto_master/secp256k1/signatures/ring_sha256"
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

	ECPoints, err := ConvertHexAddressesToECPoints(request.Data.Attributes.PublicKeys)
	if err != nil {
		Log(r).WithError(err).Errorf("failed to convert hex public keys to EC points")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	privateKeyHash := common.HexToHash(request.Data.Attributes.PrivateKey)

	newEC := new([5]ecc_math.ECPoint)
	for i, point := range ECPoints {
		newEC[i] = point
	}

	signature := ring_sha256.DynamicSizeRingSignatureGen(request.Data.Attributes.Message, ECPoints, request.Data.Attributes.Index, *privateKeyHash.Big())
	signature5 := ring_sha256.RingSignatureGen(request.Data.Attributes.Message, *newEC, 1, *privateKeyHash.Big())

	fmt.Println(ring_sha256.DynamicSizeRingSignatureVerify(request.Data.Attributes.Message, ECPoints, signature))
	fmt.Println(ring_sha256.RingSignatureVerify(request.Data.Attributes.Message, *newEC, signature5))

	sigResp := models.NewRingSignatureResponse(signature)

	newSig := newDynamicSizeRingSignature(sigResp.Data.Attributes.I, sigResp.Data.Attributes.C, sigResp.Data.Attributes.R)

	fmt.Println(ring_sha256.DynamicSizeRingSignatureVerify(request.Data.Attributes.Message, ECPoints, newSig))

	ape.Render(w, models.NewRingSignatureResponse(signature))
	return
}

func ConvertHexAddressesToECPoints(publicKeys []string) ([]ecc_math.ECPoint, error) {
	var ECPoints = make([]ecc_math.ECPoint, len(publicKeys))
	for i, key := range publicKeys {
		if key[:2] == "0x" {
			key = key[2:]
		}

		publicECDSA, err := crypto.HexToECDSA(key)
		if err != nil {
			return nil, errors.Wrap(err, "failed to convert hex to ecdsa")
		}

		ECPoints[i] = ecc_math.ECPoint{X: publicECDSA.X, Y: publicECDSA.Y}
	}
	return ECPoints, nil
}
