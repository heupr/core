package backtest

import "testing"

func TestWalkArchive(t *testing.T) {
	tr := ReplayServer{}
	err := tr.WalkArchive("/home/forstmeier/Downloads/test")
	if err != nil {
		t.Errorf("ERROR: %v", err)
	}
}

// TODO: various unit test scenarios
// - feed in non-gz file
// - only non-IssuesEvent objects in file
