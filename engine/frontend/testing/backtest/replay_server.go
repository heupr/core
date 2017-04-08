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
	"os"
	"path/filepath"
	// "time"

	"github.com/google/go-github/github"

	"coralreefci/engine/frontend"
	// "coralreefci/models"
)

// [X] COMPLETED
const (
	secretKey = "chalmun"
	localPath = "http://localhost:8080/"
)

// [X] COMPLETED
var modelList = []*frontend.HeuprModel{}

// [X] COMPLETED
type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// [X] COMPLETED
type BacktestServer struct {
	frontend.HeuprServer
	Archive []Event
}

// [ ] COMPLETED
// - [ ] unit test
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

// [ ] COMPLETED
// - [ ] unit test
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

// [ ] COMPLETED
// - [ ] unit test
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

// [ ] COMPLETED
// - [ ] unit test
func (b *BacktestServer) backtestHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Review to see if there is anything else that needs to be done to
	//       the individual objects prior to being written to the
	//       ResponseWriter.
	half := len(b.Archive) / 2
	for _, v := range b.Archive[:half] {
		fmt.Fprint(w, v)
	}
}

// [ ] COMPLETED
// - [ ] complete source code
// - [ ] unit test
func (b *BacktestServer) backtestPredict(w http.ResponseWriter, r *http.Request) {

}

// [ ] COMPLETED
// - [ ] parse payload for repo ID
// - [ ] unit test
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
	// TODO: Parse payload to get given repo ID int
	// plStr := payload.String()
	req.Header.Set("X-Github-Event", "issues")
	req.Header.Set("X-GitHub-Delivery", "placeholder")
	req.Header.Set("content-type", "application/json")
	mac := hmac.New(sha1.New, []byte(secretKey))
	mac.Write(payload.Bytes())
	sig := "sha1=" + hex.EncodeToString(mac.Sum(nil))
	req.Header.Set("X-Hub-Signature", sig)

	id := 1
	// TODO: Pass repo ID into id
	// TODO: Check to make sure Do() is having the right values passed in
	b.Repos[id].Client.Do(req, payload)
}
