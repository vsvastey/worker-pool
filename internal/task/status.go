package task

const (
	INPROGRESS_STATE = "in progress"
	DONE_STATE       = "done"
)

type Status struct {
	State    string
	Progress int
}
