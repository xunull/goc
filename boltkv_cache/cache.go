package boltkv_cache

import (
	"github.com/xunull/goc/boltkv"
	"time"
)

type CacheElement struct {
	Data   []byte
	Key    string
	Option option `json:"option"`
}

type CacheStore struct {
	boltkv.BoltStore
	Default option
}

func NewCacheStore(store boltkv.BoltStore, defaults ...Option) *CacheStore {
	cs := &CacheStore{}
	cs.BoltStore = store

	op := getOption(defaults...)
	cs.Default = *op
	return cs
}

func (s CacheStore) SetCache(k string, v []byte, opts ...Option) error {
	op := getOption(opts...)
	ce := &CacheElement{
		Key:  k,
		Data: v,
	}

	op.CreateAt = time.Now()

	if op.ExpireDuration != 0 {
		op.ExpireAt = op.CreateAt.Add(op.ExpireDuration)
	}
	if op.ExpireAt.IsZero() {
		if s.Default.ExpireDuration != 0 {
			op.ExpireAt = op.CreateAt.Add(s.Default.ExpireDuration)
		}
	}

	ce.Option = *op
	err := s.SetUseGob(k, ce)
	return err
}

func (s CacheStore) GetCache(k string) ([]byte, error) {
	var ce CacheElement
	err := s.GetUseGob(k, &ce)
	if err != nil {
		return nil, err
	}

	if !ce.Option.ExpireAt.IsZero() && ce.Option.ExpireAt.Before(time.Now()) {
		return nil, boltkv.TargetKeyExpiredError
	}

	ce.Option.VisitAt = time.Now()
	err = s.SetUseGob(k, &ce)
	if err != nil {
		return nil, err
	} else {
		return ce.Data, nil
	}
}
