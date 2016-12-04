package bhattacharya

import (
	"coralreefci/models/issues"
	"errors"
	"strconv"
	"strings"
)

type matrix map[string]map[string]int

// TODO: refactor into an object-oriented functionality
/*
type Matrix struct {
	Elements map[string]map[string]int
}
*/

// BuildMatrix takes two arguments:
// expected - slice of issues used in testing; static data
// predicted - slice of issues the model predicted; output data
// these are the same length as the latter is just predictions of the former
func BuildMatrix(expected, predicted []issues.Issue) (matrix, error) {
	if len(expected) != len(predicted) {
		return nil, errors.New("INPUT SLICES ARE NOT EQUAL LENGTH")
	}

	outputMatrix := make(map[string]map[string]int)

	for i := 0; i < len(expected); i++ {
		exp := expected[i].Assignee
		pre := predicted[i].Assignee
		if _, ok := outputMatrix[exp]; ok {
			outputMatrix[exp][pre] += 1
		} else {
			outputMatrix[exp] = make(map[string]int)
			outputMatrix[exp][pre] = 1
		}
	}
	return outputMatrix, nil
}

func getClassTP(class string, m matrix) float64 {
	return float64(m[class][class])
}

func getClassTN(class string, m matrix) float64 {
	count := 0.0
	for columnHead := range m {
		if columnHead == class {
			continue
		}
		for rowHead := range m[columnHead] {
			if rowHead == class {
				continue
			}
			count += float64(m[columnHead][rowHead])
		}
	}
	return count
}

func getClassFP(class string, m matrix) float64 {
	count := 0.0
	for columnHead := range m {
		if columnHead == class {
			continue
		}
		count += float64(m[columnHead][class])
	}
	return count
}

func getClassFN(class string, m matrix) float64 {
	count := 0.0
	for rowHead := range m[class] {
		if rowHead == class {
			continue
		}
		count += float64(m[class][rowHead])
	}
	return count
}

func round(number float64) float64 {
	return float64(int(number*100.00)) / 100.00
}

func toString(number float64) string {
	return strconv.FormatFloat(number, 'f', 2, 64)
}

func getPrecision(class string, m matrix) float64 {
	classTP := getClassTP(class, m)
	classFP := getClassFP(class, m)
	return round(classTP / (classTP + classFP))
}

func getRecall(class string, m matrix) float64 {
	classTP := getClassTP(class, m)
	classFN := getClassFN(class, m)
	return round(classTP / (classTP + classFN))
}

func getAccuracy(m matrix) float64 {
	correct := 0.0
	total := 0.0
	for columnHead := range m {
		for rowHead := range m[columnHead] {
			if columnHead == rowHead {
				correct += float64(m[columnHead][rowHead])
			}
			total += float64(m[columnHead][rowHead])
		}
	}
	return round(float64(correct) / float64(total))
}

func getTestCount(m matrix) float64 {
	count := 0.0
	for columnHead := range m {
		for rowHead := range m[columnHead] {
			count += float64(m[columnHead][rowHead])
		}
	}
	return count
}

func fillMatrix(m matrix) matrix {
	for columnHead := range m {
		for key := range m {
			if _, ok := m[columnHead][key]; ok {
				continue
			} else {
				m[columnHead][key] = 0
			}
		}
	}
	return m
}

func ClassSummary(class string, m matrix) string {
	input := []string{"SUMMARY RESULTS FOR CLASS: ", class, "\n",
		"TRUE POSITIVES:  ", toString(getClassTP(class, m)), "\n",
		"TRUE NEGATIVES:  ", toString(getClassTN(class, m)), "\n",
		"FALSE POSITIVES: ", toString(getClassFP(class, m)), "\n",
		"FALSE NEGATIVES: ", toString(getClassFN(class, m)), "\n",
		"PRECISION:       ", toString(getPrecision(class, m)), "\n",
		"RECALL:          ", toString(getRecall(class, m)), "\n"}
	output := strings.Join(input, "")
	return output
}

func FullSummary(m matrix) string {
	input := []string{"SUMMARY RESULTS FOR FULL MATRIX\n",
		"TOTAL TESTS:    ", toString(getTestCount(m)), "\n",
		"TOTAL ACCURACY: ", toString(getAccuracy(m)), "\n"}
	output := strings.Join(input, " ")
	return output
}
