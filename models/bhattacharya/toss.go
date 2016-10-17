package bhattacharya

import (
	"coralreefci/models/issues"
	"sort"
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

func Tossing(scores []float64, top int) []int {
	scoreMap := make(map[int]float64)
	depth := 1
	for i := 0; i < len(scores); i++ {
		scoreMap[i] = scores[i]
	}
	values := []float64{}

	for _, value := range scoreMap {
		values = append(values, value)
	}
	sort.Float64s(values)
	flipScoreMap := make(map[float64]int)
	for integer, floater := range scoreMap {
		flipScoreMap[floater] = integer
	}
	topIndex := []int{}

	if (len(values) < top) {
		depth = len(values)
	} else {
		depth = top
	}
	for i := 0; i < depth; i++ {
		if _, ok := flipScoreMap[values[i]]; ok {
				topIndex = append(topIndex, flipScoreMap[values[i]])
			}
	}
	return topIndex
}
