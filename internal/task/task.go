package task

type Status struct {
	Progress int
}

type Task interface {
	Do() <-chan Status
	Name() string
}
