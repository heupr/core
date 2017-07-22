package backend

import (
	"go.uber.org/zap"

	"coralreefci/utils"
)

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
			w.Queue <- w.Work
			select {
			case repodata := <-w.Work:
				w.Repos.Lock()

				if w.Repos.Actives[repodata.RepoID] != nil {
					utils.AppLog.Error("repo not initialized before worker start, repo ID: ", zap.Int("repoID", repodata.RepoID))
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
				continue
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
