package data

type RequestsStatus string

const (
	IN_PROGRESS RequestsStatus = "in progress"
	PENDING     RequestsStatus = "pending"
	FAILED      RequestsStatus = "failed"
	SUCCESS     RequestsStatus = "success"
)

type Requests interface {
	New() Requests

	Insert(request Request) (Request, error)
	Update(request RequestToUpdate) error
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

type RequestToUpdate struct {
	Status RequestsStatus `structs:"status"`
	Error  *string        `structs:"error,omitempty"`
}
