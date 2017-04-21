package dispatcher

import (
	// "coralreefci/engine/gateway/conflation"
	"coralreefci/engine/onboarder"
	"coralreefci/engine/onboarder/retriever"
	// "coralreefci/models"
	// "coralreefci/models/bhattacharya"
)

type Worker struct {
	ID    int
	Work  chan *retriever.RepoData
	Queue chan chan *retriever.RepoData
	Repos map[int]*onboarder.ArchRepo
	Quit  chan bool
}

func NewWorker(id int, queue chan chan *retriever.RepoData) Worker {
	return Worker{
		ID:    id,
		Work:  make(chan *retriever.RepoData),
		Queue: queue,
		Repos: make(map[int]*onboarder.ArchRepo),
		Quit:  make(chan bool),
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			w.Queue <- w.Work
			select {
			case repodata := <-w.Work:
				if len(repodata.Open) != 0 {
					w.Repos[repodata.RepoID].Hive.Blender.Models[0].Conflator.SetIssueRequests(repodata.Open)
				}
				if len(repodata.Closed) != 0 {
					w.Repos[repodata.RepoID].Hive.Blender.Models[0].Conflator.SetIssueRequests(repodata.Closed)
				}
				if len(repodata.Pulls) != 0 {
					w.Repos[repodata.RepoID].Hive.Blender.Models[0].Conflator.SetPullRequests(repodata.Pulls)
				}
				w.Repos[repodata.RepoID].Hive.Blender.Models[0].Conflator.Conflate()
				// TODO: Implement learn method calls.
			case <-w.Quit:
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	go func() {
		w.Quit <- true
	}()
}
