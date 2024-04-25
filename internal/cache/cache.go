package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

type Cache interface {
	SetStruct(string, interface{}) error
	GetStruct(string, interface{}) error
	DeleteStruct(string) error
}

type service struct {
	redisClient *redis.Client
	ctx         context.Context
}

func New(addr string) Cache {
	ctx := context.Background()

	// Initialize the Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	s := &service{redisClient: redisClient, ctx: ctx}
	return s
}

func (s *service) SetStruct(key string, value interface{}) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return s.redisClient.Set(s.ctx, key, jsonData, 10*time.Minute).Err()
}

func (s *service) GetStruct(key string, dest interface{}) error {
	jsonData, err := s.redisClient.Get(s.ctx, key).Bytes()
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonData, dest)
	return err
}

func (s *service) DeleteStruct(key string) error {
	return s.redisClient.Del(s.ctx, key).Err()
}
