package models

import (
	"coralreefci/engine/gateway/conflation"
	"errors"
	"strings"
)

type matrix map[string]map[string]int

// DOC: BuildMatrix generates a new confusion matrix for evaluation.
//      BuildMatrix takes two arguments:
//      - expected - slice of issues used in testing; static data
//      - predicted - slice of issues the model predicted; output data
//      These are the same length as the latter is just predictions of the
//      former.
func (m *Model) BuildMatrix(expected, predicted []conflation.ExpandedIssue) (matrix, error) {
	if len(expected) != len(predicted) {
		return nil, errors.New("INPUT SLICES ARE NOT EQUAL LENGTH")
	}

	outputMatrix := make(map[string]map[string]int)

	for i := 0; i < len(expected); i++ {
		exp := *expected[i].Issue.Assignee.Name  // NOTE: this is a big assumption
		pre := *predicted[i].Issue.Assignee.Name // NOTE: this is a big assumption
		if _, ok := outputMatrix[exp]; ok {
			outputMatrix[exp][pre] += 1
		} else {
			outputMatrix[exp] = make(map[string]int)
			outputMatrix[exp][pre] = 1
		}
	}
	return outputMatrix, nil
}

func (m matrix) getClassTP(class string) float64 {
	return float64(m[class][class])
}

func (m matrix) getClassTN(class string) float64 {
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

func (m matrix) getClassFP(class string) float64 {
	count := 0.0
	for columnHead := range m {
		if columnHead == class {
			continue
		}
		count += float64(m[columnHead][class])
	}
	return count
}

func (m matrix) getClassFN(class string) float64 {
	count := 0.0
	for rowHead := range m[class] {
		if rowHead == class {
			continue
		}
		count += float64(m[class][rowHead])
	}
	return count
}

func (m matrix) getPrecision(class string) float64 {
	classTP := m.getClassTP(class)
	classFP := m.getClassFP(class)
	return Round(classTP / (classTP + classFP))
}

func (m matrix) getRecall(class string) float64 {
	classTP := m.getClassTP(class)
	classFN := m.getClassFN(class)
	return Round(classTP / (classTP + classFN))
}

func (m matrix) getAccuracy() float64 {
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
	return Round(float64(correct) / float64(total))
}

func (m matrix) getTestCount() float64 {
	count := 0.0
	for columnHead := range m {
		for rowHead := range m[columnHead] {
			count += float64(m[columnHead][rowHead])
		}
	}
	return count
}

func (m matrix) fillMatrix() matrix {
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

func (m matrix) ClassSummary(class string) string {
	input := []string{"SUMMARY RESULTS FOR CLASS: ", class, "\n",
		"TRUE POSITIVES:  ", ToString(m.getClassTP(class)), "\n",
		"TRUE NEGATIVES:  ", ToString(m.getClassTN(class)), "\n",
		"FALSE POSITIVES: ", ToString(m.getClassFP(class)), "\n",
		"FALSE NEGATIVES: ", ToString(m.getClassFN(class)), "\n",
		"PRECISION:       ", ToString(m.getPrecision(class)), "\n",
		"RECALL:          ", ToString(m.getRecall(class)), "\n"}
	output := strings.Join(input, "")
	return output
}

func (m matrix) FullSummary() string {
	input := []string{"SUMMARY RESULTS FOR FULL MATRIX\n",
		"TOTAL TESTS:    ", ToString(m.getTestCount()), "\n",
		"TOTAL ACCURACY: ", ToString(m.getAccuracy()), "\n"}
	output := strings.Join(input, " ")
	return output
}
