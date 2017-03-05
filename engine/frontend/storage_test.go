package frontend

import "testing"

// var (
//     testHeuprServer := HeuprServer{}
// 	testDB     = BoltDB{}
// 	testBucket = 2187
// 	testKey    = "Leia"
// 	testValue  = "Princess of Alderaan"
// )

func Test_open(t *testing.T) {
	testServer := HeuprServer{}
	// testDB     := BoltDB{}
	testBucket := 2187
	testKey := "Leia"
	testValue := "Princess of Alderaan"

	defer testServer.closeDB()
	err := testServer.openDB()
	if err != nil {
		t.Error(err) // TODO: Flesh out message
	}

	if err != nil {
		t.Errorf("Error opening new database instance; %v", err)
	}
	t.Run("store", func(t *testing.T) {
		err := testServer.Database.storeData(testBucket, testKey, testValue)
		if err != nil {
			t.Error("Error in adding data to database file")
		}
	})
	t.Run("retrieve", func(t *testing.T) {
		value, err := testServer.Database.retrieveData(testBucket, testKey)
		if err != nil {
			t.Errorf(
				"Error retrieving data from database",
				"\nExpected %v; received %v", testValue, value,
			)
		}
	})
	t.Run("delete", func(t *testing.T) {
		err := testServer.Database.deleteData(testBucket)
		if err != nil {
			t.Errorf(
				"Error deleting database entry",
				"\n", err,
			)
		}
	})
}

// func Test_storeData(t *testing.T) {
// 	err := storeData(testBucket, testKey, testValue)
// 	if err != nil {
// 		t.Error("Error in adding data to database file")
// 	}
// }

// func Test_retrieveData(t *testing.T) {
// 	value, err := retrieveData(testBucket, testKey)
// 	if err != nil {
// 		t.Errorf(
// 			"Error retrieving data from database",
// 			"\nExpected %v; received %v", testValue, value,
// 		)
// 	}
// }

// func Test_deleteData(t *testing.T) {
// 	err := deleteData(testBucket)
// 	if err != nil {
// 		t.Errorf(
// 			"Error deleting database entry",
// 			"\n", err,
// 		)
// 	}
// }
