package task

type Task interface {
	Do() <-chan Status
}
