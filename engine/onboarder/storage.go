package onboarder

import (
	"fmt"
	"strconv"

	"github.com/boltdb/bolt"
)

// DOC: Keys thus far are "hookID" and "token" in buckets named "repoID".
type BoltDB struct {
	db *bolt.DB
}

func (b *BoltDB) store(repoID int, key string, value interface{}) error {
	idBytes := []byte(strconv.Itoa(repoID))
	keyBytes := []byte(key)
	valueBytes, err := func(input interface{}) ([]byte, error) {
		switch i := input.(type) {
		case int:
			return []byte(strconv.Itoa(i)), nil
		case string:
			return []byte(i), nil
		default:
			return nil, fmt.Errorf("Invalid input type %v; accepted types: int, string", i)
		}
	}(value)
	if err != nil {
		return err
	}

	err = b.db.Update(func(tx *bolt.Tx) error {
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

func (b *BoltDB) retrieve(repoID int, key string) (interface{}, error) {
	idBytes := []byte(strconv.Itoa(repoID))
	keyBytes := []byte(key)
	valueBytes := []byte{}

	err := b.db.View(func(tx *bolt.Tx) error {
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

func (b *BoltDB) retrieveBulk(input ...string) ([][]byte, error) {
	tokens := [][]byte{}
	if err := b.db.Batch(func(tx *bolt.Tx) error {
		if err := tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			bucket := tx.Bucket(name)
			tokens = append(tokens, bucket.Get([]byte("token")))
			return nil
		}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return tokens, nil
}

// DOC: For when users remove their repos from the Heupr service.
func (b *BoltDB) delete(repoID int) error {
	idBytes := []byte(strconv.Itoa(repoID))

	err := b.db.Update(func(tx *bolt.Tx) error {
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
