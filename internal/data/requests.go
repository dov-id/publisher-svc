package data

type RequestsStatus string

const (
	RequestsStatusInProgress RequestsStatus = "in progress"
	RequestsStatusPending    RequestsStatus = "pending"
	RequestsStatusFailed     RequestsStatus = "failed"
	RequestsStatusSuccess    RequestsStatus = "success"
)

type Requests interface {
	New() Requests

	Insert(request Request) (Request, error)
	Update(request Request) error
	Delete() error
	Get() (*Request, error)
	Select() ([]Request, error)

	FilterByIds(ids ...string) Requests
}

type Request struct {
	Id     string         `json:"id" db:"id" structs:"id"`
	Status RequestsStatus `json:"status" db:"status" structs:"status"`
	Error  string         `json:"error" db:"error" structs:"error"`
}
