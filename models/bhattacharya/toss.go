package bhattacharya

import (
	"coralreefci/models/issues"
	"sort"
	// "strconv" // NOTE: temporarily excluded
	"time"
)

type Assignee struct {
	Name          string
	LastActive    time.Time
	Profile       []string
	Contributions int
}

type Assignees map[string]*Assignee

type TossingGraph struct {
	Assignees  []string
	GraphDepth int
	Logger     *CoralReefLogger
}

// NOTE: currently the setup of the logic is the NBClassifier is the "gatekeeper" for the
// distinct list of assigness which is then passed into the TossingGraph; however, going
// forward, this will ultimately be swapped so that TossingGraph is first passed the
// complete list of assingees which it then performs distinct assignee logic on so that
// it is the new "gatekeeper" for the program.
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
	for index, _ := range profiles {
		cleaned := profileFilter(profiles[index].Profile)
		profiles[index].Profile = cleaned
	}
	return profiles
}

func profileFilter(input []string) []string {
	found := make(map[string]bool)
	clean := []string{}
	for i := 0; i < len(input); i++ {
		if found[input[i]] != true {
			found[input[i]] = true
			clean = append(clean, input[i])
		}
	}
	return clean
}

func (c *TossingGraph) Tossing(scores []float64) []int {
	scoreMap := make(map[int]float64)
	for i := 0; i < len(scores); i++ {
		scoreMap[i] = scores[i]
	}

	scoreValues := []float64{}
	for _, value := range scoreMap {
		scoreValues = append(scoreValues, value)
	}

	sort.Sort(sort.Reverse(sort.Float64Slice(scoreValues))) // TODO: refactor and check to see if this is appropriate (probably duplicated)

	flipScoreMap := make(map[float64]int)
	for integer, floater := range scoreMap {
		flipScoreMap[floater] = integer
	}

	topIndex := []int{}

	top := 0
	if len(scores) < c.GraphDepth {
		top = len(scores)
	} else {
		top = c.GraphDepth
	}

	for i := 0; i < top; i++ {
		if _, ok := flipScoreMap[scoreValues[i]]; ok {
			topIndex = append(topIndex, flipScoreMap[scoreValues[i]])
		}
	}

	// NOTE: temporarily excluding various logging applications.
	// TODO: include a "logging flag" that would fire this particular section
	// of the logging - do not focus on it currently and have the program log
	/*
		logOutput := []string{}
		for i := 0; i < len(scores); i++ {
			if _, ok := flipScoreMap[scoreValues[i]]; ok {
				logOutput = append(logOutput, strconv.FormatFloat(scoreValues[i], 'f', -1, 32) + "," + c.Assignees[flipScoreMap[scoreValues[i]]] )
			}
		}
		c.Logger.Log(logOutput)
	*/
	// TODO: logger generation:
	// List of all contributors ranked by their logScores
	// ^ this needs the location within the input slice of the contributors
	// ^ ths is used to match against the assignees string value slice

	return topIndex
}
