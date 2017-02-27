package frontend

import (
	"fmt"
	"strconv"

	// "golang.org/x/crypto/bcrypt"
	"github.com/boltdb/bolt"
)

func storeData(repoID int, key string, value interface{}) error {
	db, err := bolt.Open("storage.db", 0644, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	idBytes := []byte(strconv.Itoa(repoID))
	keyBytes := []byte(key)
	valueBytes, err := func(input interface{}) ([]byte, error) {
		switch i := input.(type) {
		case int:
			return []byte(strconv.Itoa(i)), nil
		case string:
			return []byte(i), nil
		default:
			return nil, fmt.Errorf("Invalid input type %v\nAccepted types: int, string", i)
		}
	}(value)
	if err != nil {
		return err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(idBytes)
		if err != nil {
			return err
		}

		err = bucket.Put(keyBytes, valueBytes)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func retrieveData(repoID int, key string) (interface{}, error) {
	db, err := bolt.Open("storage.db", 0644, nil)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	idBytes := []byte(strconv.Itoa(repoID))
	keyBytes := []byte(key)
	valueBytes := []byte{}

	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(idBytes)
		if bucket == nil {
			return fmt.Errorf("Repository %v not found in database", repoID)
		}
		valueBytes = bucket.Get(keyBytes)
		if valueBytes == nil {
			return fmt.Errorf("Key %v not found in database for repository %v", key, repoID)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	result, err := func(input []byte) (interface{}, error) {
		if num, err := strconv.ParseUint(string(input), 10, 8); err == nil {
			return int(num), nil
		} else {
			return string(input), nil
		}
	}(valueBytes)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func deleteData(repoID int) error {
	db, err := bolt.Open("storage.db", 0644, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	idBytes := []byte(strconv.Itoa(repoID))

	err = db.Update(func(tx *bolt.Tx) error {
		if err := tx.DeleteBucket(idBytes); err != nil {
			return err
		}
		if tx.Bucket(idBytes) != nil {
			return fmt.Errorf("Target repo %v not removed", repoID)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
