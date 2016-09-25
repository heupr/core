package bhattacharya

import (
	"coralreef-ci/models/issues"
	"math/rand"
)

func Shuffle(issueList []issues.Issue, seed int64) {
	rand.Seed(seed)
	for i := 0; i < len(issueList); i++ {
		r := rand.Intn(i + 1)
		if i != r {
			issueList[r].Body, issueList[i].Body = issueList[i].Body, issueList[r].Body
		}
	}
}
