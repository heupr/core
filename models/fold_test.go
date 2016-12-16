package models

import (
    "coralreefci/models/bhattacharya"
	"coralreefci/models/issues"
	"strconv"
	"testing"
)

var testingIssues = []issues.Issue{
	{Assignee: "Mike", Body: "I hail from the state of Montana"},
	{Assignee: "Woz", Body: "Boys, I'm wild and wooly and rough"},
	{Assignee: "John", Body: "I ride the broncos drink and smoke"},
	{Assignee: "Mike", Body: "Do everything that's tough"},
	{Assignee: "Woz", Body: "I'm always lighthearted and free from care"},
	{Assignee: "John", Body: "to a friend I'm always true"},
	{Assignee: "Mike", Body: "You'll always find me ready to fight"},
	{Assignee: "Woz", Body: "For dear old Sigma Nu!"},
	{Assignee: "Mike", Body: "Who am I sir"},
	{Assignee: "John", Body: "A fraternity man am I"},
	{Assignee: "Woz", Body: "A Sigma Nu sir"},
	{Assignee: "Mike", Body: "And will be until I die"},
	{Assignee: "John", Body: "Hi-rickety whoopty do"},
	{Assignee: "Woz", Body: "What's the matter with Sigma Nu"},
	{Assignee: "Mike", Body: "Hallabaloo, terrikahoo"},
	{Assignee: "John", Body: "All together for Sigma Nu!"},
	{Assignee: "Mike", Body: "The Sigma Nu's they like to live"},
	{Assignee: "Woz", Body: "But when it comes to die"},
	{Assignee: "John", Body: "You'll never hear them groan or moan"},
	{Assignee: "Mike", Body: "You'll never hear them sigh"},
	{Assignee: "Woz", Body: "They march right up to the pearly gate"},
	{Assignee: "John", Body: "You bet your life they do"},
	{Assignee: "Mike", Body: "For at the gate, the meet Saint Pete"},
	{Assignee: "Woz", Body: "And he's a Sigma Nu!"},
}

func TestFold(t *testing.T) {
	nbModel := Model{Algorithm: &bhattacharya.NBClassifier{}}
	result, _ := nbModel.JohnFold(testingIssues)
	number, _ := strconv.ParseFloat(result, 64)
	if number < 0.00 && number > 1.00 {
		t.Error(
			"\nRESULT IS OUTSIDE ACCEPTABLE RANGE - JOHN FOLD",
			"\nEXPECTED BETWEEN 0.00 AND 1.00",
			"\nACTUAL: %f", number,
		)
	}
}
