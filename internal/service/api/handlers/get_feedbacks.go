package handlers

import (
	"net/http"

	"github.com/dov-id/publisher-svc/internal/helpers"
	"github.com/dov-id/publisher-svc/internal/service/api/requests"
	"github.com/dov-id/publisher-svc/internal/service/api/responses"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetFeedbacks(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetFeedbacksRequest(r)
	if err != nil {
		Log(r).WithError(err).Debug("failed to parse get feedbacks request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	totalCountStmt := FeedbacksQ(r).Count()
	feedbacksStmt := FeedbacksQ(r)

	if request.Course != nil {
		totalCountStmt = totalCountStmt.FilterByCourses(*request.Course)
		feedbacksStmt = feedbacksStmt.FilterByCourses(*request.Course)
	}

	totalCount, err := totalCountStmt.GetTotalCount()
	if err != nil {
		Log(r).WithError(err).Debugf("failed to get total count from db")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	feedbacks, err := feedbacksStmt.Page(request.OffsetPageParams).Select()
	if err != nil {
		Log(r).WithError(err).Debugf("failed to select feedbacks from db")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	response := responses.NewFeedbackListResponse(
		feedbacks,
		int64(request.OffsetPageParams.PageNumber*request.OffsetPageParams.Limit),
	)
	response.Meta.TotalCount = totalCount
	response.Links = helpers.GetOffsetLinksForPGParams(r, request.OffsetPageParams)
	ape.Render(w, response)
}
