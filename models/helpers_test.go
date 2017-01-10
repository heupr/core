package models

import "testing"

func TestRound(t *testing.T) {
	roundValue := Round(3.14159265359)
	if roundValue != 3.1416 {
		t.Error(
			"ROUND HELPER FUNCTION FAILING",
			"\n", roundValue,
		)
	}
}

func TestToString(t *testing.T) {
	toStringValue := ToString(3.1416)
	if len(toStringValue) != 6 {
		t.Error(
			"\nTOSTRING HELPER FUNCTION FAILING",
			"\n", toStringValue,
		)
	}
}
