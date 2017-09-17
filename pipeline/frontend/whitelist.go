package frontend

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/boltdb/bolt"
	"github.com/google/go-github/github"
	// "github.com/spf13/viper"
)

var (
	databaseName = "whitelist.db"
	bucketName   = "active-repos"
	whitelist    = "whitelist.toml"
)

type User struct {
	Name  string
	Repo  string
	Token string
}

type Config struct {
	Title   string
	Maximum int
	user    []User
}

func (fs *FrontendServer) AutomaticWhitelist(repo github.Repository, token []byte) error {
	boltDB, err := bolt.Open(databaseName, 0644, nil)
	if err != nil {
		return err
	}
	defer boltDB.Close()

	boltDB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}
		count := bucket.Stats().KeyN
		config := Config{}
		if _, err := toml.DecodeFile(whitelist, &config); err != nil {
			return err
		}

		if count <= config.Maximum {
			content := []string{*repo.Owner.Login, *repo.Name}
			buf := &bytes.Buffer{}
			gob.NewEncoder(buf).Encode(content)
			info := buf.Bytes()

			err := bucket.Put([]byte(bucketName), info)
			if err != nil {
				return err
			}
		} else {
			e := fmt.Errorf("Maximum allowed beta users reached - see %v to adjust cap", whitelist)
			fmt.Println(e)
			return e
		}
		return nil
	})
	return nil
}

func (fs *FrontendServer) ManualWhitelist() {

}
