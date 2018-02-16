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
	Actives map[int64]*ArchRepo
}

type Server struct {
	Server   http.Server
	Database MemSQL
	Repos    *ActiveRepos
}

// HeuprInstallationEvent is a workaround a Github API limitation. This is
// required to wrap HeuprInstallation.
type HeuprInstallationEvent struct {
	// The action that was performed. Can be either "created" or "deleted".
	Action            *string            `json:"action,omitempty"`
	Sender            *github.User       `json:"sender,omitempty"`
	HeuprInstallation *HeuprInstallation `json:"installation,omitempty"`
	Repositories      []HeuprRepository  `json:"repositories,omitempty"`
}

// HeuprInstallation is a workaround a Github API limitation. The go-github
// library is missing the repositories field.
type HeuprInstallation struct {
	ID              *int         `json:"id,omitempty"`
	Account         *github.User `json:"account,omitempty"`
	AppID           *int         `json:"app_id,omitempty"`
	AccessTokensURL *string      `json:"access_tokens_url,omitempty"`
	RepositoriesURL *string      `json:"repositories_url,omitempty"`
	HTMLURL         *string      `json:"html_url,omitempty"`
}

type HeuprRepository struct {
	ID       *int64  `json:"id,omitempty"`
	Name     *string `json:"name,omitempty"`
	FullName *string `json:"full_name,omitempty"`
}

func (bs *Server) activateHandler(w http.ResponseWriter, r *http.Request) {
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
			labels, err := bs.Database.ReadLabels(*repository.ID)
			if err != nil {
				http.Error(w, "error fetching issue labels", 500)
			}
			bs.NewArchRepo(*repository.ID, settings, labels)
			bs.NewClient(*repository.ID, *integration.HeuprInstallation.AppID, *integration.HeuprInstallation.ID)
			bs.NewModel(*repository.ID)
		}
	}
}

func (bs *Server) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/activate-ingestor-backend", bs.activateHandler)
	bs.Server = http.Server{
		Addr:    utils.Config.BackendServerAddress,
		Handler: mux,
	}
	bs.OpenSQL()
	defer bs.CloseSQL()

	bs.Repos = &ActiveRepos{Actives: make(map[int64]*ArchRepo)}

	integrations, err := bs.Database.ReadIntegrations()
	if err != nil {
		utils.AppLog.Error("retrieve bulk tokens on ingestor restart", zap.Error(err))
	}

	for _, integration := range integrations {
		if _, ok := bs.Repos.Actives[integration.RepoID]; !ok {
			settings, err := bs.Database.ReadHeuprConfigSettingsByRepoID(integration.RepoID)
			if err != nil {
				panic(err)
			}
			if settings.StartTime.IsZero() {
				settings.StartTime = time.Now()
			}
			labels, err := bs.Database.ReadLabels(integration.RepoID)
			if err != nil {
				panic(err)
			}
			bs.NewArchRepo(integration.RepoID, settings, labels)
			bs.NewClient(integration.RepoID, integration.AppID, integration.InstallationID)
			bs.NewModel(integration.RepoID)
		}
	}

	// Keeping this channel to implement graceful shutdowns if needed.
	wiggin := make(chan bool)
	bs.Timer(wiggin)
	bs.Server.ListenAndServe()
}

func (bs *Server) OpenSQL() {
	bs.Database.Open()
}

func (bs *Server) CloseSQL() {
	bs.Database.Close()
}

// Timer conducts periodic pulldowns from the MemSQL database for processing.
func (bs *Server) Timer(ender chan bool) {
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
				collector(data)
			case <-ender:
				ticker.Stop()
				close(ender)
				return
			}
		}
	}()
}
