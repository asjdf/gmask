package gmask

import (
	"fmt"
	"testing"
)

func TestExample(t *testing.T) {
	demo := struct {
		PubInfo        string `json:"pub"`
		StrChar        string `json:"secretChar" mask:"char"`
		StrChar3       string `json:"secretChar8" mask:"char,3"`
		StrCharSameLen string `json:"secretCharSameLen" mask:"char,-1"`
		StrCharDash    string `json:"secretCharSameStr" mask:"char,,-"`
		StrRand        string `json:"secretRand" mask:"rand"`
		StrHash        string `json:"secretHash" mask:"hash"`
		StrZero        string `json:"secretZero" mask:"zero"`
	}{
		PubInfo:        "This field won't be masked",
		StrChar:        "This field will be ********",
		StrChar3:       "This field will be ***",
		StrCharSameLen: "This field will be same length *",
		StrCharDash:    "This field will be --------",
		StrRand:        "This field will be a random string",
		StrHash:        "This field will be a hash string",
		StrZero:        "This field will be a empty string",
	}
	fmt.Printf("%+v\n", demo)
	demoMasked, _ := Mask(demo)
	fmt.Printf("%+v\n", demoMasked)
}
