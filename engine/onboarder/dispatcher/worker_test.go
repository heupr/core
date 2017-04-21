package dispatcher

import (
	"testing"

	"coralreefci/engine/onboarder/retriever"
)

func TestNewWorker(t *testing.T) {
	testID := 1
	channel := make(chan chan *retriever.RepoData)
	testWorker := NewWorker(testID, channel)
	if testWorker.ID != testID {
		t.Error("Failure creating new worker object")
	}
}
