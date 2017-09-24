package frontend

import (
	"io/ioutil"
	"os"
	"strconv"
	"testing"

	"github.com/boltdb/bolt"
)

func Test_open(t *testing.T) {
	filename := "test-storage.db"
	file, err := ioutil.TempFile("", filename)
	if err != nil {
		t.Errorf("generate storage test file: %v", err)
	}
	file.Close()
	defer os.Remove(filename)

	testDB, err := bolt.Open(filename, 0644, nil)
	if err != nil {
		t.Errorf("error opening test database: %v", err)
	}

	testServer := FrontendServer{Database: BoltDB{DB: testDB}}
	testBucket := "hook"
	testKey := 7
	testValue := []byte(strconv.Itoa(2224))

	t.Run("init", func(t *testing.T) {
		if err := testServer.Database.Initialize(); err != nil {
			t.Errorf("error starting new in-memory database: %v", err)
		}
		err = testServer.Database.DB.Update(func(tx *bolt.Tx) error {
			if hooks := tx.Bucket([]byte("hook")); hooks == nil {
				t.Errorf("test bucket error: %v", err)
			}
			if tokens := tx.Bucket([]byte("token")); tokens == nil {
				t.Errorf("test bucket error: %v", err)
			}
			return nil
		})
	})

	t.Run("store", func(t *testing.T) {
		if err := testServer.Database.Store(testBucket, testKey, testValue); err != nil {
			t.Errorf("error in adding data to database file: %v", err)
		}
	})

	t.Run("retrieve", func(t *testing.T) {
		value, err := testServer.Database.Retrieve(testBucket, testKey)
		if err != nil {
			t.Errorf("error retrieving data from database - expected %v; received %v", testValue, value)
		}
	})

	t.Run("bulk", func(t *testing.T) {
		_, _, err := testServer.Database.RetrieveBulk(testBucket)
		if err != nil {
			t.Errorf("error pulling all data from bucket: %v", err)
		}
	})

	t.Run("delete", func(t *testing.T) {
		err := testServer.Database.Delete(testBucket, testKey)
		if err != nil {
			t.Errorf("error deleting database entry: %v", err)
		}
	})
}
