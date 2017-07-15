package backend

import "fmt"

type Worker struct {
	ID    int
	Work  chan *RepoData
	Queue chan chan *RepoData
	Repos *ActiveRepos
	Quit  chan bool
}

func (bs *BackendServer) NewWorker(workerID int, queue chan chan *RepoData) Worker {
	return Worker{
		ID:    workerID,
		Work:  make(chan *RepoData),
		Queue: queue,
		Repos: bs.Repos,
		Quit:  make(chan bool),
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			fmt.Println("AFTER FOR")
			w.Queue <- w.Work
			fmt.Println("QUEUED")
			select {
			case repodata := <-w.Work:
				fmt.Println("WORK")
				w.Repos.Lock()

				if w.Repos.Actives[repodata.RepoID] != nil {
					// Generate warning
					// Exit loop / redirect back to initiate ArchRepo & etc.
				}

				if len(repodata.Open) != 0 {
					w.Repos.Actives[repodata.RepoID].Hive.Blender.Conflator.SetIssueRequests(repodata.Open)
				}
				if len(repodata.Closed) != 0 {
					w.Repos.Actives[repodata.RepoID].Hive.Blender.Conflator.SetIssueRequests(repodata.Closed)
				}
				if len(repodata.Pulls) != 0 {
					w.Repos.Actives[repodata.RepoID].Hive.Blender.Conflator.SetPullRequests(repodata.Pulls)
				}
				w.Repos.Actives[repodata.RepoID].Hive.Blender.Conflator.Conflate()

				w.Repos.Actives[repodata.RepoID].Hive.Blender.TrainModels()
				w.Repos.Actives[repodata.RepoID].TriageOpenIssues()

				w.Repos.Unlock()
				fmt.Println("WORK COMPLETED")
				continue //Try with this first. If this doesn't work then remove default: and try again
			case <-w.Quit:
				fmt.Println("QUITTER")
				return
				// default:
				//     fmt.Println("FALL THROUGH")
			}
		}
		fmt.Println("AFTER GOROUTINE")
	}()
}

func (w *Worker) Stop() {
	fmt.Println("STOPPED")
	go func() {
		w.Quit <- true
	}()
}
