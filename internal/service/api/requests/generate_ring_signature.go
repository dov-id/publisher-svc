package requests

import (
	"encoding/json"
	"net/http"

	"github.com/dov-id/publisher-svc/resources"
	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type GenerateRingSigRequest struct {
	Data resources.GenRingSig `json:"data"`
}

func NewGenerateRingSignatureRequest(r *http.Request) (GenerateRingSigRequest, error) {
	var request GenerateRingSigRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, request.validate()
}

func (r *GenerateRingSigRequest) validate() error {
	return validation.Errors{
		"message":     validation.Validate(&r.Data.Attributes.Message, validation.Required),
		"public_key":  validation.Validate(&r.Data.Attributes.PublicKeys, validation.Required),
		"index":       validation.Validate(&r.Data.Attributes.Index, validation.Required),
		"private_key": validation.Validate(&r.Data.Attributes.PrivateKey, validation.Required),
	}.Filter()
}
