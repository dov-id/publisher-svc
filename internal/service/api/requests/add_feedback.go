package requests

import (
	"encoding/json"
	"net/http"

	"github.com/dov-id/publisher-svc/internal/data"
	"github.com/dov-id/publisher-svc/resources"
	"github.com/ethereum/go-ethereum/common"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/ipfs/go-cid"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type AddFeedbackRequest resources.AddFeedbackRequest

func NewAddFeedbackRequest(r *http.Request) (AddFeedbackRequest, error) {
	var request AddFeedbackRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, request.validate()
}

func (r *AddFeedbackRequest) validate() error {
	return validation.Errors{
		"course": validation.Validate(
			&r.Data.Attributes.Course, validation.Required, validation.By(MustBeValidEthAddress),
		),
		"network": validation.Validate(
			&r.Data.Attributes.Network,
			validation.In(data.EthereumNetwork.String(), data.PolygonNetwork.String(), data.QNetwork.String()),
		),
		"feedback": validation.Validate(
			&r.Data.Attributes.Feedback, validation.Required, validation.By(MustBeValidCID),
		),
		"zk_proof": validation.Validate(&r.Data.Attributes.ZkProof, validation.Required),
	}.Filter()
}

func MustBeValidEthAddress(src interface{}) error {
	raw, ok := src.(*string)
	if !ok {
		return data.ErrNotString
	}
	if !common.IsHexAddress(*raw) {
		return data.ErrInvalidEthAddress
	}

	return nil
}

func MustBeValidCID(src interface{}) error {
	raw, ok := src.(*string)
	if !ok {
		return data.ErrNotString
	}

	_, err := cid.Decode(*raw)
	if err != nil {
		return err
	}

	return nil
}
