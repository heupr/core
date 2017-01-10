package conflation

import (
	"github.com/google/go-github/github"
	"testing"
)

var TestStruct7 = Scenario7{}

func TestFilter7(t *testing.T) {
	number := 2
	bodyNaked := "The Mos Espa Grand Arena"
	bodyClothed := "Gasgano's podracer is Fixed #2"
	pullNaked := github.PullRequest{Number: &number, Body: &bodyNaked}
	pullClothed := github.PullRequest{Number: &number, Body: &bodyClothed}
	pullWith := &ExpandedIssue{PullRequest: CRPullRequest{pullNaked, []int{}, []CRIssue{}}}
	pullWithout := &ExpandedIssue{PullRequest: CRPullRequest{pullClothed, []int{}, []CRIssue{}}}

	if !TestStruct7.Filter(pullWith) {
		t.Error(
			"\nNAKED PULL REQUEST INCORRECTLY REMOVED",
			"\nTEST PULL REQUEST : ", *pullWith,
		)
	}
	if TestStruct7.Filter(pullWithout) {
		t.Error(
			"\nCLOSING PULL REQUEST INCORRECTLY INCLUDED",
			"\nTEST PULL REQUEST : ", *pullWithout,
		)
	}
}
