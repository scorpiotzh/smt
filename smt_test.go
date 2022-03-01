package smt

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/nervosnetwork/ckb-sdk-go/crypto/blake2b"
	"testing"
	"time"
)

func TestSparseMerkleTree(t *testing.T) {
	tree := NewSparseMerkleTree("", nil)
	key := Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000")
	value := Hex2Bytes("00ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	_ = tree.Update(key, value)
	fmt.Println("root:", tree.Root.String())

	key = Hex2Bytes("0100000000000000000000000000000000000000000000000000000000000000")
	value = Hex2Bytes("11ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	_ = tree.Update(key, value)
	fmt.Println("root:", tree.Root.String())

	key = Hex2Bytes("0200000000000000000000000000000000000000000000000000000000000000")
	value = Hex2Bytes("22ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	_ = tree.Update(key, value)
	fmt.Println("root:", tree.Root.String())

	key = Hex2Bytes("0300000000000000000000000000000000000000000000000000000000000000")
	value = Hex2Bytes("33ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	_ = tree.Update(key, value)
	fmt.Println("root:", tree.Root.String())
}

func TestMerkleProof(t *testing.T) {
	tree := NewSparseMerkleTree("", nil)
	key := Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000")
	value := Hex2Bytes("00ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	_ = tree.Update(key, value)

	key = Hex2Bytes("0100000000000000000000000000000000000000000000000000000000000000")
	value = Hex2Bytes("11ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	_ = tree.Update(key, value)

	key = Hex2Bytes("0300000000000000000000000000000000000000000000000000000000000000")
	value = Hex2Bytes("33ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	_ = tree.Update(key, value)
	fmt.Println("root:", tree.Root.String())

	var keys, values []H256
	keys = append(keys, Hex2Bytes("0400000000000000000000000000000000000000000000000000000000000000"))
	values = append(values, H256Zero())
	proof, err := tree.MerkleProof(keys, values)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("proof:", proof.String())
	fmt.Println(Verify(tree.Root, proof, keys, values))
}

func TestMerge(t *testing.T) {
	nodeKey := Hex2Bytes("0x00000000000000000000000000000000000000000000000000000000000000d0")
	lhs := MergeValueZero{
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
	res := Merge(255, nodeKey, &lhs, rhs)
	fmt.Println(res.String())
}

func TestMerkleProof2(t *testing.T) {
	tree := NewSparseMerkleTree("", nil)
	key1 := Hex2Bytes("0x88e6966ee9d691a6befe0664bb54c7b45bbda274a7cf5fa8cd07f56d94741223")
	value1 := Hex2Bytes("0xe19e9083ca4dbbee50e56c9825eed7fd750c1982f86412275c9efedf3440f83b")
	_ = tree.Update(key1, value1)
	fmt.Println("root:", tree.Root, tree.Root.String())

	key2 := Hex2Bytes("0x26938181394b731558f2bcc40926fc1e38c3036f249319cf7ba845a4bbd769d3")
	value2 := Hex2Bytes("0x358305052f73809142a8b4c11f7becbef9d15ac718f86e0e7d51a7f1f2383718")
	_ = tree.Update(key2, value2)
	fmt.Println("root:", tree.Root, tree.Root.String())

	key := Hex2Bytes("0x3b66d16df3f793044f09494c7d3fd540be1a94a3ed4c8a686c595cf144703e64")
	value := Hex2Bytes("0xb0aa5768b4893807d97ca1785b352fb8ec8f9b5521f549b0a285578b8a57ea97")
	_ = tree.Update(key, value)
	fmt.Println("root:", tree.Root, tree.Root.String())

	key = Hex2Bytes("0x95377d6ba3f39fbfdd93f2fe7bb29ff1a52aa4baf0d1ae86d73f7ac9f5de31df")
	value = Hex2Bytes("0x61e0ef0afc4eaabbe68ee97bd54e629aa2914085ab3676f0f6df63eb233cc07b")
	_ = tree.Update(key, value)
	fmt.Println("root:", tree.Root, tree.Root.String())

	var keys, values []H256
	k := H256Zero()
	v := H256Zero()
	k[0] = '1'
	keys = append(keys, k)
	values = append(values, v)
	//keys = appen(keys, key2)
	//values = append(values, value2)

	proof, err := tree.MerkleProof(keys, values)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(proof)
	fmt.Println(Verify(tree.Root, proof, keys, values))
}

func TestSmt(t *testing.T) {
	fmt.Println(time.Now().String())
	tree := NewSparseMerkleTree("", nil)
	count := 100
	for i := 0; i < count; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("value-%d", i)
		k, _ := blake2b.Blake256([]byte(key))
		v, _ := blake2b.Blake256([]byte(value))
		//fmt.Println("k:",common.Bytes2Hex(k))
		//fmt.Println("v:",common.Bytes2Hex(v))
		_ = tree.Update(k, v)
	}
	fmt.Println(time.Now().String())
	for i := 0; i < count; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("value-%d", i)
		var keys, values []H256
		k1, _ := blake2b.Blake256([]byte(key))
		keys = append(keys, k1)
		v1, _ := blake2b.Blake256([]byte(value))
		values = append(values, v1)
		proof, err := tree.MerkleProof(keys, values)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(Verify(tree.Root, proof, keys, values))
	}
	fmt.Println(time.Now().String())

}

func TestRedisStore(t *testing.T) {
	// 10000 5min 900M
	red := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	s := NewRedisStore(red)
	fmt.Println(time.Now().String())
	tree := NewSparseMerkleTree("test", s)
	count := 10000
	for i := 0; i < count; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("value-%d", i)
		k, _ := blake2b.Blake256([]byte(key))
		v, _ := blake2b.Blake256([]byte(value))
		//fmt.Println("k:",common.Bytes2Hex(k))
		//fmt.Println("v:",common.Bytes2Hex(v))
		if err := tree.Update(k, v); err != nil {
			t.Fatal(err)
		}
	}
	fmt.Println(time.Now().String())

	var keys, values []H256
	key := fmt.Sprintf("key-%d", 1)
	value := fmt.Sprintf("value-%d", 1)
	k, _ := blake2b.Blake256([]byte(key))
	v, _ := blake2b.Blake256([]byte(value))
	keys = append(keys, k)
	values = append(values, v)

	proof, err := tree.MerkleProof(keys, values)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(proof)
	fmt.Println(Verify(tree.Root, proof, keys, values))
}

func TestBranchNode(t *testing.T) {
	node := BranchNode{
		Left: &MergeValueH256{Value: H256Zero()},
		Right: &MergeValueZero{
			BaseNode:  H256Zero(),
			ZeroBits:  H256Zero(),
			ZeroCount: 0,
		},
	}
	res, err := json.Marshal(&node)
	fmt.Println(string(res), err)
	var data map[string]interface{}
	err = json.Unmarshal(res, &data)
	fmt.Println(err, data)
}
