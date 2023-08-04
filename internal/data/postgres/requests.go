package postgres

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/dov-id/publisher-svc/internal/data"
	"github.com/fatih/structs"
	pkgErrors "github.com/pkg/errors"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const (
	requestsTableName = "requests"
	requestsIdColumn  = requestsTableName + ".id"
)

type RequestsQ struct {
	db            *pgdb.DB
	selectBuilder sq.SelectBuilder
	updateBuilder sq.UpdateBuilder
	deleteBuilder sq.DeleteBuilder
}

func NewRequestsQ(db *pgdb.DB) data.Requests {
	return &RequestsQ{
		db:            db,
		selectBuilder: sq.Select("*").From(requestsTableName),
		updateBuilder: sq.Update(requestsTableName),
		deleteBuilder: sq.Delete(requestsTableName),
	}
}

func (r RequestsQ) New() data.Requests {
	return NewRequestsQ(r.db.Clone())
}

func (r RequestsQ) Get() (*data.Request, error) {
	var result data.Request
	err := r.db.Get(&result, r.selectBuilder)

	if pkgErrors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return &result, err
}

func (r RequestsQ) Select() ([]data.Request, error) {
	var result []data.Request

	err := r.db.Select(&result, r.selectBuilder)

	return result, err
}

func (r RequestsQ) Insert(request data.Request) (data.Request, error) {
	var result data.Request

	err := r.db.Get(
		&result,
		sq.Insert(requestsTableName).
			SetMap(structs.Map(request)).
			Suffix("RETURNING *"),
	)

	return result, err
}

func (r RequestsQ) Update(request data.Request) error {
	r.updateBuilder = r.updateBuilder.
		SetMap(structs.Map(request)).
		Where(sq.Eq{requestsIdColumn: request.Id})

	return r.db.Exec(r.updateBuilder)
}

func (r RequestsQ) Delete() error {
	var deleted []data.Request

	err := r.db.Select(&deleted, r.deleteBuilder.Suffix("RETURNING *"))
	if err != nil {
		return err
	}

	if len(deleted) == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r RequestsQ) FilterByIds(ids ...string) data.Requests {
	equalIds := sq.Eq{requestsIdColumn: ids}

	r.selectBuilder = r.selectBuilder.Where(equalIds)
	r.updateBuilder = r.updateBuilder.Where(equalIds)
	r.deleteBuilder = r.deleteBuilder.Where(equalIds)

	return r
}
