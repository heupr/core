package backtest

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
)

const secretKey = "chalmun"

type ReplayServer struct {
	client http.Client
}

func CreateReplayServer() *ReplayServer {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	return &ReplayServer{client: http.Client{Transport: tr}}
}

// TODO: ngrok url is now located here and in hooker.go (lets fix that with an
// env variable. Fortunately ngrok is written in Golang (so that helps))
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

	r.client.Do(req)
}

/*
TODO lets figure out where to put this one
func (r *ReplayServer) LoadGithubArchive() {

}*/
