package smt

import (
	"encoding/json"
	"github.com/go-redis/redis"
)

type RedisStore struct {
	smtName string
	red     *redis.Client
}

func NewRedisStore(smtName string, red *redis.Client) *RedisStore {
	return &RedisStore{smtName: smtName, red: red}
}

func (s *RedisStore) UpdateRoot(root H256) error {
	return s.red.HSet(s.smtName, "root", root.String()).Err()
}

func (s *RedisStore) Root() (H256, error) {
	if res, err := s.red.HGet(s.smtName, "root").Result(); err != nil {
		if err == redis.Nil {
			return H256Zero(), nil
		}
		return nil, err
	} else {
		return Hex2Bytes(res), nil
	}
}

func (s *RedisStore) GetBranch(key BranchKey) (*BranchNode, error) {
	keyHash := key.GetHash()
	if res, err := s.red.HGet(s.smtName, keyHash).Result(); err != nil {
		if err == redis.Nil {
			return nil, StoreErrorNotExist
		}
		return nil, err
	} else {
		var node BranchNode
		if err = json.Unmarshal([]byte(res), &node); err != nil {
			return nil, err
		} else {
			return &node, nil
		}
	}
}

func (s *RedisStore) InsertBranch(key BranchKey, node BranchNode) error {
	keyHash := key.GetHash()
	data, err := json.Marshal(&node)
	if err != nil {
		return err
	}
	if _, err := s.red.HSet(s.smtName, keyHash, data).Result(); err != nil {
		return err
	}
	return nil
}

func (s *RedisStore) RemoveBranch(key BranchKey) error {
	keyHash := key.GetHash()
	if _, err := s.red.HDel(s.smtName, keyHash).Result(); err != nil {
		return err
	}
	return nil
}
