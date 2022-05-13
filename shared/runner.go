package shared

import "sync"

type Runner struct {
	wg sync.WaitGroup
}

func (runner *Runner) RunParallel(closure func()) {
	runner.wg.Add(1)
	go func() {
		closure()
		runner.wg.Done()
	}()
}

func (runner *Runner) Wait() {
	runner.wg.Wait()
}
