package frontend

import (
	"testing"

	"github.com/google/go-github/github"
)

var input = make(chan chan github.IssuesEvent)

func TestNewWorker(t *testing.T) {
	output := NewWorker(1, input)
	comp := new(Worker)
	if output.ID == comp.ID {
		t.Errorf("\nNewWorker is failing to generate an initalized Worker\nGenerated struct: %+v", output)
	}
}
