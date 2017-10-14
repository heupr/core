package gateway

import (
	"core/utils"
	"encoding/gob"
	"os"
)

type DiskCache struct {
}

func (d *DiskCache) Set(key string, values interface{}) (err error) {
	file, err := os.OpenFile(utils.Config.DataCachesPath+key, os.O_WRONLY|os.O_CREATE, 0644)
	defer file.Close()
	if err != nil {
		return err
	}
	enc := gob.NewEncoder(file)
	err = enc.Encode(values)
	return
}

func (d *DiskCache) TryGet(key string, values interface{}) (err error) {
	if _, err = os.Stat(utils.Config.DataCachesPath + key); err == nil {
		return d.Get(key, values)
	}
	return
}

func (d *DiskCache) Get(key string, values interface{}) (err error) {
	file, err := os.Open(utils.Config.DataCachesPath + key)
	defer file.Close()
	if err != nil {
		return err
	}
	dec := gob.NewDecoder(file)
	err = dec.Decode(values)
	return
}
