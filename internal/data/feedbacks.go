package data

type Feedbacks interface {
	New() Feedbacks

	Insert(feedback Feedback) error
	Delete() error
	Get() (*Feedback, error)
	Select() ([]Feedback, error)

	FilterByCourses(courses ...string) Feedbacks
}

type Feedback struct {
	Course  string `json:"course" db:"course" structs:"course"`
	Content string `json:"content" db:"content" structs:"content"`
}
