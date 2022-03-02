package smt

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

func TestSparseMerkleTree(t *testing.T) {
	tree := NewSparseMerkleTree(nil)
	key := Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000")
	value := Hex2Bytes("00ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	_ = tree.Update(key, value)
	fmt.Println(tree.Root())

	key = Hex2Bytes("0100000000000000000000000000000000000000000000000000000000000000")
	value = Hex2Bytes("11ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	_ = tree.Update(key, value)
	fmt.Println(tree.Root())

	key = Hex2Bytes("0200000000000000000000000000000000000000000000000000000000000000")
	value = Hex2Bytes("22ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	_ = tree.Update(key, value)
	fmt.Println(tree.Root())

	key = Hex2Bytes("0300000000000000000000000000000000000000000000000000000000000000")
	value = Hex2Bytes("33ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	_ = tree.Update(key, value)
	fmt.Println(tree.Root())
}

func TestMerkleProof(t *testing.T) {
	tree := NewSparseMerkleTree(nil)
	key := Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000")
	value := Hex2Bytes("00ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	_ = tree.Update(key, value)

	key = Hex2Bytes("0100000000000000000000000000000000000000000000000000000000000000")
	value = Hex2Bytes("11ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	_ = tree.Update(key, value)

	key = Hex2Bytes("0300000000000000000000000000000000000000000000000000000000000000")
	value = Hex2Bytes("33ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	_ = tree.Update(key, value)
	fmt.Println(tree.Root())

	var keys, values []H256
	keys = append(keys, Hex2Bytes("0400000000000000000000000000000000000000000000000000000000000000"))
	values = append(values, H256Zero())
	proof, err := tree.MerkleProof(keys, values)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("proof:", proof.String())
	root, err := tree.Root()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(Verify(root, proof, keys, values))
}

func TestMerge(t *testing.T) {
	nodeKey := Hex2Bytes("0x00000000000000000000000000000000000000000000000000000000000000d0")
	lhs := MergeValue{
		Value:     nil,
		BaseNode:  Hex2Bytes("0x9180ed6242e737f554d3c4f7b8f8f810581d810bcf3c1075070b45a6104d5ff8"),
		ZeroBits:  Hex2Bytes("0x26938181394b731558f2bcc40926fc1e38c3036f249319cf7ba845a4bbd76903"),
		ZeroCount: 251,
	}
	rhs := MergeValueFromZero()
	//rhs := smt.MergeValueZero{
	//	BaseNode:  common.Hex2Bytes("0x9180ed6242e737f554d3c4f7b8f8f810581d810bcf3c1075070b45a6104d5ff8"),
	//	ZeroBits:  common.Hex2Bytes("0x26938181394b731558f2bcc40926fc1e38c3036f249319cf7ba845a4bbd76953"),
	//	ZeroCount: 255,
	//}
	res := Merge(255, nodeKey, lhs, rhs)
	fmt.Println(res.String())
}

func TestMerkleProof2(t *testing.T) {
	tree := NewSparseMerkleTree(nil)
	key1 := Hex2Bytes("0x88e6966ee9d691a6befe0664bb54c7b45bbda274a7cf5fa8cd07f56d94741223")
	value1 := Hex2Bytes("0xe19e9083ca4dbbee50e56c9825eed7fd750c1982f86412275c9efedf3440f83b")
	_ = tree.Update(key1, value1)
	fmt.Println(tree.Root())

	key2 := Hex2Bytes("0x26938181394b731558f2bcc40926fc1e38c3036f249319cf7ba845a4bbd769d3")
	value2 := Hex2Bytes("0x358305052f73809142a8b4c11f7becbef9d15ac718f86e0e7d51a7f1f2383718")
	_ = tree.Update(key2, value2)
	fmt.Println(tree.Root())

	key := Hex2Bytes("0x3b66d16df3f793044f09494c7d3fd540be1a94a3ed4c8a686c595cf144703e64")
	value := Hex2Bytes("0xb0aa5768b4893807d97ca1785b352fb8ec8f9b5521f549b0a285578b8a57ea97")
	_ = tree.Update(key, value)
	fmt.Println(tree.Root())

	key = Hex2Bytes("0x95377d6ba3f39fbfdd93f2fe7bb29ff1a52aa4baf0d1ae86d73f7ac9f5de31df")
	value = Hex2Bytes("0x61e0ef0afc4eaabbe68ee97bd54e629aa2914085ab3676f0f6df63eb233cc07b")
	_ = tree.Update(key, value)
	fmt.Println(tree.Root())

	var keys, values []H256
	k := H256Zero()
	v := H256Zero()
	k[0] = '1'
	keys = append(keys, k)
	values = append(values, v)

	proof, err := tree.MerkleProof(keys, values)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(proof)
	root, err := tree.Root()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(Verify(root, proof, keys, values))
}

func TestSmt(t *testing.T) {
	// 10000 4s
	// 100000 2min
	fmt.Println(time.Now().String())
	tree := NewSparseMerkleTree(nil)
	count := 100
	for i := 0; i < count; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("value-%d", i)
		k := Sha256(key)
		v := Sha256(value)
		_ = tree.Update(k, v)
	}
	fmt.Println(time.Now().String())
	for i := 0; i < count; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("value-%d", i)
		var keys, values []H256
		k1 := Sha256(key)
		keys = append(keys, k1)
		v1 := Sha256(value)
		values = append(values, v1)
		proof, err := tree.MerkleProof(keys, values)
		if err != nil {
			t.Fatal(err)
		}
		root, err := tree.Root()
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(Verify(root, proof, keys, values))
	}
	fmt.Println(time.Now().String())

}

func TestRedisStore(t *testing.T) {
	// 10000 4min 800M
	red := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	s := NewRedisStore("test", red)
	fmt.Println(time.Now().String())
	tree := NewSparseMerkleTree(s)
	count := 100
	//for i := 0; i < count; i++ {
	//	key := fmt.Sprintf("key-%d", i)
	//	value := fmt.Sprintf("value-%d", i)
	//	k := Sha256(key)
	//	v := Sha256(value)
	//	if err := tree.Update(k, v); err != nil {
	//		t.Fatal(err)
	//	}
	//}
	fmt.Println(time.Now().String())
	for i := 0; i < count; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("value-%d", i)
		var keys, values []H256
		k1 := Sha256(key)
		keys = append(keys, k1)
		v1 := Sha256(value)
		values = append(values, v1)
		proof, err := tree.MerkleProof(keys, values)
		if err != nil {
			t.Fatal(err)
		}
		root, err := tree.Root()
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(Verify(root, proof, keys, values))
	}
	fmt.Println(time.Now().String())
}

func TestMongodbStoreDB(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatal(err)
	}
	s := NewMongoStore(ctx, client, "smt", "test")
	collection := s.client.Database("smt").Collection("test")
	if err := collection.Drop(s.ctx); err != nil {
		t.Fatal(err)
	}
	//key := BranchKey{
	//	Height:    0,
	//	NodeKey:   H256Zero(),
	//}
	//node := BranchNode{
	//	Left:  MergeValueFromZero(),
	//	Right: MergeValueFromZero(),
	//}
	//err = s.InsertBranch(key, node)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//res, err := s.GetBranch(key)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println(res)
	//fmt.Println(s.UpdateRoot(H256Zero()))
	//fmt.Println(s.Root())
}

func TestMongodbStore(t *testing.T) {
	// 1000 2min 30M
	// 10000 16min 330M
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatal(err)
	}
	s := NewMongoStore(ctx, client, "smt", "test")
	fmt.Println(time.Now().String())
	tree := NewSparseMerkleTree(s)
	count := 100
	//for i := 0; i < count; i++ {
	//	key := fmt.Sprintf("key-%d", i)
	//	value := fmt.Sprintf("value-%d", i)
	//	k := Sha256(key)
	//	v := Sha256(value)
	//	if err := tree.Update(k, v); err != nil {
	//		t.Fatal(err)
	//	}
	//}
	fmt.Println(time.Now().String())

	for i := 0; i < count; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("value-%d", i)
		var keys, values []H256
		k1 := Sha256(key)
		keys = append(keys, k1)
		v1 := Sha256(value)
		values = append(values, v1)
		proof, err := tree.MerkleProof(keys, values)
		if err != nil {
			t.Fatal(err)
		}
		root, err := tree.Root()
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(Verify(root, proof, keys, values))
	}
}

func TestDelete(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatal(err)
	}
	s := NewMongoStore(ctx, client, "smt", "test")

	//red := redis.NewClient(&redis.Options{
	//	Addr:     "127.0.0.1:6379",
	//	Password: "",
	//	DB:       0,
	//})
	//s := NewRedisStore(red)
	tree := NewSparseMerkleTree(s)

	//tree := NewSparseMerkleTree("", nil)
	key := H256Zero()
	value := H256Zero()
	value[0] = '1'
	err = tree.Update(key, value)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(tree.Root())
	value = H256Zero()
	err = tree.Update(key, value)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(tree.Root())
}
