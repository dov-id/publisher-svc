package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/urlval"
)

type GetReqRequest struct {
	RequestId *string `filter:"request_id"`
}

func NewGetReqRequest(r *http.Request) (GetReqRequest, error) {
	var request GetReqRequest

	err := urlval.Decode(r.URL.Query(), &request)

	return request, err
}
