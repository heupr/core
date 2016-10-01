package bhattacharya

import (
    "coralreef-ci/models/issues"
    "errors"
)

func BuildMatrix(predicted, expected []issues.Issue) (map[string]map[string]int, error) {
    matrix := make(map[string]map[string]int)

    if len(predicted) != len(expected) {
        return errors.New("SLICE LENGTH MISMATCH")
    }
}

// arguments: predicted fixer list, actual fixer list
// return: map of string of map of string of int
// - matrix := make()

// structure of the algorithm:
// for the actual assignee
// -
