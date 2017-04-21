package dispatcher

import (
	"coralreefci/engine/onboarder"
	"coralreefci/engine/onboarder/retriever"
)

var Workers chan chan *retriever.RepoData

type Dispatcher struct {
	Repos map[int]*onboarder.ArchRepo
}

func (d *Dispatcher) Start(count int) {
	Workers = make(chan chan *retriever.RepoData, count)
	for i := 0; i < count; i++ {
		worker := NewWorker(i+1, Workers)
		worker.Repos = d.Repos // TODO: Move into NewWorker function as argument.
		worker.Start()
	}

	go func() {
		for {
			work := <-Workload
			go func() {
				workers := <-Workers
				workers <- work
			}()
		}
	}()
}
