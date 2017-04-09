package backtest

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	// "reflect"
	"os"
	"path/filepath"

	"github.com/google/go-github/github"

	"coralreefci/engine/frontend"
	// "coralreefci/models"
)

const (
	secretKey = "chalmun"
	localPath = "http://localhost:8080/"
)

var modelList = []*frontend.HeuprModel{}

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type BacktestServer struct {
	frontend.HeuprServer
	Archive []Event
}

func (b *BacktestServer) loadArchive() {
	args := os.Args
	for _, arg := range args[1:] {
		fn, _ := os.Stat(arg)
		switch mode := fn.Mode(); {
		case mode.IsDir():
			err := filepath.Walk(arg, func(path string, f os.FileInfo, err error) error {
				b.parseFile(path)
				return nil
			})
			if err != nil {
				fmt.Printf("Error walking directory: %v", err)
			}
		case mode.IsRegular():
			b.parseFile(arg)
		default:
			fmt.Println("Unrecognized argument; provide a file or directory")
		}
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
		e := Event{}
		if err := jd.Decode(&e); err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		b.Archive = append(b.Archive, e)
	}
	return nil
}

func (b *BacktestServer) loadRepos() {
	// TODO: parsing some input (likely a file w/ JSON content)
	// - defines the desired repos to run (specific GitHub Repositories)
	repos := []github.Repository{} // TEMPORARY

	client := github.NewClient(nil)
	u, _ := url.Parse(localPath)
	client.BaseURL = u
	client.UploadURL = u

	for i := 0; i < len(repos); i++ {
		b.Repos[i] = &frontend.HeuprRepo{
			Repo:   &repos[i],
			Hive:   &frontend.HeuprHive{Models: modelList},
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
}

func (b *BacktestServer) backtestPredict(w http.ResponseWriter, r *http.Request) {

}

// TODO: ngrok url is now located here and in hooker.go (lets fix that with an
//       env variable. Fortunately ngrok is written in Golang (so that helps))
// TODO: Per Gor Replay File Add Missing HTTP Headers (File in Slack Channel - requests_0.gor)
// TODO: (see unit test file for more TODOS)
// TODO: Perf: Reuse Http Request objects
func (b *BacktestServer) HTTPPost(payload *bytes.Buffer) {
	req, err := http.NewRequest("POST", localPath, payload)
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

	repo := github.Repository{}
	if err := json.Unmarshal(payload.Bytes(), &repo); err != nil {
		fmt.Println("Error unmarshalling JSON into repo")
	}
	id := *repo.ID
	b.Repos[id].Client.Do(req, payload)
}
