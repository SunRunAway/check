package check

type runner struct {
	concurrency int

	ch chan func()
	done chan struct{}
}

func (r *runner) start() {
	r.ch = make(chan func(), 10000)
	r.done = make(chan struct{})
	for i := 0; i < r.concurrency; i++ {
		go func() {
			for {
				select {
				case do := <-r.ch:
					do()
				case <-r.done:
					return
				}
			}
		}()
	}
}

func (r *runner) run(do func()) {
	r.ch <- do
}

func (r *runner) stop() {
	close(r.done)
}