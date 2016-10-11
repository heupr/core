package bhattacharya

import (
    "coralreef-ci/models/issues"
    "time"
)

type Assignee struct {
	Name         string
	LastActive   time.Time
	Profile      []string
}

func BuildProfiles(issues []issues.Issue) map[string][]Assignee {
    profiles := make(map[string][]Assignee)

    for i := 0; i < len(issues); i ++ {
        name := issues[i].Assignee
        active := assignees[i].Resolved
        profile := issues[i].Labels

        if _, ok := profiles[name]; ok {
            // if true:
            // - append new labels to the Profile value string slice
            // - compare activity to LastActivity and update if more recent
        } else {
            profiles[name].Name = name
            profiles[name].LastActive = active
            profiles[name].Profile = profile
        }

        // if _, ok := profiles[assignee]; ok {
        //     profiles[assignee] = append(profiles[assignee], labels...)
        // } else {
        //     profiles[assignee] = labels
        // }
        // TODO: LastActivity
        // build functionality to find most recent activity
        // should go here
        // compare incoming date to existing date for given Name
        // - if more recent, replace it; else ignore and move on
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
