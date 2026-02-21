package types

type Status string

const (
	Pending   Status = "Pending"
	Completed Status = "Completed"
)

type Task struct {
	ID         int
	Desc       string
	CurrStatus Status
}
