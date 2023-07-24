package handlers

import (
	"fmt"
	"math/big"
	"net/http"

	"github.com/dov-id/publisher-svc/contracts"
	"github.com/dov-id/publisher-svc/crypto_master/secp256k1/signatures/ring_sha256"
	"github.com/dov-id/publisher-svc/internal/data"
	"github.com/dov-id/publisher-svc/internal/helpers"
	"github.com/dov-id/publisher-svc/internal/service/api/requests"
	"github.com/dov-id/publisher-svc/internal/service/api/responses"
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

	isVerifiedSig := ring_sha256.DynamicSizeRingSignatureVerifyBytes(request.Data.Attributes.Feedback, ECPoints, signature)
	if !isVerifiedSig {
		Log(r).Errorf("signature is not verified")
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
		errorMsg := err.Error()

		err = RequestsQ(r).FilterByIds(dbRequest.Id).Update(data.RequestToUpdate{Status: data.FAILED, Error: &errorMsg})
		if err != nil {
			err = errors.Wrap(err, "failed to update request status")
		}

		Log(r).WithError(err).Errorf("failed to send feedback to contract")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusAccepted)
	ape.Render(w, responses.NewRequestResponse(dbRequest))
	return
}

type addFeedbackParams struct {
	requestId        string
	feedbackRegistry *contracts.FeedbackRegistry
	client           *ethclient.Client
	auth             *bind.TransactOpts
	course           common.Address
	i                *big.Int
	c                []*big.Int
	r                []*big.Int
	publicKeysX      []*big.Int
	publicKeysY      []*big.Int
	merkleTreeProofs [][][32]byte
	keys             [][32]byte
	values           [][32]byte
	ipfsHash         string
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

	params, err := prepareAddFeedbackParams(req)
	if err != nil {
		return errors.Wrap(err, "failed to prepare add feedback params")
	}

	params.client = client
	params.auth = auth
	params.requestId = requestId
	params.feedbackRegistry = feedbackRegistry

	err = makeAddFeedbackTx(*params, r)
	if err != nil {
		return errors.Wrap(err, "failed to make add feedback transaction")
	}

	return nil
}

func makeAddFeedbackTx(params addFeedbackParams, r *http.Request) error {
	transaction, err := params.feedbackRegistry.AddFeedback(
		params.auth,
		params.course,
		params.i,
		params.c,
		params.r,
		params.publicKeysX,
		params.publicKeysY,
		params.merkleTreeProofs,
		params.keys,
		params.values,
		params.ipfsHash,
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

func newDynamicSizeRingSignature(i string, cHexArr []string, rHexArr []string) ring_sha256.DynamicSizeRingSignature {
	var signature ring_sha256.DynamicSizeRingSignature
	arraysLen := len(cHexArr)

	signature.I = *helpers.StringToBigInt(i, 10)

	c := make([]big.Int, arraysLen)
	r := make([]big.Int, arraysLen)

	for i := 0; i < arraysLen; i++ {
		c[i] = *helpers.StringToBigInt(cHexArr[i], 10)
		r[i] = *helpers.StringToBigInt(rHexArr[i], 10)
	}

	signature.R = r
	signature.C = c

	return signature
}

func prepareAddFeedbackParams(req requests.AddFeedbackRequest) (*addFeedbackParams, error) {
	var params addFeedbackParams

	params.course = common.HexToAddress(req.Data.Attributes.Course)

	params.i = helpers.StringToBigInt(req.Data.Attributes.Signature.I, 10)
	params.r = helpers.StringArrToBigIntArr(req.Data.Attributes.Signature.R, 10)
	params.c = helpers.StringArrToBigIntArr(req.Data.Attributes.Signature.C, 10)

	ECPoints, err := helpers.ConvertHexKeysToECPoints(req.Data.Attributes.PublicKeys)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert hex keys to ec points")
	}

	for _, point := range ECPoints {
		params.publicKeysX = append(params.publicKeysX, point.X)
		params.publicKeysY = append(params.publicKeysY, point.Y)
	}

	for _, proof := range req.Data.Attributes.Proofs {
		params.keys = append(params.keys, helpers.StringToByte32(proof.NodeKey))
		params.values = append(params.values, helpers.StringToByte32(proof.NodeValue))
		params.merkleTreeProofs = append(params.merkleTreeProofs, helpers.StringArrToByte32Arr(proof.Proof))
	}

	params.ipfsHash = req.Data.Attributes.Feedback

	return &params, nil
}
