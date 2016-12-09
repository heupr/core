package gateway

import (
	"encoding/gob"
	"os"
)

type DiskCache struct {
}

func (d *DiskCache) Set(key string, values interface{}) (err error) {
	file, err := os.OpenFile("../../data/caches/"+key, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	enc := gob.NewEncoder(file)
	err = enc.Encode(values)
	return
}

func (d *DiskCache) TryGet(key string, values interface{}) (err error) {
	if _, err = os.Stat(key); err == nil {
		return d.Get(key, values)
	}
	return
}

func (d *DiskCache) Get(key string, values interface{}) (err error) {
	file, err := os.Open(key)
	if err != nil {
		return err
	}
	dec := gob.NewDecoder(file)
	err = dec.Decode(values)
	return
}
