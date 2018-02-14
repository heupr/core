package backend

import (
	"context"
	"net/url"
	"testing"
	"time"

	language "cloud.google.com/go/language/apiv1"
	"github.com/google/go-github/github"

	"core/models/labelmaker"
	"core/pipeline/gateway/conflation"
)

func TestWorker(t *testing.T) {
	repoID := int64(23)

	ctx := context.Background()
	lngClient, err := language.NewClient(ctx)
	if err != nil {
		t.Error("error creating language client", err)
	}

	client := github.NewClient(nil)
	url, _ := url.Parse("http://localhost:8000/")
	client.BaseURL = url
	client.UploadURL = url

	bs := new(Server)
	bs.Repos = new(ActiveRepos)
	bs.Repos.Actives = make(map[int64]*ArchRepo)

	created := time.Now()
	started := created.AddDate(0, 0, -1)

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
		Labelmaker: &labelmaker.LBModel{
			Classifier: &labelmaker.LBClassifier{
				Client: lngClient,
				Gateway: labelmaker.CachedNlpGateway{
					NlpGateway: &labelmaker.NlpGateway{
						Client: lngClient,
					},
				},
				Ctx: ctx,
			},
		},
		Settings: HeuprConfigSettings{
			StartTime: started,
		},
	}

	fullname := "skywalker/t-16"
	i := []*github.Issue{
		&github.Issue{
			ID:     github.Int64(1),
			Number: github.Int(1),
			Title:  github.String("Begger's Canyon"),
			Body:   github.String("You'll be a damp mark on the dark side of a canyon wall."),
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
				ID:     github.Int64(2),
				Number: github.Int(2),
				Title:  github.String("Threading the needle"),
				Body:   github.String("You'll be a damp mark on the dark side of a canyon wall."),
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
			// NOTE: Using Sleep to allow the Stop method / Quit selection to
			// complete in the unit test (avoiding race).
			break
		}
	}
}
