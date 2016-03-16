package Utils

import (
	"math/rand"
	"crypto/sha256"
)

func GenSalt() []byte {
	result := make([]byte, 10)
	for i := 0; i < 10; i++ {
		result[i] = (byte)(rand.Uint32() % 256)
	}
	return result[:]
}

func Sha256(bs []byte) []byte {
	hasher := sha256.New()
	hasher.Write(bs)
	return hasher.Sum(nil)
}