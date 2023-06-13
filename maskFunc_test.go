package gmask

import (
	"github.com/stretchr/testify/assert"
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
