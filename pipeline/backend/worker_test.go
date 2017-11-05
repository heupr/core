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
	i := []*github.Issue{
		&github.Issue{
			ID:     github.Int(2187),
			Number: github.Int(1),
			User: &github.User{
				Login: github.String("luke"),
			},
			URL: github.String("fake-url"),
			Repository: &github.Repository{
				FullName: &fullname,
			},
			CreatedAt: &created,
		},
	}

	bs.Repos.Actives[repoID].Hive.Blender.Conflator.SetIssueRequests(i)
	bs.Repos.Actives[repoID].Limit = time.Now() //.AddDate(0, 0, -1)

	work := &RepoData{
		RepoID: repoID,
		Open: []*github.Issue{
			&github.Issue{
				ID:     github.Int(2188),
				Number: github.Int(2),
				User: &github.User{
					Login: github.String("luke"),
				},
				URL:       github.String("fake-url"),
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
