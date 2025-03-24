package utils

import (
	"testing"

	compat "github.com/bitcoin-sv/go-sdk/compat/bip32"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testXPub = "xpub661MyMwAqRbcFrBJbKwBGCB7d3fr2SaAuXGM95BA62X41m6eW2ehRQGW4xLi9wkEXUGnQZYxVVj4PxXnyrLk7jdqvBAs1Qq9gf6ykMvjR7J"
)

func Test_DeriveAddresses(t *testing.T) {
	xPub, errX := compat.NewKeyFromString(testXPub)
	require.NoError(t, errX)

	t.Run("DeriveAddresses 1", func(t *testing.T) {
		external, internal, err := DeriveAddresses(xPub, 1)
		require.NoError(t, err)
		assert.Equal(t, "1PQW54xMn5KA6uK7wgfzN4y7ZXMi6o7Qtm", internal)
		assert.Equal(t, "16fq7PmmXXbFUG5maT5Xvr2zDBUgN1xdMF", external)
	})

	t.Run("DeriveAddresses 2", func(t *testing.T) {
		external, internal, err := DeriveAddresses(xPub, 2)
		require.NoError(t, err)
		assert.Equal(t, "13St2SHkw1b8ZuaExyMf6ZMEzNjYbWRqL4", internal)
		assert.Equal(t, "19jswATg9vBFta1aRnEjPHa2KMwafkmANj", external)
	})
}

func Test_DeriveAddress(t *testing.T) {
	xPub, errX := compat.NewKeyFromString(testXPub)
	require.NoError(t, errX)

	t.Run("DeriveAddresses 1", func(t *testing.T) {
		internal, err := DeriveAddress(xPub, ChainInternal, 1)
		require.NoError(t, err)
		var external string
		external, err = DeriveAddress(xPub, ChainExternal, 1)
		require.NoError(t, err)
		assert.Equal(t, "1PQW54xMn5KA6uK7wgfzN4y7ZXMi6o7Qtm", internal)
		assert.Equal(t, "16fq7PmmXXbFUG5maT5Xvr2zDBUgN1xdMF", external)
	})

	t.Run("DeriveAddresses 2", func(t *testing.T) {
		internal, err := DeriveAddress(xPub, ChainInternal, 2)
		require.NoError(t, err)
		var external string
		external, err = DeriveAddress(xPub, ChainExternal, 2)
		require.NoError(t, err)
		assert.Equal(t, "13St2SHkw1b8ZuaExyMf6ZMEzNjYbWRqL4", internal)
		assert.Equal(t, "19jswATg9vBFta1aRnEjPHa2KMwafkmANj", external)
	})
}

// Benchmark_DeriveAddresses will benchmark the method DeriveAddresses()
func Benchmark_DeriveAddresses(b *testing.B) {
	xPub, errX := compat.NewKeyFromString(testXPub)
	if errX != nil {
		b.Fail()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = DeriveAddresses(xPub, uint32(i))
	}
}

// Benchmark_DeriveAddress will benchmark the method DeriveAddress()
func Benchmark_DeriveAddress(b *testing.B) {
	xPub, errX := compat.NewKeyFromString(testXPub)
	if errX != nil {
		b.Fail()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = DeriveAddress(xPub, ChainInternal, uint32(i))
	}
}
