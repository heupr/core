package syntaxnet

import (
	"fmt"
	"testing"
)

/*
type ConllWord struct {
	Id      int
	Form    string
	Lemma   string
	Upostag string
	Xpostag string
	Feats   string
	Head    int
	Deprel  string
	Deps    string
	Misc    string
}*/

func TestSyntaxNet(t *testing.T) {

	words := SyntaxTree("We need to expose ResourceSet.Table but this would mean exposing Hashtable from corelib .")

	for _, word := range words {
		fmt.Println(word)
		fmt.Println("Id:", word.Id)
		fmt.Println("Form:", word.Form)
		fmt.Println("Lemma:", word.Lemma)
		fmt.Println("Upostag:", word.Upostag)
		fmt.Println("Xpostag:", word.Xpostag)
		fmt.Println("Feats:", word.Feats)
		fmt.Println("Head:", word.Head)
		fmt.Println("Deprel:", word.Deprel)
		fmt.Println("Deps:", word.Deps)
		fmt.Println("Misc:", word.Misc)
		fmt.Println("---------------------------------")
	}
	t.Error("Test Error")
}
