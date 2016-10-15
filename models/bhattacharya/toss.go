package bhattacharya

import (
	"coralreefci/models/issues"
	"time"
)

type Assignee struct {
	Name          string
	LastActive    time.Time
	Profile       []string
	Contributions int
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
			profiles[name].Contributions += 1
		} else {
			profiles[name] = &Assignee{
				Name:          name,
				LastActive:    active,
				Profile:       labels,
				Contributions: 1,
			}
		}
	}

    // for each item in the profiles
    // - filter out profile labels
    // - insert new filter of labels

    for index, _ := range profiles {
        cleaned := profileFilter(profiles[index].Profile)
        profiles[index].Profile = cleaned
    }

	return profiles
}

func profileFilter(input []string) []string {
    found := make(map[string]bool)
    clean := []string{}
    for i := 0; i < len(input); i ++ {
        if found[input[i]] != true {
            found[input[i]] = true
            clean = append(clean, input[i])
        }
    }
    return clean
}

// taking score from scores in the LogScores
// in the nb_classifier.go file
// scores -> []float64
// look at the topTrhee function
// primitive tossing graph

// functions:
// tossing function
// - actually acounts for given number of possible assignees
// ranking function
// - provides logic for pruning given list
// profile builder
// - constructs slice of labels for each assignee
