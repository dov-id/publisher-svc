package handlers

import (
	"context"
	"math/big"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dov-id/publisher-svc/internal/config"
	"github.com/dov-id/publisher-svc/internal/data"
	"github.com/dov-id/publisher-svc/internal/helpers"
	"github.com/dov-id/publisher-svc/internal/service/api/requests"
	"github.com/dov-id/publisher-svc/internal/service/api/responses"
	"github.com/dov-id/publisher-svc/resources"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/iden3/go-rapidsnark/types"
	"github.com/iden3/go-rapidsnark/verifier"
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

	err = verifyCircomProof(request.Data.Attributes.ZkProof)
	if err != nil {
		Log(r).WithError(err).Debugf("failed to verify circom proof")
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
		request:   dbRequest,
		requestsQ: RequestsQ(r),
		logger:    Log(r),
		wallet:    Cfg(r).Networks().Networks[data.Network(request.Data.Attributes.Network)].WalletCfg,
	})
	if err != nil {
		Log(r).WithError(err).Debugf("failed to send feedback to contract")

		dbRequest.Error = err.Error()
		dbRequest.Status = data.RequestsStatusFailed
		err = RequestsQ(r).Update(dbRequest)
		if err != nil {
			err = errors.Wrap(err, "failed to update request status")
		}

		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusAccepted)
	ape.Render(w, responses.NewRequestResponse(dbRequest))
	return
}

func verifyCircomProof(zkProof resources.ZkProof) error {
	verificationKey, err := getVerificationKey()
	if err != nil {
		return errors.Wrap(err, "failed to get verification key")
	}

	err = verifier.VerifyGroth16(*convertZKProof(zkProof), verificationKey)
	if err != nil {
		return errors.Wrap(err, "failed to verify groth16 proof")
	}

	return nil
}

func getVerificationKey() ([]byte, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get current directory path")
	}

	keyFile := filepath.Join(currentDir, "verification_key.json")

	verificationKeyBytes, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read verification key file")
	}

	return verificationKeyBytes, nil
}

func convertZKProof(resProof resources.ZkProof) *types.ZKProof {
	proof := new(types.ZKProof)

	proof.PubSignals = resProof.PubSignals
	proof.Proof = &types.ProofData{
		A:        resProof.Proof.PiA,
		B:        resProof.Proof.PiB,
		C:        resProof.Proof.PiC,
		Protocol: resProof.Proof.Protocol,
	}

	return proof
}

type addFeedbackParams struct {
	network   data.Network
	request   data.Request
	requestsQ data.Requests
	logger    *logan.Entry
	wallet    *config.WalletCfg

	auth        *bind.TransactOpts
	course      common.Address
	ipfsHash    string
	pA_         [2]*big.Int
	pB_         [2][2]*big.Int
	pC_         [2]*big.Int
	pubSignals_ [11]*big.Int
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

func prepareAddFeedbackParams(
	req requests.AddFeedbackRequest,
	params *addFeedbackParams,
) error {
	var ok bool

	for i := 0; i < 2; i++ {
		params.pA_[i], ok = new(big.Int).SetString(req.Data.Attributes.ZkProof.Proof.PiA[i], 10)
		if !ok {
			return data.ErrFailedToSetString
		}

		j := 0
		l := 1
		for ; j < 2; j++ {
			params.pB_[i][j], ok = new(big.Int).SetString(req.Data.Attributes.ZkProof.Proof.PiB[i][l], 10)
			if !ok {
				return data.ErrFailedToSetString
			}
			l--
		}

		params.pC_[i], ok = new(big.Int).SetString(req.Data.Attributes.ZkProof.Proof.PiC[i], 10)
		if !ok {
			return data.ErrFailedToSetString
		}
	}

	for i := 0; i < 11; i++ {
		params.pubSignals_[i], ok = new(big.Int).SetString(req.Data.Attributes.ZkProof.PubSignals[i], 10)
		if !ok {
			return data.ErrFailedToSetString
		}
	}
	
	params.course = common.HexToAddress(req.Data.Attributes.Course)

	params.ipfsHash = req.Data.Attributes.Feedback

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
		params.ipfsHash,
		params.pA_,
		params.pB_,
		params.pC_,
		params.pubSignals_,
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
