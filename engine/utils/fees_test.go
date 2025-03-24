package utils

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/stretchr/testify/assert"
)

func TestGetInputSizeForType(t *testing.T) {
	t.Parallel()

	t.Run("valid input type", func(t *testing.T) {
		assert.Equal(t, uint64(148), GetInputSizeForType(ScriptTypePubKeyHash))
	})

	t.Run("unknown input type", func(t *testing.T) {
		assert.Equal(t, uint64(500), GetInputSizeForType("unknown"))
	})
}

func TestGetOutputSizeForType(t *testing.T) {
	t.Parallel()

	t.Run("valid output type", func(t *testing.T) {
		assert.Equal(t, uint64(34), GetOutputSize("76a914a7bf13994cb80a6c17ca3624cae128bf1ff4c57b88ac"))
	})

	t.Run("unknown input type", func(t *testing.T) {
		assert.Equal(t, uint64(500), GetOutputSize(""))
	})
}

func TestIsLowerThan(t *testing.T) {
	t.Run("same satoshis, different bytes", func(t *testing.T) {
		one := bsv.FeeUnit{
			Satoshis: 1,
			Bytes:    1000,
		}
		two := bsv.FeeUnit{
			Satoshis: 1,
			Bytes:    20,
		}
		assert.True(t, one.IsLowerThan(&two))
		assert.False(t, two.IsLowerThan(&one))
	})
	t.Run("same bytes, different satoshis", func(t *testing.T) {
		one := bsv.FeeUnit{
			Satoshis: 1,
			Bytes:    20,
		}
		two := bsv.FeeUnit{
			Satoshis: 2,
			Bytes:    20,
		}
		assert.True(t, one.IsLowerThan(&two))
		assert.False(t, two.IsLowerThan(&one))
	})

	t.Run("zero as bytes in denominator", func(t *testing.T) {
		one := bsv.FeeUnit{
			Satoshis: 1,
			Bytes:    0,
		}
		two := bsv.FeeUnit{
			Satoshis: 2,
			Bytes:    0,
		}
		assert.False(t, one.IsLowerThan(&two))
		assert.False(t, two.IsLowerThan(&one))
	})
}

func TestLowestFee(t *testing.T) {
	initTest := func() (feeList []bsv.FeeUnit, defaultFee bsv.FeeUnit) {
		feeList = []bsv.FeeUnit{
			{
				Satoshis: 1,
				Bytes:    20,
			},
			{
				Satoshis: 2,
				Bytes:    20,
			},
			{
				Satoshis: 3,
				Bytes:    20,
			},
		}
		defaultFee = bsv.FeeUnit{
			Satoshis: 4,
			Bytes:    20,
		}
		return
	}

	t.Run("lowest fee among feeList elements, despite defaultValue", func(t *testing.T) {
		feeList, defaultFee := initTest()
		defaultFee.Satoshis = 1
		defaultFee.Bytes = 50
		assert.Equal(t, feeList[0], *LowestFee(feeList, &defaultFee))
	})

	t.Run("lowest fee as first value", func(t *testing.T) {
		feeList, defaultFee := initTest()
		assert.Equal(t, feeList[0], *LowestFee(feeList, &defaultFee))
	})

	t.Run("lowest fee as middle value", func(t *testing.T) {
		feeList, defaultFee := initTest()
		feeList[1].Bytes = 50
		assert.Equal(t, feeList[1], *LowestFee(feeList, &defaultFee))
	})

	t.Run("lowest fee as defaultValue", func(t *testing.T) {
		_, defaultFee := initTest()
		feeList := []bsv.FeeUnit{}
		assert.Equal(t, defaultFee, *LowestFee(feeList, &defaultFee))
	})
}
