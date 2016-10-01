package bhattacharya

import (
	"testing"
	"coralreef-ci/models/issues"
)

func TestModel(t *testing.T) {
	v := 1
	if v != 1 {
		t.Error(
			"For", v,
			"expected", 1,
			"got", v,
		)
	}
}

func CreateTrainingIssues() []issues.Issue {
	return []issues.Issue{
		issues.Issue { Body: "parallel test body", Assignee: "Mike"},
		issues.Issue { Body: "fxcore test body and", Assignee: "John"},
		issues.Issue { Body: "fxcore test body or", Assignee: "John"}}
}

func CreateValidationIssues() []issues.Issue {
	return []issues.Issue{
		issues.Issue { Body: "parallel test code", Assignee: "Mike"},
		issues.Issue { Body : "fxcore test code and", Assignee: "John"},
		issues.Issue { Body: "fxcore test code or", Assignee: "John"}}
	}

func CreateUnassignedIssues() []issues.Issue {
	return []issues.Issue{
		issues.Issue { Body: "parallel test code"},
		issues.Issue { Body : "fxcore test code and"},
		issues.Issue { Body: "fxcore test code or"}}
}

func TestLearn(t *testing.T) {
	nbModel := Model{ classifier: &NbClassifer{}}
	trainingSet := CreateTrainingIssues()
	validationSet := CreateValidationIssues()

	nbModel.Learn(trainingSet)

	for i := 0; i < len(validationSet); i++ {
		input := validationSet[i].Body
		expected := validationSet[i].Assignee
		actual := nbModel.Predict(validationSet[i])
		Assert(t, expected, actual, input)
	}
}

func TestStopWords(t *testing.T) {
	output := CreateTrainingIssues()
	originalInput := CreateTrainingIssues()
	removeStopWords(output)

	Assert(t, "parallel test body", output[0].Body, originalInput[0].Body)
	Assert(t, "fxcore test body", output[1].Body, originalInput[1].Body)
	Assert(t, "fxcore test body", output[2].Body, originalInput[2].Body)
}

func Assert(t *testing.T, expected string, actual string, input string) {
	if (actual != expected) {
		t.Error(
			"\nFOR     : ", input,
			"\nEXPECTED: ", expected,
			"\nGOT     : ", actual,
		)
	}
}
