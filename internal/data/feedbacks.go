package data

import "gitlab.com/distributed_lab/kit/pgdb"

type Feedbacks interface {
	New() Feedbacks

	Insert(feedback Feedback) error
	Delete() error
	Get() (*Feedback, error)
	Select() ([]Feedback, error)

	Count() Feedbacks
	GetTotalCount() (int64, error)

	FilterByCourses(courses ...string) Feedbacks
	Page(pageParams pgdb.OffsetPageParams) Feedbacks
}

type Feedback struct {
	Course  string `json:"course" db:"course" structs:"course"`
	Content string `json:"content" db:"content" structs:"content"`
}
