package models

import (
	"testing"
)

var exp = []string{"John", "Mike", "Woz", "John", "Mike", "Woz", "John", "Mike", "Woz"}
var pre = []string{"John", "John", "Mike", "Woz", "Woz", "Mike", "Mike", "Mike", "John"}

var metrics = map[string]float64{
	"MikeTP":    1.0,
	"MikeTN":    3.0,
	"MikeFP":    3.0,
	"MikeFN":    2.0,
	"FullCount": 9.0,
	"Precision": 0.25,
	"Recall":    0.3333,
	"Accuracy":  0.2222,
}

func TestBuildMatrix(t *testing.T) {
	nbModel := Model{}
	matrix, _, _ := nbModel.BuildMatrix(exp, pre)

	if len(matrix) == 0 {
		t.Error(
			"\nEMPTY MATRIX",
			"\nCONTENTS", matrix)
	}

	countTP := matrix.getClassTP("Mike")
	if metrics["MikeTP"] != countTP {
		t.Error(
			"\nCLASS TRUE POSITIVE MISCOUNT",
			"\nEXPECTED:  ", metrics["MikeTP"],
			"\nACTUAL:    ", countTP)
	}

	countNP := matrix.getClassTN("Mike")
	if metrics["MikeTN"] != countNP {
		t.Error(
			"\nCLASS TRUE NEGATIVE MISCOUNT",
			"\nEXPECTED:  ", metrics["MikeTN"],
			"\nACTUAL:    ", countNP)
	}

	countFP := matrix.getClassFP("Mike")
	if metrics["MikeFP"] != countFP {
		t.Error(
			"\nCLASS FALSE POSITIVE MISCOUNT",
			"\nEXPECTED:  ", metrics["MikeFP"],
			"\nACTUAL:    ", countFP)
	}

	countFN := matrix.getClassFN("Mike")
	if metrics["MikeFN"] != countFN {
		t.Error(
			"\nCLASS FALSE NEGATIVE MISCOUNT",
			"\nEXPECTED:  ", metrics["MikeFN"],
			"\nACTUAL:    ", countFN)
	}

	classPrecision := matrix.getPrecision("Mike")
	if metrics["Precision"] != classPrecision {
		t.Error(
			"\nCLASS PRECISION MISCALCULATED",
			"\nEXPECTED:  ", metrics["Precision"],
			"\nACTUAL:    ", classPrecision)
	}

	classRecall := matrix.getRecall("Mike")
	if metrics["Recall"] != classRecall {
		t.Error(
			"\nCLASS RECALL MISCALCULATED",
			"\nEXPECTED:  ", metrics["Recall"],
			"\nACTUAL:    ", classRecall)
	}

	fullAccuracy := matrix.getAccuracy()
	if metrics["Accuracy"] != fullAccuracy {
		t.Error(
			"\nFULL MATRIX ACCURACY MISCALCULATED",
			"\nEXPECTED:  ", metrics["Accuracy"],
			"\nACTUAL:    ", fullAccuracy)
	}

	fullCount := matrix.getTestCount()
	if metrics["FullCount"] != fullCount {
		t.Error(
			"\nFULL MATRIX RESULT COUNT MISCALCULATED",
			"\nEXPECTED:  ", metrics["FullCount"],
			"\nACTUAL:    ", fullCount)
	}

	fullMatrix := matrix.fillMatrix()
	for key := range fullMatrix {
		if len(fullMatrix) != len(fullMatrix[key]) {
			t.Error(
				"\nMATRIX IS NOT EQUAL IN DIMENSIONS",
				"\nEXPECTED LENGTH:  ", len(fullMatrix),
				"\nACTUAL LENGTH:    ", len(fullMatrix[key]))
		}
	}

	classOutput := fullMatrix.ClassSummary("John")
	if classOutput == "" {
		t.Error(
			"\nNO OUTPUT STRING FROM CLASS SUMMARY",
		)
	}

	fullOutput := fullMatrix.FullSummary()
	if fullOutput == "" {
		t.Error(
			"\nNO OUTPUT STRING FROM FULL SUMMARY",
		)
	}
}
