package replay

import (
	"bytes"
	"coralreefci/engine/ingestor"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	gzip "github.com/klauspost/pgzip"
	"io"
	"io/ioutil"
	"net/http"
	//"net/url"
	// "reflect"
	"os"
	"path/filepath"

	"github.com/google/go-github/github"
	"github.com/gorilla/mux"
	//"coralreefci/engine/onboarder"
	// "coralreefci/models"
	"time"
)

const (
	secretKey = "test"
	localPath = "http://localhost:8000/"
)

//var modelList = []*onboarder.ArchModel{}

var webhooksplit float32 = 1

type Issue struct {
	Payload json.RawMessage `json:issue`
}

type BacktestServer struct {
	client http.Client
	DB     *ingestor.Database
	server http.Server
	//onboarder.RepoServer
	events        []ingestor.Event
	WebhookEvents []ingestor.Event
}

func (b *BacktestServer) routes() *mux.Router {
	gorilla := mux.NewRouter()
	gorilla.HandleFunc("/repos/{org}/{repo}/issues", b.getIssues)
	return gorilla
}

func (b *BacktestServer) Start() {
	b.server = http.Server{Addr: "127.0.0.1:8000", Handler: b.routes()}
	err := b.server.ListenAndServe()
	if err != nil {
		fmt.Println("BacktestServer", err)
	}

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	b.client = http.Client{Transport: tr}
}

func (b *BacktestServer) LoadArchive(path string) {
	fn, _ := os.Stat(path)
	switch mode := fn.Mode(); {
	case mode.IsDir():
		files, _ := ioutil.ReadDir(path)
		loadedFiles := 0
		totalFiles := len(files)
		err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
			parseErr := b.parseFile(path)
			if parseErr == nil {
				loadedFiles++
				fmt.Printf("Loaded %d out of %d files\n", loadedFiles, totalFiles)
			} else {
				fmt.Println(parseErr)
			}
			if loadedFiles > 0 && loadedFiles%15 == 0 {
				b.DB.BulkInsertBacktestEvents(b.events)
				b.events = []ingestor.Event{}
				fmt.Printf("Inserted %d out of %d files\n", loadedFiles, totalFiles)
			}
			return nil
		})
		if err != nil {
			fmt.Printf("Error walking directory: %v", err)
		}
	case mode.IsRegular():
		b.parseFile(path)
	default:
		fmt.Println("Unrecognized argument; provide a file or directory")
	}
	if len(b.events) > 0 {
		b.DB.BulkInsertBacktestEvents(b.events)
		fmt.Printf("Inserted remaining records", len(b.events))
		b.events = []ingestor.Event{}
	}
}

func (b *BacktestServer) parseFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	gr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gr.Close()
	jd := json.NewDecoder(gr)
	for {
		e := ingestor.Event{}
		if err := jd.Decode(&e); err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		switch e.Type {
		case "IssuesEvent":
			//case "PullRequestEvent":
			b.events = append(b.events, e)
		}
	}
	return nil
}

/*
func (b *BacktestServer) loadRepos() {
	// TODO: parsing some input (likely a file w/ JSON content)
	// - defines the desired repos to run (specific GitHub Repositories)
	repos := []github.Repository{} // TEMPORARY

	client := github.NewClient(nil)
	u, _ := url.Parse(localPath)
	client.BaseURL = u
	client.UploadURL = u

	for i := 0; i < len(repos); i++ {
		b.Repos[i] = &onboarder.ArchRepo{
			Repo:   &repos[i],
			Hive:   &onboarder.ArchHive{Models: modelList},
			Client: client,
		}
	}
}

func (b *BacktestServer) backtestHandler(w http.ResponseWriter, r *http.Request) {
	jsonRepo := r.FormValue("repository")
	repo := github.Repository{}
	if err := json.Unmarshal([]byte(jsonRepo), &repo); err != nil {
		fmt.Println("Error unmarshalling JSON into repo")
	}

	// TODO: There still needs to be logic in what is pulled out of the Archive
	//       and how best to return it to the writer (likely via fmt.Fprint).
	//       The archive has potential to be relatively large and looping
	//       through it could be expensive particularly if this handler is
	//       being hit for every repo in a given backtest run.
}*/

func (b *BacktestServer) backtestPredict(w http.ResponseWriter, r *http.Request) {

}

func (b *BacktestServer) getIssues(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	org := vars["org"]
	repo := vars["repo"]
	events, _ := b.DB.ReadBacktestEvents(org + "/" + repo)
	issues := make([]*github.Issue, len(events))
	if webhooksplit == 1 {
		for i := 0; i < len(events); i++ {
			issues[i] = events[i].Payload.Issue
		}
	} else {
		for i := 0; i < int(float32(len(events))/webhooksplit); i++ {
			issues[i] = events[i].Payload.Issue
		}
		for i := int(float32(len(events)) / webhooksplit); i < len(events); i++ {
			b.WebhookEvents = append(b.WebhookEvents, events[i])
		}
	}
	payload, _ := json.Marshal(&issues)
	w.Write(payload)
}

func (b *BacktestServer) HTTPPost(payload *bytes.Buffer) {
	req, err := http.NewRequest("POST", "http://localhost:8080/hook", payload)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("X-Github-Event", "issues")
	req.Header.Set("X-GitHub-Delivery", "placeholder")
	req.Header.Set("content-type", "application/json")
	mac := hmac.New(sha1.New, []byte(secretKey))
	mac.Write(payload.Bytes())
	sig := "sha1=" + hex.EncodeToString(mac.Sum(nil))
	req.Header.Set("X-Hub-Signature", sig)

	b.client.Do(req)
}
