package handlers

import (
	"net/http"

	"github.com/dov-id/publisher-svc/internal/service/api/requests"
	"github.com/dov-id/publisher-svc/internal/service/api/responses"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetRequest(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetReqRequest(r)
	if err != nil {
		Log(r).WithError(err).Error("failed to parse get request request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if request.RequestId == nil {
		Log(r).WithError(err).Error("request id is empty")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	req, err := RequestsQ(r).FilterByIds(*request.RequestId).Get()
	if err != nil {
		Log(r).WithError(err).Error("failed to get request")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if req == nil {
		Log(r).Errorf("request not found")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	ape.Render(w, responses.NewRequestResponse(*req))
	return
}
