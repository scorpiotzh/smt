package smt

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongodbStore struct {
	ctx    context.Context
	client *mongo.Client
}

func NewMongoStore(ctx context.Context, client *mongo.Client) *MongodbStore {
	return &MongodbStore{ctx: ctx, client: client}
}

type MongoBranchNode struct {
	LeftValue      H256
	LeftBaseNode   H256
	LeftZeroBits   H256
	LeftZeroCount  byte
	RightValue     H256
	RightBaseNode  H256
	RightZeroBits  H256
	RightZeroCount byte
}

func MongoBranchNodeToMergeValue(node MongoBranchNode) BranchNode {
	var data BranchNode
	if node.LeftValue != nil {
		data.Left = &MergeValueH256{Value: node.LeftValue}
	} else {
		data.Left = &MergeValueZero{
			BaseNode:  node.LeftBaseNode,
			ZeroBits:  node.LeftZeroBits,
			ZeroCount: node.LeftZeroCount,
		}
	}
	if node.RightValue != nil {
		data.Right = &MergeValueH256{Value: node.RightValue}
	} else {
		data.Right = &MergeValueZero{
			BaseNode:  node.RightBaseNode,
			ZeroBits:  node.RightZeroBits,
			ZeroCount: node.RightZeroCount,
		}
	}
	return data
}

func MergeValueToMongoBranchNode(data BranchNode) MongoBranchNode {
	var res MongoBranchNode
	if v, ok := (data.Left).(*MergeValueH256); ok {
		res.LeftValue = v.Value
	} else if z, ok := (data.Left).(*MergeValueZero); ok {
		res.LeftZeroBits = z.ZeroBits
		res.LeftBaseNode = z.BaseNode
		res.LeftZeroCount = z.ZeroCount
	}
	if v, ok := (data.Right).(*MergeValueH256); ok {
		res.RightValue = v.Value
	} else if z, ok := (data.Right).(*MergeValueZero); ok {
		res.RightZeroBits = z.ZeroBits
		res.RightBaseNode = z.BaseNode
		res.RightZeroCount = z.ZeroCount
	}

	return res
}

var countNum = 0

func (m *MongodbStore) GetBranch(key BranchKey) (*BranchNode, error) {
	keyHash := key.GetHash()
	var data MongoBranchNode
	collection := m.client.Database("smt").Collection(key.NameSpace)
	err := collection.FindOne(m.ctx, bson.M{"_id": keyHash}).Decode(&data)
	if err != nil {
		return nil, err
	}
	node := MongoBranchNodeToMergeValue(data)
	if key.Height == 255 {
		countNum++
		fmt.Println("GetBranch:", countNum, time.Now().String())
	}
	return &node, nil
}

func (m *MongodbStore) InsertBranch(key BranchKey, node BranchNode) error {
	keyHash := key.GetHash()
	newNode := MergeValueToMongoBranchNode(node)
	collection := m.client.Database("smt").Collection(key.NameSpace)
	update := bson.M{"$set": newNode}
	updateOpts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(context.Background(), bson.M{"_id": keyHash}, update, updateOpts)
	if err != nil {
		return err
	}
	return nil
}

func (m *MongodbStore) RemoveBranch(key BranchKey) error {
	keyHash := key.GetHash()
	collection := m.client.Database("smt").Collection(key.NameSpace)
	_, err := collection.DeleteOne(m.ctx, bson.M{"_id": keyHash})
	if err != nil {
		return err
	}
	return nil
}
