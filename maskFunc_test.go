package gmask

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestMaskZero(t *testing.T) {
	str := "should be emptyStr"
	emptyStr, err := MaskZero(str)
	assert.NoError(t, err)
	assert.Equal(t, "", emptyStr.(string))
	assert.Equal(t, "should be emptyStr", str)

	num := 57128
	zero, err := MaskZero(num)
	assert.NoError(t, err)
	assert.Equal(t, 0, zero.(int))
	assert.Equal(t, 57128, num)
}

func TestMaskCharString(t *testing.T) {
	originalStr := "abcdef"
	starMasker := MaskCharString()

	// default
	maskedStr, err := starMasker(originalStr)
	assert.NoError(t, err)
	assert.Equal(t, "********", maskedStr)
	assert.Equal(t, "abcdef", originalStr)

	// Set length of masked str to 3
	maskedStr, err = starMasker(originalStr, "3")
	assert.NoError(t, err)
	assert.Equal(t, "***", maskedStr)
	assert.Equal(t, "abcdef", originalStr)

	// Set length of masked str to length of input str
	maskedStr, err = starMasker(originalStr, "-1")
	assert.NoError(t, err)
	assert.Equal(t, "******", maskedStr)
	assert.Equal(t, "abcdef", originalStr)

	// Set mask char to dash
	maskedStr, err = starMasker(originalStr, "", "-")
	assert.NoError(t, err)
	assert.Equal(t, "--------", maskedStr)
	assert.Equal(t, "abcdef", originalStr)

	// Set default mask char to dash
	dashMasker := MaskCharString("-")

	maskedStr, err = dashMasker(originalStr, "", "*")
	assert.NoError(t, err)
	assert.Equal(t, "********", maskedStr)
	assert.Equal(t, "abcdef", originalStr)

	maskedStr, err = dashMasker(originalStr)
	assert.NoError(t, err)
	assert.Equal(t, "--------", maskedStr)
	assert.Equal(t, "abcdef", originalStr)
}

func TestMaskRandString(t *testing.T) {
	originalStr := "abcdef"

	maskedStr, err := MaskRandString(originalStr)
	assert.NoError(t, err)
	assert.Equal(t, "abcdef", originalStr)
	assert.NotEqual(t, originalStr, maskedStr)
	assert.Len(t, maskedStr, 8)

	maskedStr, err = MaskRandString(originalStr, "3")
	assert.NoError(t, err)
	assert.Equal(t, "abcdef", originalStr)
	assert.NotEqual(t, originalStr, maskedStr)
	assert.Len(t, maskedStr, 3)

	maskedStr, err = MaskRandString(originalStr, "-1")
	assert.NoError(t, err)
	assert.Equal(t, "abcdef", originalStr)
	assert.NotEqual(t, originalStr, maskedStr)
	assert.Len(t, maskedStr, len(originalStr))
}

func TestMaskHashString(t *testing.T) {
	originalStr := "abcdef"

	for _, a := range []string{"md5", "sha1", "sha256"} {
		maskedStr, err := MaskHashString(originalStr, a)
		assert.NoError(t, err)
		assert.Equal(t, "abcdef", originalStr)
		assert.Greater(t, len(maskedStr), len(originalStr))
		assert.NotEqual(t, maskedStr, originalStr)
	}
}

func TestMaskRandFloat64(t *testing.T) {
	var origin float64 = -1

	masked, err := MaskRandFloat64(origin)
	assert.NoError(t, err)
	assert.NotEqual(t, masked, origin)
	assert.True(t, 0 <= masked && masked < 1)

	masked, err = MaskRandFloat64(origin, "10")
	assert.NoError(t, err)
	assert.NotEqual(t, masked, origin)
	assert.True(t, 0 <= masked && masked < 10)

	masked, err = MaskRandFloat64(origin, "20", "10")
	assert.NoError(t, err)
	assert.NotEqual(t, masked, origin)
	assert.True(t, 10 <= masked && masked < 20)

	masked, err = MaskRandFloat64(origin, "20", "10", "3")
	assert.NoError(t, err)
	assert.NotEqual(t, masked, origin)
	assert.True(t, 10 <= masked && masked < 20)
	assert.Zero(t, int(masked*math.Pow10(4))%10)
	assert.NotZero(t, int(masked*math.Pow10(3))%10)
}

func TestMaskRandInt(t *testing.T) {
	var origin = -1

	masked, err := MaskRandInt(origin)
	assert.NoError(t, err)
	assert.NotEqual(t, masked, origin)
	assert.True(t, 0 <= masked && masked < math.MaxInt)

	masked, err = MaskRandInt(origin, "10")
	assert.NoError(t, err)
	assert.NotEqual(t, masked, origin)
	assert.True(t, 0 <= masked && masked < 10)

	masked, err = MaskRandInt(origin, "20", "10")
	assert.NoError(t, err)
	assert.NotEqual(t, masked, origin)
	assert.True(t, 10 <= masked && masked < 20)
}

func TestMaskRandUint(t *testing.T) {
	var origin uint = 0

	masked, err := MaskRandUint(origin)
	assert.NoError(t, err)
	assert.NotEqual(t, masked, origin)
	assert.True(t, masked < math.MaxUint)

	masked, err = MaskRandUint(origin, "10")
	assert.NoError(t, err)
	assert.NotEqual(t, masked, origin)
	assert.True(t, masked < 10)

	masked, err = MaskRandUint(origin, "20", "10")
	assert.NoError(t, err)
	assert.NotEqual(t, masked, origin)
	assert.True(t, 10 <= masked && masked < 20)
}
