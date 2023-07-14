package postgres

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/dov-id/publisher-svc/internal/data"
	"github.com/fatih/structs"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const (
	feedbacksTableName    = "feedbacks"
	feedbacksCourseColumn = feedbacksTableName + ".course"
)

type FeedbacksQ struct {
	db            *pgdb.DB
	selectBuilder sq.SelectBuilder
	updateBuilder sq.UpdateBuilder
	deleteBuilder sq.DeleteBuilder
}

func NewFeedbacksQ(db *pgdb.DB) data.Feedbacks {
	return &FeedbacksQ{
		db:            db,
		selectBuilder: sq.Select("*").From(feedbacksTableName),
		updateBuilder: sq.Update(feedbacksTableName),
		deleteBuilder: sq.Delete(feedbacksTableName),
	}
}

func (q FeedbacksQ) New() data.Feedbacks {
	return NewFeedbacksQ(q.db)
}

func (q FeedbacksQ) Get() (*data.Feedback, error) {
	var result data.Feedback
	err := q.db.Get(&result, q.selectBuilder)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

func (q FeedbacksQ) Select() ([]data.Feedback, error) {
	var result []data.Feedback

	err := q.db.Select(&result, q.selectBuilder)

	return result, err
}

func (q FeedbacksQ) Insert(user data.Feedback) error {
	query := sq.Insert(feedbacksTableName).
		SetMap(structs.Map(user)).
		Suffix("ON CONFLICT (course, content) DO NOTHING")

	return q.db.Exec(query)
}

func (q FeedbacksQ) Delete() error {
	var deleted []data.Feedback

	err := q.db.Select(&deleted, q.deleteBuilder.Suffix("RETURNING *"))
	if err != nil {
		return err
	}

	if len(deleted) == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (q FeedbacksQ) FilterByCourses(courses ...string) data.Feedbacks {
	equalCourses := sq.Eq{feedbacksCourseColumn: courses}

	q.selectBuilder = q.selectBuilder.Where(equalCourses)
	q.updateBuilder = q.updateBuilder.Where(equalCourses)
	q.deleteBuilder = q.deleteBuilder.Where(equalCourses)

	return q
}

func (q FeedbacksQ) Count() data.Feedbacks {
	q.selectBuilder = sq.Select("COUNT (*)").From(feedbacksTableName)

	return q
}

func (q FeedbacksQ) GetTotalCount() (int64, error) {
	var count int64
	err := q.db.Get(&count, q.selectBuilder)

	return count, err
}

func (q FeedbacksQ) Page(pageParams pgdb.OffsetPageParams) data.Feedbacks {
	q.selectBuilder = pageParams.ApplyTo(q.selectBuilder, feedbacksCourseColumn)

	return q
}
