package bhattacharya

import (
	"coralreefci/models/issues"
	"time"
)

type Assignee struct {
	Name       string
	LastActive time.Time
	Profile    []string
}

type Assignees map[string]*Assignee

func BuildProfiles(issues []issues.Issue) Assignees {
	profiles := make(Assignees)

	for i := 0; i < len(issues); i++ {
		name := issues[i].Assignee
		active := issues[i].Resolved
		labels := issues[i].Labels

		if _, ok := profiles[name]; ok {
            if active.After(profiles[name].LastActive) {
                profiles[name].LastActive = active
            }
            profiles[name].Profile = append(profiles[name].Profile, labels...)
		} else {
			profiles[name] = &Assignee{
				Name:       name,
				LastActive: active,
				Profile:    labels,
			}
		}
    }

	// insert a filter for profiles here
	// stand alone function

	return profiles

}

// taking score from scores in the LogScores
// in the nb_classifier.go file
// scores -> []float64
// look at the topTrhee function
// primitive tossing graph

// new data from GH:
// username column - string type
// - include last activity associated
// last activity - time type
// labels column - string slice type
// - associate with each username (fixer of the given issue)

// functions:
// tossing function
// - actually acounts for given number of possible assignees
// ranking function
// - provides logic for pruning given list
// profile builder
// - constructs slice of labels for each assignee
