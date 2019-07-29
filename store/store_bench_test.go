package store

import (
	"math/rand"
	"os"
	"testing"
)

func randBytes(length int) []byte {
	key := make([]byte, length)
	rand.Read(key)
	return key
}

func BenchmarkInserts(b *testing.B) {
	// Setup
	defer func() {
		os.RemoveAll("mstate.db")
	}()
	st := NewStateStore(".")
	dcache := NewCache(st)
	b.ResetTimer()

	for i := 1; i <= b.N; i++ {
		dcache.Set(randBytes(32), randBytes(10000))

	}
	dcache.ApplyToState()
}
