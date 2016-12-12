package models

import (
    "coralreefci/models/confuse"
    "coralreefci/models/fold"
)

type Model struct {
    Algorithm   Algorithm
    Utilities Utilities
}

type Algorithm interface {
    Learn(input interface{})
    Predict(input interface{}) []string
}

type Utilities struct {
    Fold    fold.Fold
    Confuse confuse.Confuse
}





/*
Model{
    Algorithm: bhattacharya.NBClassifier
    - Learn
    - Predict
    Utilities: Utilities
}
*/

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
