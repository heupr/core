package main

import (
	"core/pipeline/gateway/conflation"
	"math/rand"
)

// DOC: Shuffle provides object shuffling for training / backtesting purposes.
//      This helper function returns a new instance of a shuffled list.
func Shuffle(issues []conflation.ExpandedIssue, seed int64) {
	rand.Seed(seed)
	for i := 0; i < len(issues); i++ {
		r := rand.Intn(len(issues) - 1)
		if i != r {
			issues[r], issues[i] = issues[i], issues[r]
		}
	}
}
