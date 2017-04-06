package backtest

import (
	"bytes"
    "compress/gzip"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
    "encoding/json"
	"fmt"
	"net/http"
    "net/url"
    "os"
    "path/filepath"
	"time"

    "github.com/google/go-github/github"

    "coralreefci/models"
)

const (
    secretKey = "chalmun"
    localPath = "http://localhost:8080/"
    modelList = []models.Model{}
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type BacktestServer struct {
    HeuprServer
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
        case mode.IsRegular():
            b.parseFile(arg)
        default:
            fmt.Println("Unrecognized argument; provide a file or directory")
        }
    }
}

func (b *BacktestServer) parseFile(filename string) {
    f, _ := os.Open(filename)
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
}

func (b *BacktestServer) loadRepos() {
    // TODO: parsing some input
    // defines the desired repos to run
    // specific GitHub Repositories
    repos = []github.Repositories{} // TEMPORARY

    client := github.NewClient(nil)
    u, _ := url.Parse(localPath)
    client.BaseURL = u
    client.UploadURL = u

    for i := 0; i < len(repos); i ++ {
        b.Repos[i] = &HeuprRepo{
            Repo: &repos[i]
            Hive: &HeuprHive{Models: &models}
            Client: &client,
        }
    }
}

func (b *BacktestServer) backtestHandler(w *http.ResponseWriter, r *http.Request) {
    // TODO: Review to see if there is anything else that needs to be done to
    //       the individual objects prior to being written to the
    //       ResponseWriter.
    for _, v := range h.Archive {
        fmt.Fprint(w, v)
    }
}

// TODO: ngrok url is now located here and in hooker.go (lets fix that with an
//       env variable. Fortunately ngrok is written in Golang (so that helps))
// TODO: Per Gor Replay File Add Missing HTTP Headers (File in Slack Channel - requests_0.gor)
// TODO: (see unit test file for more TODOS)
// TODO: Perf: Reuse Http Request objects
func (r *ReplayServer) HTTPPost(payload *bytes.Buffer) {
    req, err := http.NewRequest("POST", "http://5b0f0030.ngrok.io/hook", payload)
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

    // TODO: Note that this will need to change if the default address needs
    // to be overwritten (e.g. r.client.NewRequest())
    r.client.Do(req)
}
