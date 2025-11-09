package taskScheduler

import "time"

type Task struct {
	Action func()
}
type Scheduler struct {
	Tasks []Task
}

func NewTaskScheduler(Tasks []Task) *Scheduler {
	return &Scheduler{Tasks: Tasks}
}

func NewTask(f func()) Task {
	return Task{f}
}

func (s *Scheduler) Run() {
	go func() {
		t := time.NewTicker(time.Hour * 24)
		for {
			select {
			case <-t.C:
				for _, task := range s.Tasks {
					task.Action()
				}
			}
		}
	}()
}
