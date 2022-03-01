package smt

import (
	"encoding/json"
	"fmt"
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
	if res, err := s.red.HGet(key.NameSpace, keyHash).Result(); err != nil {
		return nil, err
	} else {
		var list []RedisBranchNode
		if err = json.Unmarshal([]byte(res), &list); err != nil {
			return nil, err
		} else {
			var node BranchNode
			node.Left = RedisBranchNodeToMergeValue(list[0])
			node.Right = RedisBranchNodeToMergeValue(list[1])
			return &node, nil
		}
	}
	return nil, StoreErrorNotExistBranch
}

type RedisBranchNode struct {
	Value     H256
	BaseNode  H256
	ZeroBits  H256
	ZeroCount byte
}

func RedisBranchNodeToMergeValue(node RedisBranchNode) MergeValue {
	if node.Value != nil {
		return &MergeValueH256{Value: node.Value}
	} else {
		return &MergeValueZero{
			BaseNode:  node.BaseNode,
			ZeroBits:  node.ZeroBits,
			ZeroCount: node.ZeroCount,
		}
	}
}

func MergeValueToRedisBranchNode(data MergeValue) RedisBranchNode {
	var res RedisBranchNode
	if v, ok := (data).(*MergeValueH256); ok {
		res.Value = v.Value
	} else if z, ok := (data).(*MergeValueZero); ok {
		res.ZeroBits = z.ZeroBits
		res.BaseNode = z.BaseNode
		res.ZeroCount = z.ZeroCount
	}
	return res
}

func (s *RedisStore) InsertBranch(key BranchKey, node BranchNode) error {
	keyHash := key.GetHash()
	var list []RedisBranchNode
	list = append(list, MergeValueToRedisBranchNode(node.Left))
	list = append(list, MergeValueToRedisBranchNode(node.Right))
	data, err := json.Marshal(&list)
	if err != nil {
		return err
	}
	if _, err := s.red.HSet(key.NameSpace, keyHash, data).Result(); err != nil {
		return err
	}
	return nil
}

func (s *RedisStore) RemoveBranch(key BranchKey) error {
	fmt.Println("RemoveBranch")
	keyHash := key.GetHash()
	if _, err := s.red.HDel(key.NameSpace, keyHash).Result(); err != nil {
		return err
	}
	return nil
}
