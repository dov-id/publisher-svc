package models

import (
	"github.com/dov-id/publisher-svc/internal/data"
	"github.com/dov-id/publisher-svc/resources"
)

func newFeedback(feedback data.Feedback, index int64) resources.Feedback {
	return resources.Feedback{
		Key: resources.NewKeyInt64(index, resources.FEEDBACK),
		Attributes: resources.FeedbackAttributes{
			Course:  feedback.Course,
			Content: feedback.Content,
		},
	}
}

func newFeedbackList(feedbacks []data.Feedback, offset int64) []resources.Feedback {
	var list = make([]resources.Feedback, 0)

	for i, feedback := range feedbacks {
		list = append(list, newFeedback(feedback, int64(i)+offset))
	}

	return list
}

func NewFeedbackListResponse(feedbacks []data.Feedback, offset int64) FeedbacksListResponse {
	return FeedbacksListResponse{
		Data: newFeedbackList(feedbacks, offset),
	}
}

type Meta struct {
	TotalCount int64 `json:"total_count"`
}

type FeedbacksListResponse struct {
	Meta  Meta                 `json:"meta"`
	Data  []resources.Feedback `json:"data"`
	Links *resources.Links     `json:"links"`
}
