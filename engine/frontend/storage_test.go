package frontend

import "testing"

var (
	testBucket = 2187
	testKey    = "Leia"
	testValue  = "Princess of Alderaan"
)

func TeststoreData(t *testing.T) {
	err := storeData(testBucket, testKey, testValue)
	if err != nil {
		t.Error("Error in adding data to database file")
	}
}

func TestretrieveData(t *testing.T) {
	value, err := retrieveData(testBucket, testKey)
	if err != nil {
		t.Errorf(
			"Error retrieving data from database",
			"\nExpected %v; received %v", testValue, value,
		)
	}
}

func TestdeleteData(t *testing.T) {
	err := deleteData(testBucket)
	if err != nil {
		t.Errorf(
			"Error deleting database entry",
			"\n", err,
		)
	}
}
