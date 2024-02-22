/*
Package utils is used for generic methods and values that are used across all packages
*/
package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash/adler32"
	"math"
	"strconv"

	"github.com/libsv/go-bt/v2"
)

const (
	// XpubKeyLength is the length of an xPub string key
	XpubKeyLength = 111

	// ChainInternal internal chain num
	ChainInternal = uint32(1)

	// ChainExternal external chain num
	ChainExternal = uint32(0)

	// MaxInt32 max integer for int32
	MaxInt32 = int64(1<<(32-1) - 1)
)

// Hash will generate a hash of the given string (used for xPub:hash)
func Hash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// RandomHex returns a random hex string and error
func RandomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// GetChildNumsFromHex get an array of uint32 numbers from the hex string
func GetChildNumsFromHex(hexHash string) ([]uint32, error) {
	strLen := len(hexHash)
	size := 8
	splitLength := int(math.Ceil(float64(strLen) / float64(size)))
	childNums := make([]uint32, 0)
	for i := 0; i < splitLength; i++ {
		start := i * size
		stop := start + size
		if stop > strLen {
			stop = strLen
		}
		num, err := strconv.ParseInt(hexHash[start:stop], 16, 64)
		if err != nil {
			return nil, err
		}
		if num > MaxInt32 {
			num = num - MaxInt32
		}
		childNums = append(childNums, uint32(num)) // todo: re-work to remove casting (possible cutoff)
	}

	return childNums, nil
}

// StringInSlice check whether the string already is in the slice
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// GetTransactionIDFromHex get the transaction ID from the given transaction hex
func GetTransactionIDFromHex(hex string) (string, error) {
	parsedTx, err := bt.NewTxFromString(hex)
	if err != nil {
		return "", err
	}
	return parsedTx.TxID(), nil
}

// LittleEndianBytes64 returns a byte array in little endian from an unsigned integer of 64 bytes.
func LittleEndianBytes64(value uint64, resultLength uint32) []byte {
	buf := make([]byte, resultLength)
	binary.LittleEndian.PutUint64(buf, value)

	return buf
}

// HashAdler32 returns computed string calculated with Adler32 function.
func HashAdler32(input string) (string, error) {
	if input == "" {
		return "", fmt.Errorf("input string is empty - cannot apply adler32 hash function")
	}
	data := []byte(input)
	hasher := adler32.New()
	_, err := hasher.Write(data)
	if err != nil {
		return "", err
	}

	sum := hasher.Sum32()

	return fmt.Sprintf("%08x", sum), nil
}

// SafeAssign - Assigns value (not pointer) the src to dest if src is not nil
func SafeAssign[T any](dest *T, src *T) {
	if src != nil {
		*dest = *src
	}
}
