package backend

import (
	"testing"
)

func Testcollector(t *testing.T) {
	repodataMap := map[int64]*RepoData{
		1: &RepoData{
			RepoID: 66,
			EligibleAssignees: map[string]int{
				"Cody": 2224,
			},
		},
		2: &RepoData{
			RepoID: 66,
			EligibleAssignees: map[string]int{
				"Bacara": 1138,
			},
		},
		3: &RepoData{
			RepoID: 66,
			EligibleAssignees: map[string]int{
				"Bly": 5024,
			},
		},
		4: &RepoData{
			RepoID: 66,
			EligibleAssignees: map[string]int{
				"Neyo": 8826,
			},
		},
	}
	collector(repodataMap)
	if len(workload) != len(repodataMap) {
		t.Errorf(
			"collector incorrectly populating workload; wanted %v, received %v",
			len(repodataMap),
			len(workload),
		)
	}
}
