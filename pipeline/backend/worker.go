package backend

import (
	"core/utils"
	"go.uber.org/zap"
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
				if w.Repos.Actives[repodata.RepoID] == nil {
					utils.AppLog.Error("repo not initialized before worker start", zap.Int("RepoID", repodata.RepoID))
					continue
				}

				w.Repos.RLock()
				repo := w.Repos.Actives[repodata.RepoID]
				w.Repos.RUnlock()

				repo.Lock()
				if len(repodata.Open) != 0 {
					repo.Hive.Blender.Conflator.SetIssueRequests(repodata.Open)
					issues := repo.Hive.Blender.Conflator.Context.Issues
					utils.AppLog.Info("Events", zap.Int("Open", len(repodata.Open)), zap.Int("Total", len(issues)), zap.Int("RepoID", repodata.RepoID))
				}
				if len(repodata.Closed) != 0 {
					repo.Hive.Blender.Conflator.SetIssueRequests(repodata.Closed)
					issues := repo.Hive.Blender.Conflator.Context.Issues
					utils.AppLog.Info("Events", zap.Int("Closed", len(repodata.Closed)), zap.Int("Total", len(issues)), zap.Int("RepoID", repodata.RepoID))
				}
				if len(repodata.Pulls) != 0 {
					repo.Hive.Blender.Conflator.SetPullRequests(repodata.Pulls)
					issues := repo.Hive.Blender.Conflator.Context.Issues
					utils.AppLog.Info("Events", zap.Int("Pulls", len(repodata.Pulls)), zap.Int("Total", len(issues)), zap.Int("RepoID", repodata.RepoID))
				}
				utils.AppLog.Info("Conflator.Conflate() ", zap.Int("RepoID", repodata.RepoID))
				repo.Hive.Blender.Conflator.Conflate()

				utils.AppLog.Info("Blender.TrainModels() ", zap.Int("RepoID", repodata.RepoID))
				repo.Hive.Blender.TrainModels()

				repo.AssigneeAllocations = repodata.AssigneeAllocations
				repo.EligibleAssignees = repodata.EligibleAssignees

				utils.AppLog.Info("TriageOpenIssues() - Begin ", zap.Int("RepoID", repodata.RepoID))
				repo.TriageOpenIssues()
				utils.AppLog.Info("TriageOpenIssues() - Complete ", zap.Int("RepoID", repodata.RepoID))
				repo.Unlock()
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
