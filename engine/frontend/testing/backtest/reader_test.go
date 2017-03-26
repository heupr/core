package backtest

import "testing"

func TestDirectoryWalk(t *testing.T) {
    tr := ReplayServer{}
    err := tr.DirectoryWalk()
    if err != nil {
        t.Errorf("ERROR: %v", err)
    }
}
