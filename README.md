# gmask
gmask is a Go library for masking sensitive information in variable.

## Getting gmask

With [Go module](https://github.com/golang/go/wiki/Modules) support, simply add the following import

```go
import "github.com/asjdf/gmask"
```

to your code, and then `go [build|run|test]` will automatically fetch the necessary dependencies.

Otherwise, run the following Go command to install the gmask package:

```shell
go get -u github.com/asjdf/gmask
```

## Demo

You can get more demo in [gmask_test.go](./gmask_test.go)

```go
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
```

output:

```shell
{PubInfo:This field won't be masked StrChar:This field will be ******** StrChar3:This field will be *** StrCharSameLen:This field will be same length * StrCharDash:This field will be -------- StrRand:This field will be a random string StrHash:This field will be a hash string StrZero:This field will be a empty string}
{PubInfo:This field won't be masked StrChar:******** StrChar3:*** StrCharSameLen:******************************** StrCharDash:-------- StrRand:vOvMXGbu StrHash:469e0aae1b45c13042c0f95e4a5bea77a2696bd9b7d8694a6023f1ad1b3479f6 StrZero:}
```

## Available tags

## How to Contribute

If you are interested in this project, you can contribute in the following ways:

- Submitting issues: If you find any problems or have any suggestions, please submit an issue on GitHub.
- Submitting pull requests: If you have solved a problem or implemented a new feature, please submit a pull request.
- Improving documentation: We greatly need contributors to improve the documentation. If you find any deficiencies or errors in the documentation, please submit a pull request.

We appreciate any contributions to this project, and we thank you for your support!