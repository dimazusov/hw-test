package scheduler

import "time"

type task struct {
	f func() error
	t time.Duration
}

type Task interface {
	Run() error
	Sleep()
}

func newTask(f func() error, t time.Duration) Task {
	return &task{
		f: f,
		t: t,
	}
}

func (m *task) Run() error {
	return m.f()
}

func (m *task) Sleep() {
	time.Sleep(m.t)
}
