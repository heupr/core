package models

type Model struct {
    Algorithm   Algorithm
}

// NOTE: Algorithm interface may ultimately take ExpandedIssue as its methods
//       parameter type; this will necessitate an import from engine/

type Algorithm interface {
    Learn(input ...interface{})
    Predict(input ...interface{}) []string
}

func (m *Model) Learn(input ...interface{}) {
    m.Algorithm.Learn(input)
}

func (m *Model) Predict(input ...interface{}) []string {
    return m.Algorithm.Predict(input)
}

// nbModel := bhattacharya.Model{Classifier: &bhattacharya.NBClassifier{Logger: &logger}, Logger: &logger}
// nbModel := models.Model{Algorithm: &bhattacharya.NBClassifier{Logger: &logger}}

/*
So fold and confuse can be new methods on the new Model struct
Those methods on the new model struct would call the confuse or fold class
Toss... is currently tightly coupled with battacharya and that cannot be moved out as easily
So on the new Model struct you will also need FoldImpl and ConfuseImpl fields..... Basically confuse and fold will need to be plugged in.
They could both be fields inside a new Grades Field/Type

func (m *Model) Fold () {
    m.Grades.Fold(...) etc
}

type interface Algorithm {
    train
    predict
}
*/
