package handlers

import (
	"fmt"
	"math/big"
	"net/http"

	"github.com/dov-id/publisher-svc/contracts"
	"github.com/dov-id/publisher-svc/crypto_master/secp256k1/signatures/ring_sha256"
	"github.com/dov-id/publisher-svc/internal/data"
	"github.com/dov-id/publisher-svc/internal/helpers"
	"github.com/dov-id/publisher-svc/internal/service/api/models"
	"github.com/dov-id/publisher-svc/internal/service/api/requests"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func AddFeedback(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewAddFeedbackRequest(r)
	if err != nil {
		Log(r).WithError(err).Error("failed to parse add feedback request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	ECPoints, err := helpers.ConvertHexKeysToECPoints(request.Data.Attributes.PublicKeys)
	if err != nil {
		Log(r).WithError(err).Errorf("failed to convert hex public keys to EC points")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	signature := newDynamicSizeRingSignature(request.Data.Attributes.Signature.I, request.Data.Attributes.Signature.C, request.Data.Attributes.Signature.R)

	isVerifiedSig := ring_sha256.DynamicSizeRingSignatureVerify(request.Data.Attributes.Feedback, ECPoints, signature)
	if !isVerifiedSig {
		Log(r).WithError(err).Error("signature is not verified")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	dbRequest, err := RequestsQ(r).Insert(data.Request{
		Id:     uuid.New().String(),
		Status: data.PENDING,
		Error:  "",
	})
	if err != nil {
		Log(r).WithError(err).Errorf("failed to insert new request")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	err = sendFeedbackToContract(request.Data.Attributes.Network, r, request, dbRequest.Id)
	if err != nil {
		Log(r).WithError(err).Errorf("failed to send feedback to contract")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusAccepted)
	ape.Render(w, models.NewRequestResponse(dbRequest))
	return
}

type addFeedbackParams struct {
	requestId        string
	feedbackRegistry *contracts.FeedbackRegistry
	client           *ethclient.Client
	auth             *bind.TransactOpts
	//Set other when finish contract
}

func sendFeedbackToContract(network string, r *http.Request, req requests.AddFeedbackRequest, requestId string) error {
	client, err := helpers.CreateNetworkClient(network, Cfg(r).Networks().Networks[network], Cfg(r).Networks().Networks[data.InfuraNetwork])
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to create `%s` client", network))
	}

	auth, err := helpers.GetAuth(client, Cfg(r).Networks().Networks[data.MetamaskNetwork].Key)
	if err != nil {
		return errors.Wrap(err, "failed to get auth options")
	}

	feedbackRegistry, err := contracts.NewFeedbackRegistry(common.HexToAddress(Cfg(r).FeedbackRegistry().Addresses[network]), client)
	if err != nil {
		return errors.Wrap(err, "failed to create new feedback registry instance")
	}

	//TODO: finish contract and then make tx
	err = makeAddFeedbackTx(
		addFeedbackParams{
			requestId:        requestId,
			feedbackRegistry: feedbackRegistry,
			client:           client,
			auth:             auth,
		},
		r,
	)
	if err != nil {
		return errors.Wrap(err, "failed to make add feedback transaction")
	}

	return nil
}

func makeAddFeedbackTx(params addFeedbackParams, r *http.Request) error {
	transaction, err := params.feedbackRegistry.AddFeedback(
		params.auth,
		//rest of params
	)
	if err != nil {
		if err.Error() == data.ReplacementTxUnderpricedErr {
			params.auth.Nonce = big.NewInt(params.auth.Nonce.Int64() + 1)
			return makeAddFeedbackTx(params, r)
		}

		return errors.Wrap(err, "failed to add feedback")
	}

	err = RequestsQ(r).FilterByIds(params.requestId).Update(data.RequestToUpdate{Status: data.IN_PROGRESS})
	if err != nil {
		return errors.Wrap(err, "failed to update request status")
	}

	helpers.WaitForTransactionMined(params.client, transaction, Log(r), params.requestId, RequestsQ(r))

	return nil
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
