package backtest

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

func (r *ReplayServer) WalkArchive(dir string) error {
	err := filepath.Walk(dir, func(fp string, fi os.FileInfo, err error) error {
		if !fi.IsDir() {
			f, err := os.Open(fp)
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
				switch e.Type {
				case "IssuesEvent":
					buf := bytes.NewBufferString(string(e.Payload))
					r.HTTPPost(buf)
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
