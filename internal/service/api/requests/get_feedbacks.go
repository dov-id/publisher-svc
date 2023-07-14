package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
)

type GetFeedbacksRequest struct {
	pgdb.OffsetPageParams

	Course *string `filter:"course"`
}

func NewGetFeedbacksRequest(r *http.Request) (GetFeedbacksRequest, error) {
	var request GetFeedbacksRequest

	err := urlval.Decode(r.URL.Query(), &request)

	return request, err
}
