package onboarder

import "testing"

func Test_open(t *testing.T) {
	testServer := RepoServer{}
	testBucket := 2187
	testKey := "Leia"
	testValue := "Princess of Alderaan"

	defer testServer.CloseDB()
	err := testServer.OpenDB()
	if err != nil {
		t.Errorf("Error opening test database: ", err)
	}

	if err != nil {
		t.Errorf("Error opening new database instance: %v", err)
	}
	t.Run("store", func(t *testing.T) {
		err := testServer.Database.store(testBucket, testKey, testValue)
		if err != nil {
			t.Error("Error in adding data to database file")
		}
	})
	t.Run("retrieve", func(t *testing.T) {
		value, err := testServer.Database.retrieve(testBucket, testKey)
		if err != nil {
			t.Errorf("Error retrieving data from database - expected %v; received %v", testValue, value)
		}
	})
	t.Run("bulk", func(t *testing.T) {
		_, err := testServer.Database.retrieveBulk()
		if err != nil {
			t.Errorf("Error pulling all data from default database: %v", err)
		}
	})
	t.Run("delete", func(t *testing.T) {
		err := testServer.Database.delete(testBucket)
		if err != nil {
			t.Errorf("Error deleting database entry: %v", err)
		}
	})
}
