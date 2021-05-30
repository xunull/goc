package boltkv

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	bolt "go.etcd.io/bbolt"
)

type BoltStore struct {
	db         *bolt.DB
	bucketName string
}

func InitBoltStore(name string, dir string) (*BoltStore, error) {
	bs := &BoltStore{}
	bs.bucketName = name

	exist, err := IsDbExist(name, dir)
	if err != nil {
		return nil, err
	}
	if !exist {
		// create
	}
	db, err := bolt.Open(getDbFilePath(name, dir), 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	bs.db = db
	return bs, nil
}

func (s BoltStore) Close() error {
	return s.db.Close()
}

func (s BoltStore) getBucketName() []byte {
	return []byte(s.bucketName)
}

func (s BoltStore) Set(k string, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	err = s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(s.getBucketName())
		return bucket.Put([]byte(k), b)
	})
	if err != nil {
		return err
	}
	return nil
}

func (s BoltStore) Get(k string, v interface{}) error {
	var data []byte
	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(s.getBucketName())
		t := bucket.Get([]byte(k))
		if t != nil {
			data = make([]byte, len(t))
			copy(data, t)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if data == nil {
		return TargetKeyNotFoundError
	}
	return json.Unmarshal(data, v)
}

// ---------------------------------------------------------------------------------------------------------------------

func (s BoltStore) SetUseGob(k string, v interface{}) error {
	var bf bytes.Buffer
	enc := gob.NewEncoder(&bf)

	err := enc.Encode(v)
	if err != nil {
		return err
	}
	err = s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(s.getBucketName())
		return bucket.Put([]byte(k), bf.Bytes())
	})
	if err != nil {
		return err
	}
	return nil
}

func (s BoltStore) GetUseGob(k string, v interface{}) error {
	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(s.getBucketName())
		t := bucket.Get([]byte(k))
		br := bytes.NewReader(t)
		dec := gob.NewDecoder(br)
		if t != nil {
			err := dec.Decode(v)
			return err
		} else {
			return TargetKeyNotFoundError
		}
	})
	if err != nil {
		return err
	}
	return nil
}

func (s BoltStore) Delete(k string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(s.getBucketName())
		return bucket.Delete([]byte(k))
	})
}
