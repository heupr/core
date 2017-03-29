package backtest

import (
    "bytes"
    "compress/gzip"
    "encoding/json"
	"io/ioutil"
    "os"
	"testing"
)

func TestWalkArchive(t *testing.T) {
    tr := ReplayServer{}
    testDir, err := ioutil.TempDir(os.TempDir(), "tatooine")
    if err != nil {
        t.Errorf("Error generating test folder: %v", err)
    }
    err = tr.WalkArchive(testDir)
    if err != nil {
        t.Errorf("Error handling empty directory: %v", err)
    }
    textFile, err := ioutil.TempFile(testDir, "mos_espa.txt")
    if err != nil {
        t.Errorf("Error handling non-JSON file: %v", err)
        t.Errorf("Error writing text test file: %v", err)
    }
    textFile.Write([]byte("Watto's Junkshop"))
    err = tr.WalkArchive(testDir)
    if err != nil {
    }
    os.Remove(textFile.Name())
    content := Event{"IssuesEvent",json.RawMessage("Chalmun's Spaceport Cantina")}
    jsonFile, err := json.Marshal(content)
    if err != nil {
        t.Errorf("Error writing JSON test file: %v", err)
    }
    gz := bytes.Buffer{}
    zipper := gzip.NewWriter(&gz)
    _, err = zipper.Write(jsonFile)
    if err != nil {
        t.Errorf("Error creating gzip file: %v", err)
    }
    err = zipper.Close()
    if err != nil {
        t.Errorf("Error closing gzip writer: %v", err)
    }
    err = tr.WalkArchive(testDir)
    if err != nil {
        t.Errorf("Error handling JSON gzip file: %v", err)
    }
    os.Remove(textFile.Name())
    os.Remove(testDir)
}
