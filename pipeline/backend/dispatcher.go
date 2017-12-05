package backend

var Workers chan chan *RepoData

func (bs *BackendServer) Dispatcher(count int) {
	Workers = make(chan chan *RepoData, count)
	for i := 0; i < count; i++ {
		worker := bs.NewWorker(i+1, Workers)
		worker.Start()
	}

	go func() {
		for {
			work := <-workload
			go func() {
				workers := <-Workers
				workers <- work
			}()
		}
	}()
}
