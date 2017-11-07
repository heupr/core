package backend

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/google/go-github/github"
	"go.uber.org/zap"

	"core/utils"
)

type ActiveRepos struct {
	sync.RWMutex
	Actives map[int]*ArchRepo
}

type BackendServer struct {
	Server   http.Server
	Database MemSQL
	Repos    *ActiveRepos
}

// Workaround Github API limitation. This is required to wrap HeuprInstallation
type HeuprInstallationEvent struct {
	// The action that was performed. Can be either "created" or "deleted".
	Action            *string            `json:"action,omitempty"`
	Sender            *github.User       `json:"sender,omitempty"`
	HeuprInstallation *HeuprInstallation `json:"installation,omitempty"`
	Repositories      []HeuprRepository  `json:"repositories,omitempty"`
}

// Workaround Github API limitation. go-github is missing repositories field
type HeuprInstallation struct {
	ID              *int         `json:"id,omitempty"`
	Account         *github.User `json:"account,omitempty"`
	AppID           *int         `json:"app_id,omitempty"`
	AccessTokensURL *string      `json:"access_tokens_url,omitempty"`
	RepositoriesURL *string      `json:"repositories_url,omitempty"`
	HTMLURL         *string      `json:"html_url,omitempty"`
}

type HeuprRepository struct {
	ID       *int    `json:"id,omitempty"`
	Name     *string `json:"name,omitempty"`
	FullName *string `json:"full_name,omitempty"`
}

func (bs *BackendServer) activateHandler(w http.ResponseWriter, r *http.Request) {
	var activationParams struct {
		InstallationEvent HeuprInstallationEvent `json:"installation_event,omitempty"`
		Limit             time.Time              `json:"limit,omitempty"`
	}
	err := json.NewDecoder(r.Body).Decode(&activationParams)
	if err != nil {
		utils.AppLog.Error("unable to decode json message", zap.Error(err))
	}
	integration := activationParams.InstallationEvent
	for _, repository := range integration.Repositories {
		if _, ok := bs.Repos.Actives[*repository.ID]; !ok {
			settings := HeuprConfigSettings{StartTime: time.Now(), IgnoreLabels: make(map[string]bool), IgnoreUsers: make(map[string]bool)}
			bs.NewArchRepo(*repository.ID, settings)
			bs.NewClient(*repository.ID, *integration.HeuprInstallation.AppID, *integration.HeuprInstallation.ID)
			bs.NewModel(*repository.ID)
		}
	}
}

func (bs *BackendServer) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/activate-ingestor-backend", bs.activateHandler)
	bs.Server = http.Server{
		Addr:    utils.Config.BackendServerAddress,
		Handler: mux,
	}
	bs.OpenSQL()
	defer bs.CloseSQL()

	bs.Repos = &ActiveRepos{Actives: make(map[int]*ArchRepo)}

	integrations, err := bs.Database.ReadIntegrations()
	if err != nil {
		utils.AppLog.Error("retrieve bulk tokens on ingestor restart", zap.Error(err))
	}

	for _, integration := range integrations {
		if _, ok := bs.Repos.Actives[integration.RepoId]; !ok {
			settings, err := bs.Database.ReadHeuprConfigSettingsByRepoId(integration.RepoId)
			if err != nil {
				panic(err)
			}
			if settings.StartTime.IsZero() {
				settings.StartTime = time.Now()
			}
			bs.NewArchRepo(integration.RepoId, settings)
			bs.NewClient(integration.RepoId, integration.AppId, integration.InstallationId)
			bs.NewModel(integration.RepoId)
		}
	}

	// Keeping this channel to implement graceful shutdowns if needed.
	wiggin := make(chan bool)
	bs.Timer(wiggin)
	bs.Server.ListenAndServe()
}

func (bs *BackendServer) OpenSQL() {
	bs.Database.Open()
}

func (bs *BackendServer) CloseSQL() {
	bs.Database.Close()
}

// Periodically conducts pulldowns from the MemSQL database for processing.
func (bs *BackendServer) Timer(ender chan bool) {
	ticker := time.NewTicker(time.Second * 5)

	bs.Dispatcher(10)

	go func() {
		for {
			select {
			case <-ticker.C:
				data, err := bs.Database.Read()
				if err != nil {
					utils.AppLog.Error("backend timer", zap.Error(err))
				}
				Collector(data)
			case <-ender:
				ticker.Stop()
				close(ender)
				return
			}
		}
	}()
}
