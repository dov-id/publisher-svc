package handlers

import (
	"context"
	"math/big"
	"net/http"

	"github.com/dov-id/publisher-svc/crypto_master/secp256k1/ecc_math"
	"github.com/dov-id/publisher-svc/crypto_master/secp256k1/signatures/ring_sha256"
	"github.com/dov-id/publisher-svc/internal/config"
	"github.com/dov-id/publisher-svc/internal/data"
	"github.com/dov-id/publisher-svc/internal/helpers"
	"github.com/dov-id/publisher-svc/internal/service/api/requests"
	"github.com/dov-id/publisher-svc/internal/service/api/responses"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	pkgErrors "github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func AddFeedback(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewAddFeedbackRequest(r)
	if err != nil {
		Log(r).WithError(err).Debug("failed to parse add feedback request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	ECPoints, err := ConvertHexKeysToECPoints(request.Data.Attributes.PublicKeys)
	if err != nil {
		Log(r).WithError(err).Debugf("failed to convert hex public keys to EC points")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	signature, err := newDynamicSizeRingSignature(request.Data.Attributes.Signature.I, request.Data.Attributes.Signature.C, request.Data.Attributes.Signature.R)
	if err != nil {
		Log(r).WithError(err).Debugf("failed to convert hex signature to big int values")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	isVerifiedSig := ring_sha256.RingSignatureVerify(request.Data.Attributes.Feedback, ECPoints, *signature)
	if !isVerifiedSig {
		Log(r).Debugf("signature is not verified")
		ape.RenderErr(w, problems.NotAllowed())
		return
	}

	dbRequest, err := RequestsQ(r).Insert(data.Request{
		Id:     uuid.NewString(),
		Status: data.RequestsStatusPending,
	})
	if err != nil {
		Log(r).WithError(err).Debugf("failed to insert new request")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	err = sendFeedbackToContract(ParentCtx(r), request, addFeedbackParams{
		request:    dbRequest,
		requestsQ:  RequestsQ(r),
		logger:     Log(r),
		signature:  *signature,
		publicKeys: ECPoints,
		wallet:     Cfg(r).Networks().Networks[data.Network(request.Data.Attributes.Network)].WalletCfg,
	})
	if err != nil {
		dbRequest.Error = err.Error()
		dbRequest.Status = data.RequestsStatusFailed
		err = RequestsQ(r).Update(dbRequest)
		if err != nil {
			err = errors.Wrap(err, "failed to update request status")
		}

		Log(r).WithError(err).Debugf("failed to send feedback to contract")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusAccepted)
	ape.Render(w, responses.NewRequestResponse(dbRequest))
	return
}

type addFeedbackParams struct {
	network   data.Network
	request   data.Request
	requestsQ data.Requests
	logger    *logan.Entry
	wallet    *config.WalletCfg

	signature  ring_sha256.RingSignature
	publicKeys []ecc_math.ECPoint

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

func sendFeedbackToContract(ctx context.Context, req requests.AddFeedbackRequest, params addFeedbackParams) error {
	network := data.Network(req.Data.Attributes.Network)

	clients, err := helpers.GetNetworkClientsFromCtx(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get network clients from context")
	}

	auth, err := helpers.GetAuth(ctx, clients[network], params.wallet)
	if err != nil {
		return errors.Wrap(err, "failed to get auth options")
	}

	err = prepareAddFeedbackParams(req, &params)
	if err != nil {
		return errors.Wrap(err, "failed to prepare add feedback params")
	}

	params.network = network
	params.auth = auth

	err = makeAddFeedbackTx(ctx, params)
	if err != nil {
		return errors.Wrap(err, "failed to make add feedback transaction")
	}

	return nil
}

func makeAddFeedbackTx(ctx context.Context, params addFeedbackParams) error {
	feedbackRegistries, err := helpers.GetFeedbackRegistriesFromCtx(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get feedback registries from context")
	}

	clients, err := helpers.GetNetworkClientsFromCtx(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get network clients from context")
	}

	transaction, err := feedbackRegistries[params.network].AddFeedback(
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
		if pkgErrors.Is(err, data.ErrReplacementTxUnderpriced) {
			params.auth.Nonce = big.NewInt(params.auth.Nonce.Int64() + 1)
			return makeAddFeedbackTx(ctx, params)
		}

		return errors.Wrap(err, "failed to add feedback")
	}

	params.request.Status = data.RequestsStatusInProgress
	err = params.requestsQ.Update(params.request)
	if err != nil {
		return errors.Wrap(err, "failed to update request status")
	}

	helpers.WaitForTransactionMined(
		ctx,
		clients[params.network],
		transaction,
		params.logger,
		params.request,
		params.requestsQ,
	)

	return nil
}

func newDynamicSizeRingSignature(i string, cHexArr []string, rHexArr []string) (*ring_sha256.RingSignature, error) {
	var signature ring_sha256.RingSignature

	arraysLen := len(cHexArr)

	tempBigInt, ok := new(big.Int).SetString(i, 10)
	if !ok {
		return nil, data.ErrFailedToSetString
	}
	signature.I = *tempBigInt

	c := make([]big.Int, arraysLen)
	r := make([]big.Int, arraysLen)

	for i := 0; i < arraysLen; i++ {
		tempBigInt, ok = new(big.Int).SetString(cHexArr[i], 10)
		if !ok {
			return nil, data.ErrFailedToSetString
		}
		c[i] = *tempBigInt

		tempBigInt, ok = new(big.Int).SetString(rHexArr[i], 10)
		if !ok {
			return nil, data.ErrFailedToSetString
		}
		r[i] = *tempBigInt
	}

	signature.R = r
	signature.C = c

	return &signature, nil
}

func prepareAddFeedbackParams(
	req requests.AddFeedbackRequest,
	params *addFeedbackParams,
) error {
	var ok bool

	params.course = common.HexToAddress(req.Data.Attributes.Course)

	params.i, ok = new(big.Int).SetString(req.Data.Attributes.Signature.I, 10)
	if !ok {
		return data.ErrFailedToSetString
	}

	err := prepareRingSigParams(params)
	if err != nil {
		return errors.Wrap(err, "failed to prepare ring signature params")
	}

	err = prepareMTProofs(req, params)
	if err != nil {
		return errors.Wrap(err, "failed to prepare merkle tree proofs")
	}

	params.ipfsHash = req.Data.Attributes.Feedback

	return nil
}

func prepareRingSigParams(params *addFeedbackParams) error {
	length := len(params.signature.R)

	params.r = make([]*big.Int, length)
	params.c = make([]*big.Int, length)
	params.publicKeysX = make([]*big.Int, length)
	params.publicKeysY = make([]*big.Int, length)

	params.i = &params.signature.I

	for i := 0; i < length; i++ {
		params.r[i] = &params.signature.R[i]
		params.c[i] = &params.signature.C[i]

		params.publicKeysX[i] = params.publicKeys[i].X
		params.publicKeysY[i] = params.publicKeys[i].Y
	}

	return nil
}

func prepareMTProofs(req requests.AddFeedbackRequest, params *addFeedbackParams) error {
	for _, proof := range req.Data.Attributes.Proofs {
		var tmp [32]byte

		decoded, err := hexutil.Decode(proof.NodeKey)
		if err != nil {
			return errors.Wrap(err, "failed to decode proof node key")
		}

		copy(tmp[:], decoded[:])
		params.keys = append(params.keys, tmp)

		decoded, err = hexutil.Decode(proof.NodeValue)
		if err != nil {
			return errors.Wrap(err, "failed to decode proof node value")
		}

		copy(tmp[:], decoded[:])
		params.values = append(params.values, tmp)

		var mtp [][32]byte
		for _, element := range proof.Proof {
			decoded, err = hexutil.Decode(element)
			if err != nil {
				return errors.Wrap(err, "failed to decode proof element")
			}

			copy(tmp[:], decoded[:])
			mtp = append(mtp, tmp)
		}

		params.merkleTreeProofs = append(params.merkleTreeProofs, mtp)
	}

	return nil
}

func ConvertHexKeysToECPoints(publicKeys []string) ([]ecc_math.ECPoint, error) {
	var ecPoints = make([]ecc_math.ECPoint, len(publicKeys))
	for i, key := range publicKeys {
		decodedKey, err := hexutil.Decode(key)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode public key")
		}

		publicECDSA, err := crypto.UnmarshalPubkey(decodedKey)
		if err != nil {
			return nil, errors.Wrap(err, "failed to convert hex to ecdsa")
		}

		ecPoints[i] = ecc_math.ECPoint{X: publicECDSA.X, Y: publicECDSA.Y}
	}
	return ecPoints, nil
}
