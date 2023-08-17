package gmask

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testStruct struct {
	PubInfo        string `json:"pub"`
	StrChar        string `json:"secretChar" mask:"char"`
	StrChar3       string `json:"secretChar8" mask:"char,3"`
	StrCharSameLen string `json:"secretCharSameLen" mask:"char,-1"`
	StrCharDash    string `json:"secretCharSameStr" mask:"char,,-"`
	StrRand        string `json:"secretRand" mask:"rand"`
	StrHash        string `json:"secretHash" mask:"hash"`

	StrZero   string  `json:"secretZero" mask:"zero"`
	IntZero   int     `json:"intZero" mask:"zero"`
	UintZero  uint    `json:"uintZero" mask:"zero"`
	FloatZero float64 `json:"floatZero" mask:"zero"`
	AnyZero   any     `json:"anyZero" mask:"zero"`
}

func TestExample(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	demo := testStruct{
		PubInfo:        "This field won't be masked",
		StrChar:        "This field will be ********",
		StrChar3:       "This field will be ***",
		StrCharSameLen: "This field will be same length *",
		StrCharDash:    "This field will be --------",
		StrRand:        "This field will be a random string",
		StrHash:        "This field will be a hash string",
		StrZero:        "This field will be a empty string",
		IntZero:        57128, // This field will be 0
		UintZero:       57128, // This field will be 0
		FloatZero:      57128, // This field will be 0
	}
	fmt.Printf("%+v\n", demo)
	demoMasked, _ := Mask(demo)
	fmt.Printf("%+v\n", demoMasked)
}

func TestMasker_Reg(t *testing.T) {
	New().
		RegMaskAnyFunc(MaskTypeZero, MaskZero).
		RegMaskStringFunc(MaskTypeChar, MaskCharString()).
		RegMaskStringFunc(MaskTypeRandom, MaskRandString).
		RegMaskStringFunc(MaskTypeHash, MaskHashString).
		RegMaskIntFunc(MaskTypeRandom, MaskRandInt).
		RegMaskFloat64Func(MaskTypeRandom, MaskRandFloat64).
		RegMaskUintFunc(MaskTypeRandom, MaskRandUint)
}

func TestMasker_Mask(t *testing.T) {
	demo := testStruct{
		StrZero:   "This field will be a empty string",
		IntZero:   57128, // This field will be 0
		UintZero:  57128, // This field will be 0
		FloatZero: 57128, // This field will be 0
		AnyZero:   "This field will be empty",
	}
	demoMasked, err := New().RegMaskAnyFunc(MaskTypeZero, MaskZero).Mask(demo)
	assert.NoError(t, err)
	masked := demoMasked.(testStruct)
	assert.Zero(t, masked.StrZero)
	assert.Zero(t, masked.IntZero)
	assert.Zero(t, masked.UintZero)
	assert.Zero(t, masked.FloatZero)
	assert.Zero(t, masked.AnyZero)

	// check the original demo didn't change
	assert.Equal(t, testStruct{
		StrZero:   "This field will be a empty string",
		IntZero:   57128, // This field will be 0
		UintZero:  57128, // This field will be 0
		FloatZero: 57128, // This field will be 0
		AnyZero:   "This field will be empty",
	}, demo)
}

func TestMasker_Float64(t *testing.T) {
	masked, err := New().
		RegMaskAnyFunc(MaskTypeZero, MaskZero).
		Float64(57128, MaskTypeZero)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), masked)
}

func TestMasker_String(t *testing.T) {
	masked, err := New().
		RegMaskAnyFunc(MaskTypeZero, MaskZero).
		String("This field will be a empty string", MaskTypeZero)
	assert.NoError(t, err)
	assert.Equal(t, "", masked)
}

func TestMasker_Int(t *testing.T) {
	masked, err := New().
		RegMaskAnyFunc(MaskTypeZero, MaskZero).
		Int(57128, MaskTypeZero)
	assert.NoError(t, err)
	assert.Equal(t, 0, masked)
}

func TestMasker_Uint(t *testing.T) {
	masked, err := New().
		RegMaskAnyFunc(MaskTypeZero, MaskZero).
		Uint(57128, MaskTypeZero)
	assert.NoError(t, err)
	assert.Equal(t, uint(0), masked)
}

func TestMasker_Any(t *testing.T) {
	// string
	_, masked, err := New().
		RegMaskAnyFunc(MaskTypeZero, MaskZero).
		Any("string", MaskTypeZero)
	assert.NoError(t, err)
	assert.Equal(t, "", masked)

	// int
	_, masked, err = New().
		RegMaskAnyFunc(MaskTypeZero, MaskZero).
		Any(57128, MaskTypeZero)
	assert.NoError(t, err)
	assert.Equal(t, 0, masked)

	// uint
	_, masked, err = New().
		RegMaskAnyFunc(MaskTypeZero, MaskZero).
		Any(uint(57128), MaskTypeZero)
	assert.NoError(t, err)
	assert.Equal(t, uint(0), masked)

	// float64
	_, masked, err = New().
		RegMaskAnyFunc(MaskTypeZero, MaskZero).
		Any(float64(57128), MaskTypeZero)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), masked)

	// any
	_, masked, err = New().
		RegMaskAnyFunc(MaskTypeZero, MaskZero).
		Any(any("string"), MaskTypeZero)
	assert.NoError(t, err)
	assert.Equal(t, any(""), masked)

	// map
	_, masked, err = New().
		RegMaskAnyFunc(MaskTypeZero, MaskZero).
		Any(map[string]string{"foo": "bar"}, MaskTypeZero)
	assert.NoError(t, err)
	assert.Equal(t, (map[string]string)(nil), masked)
}
