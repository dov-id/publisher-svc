package handlers

import (
	"net/http"

	"github.com/dov-id/publisher-svc/internal/data"
	"github.com/dov-id/publisher-svc/internal/service/api/models"
	"github.com/dov-id/publisher-svc/internal/service/api/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetFeedbacks(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetFeedbacksRequest(r)
	if err != nil {
		Log(r).WithError(err).Error("failed to parse get feedbacks request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	totalCount, err := FeedbacksQ(r).Count().FilterByCourses().GetTotalCount()
	if err != nil {
		Log(r).WithError(err).Errorf("failed to get total count from db")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	feedbacks, err := FeedbacksQ(r).FilterByCourses().Page(request.OffsetPageParams).Select()
	if err != nil {
		Log(r).WithError(err).Errorf("failed to select feedbacks from db")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	response := models.NewFeedbackListResponse(feedbacks, int64(request.OffsetPageParams.PageNumber*request.OffsetPageParams.Limit))
	response.Meta.TotalCount = totalCount
	response.Links = data.GetOffsetLinksForPGParams(r, request.OffsetPageParams)
	ape.Render(w, response)
}
