package testutil

import (
	"bytes"
	"crypto/rand"
	"math/big"
)

func RandomString(n int) string {
	const srcStrings = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789/"
	var buf bytes.Buffer
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(srcStrings)-1)))
		if err != nil {
			panic(err)
		}
		buf.WriteByte(srcStrings[num.Int64()])
	}
	return buf.String()
}
