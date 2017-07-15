package models

import (
	"errors"
	"fmt"
	"strings"
)

type matrix map[string]map[string]int

func (m *Model) BuildMatrix(expected, predicted []string) (matrix, []string, error) {
	if len(expected) != len(predicted) {
		return nil, nil, errors.New("Input slices are not equal length; expected: " + string(len(expected)) + ", predicted: " + string(len(predicted)))
	}

	all := append(expected, predicted...)
	distinctAssignees := []string{}
	j := 0
	for i := 0; i < len(all); i++ {
		for j = 0; j < len(distinctAssignees); j++ {
			if all[i] == distinctAssignees[j] {
				break
			}
		}
		if j == len(distinctAssignees) {
			distinctAssignees = append(distinctAssignees, all[i])
		}
	}

	outputMatrix := make(map[string]map[string]int)

	for i := 0; i < len(expected); i++ {
		exp := expected[i]
		pre := predicted[i]
		if _, ok := outputMatrix[exp]; ok {
			outputMatrix[exp][pre] += 1
		} else {
			outputMatrix[exp] = make(map[string]int)
			outputMatrix[exp][pre] = 1
		}
	}
	return outputMatrix, distinctAssignees, nil
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

func (m matrix) getClassF1(class string) float64 {
	p := m.getClassTP(class) / (m.getClassTP(class) + m.getClassFP(class))
	r := m.getClassTP(class) / (m.getClassTP(class) + m.getClassFN(class))
	output := (2 * p * r) / (p + r)
	return output
}

func (m matrix) classesEvaluation(classes []string) {
	for i := 0; i < len(classes); i++ {
		fmt.Println("Class:", classes[i], "\n", m.ClassSummary(classes[i]))
		//TODO: Fix
		//utils.ModelLog.Debug("Class: " + classes[i] + "\n" + m.ClassSummary(classes[i]))
	}
}

func (m matrix) ClassSummary(class string) string {
	input := []string{"Summary results for class: ", class, "\n",
		"True positives:  ", ToString(m.getClassTP(class)), "\n",
		"True negatives:  ", ToString(m.getClassTN(class)), "\n",
		"False positives: ", ToString(m.getClassFP(class)), "\n",
		"False negatives: ", ToString(m.getClassFN(class)), "\n",
		"Precision:       ", ToString(m.getPrecision(class)), "\n",
		"Recall:          ", ToString(m.getRecall(class)), "\n",
		"F1 score:        ", ToString(m.getClassF1(class)), "\n",
	}
	output := strings.Join(input, "")
	return output
}

func (m matrix) FullSummary() string {
	input := []string{"Summary results for full matrix\n",
		"Total tests:    ", ToString(m.getTestCount()), "\n",
		"Total accuracy: ", ToString(m.getAccuracy()), "\n",
	}
	output := strings.Join(input, " ")
	return output
}
