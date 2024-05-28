package cache

import (
	"encoding/json"
	"errors"
	"fmt"

	"time"

	"github.com/patrickmn/go-cache"
)

type Cache interface {
	SetStruct(string, interface{}) error
	GetStruct(string, interface{}) error
	DeleteStruct(string)
}

type service struct {
	cache *cache.Cache
}

var ErrMissingKey = errors.New("missing key")

func New() Cache {
	s := service{
		cache: cache.New(5*time.Minute, 10*time.Minute),
	}
	return &s
}

func (s *service) SetStruct(key string, value interface{}) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	s.cache.Set(key, jsonData, cache.DefaultExpiration)
	return nil
}

func (s *service) GetStruct(key string, dest interface{}) error {
	data, found := s.cache.Get(key)
	if !found {
		return ErrMissingKey
	}
	byteData, ok := data.([]byte)
	if !ok {
		return fmt.Errorf("data not in byte format")
	}
	err := json.Unmarshal(byteData, dest)
	return err
}

func (s *service) DeleteStruct(key string) {
	s.cache.Delete(key)
}
