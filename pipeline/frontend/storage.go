package frontend

import (
	"fmt"
	"strconv"

	"github.com/boltdb/bolt"
)

// Buckets thus far are "hook" and "token" with keys named "repoID".
type BoltDB struct {
	DB *bolt.DB
}

func (b *BoltDB) Initialize() error {
	err := b.DB.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte("token")); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists([]byte("hook")); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (b *BoltDB) Store(name string, key int, value []byte) error {
	byteName := []byte(name)
	byteKey := []byte(strconv.Itoa(key))
	err := b.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(byteName)
		if err != nil {
			return err
		}

		err = bucket.Put(byteKey, value)
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

func (b *BoltDB) Retrieve(name string, key int) ([]byte, error) {
	byteName := []byte(name)
	byteKey := []byte(strconv.Itoa(key))
	byteValue := []byte{}

	err := b.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(byteName)
		if bucket == nil {
			return fmt.Errorf("Bucket %v not found in database", name)
		}
		byteValue = bucket.Get(byteKey)
		if byteValue == nil {
			return fmt.Errorf("Repo %v not found in database for bucket %v", key, name)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return byteValue, nil
}

func (b *BoltDB) RetrieveBulk(name string) ([][]byte, [][]byte, error) {
	keys := [][]byte{}
	values := [][]byte{}

	if err := b.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(name))
		bucket.ForEach(func(key, value []byte) error {
			keys = append(keys, key)
			values = append(values, value)
			return nil
		})
		return nil
	}); err != nil {
		return nil, nil, err
	}
	return keys, values, nil
}

func (b *BoltDB) Delete(name string, key int) error {
	byteName := []byte(name)
	byteKey := []byte(strconv.Itoa(key))

	if err := b.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(byteName)
		if err := bucket.Delete(byteKey); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
