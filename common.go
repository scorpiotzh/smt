package smt

import "encoding/hex"

const (
	ByteSize     = 8
	MaxU8        = 255
	MaxStackSize = 257
)

const (
	MergeNormal            byte = 1
	MergeZeros             byte = 2
	PersonSparseMerkleTree      = "sparsemerkletree"
)

func Has0xPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

func Hex2Bytes(s string) []byte {
	if Has0xPrefix(s) {
		s = s[2:]
	}
	if len(s)%2 == 1 {
		s = "0" + s
	}
	h, _ := hex.DecodeString(s)
	return h
}

func Bytes2Hex(b []byte) string {
	h := hex.EncodeToString(b)
	if len(h) == 0 {
		h = "0"
	}
	return "0x" + h
}
