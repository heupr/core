package bhattacharya

import (
	"coralreef-ci/models/issues"
	"errors"
	"fmt"
)

type matrix map[string]map[string]int

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
	return float64(int(number*100)) / 100
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
	return round(correct / total)
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

func classSummary(class string, m matrix) {
	fmt.Printf("SUMMARY RESULTS FOR CLASS %d\n", class)
	fmt.Printf("TRUE POSITIVES:   %d\n", getClassTP(class, m))
	fmt.Printf("TRUE NEGATIVES:   %d\n", getClassTN(class, m))
	fmt.Printf("FALSE POSITIVES:  %d\n", getClassFP(class, m))
	fmt.Printf("FALSE NEGATIVES:  %d\n", getClassFN(class, m))
	fmt.Printf("PRECISION:        %d\n", getPrecision(class, m))
	fmt.Printf("RECALL:           %d\n", getRecall(class, m))
}

func fullSummary(m matrix) {
	fmt.Println("SUMMARY RESULTS FOR FULL MATRIX\n")
	fmt.Printf("TOTAL TESTS:     %d\n", getTestCount(m))
	fmt.Printf("TOTAL ACCURACY:  %d\n", getAccuracy(m))
}
