package bhattacharya

import (
	"math/rand"
)

// DOC: Shuffle provides object shuffling for training / backtesting purposes.
//      This helper function returns a new instance of a shuffled list.
func Shuffle(issueList []Issue, seed int64) {
	rand.Seed(seed)
	for i := 0; i < len(issueList); i++ {
		r := rand.Intn(len(issueList) - 1)
		if i != r {
			issueList[r], issueList[i] = issueList[i], issueList[r]
		}
	}
}
