package backend

import (
	"net/url"
	"testing"
	"time"

	"github.com/google/go-github/github"

	"core/pipeline/gateway/conflation"
)

func TestWorker(t *testing.T) {
	repoID := 23

	client := github.NewClient(nil)
	url, _ := url.Parse("http://localhost:8000/")
	client.BaseURL = url
	client.UploadURL = url

	bs := new(BackendServer)
	bs.Repos = new(ActiveRepos)
	bs.Repos.Actives = make(map[int]*ArchRepo)

	bs.Repos.Actives[repoID] = &ArchRepo{
		Client: client,
		Hive: &ArchHive{
			Blender: &Blender{
				Conflator: &conflation.Conflator{
					Context: &conflation.Context{},
					Normalizer: conflation.Normalizer{
						Context: &conflation.Context{},
					},
				},
			},
		},
	}

	fullname := "skywalker/t-16"
	created := time.Now()
	n1 := 1
	i := []*github.Issue{
		&github.Issue{
			Number: &n1,
			Repository: &github.Repository{
				FullName: &fullname,
			},
			CreatedAt: &created,
		},
	}

	bs.Repos.Actives[repoID].Hive.Blender.Conflator.SetIssueRequests(i)
	bs.Repos.Actives[repoID].Limit = time.Now() //.AddDate(0, 0, -1)

	issueID := 2187
	n2 := 2
	work := &RepoData{
		RepoID: repoID,
		Open: []*github.Issue{
			&github.Issue{
				ID:        &issueID,
				Number:    &n2,
				CreatedAt: &created,
				Repository: &github.Repository{
					FullName: &fullname,
				},
			},
		},
	}

	workerID := 1138
	ch := make(chan chan *RepoData, 2)

	worker := bs.NewWorker(workerID, ch)
	if worker.ID != workerID {
		t.Error("failure creating worker object properly")
	}

	worker.Start()
	worker.Work <- work

	for {
		if len(worker.Work) == 0 {
			worker.Stop()
			time.Sleep(1 * time.Second)
			// NOTE: Using Sleep to allow the Stop method / Quit selection to complete in the unit test (avoiding race).
			break
		}
	}
}
