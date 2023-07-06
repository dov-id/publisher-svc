package requests

import (
	"encoding/json"
	"net/http"

	"github.com/dov-id/publisher-svc/internal/data"
	"github.com/dov-id/publisher-svc/resources"
	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type AddFeedbackRequest struct {
	Data resources.AddFeedback `json:"data"`
}

func NewAddFeedbackRequest(r *http.Request) (AddFeedbackRequest, error) {
	var request AddFeedbackRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, request.validate()
}

func (r *AddFeedbackRequest) validate() error {
	return validation.Errors{
		"course":       validation.Validate(&r.Data.Attributes.Course, validation.Required),
		"network":      validation.Validate(&r.Data.Attributes.Network, validation.In(data.EthereumNetwork, data.PolygonNetwork, data.QNetwork)),
		"feedback":     validation.Validate(&r.Data.Attributes.Feedback, validation.Required),
		"signature":    validation.Validate(&r.Data.Attributes.Signature, validation.Required),
		"public_keys":  validation.Validate(&r.Data.Attributes.PublicKeys, validation.Required),
		"merkle_proof": validation.Validate(&r.Data.Attributes.Proofs, validation.Required),
	}.Filter()
}
