package backtest

import (
    "encoding/json"
    "fmt" // TEMPORARY
    "io/ioutil"
    "os"
    "path/filepath"

    "github.com/google/go-github/github"
)

func (r *ReplayServer) DirectoryWalk() error {
    rel := "/home/forstmeier/Downloads/test"
    var contents []github.Event
    // contents := make([]github.Event, 0)
    err := filepath.Walk(rel, func(fp string, fi os.FileInfo, err error) error {
        if !fi.IsDir() {

            f, err := ioutil.ReadFile(fp)
            if err != nil {
                return err
            }
            err = json.Unmarshal(f, &contents)
            if err != nil {
                return err
            }

            // fb, err := os.Open(path)
            // if err != nil {
            //     return err
            // }
            // jsonParser := json.NewDecoder(fb)
            // if err = jsonParser.Decode(&contents); err != nil {
            //     return err
            // }
            // fb, err := ioutil.ReadFile(path)
            // if err != nil {
            //     return err
            // }
            // if err = json.Unmarshal(fb, &contents); err != nil {
            //     return err
            // }
        }
        return nil
    })
    if err != nil {
        return err
    }
    fmt.Println(contents)
    return nil
}
