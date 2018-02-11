package labelmaker

import (
	"encoding/gob"
	"fmt"
	"hash"
	"hash/fnv"
	"os"

	"core/utils"
)

var fnvHash hash.Hash = fnv.New128a()

func getHash(s string) string {
	fnvHash.Write([]byte(s))
	defer fnvHash.Reset()

	return fmt.Sprintf("%x", fnvHash.Sum(nil))
}

//TODO: Factor out HashedDiskCache to common location
type HashedDiskCache struct {
}

func (h *HashedDiskCache) Set(key string, values interface{}) (err error) {
	hashedKey := getHash(key)
	file, err := os.OpenFile(utils.Config.DataCachesPath+"/"+hashedKey, os.O_WRONLY|os.O_CREATE, 0644)
	defer file.Close()
	if err != nil {
		return err
	}
	enc := gob.NewEncoder(file)
	err = enc.Encode(values)
	return
}

func (h *HashedDiskCache) TryGet(key string, values interface{}) (err error) {
	hashedKey := getHash(key)
	if _, err = os.Stat(utils.Config.DataCachesPath + "/" + hashedKey); err == nil {
		return h.Get(key, values)
	}
	return
}

func (h *HashedDiskCache) Get(key string, values interface{}) (err error) {
	hashedKey := getHash(key)
	file, err := os.Open(utils.Config.DataCachesPath + "/" + hashedKey)
	defer file.Close()
	if err != nil {
		return err
	}
	dec := gob.NewDecoder(file)
	err = dec.Decode(values)
	return
}
