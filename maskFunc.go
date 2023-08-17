package gmask

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
)

const (
	MaskTypeZero   = "zero"
	MaskTypeChar   = "char"
	MaskTypeRandom = "rand"
	MaskTypeHash   = "hash"
)

var _ MaskAnyFunc = MaskZero

// MaskZero will return the zero value of the given value
//
// Example: `mask:"zero"`
func MaskZero(value any, _ ...string) (any, error) {
	return reflect.Zero(reflect.TypeOf(value)).Interface(), nil
}

var _ MaskStringFunc = MaskCharString()

// MaskCharString will generate a MaskStringFunc which
// will return the maskChar string with given length in tag,
// if not set length then will use the default length 8,
// if set length to -1 the length will equal to original string,
// if not set maskChar, the maskChar will be *
//
// Example: `mask:"char,[length],[maskChar]"`
func MaskCharString(defaultMaskChar ...string) func(value string, arg ...string) (string, error) {
	c := "*"
	if len(defaultMaskChar) > 0 {
		c = defaultMaskChar[0]
	}
	return func(value string, arg ...string) (new string, err error) {
		length := 8
		maskChar := c

		if len(arg) == 0 {
			return strings.Repeat(maskChar, length), nil
		}
		if len(arg) >= 1 {
			if arg[0] == "-1" {
				length = len(value)
			} else if len(arg[0]) != 0 {
				length, err = strconv.Atoi(arg[0])
				if err != nil {
					return "", err
				}
			}
		}
		if len(arg) >= 2 {
			if len(arg[1]) == 1 {
				maskChar = arg[1]
			} else {
				return "", fmt.Errorf("length of maskChar must equal to 1, got %d", len(arg))
			}
		}
		return strings.Repeat(maskChar, length), nil
	}
}

var _ MaskStringFunc = MaskRandString

// MaskRandString returns a random string
// The default length is 8. If set length to -1, the length of the output
// will be equal to the length of the input string
//
// Example: `mask:"rand,[length]"`
func MaskRandString(value string, arg ...string) (new string, err error) {
	length := 8
	if len(arg) >= 1 {
		if arg[0] == "-1" {
			length = len(value)
		} else {
			length, err = strconv.Atoi(arg[0])
			if err != nil {
				return "", err
			}
		}
	}
	return randString(length), nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func randString(length int) string {
	b := make([]byte, length)
	// A randSrc.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := length-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

var _ MaskStringFunc = MaskHashString

// MaskHashString returns the hash of the given string
// supported algorithms: md5, sha1, sha256
// default algorithm is sha256
//
// Example: `mask:"hash,[algorithm]"`
func MaskHashString(value string, arg ...string) (string, error) {
	algorithm := "sha256"
	if len(arg) >= 1 {
		algorithm = arg[0]
	}

	switch algorithm {
	case "md5":
		w := md5.New()
		_, err := io.WriteString(w, value)
		if err != nil {
			return "", err
		}
		return hex.EncodeToString(w.Sum(nil)), nil
	case "sha1":
		w := sha1.New()
		_, err := io.WriteString(w, value)
		if err != nil {
			return "", err
		}
		return hex.EncodeToString(w.Sum(nil)), nil
	case "sha256":
		w := sha256.New()
		_, err := io.WriteString(w, value)
		if err != nil {
			return "", err
		}
		return hex.EncodeToString(w.Sum(nil)), nil
	default:
		return "", fmt.Errorf("%s algorithm not support", algorithm)
	}
}

var _ MaskFloat64Func = MaskRandFloat64

// MaskRandFloat64 returns a new random float64
//
// Example: `mask:"rand,[max],[min],[digit]"`
func MaskRandFloat64(_ float64, arg ...string) (float64, error) {
	var (
		max, min float64 = 1, 0
		digit            = 0
	)

	switch len(arg) {
	case 3:
		digit, _ = strconv.Atoi(arg[2])
		fallthrough
	case 2:
		min, _ = strconv.ParseFloat(arg[1], 64)
		fallthrough
	case 1:
		max, _ = strconv.ParseFloat(arg[0], 64)
	}
	dd := math.Pow10(digit)
	x := float64(int(rand.Float64()*(max-min)*dd + min*dd))
	return x / dd, nil
}

var _ MaskIntFunc = MaskRandInt

// MaskRandInt returns a new random int
//
// Example: `mask:"rand,[max],[min]"`
func MaskRandInt(_ int, arg ...string) (int, error) {
	var (
		max, min = math.MaxInt, 0
	)

	switch len(arg) {
	case 2:
		min, _ = strconv.Atoi(arg[1])
		fallthrough
	case 1:
		max, _ = strconv.Atoi(arg[0])
	}

	return rand.Intn(max-min) + min, nil
}

var _ MaskUintFunc = MaskRandUint

// MaskRandUint returns a new random uint
// if not set max value, the
//
// Example: `mask:"rand,[max],[min]"`
func MaskRandUint(_ uint, arg ...string) (uint, error) {
	var (
		max, min uint64 = math.MaxUint64, 0
	)

	switch len(arg) {
	case 2:
		min, _ = strconv.ParseUint(arg[1], 10, 64)
		fallthrough
	case 1:
		max, _ = strconv.ParseUint(arg[0], 10, 64)
	}

	return uint(rand.Uint64()%(max-min) + min), nil
}
