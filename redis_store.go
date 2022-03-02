package smt

import (
	"encoding/json"
	"github.com/go-redis/redis"
)

type RedisStore struct {
	red *redis.Client
}

func NewRedisStore(red *redis.Client) *RedisStore {
	return &RedisStore{red: red}
}

func (s *RedisStore) GetBranch(key BranchKey) (*BranchNode, error) {
	keyHash := key.GetHash()
	if res, err := s.red.HGet(key.SmtName, keyHash).Result(); err != nil {
		return nil, err
	} else {
		var node BranchNode
		if err = json.Unmarshal([]byte(res), &node); err != nil {
			return nil, err
		} else {
			return &node, nil
		}
	}
	return nil, StoreErrorNotExistBranch
}

func (s *RedisStore) InsertBranch(key BranchKey, node BranchNode) error {
	keyHash := key.GetHash()
	data, err := json.Marshal(&node)
	if err != nil {
		return err
	}
	if _, err := s.red.HSet(key.SmtName, keyHash, data).Result(); err != nil {
		return err
	}
	return nil
}

func (s *RedisStore) RemoveBranch(key BranchKey) error {
	keyHash := key.GetHash()
	if _, err := s.red.HDel(key.SmtName, keyHash).Result(); err != nil {
		return err
	}
	return nil
}
