package ingestor

var Workers chan chan interface{}

type Dispatcher struct {
	Database        DataAccess
	RepoInitializer *RepoInitializer
}

func (d *Dispatcher) Start(count int) {
	Workers = make(chan chan interface{}, count)
	for i := 0; i < count; i++ {
		worker := NewWorker(i+1, d.Database, d.RepoInitializer, Workers)
		worker.Start()
	}

	go func() {
		for {
			work := <-Workload
			go func() {
				worker := <-Workers
				worker <- work
			}()
		}
	}()
}
